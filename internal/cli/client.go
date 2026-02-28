package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// Client untuk koneksi Telnet ke OLT ZTE
type Client struct {
	host     string
	port     int
	username string
	password string
	conn     net.Conn
}

// Config untuk koneksi CLI
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// New membuat client CLI baru
func New(cfg Config) *Client {
	if cfg.Port == 0 {
		cfg.Port = 23 // Telnet default
	}
	if cfg.Username == "" {
		cfg.Username = "zte"
	}
	if cfg.Password == "" {
		cfg.Password = "zte"
	}
	return &Client{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
	}
}

// Connect melakukan koneksi Telnet ke OLT
func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("telnet connection failed: %w", err)
	}
	c.conn = conn

	// Wait for login prompt
	time.Sleep(200 * time.Millisecond)
	c.readUntil("Username:")

	// Send username
	c.send(c.username)
	time.Sleep(100 * time.Millisecond)
	c.readUntil("Password:")

	// Send password
	c.send(c.password)
	time.Sleep(200 * time.Millisecond)

	// Check if logged in (wait for prompt)
	output, _ := c.readUntil("ZXAN>", "ZXAN#")
	if strings.Contains(output, "Login invalid") || strings.Contains(output, "Access denied") {
		c.conn.Close()
		return fmt.Errorf("authentication failed")
	}

	// Enter enable mode if needed
	if strings.Contains(output, "ZXAN>") {
		c.send("enable")
		time.Sleep(100 * time.Millisecond)
		c.readUntil("Password:")
		c.send("zxr10") // Default enable password
		time.Sleep(100 * time.Millisecond)
		c.readUntil("ZXAN#")
	}

	return nil
}

// Close menutup koneksi Telnet
func (c *Client) Close() error {
	if c.conn != nil {
		c.send("exit")
		time.Sleep(50 * time.Millisecond)
		return c.conn.Close()
	}
	return nil
}

// send mengirim data ke koneksi
func (c *Client) send(data string) error {
	_, err := c.conn.Write([]byte(data + "\r\n"))
	return err
}

// readUntil membaca sampai menemukan salah satu pattern
func (c *Client) readUntil(patterns ...string) (string, error) {
	buf := make([]byte, 4096)
	result := new(bytes.Buffer)
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			return result.String(), fmt.Errorf("timeout waiting for prompt")
		default:
			c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, err := c.conn.Read(buf)
			if err != nil && err != io.EOF {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Check if we have one of the patterns
					output := result.String()
					for _, p := range patterns {
						if strings.Contains(output, p) {
							return output, nil
						}
					}
					continue
				}
				return result.String(), err
			}
			if n > 0 {
				result.Write(buf[:n])
				output := result.String()
				for _, p := range patterns {
					if strings.Contains(output, p) {
						return output, nil
					}
				}
			}
		}
	}
}

// Execute menjalankan command dan mengembalikan output
func (c *Client) Execute(ctx context.Context, cmd string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected")
	}

	// Clear buffer first
	buf := make([]byte, 4096)
	c.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	c.conn.Read(buf)

	// Send command
	if err := c.send(cmd); err != nil {
		return "", err
	}

	// Wait for response
	time.Sleep(200 * time.Millisecond)
	output, err := c.readUntil("ZXAN#", "ZXAN(config)#", "ZXAN(config-if)#", "ZXAN(gpon-onu-mng)#")
	if err != nil {
		return "", err
	}

	// Clean output
	return c.cleanOutput(output, cmd), nil
}

// cleanOutput membersihkan output dari command echo dan prompt
func (c *Client) cleanOutput(output, cmd string) string {
	lines := strings.Split(output, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip command echo
		if line == cmd {
			continue
		}

		// Skip prompts
		if strings.HasPrefix(line, "ZXAN") || strings.HasSuffix(line, "#") || strings.HasSuffix(line, ">") {
			continue
		}

		// Skip ANSI escape sequences
		if strings.Contains(line, "\x1b") {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// ExecuteMultiple menjalankan beberapa command secara berurutan
func (c *Client) ExecuteMultiple(ctx context.Context, commands []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, cmd := range commands {
		output, err := c.Execute(ctx, cmd)
		if err != nil {
			results[cmd] = fmt.Sprintf("ERROR: %v", err)
		} else {
			results[cmd] = output
		}
		time.Sleep(50 * time.Millisecond)
	}
	return results, nil
}

// IsConnected mengecek apakah client terhubung
func (c *Client) IsConnected() bool {
	if c.conn == nil {
		return false
	}
	_, err := c.conn.Write([]byte{})
	return err == nil
}

// TestConnection mengetes koneksi ke OLT
func TestConnection(cfg Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	conn.Close()
	return nil
}

// ConfigureTerminal masuk ke mode konfigurasi global
func (c *Client) ConfigureTerminal() error {
	_, err := c.Execute(context.Background(), "configure terminal")
	return err
}

// ExitConfig keluar dari mode konfigurasi
func (c *Client) ExitConfig() error {
	_, err := c.Execute(context.Background(), "exit")
	return err
}
