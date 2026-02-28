package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client untuk koneksi SSH ke OLT ZTE
type Client struct {
	host     string
	port     int
	username string
	password string
	client   *ssh.Client
	session  *ssh.Session
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
		cfg.Port = 22
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

// Connect melakukan koneksi SSH ke OLT
func (c *Client) Connect() error {
	config := &ssh.ClientConfig{
		User: c.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}

	c.client = client
	return nil
}

// Close menutup koneksi SSH
func (c *Client) Close() error {
	if c.session != nil {
		c.session.Close()
	}
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Execute menjalankan command dan mengembalikan output
func (c *Client) Execute(ctx context.Context, cmd string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Channel untuk output
	outputChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		output, err := session.CombinedOutput(cmd)
		if err != nil {
			errorChan <- err
			return
		}
		outputChan <- string(output)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case output := <-outputChan:
		return output, nil
	case err := <-errorChan:
		return "", err
	}
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
	}
	return results, nil
}

// InteractiveSession untuk session interaktif
type InteractiveSession struct {
	client   *Client
	session  *ssh.Session
	stdin    io.WriteCloser
	stdout   io.Reader
	stderr   io.Reader
	prompt   string
	enabled  bool
}

// NewInteractiveSession membuat session interaktif
func (c *Client) NewInteractiveSession() (*InteractiveSession, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}

	// Start shell
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		return nil, err
	}

	if err := session.Shell(); err != nil {
		return nil, err
	}

	return &InteractiveSession{
		client:  c,
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		prompt:  "ZXAN>",
	}, nil
}

// Close menutup session interaktif
func (s *InteractiveSession) Close() error {
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}

// Enable masuk ke mode privilege
func (s *InteractiveSession) Enable(password string) error {
	s.stdin.Write([]byte("enable\n"))
	time.Sleep(100 * time.Millisecond)
	s.stdin.Write([]byte(password + "\n"))
	time.Sleep(200 * time.Millisecond)
	s.prompt = "ZXAN#"
	s.enabled = true
	return nil
}

// ConfigureTerminal masuk ke mode konfigurasi
func (s *InteractiveSession) ConfigureTerminal() error {
	if !s.enabled {
		return fmt.Errorf("must be in enable mode first")
	}
	s.stdin.Write([]byte("configure terminal\n"))
	time.Sleep(100 * time.Millisecond)
	s.prompt = "ZXAN(config)#"
	return nil
}

// SendCommand mengirim command dan membaca output
func (s *InteractiveSession) SendCommand(cmd string) (string, error) {
	// Clear buffer
	buf := make([]byte, 4096)
	s.stdout.Read(buf)

	// Send command
	s.stdin.Write([]byte(cmd + "\n"))
	time.Sleep(200 * time.Millisecond)

	// Read output
	outputBuf := new(bytes.Buffer)
	for {
		n, err := s.stdout.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		outputBuf.Write(buf[:n])
		if strings.Contains(outputBuf.String(), s.prompt) {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return s.cleanOutput(outputBuf.String()), nil
}

// cleanOutput membersihkan output dari ANSI codes dan prompt
func (s *InteractiveSession) cleanOutput(output string) string {
	// Remove ANSI escape codes
	ansi := "\x1b\\[[0-9;]*[a-zA-Z]"
	output = strings.ReplaceAll(output, ansi, "")

	// Remove command echo and prompt
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ZXAN") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// IsConnected mengecek apakah client terhubung
func (c *Client) IsConnected() bool {
	if c.client == nil {
		return false
	}
	_, _, err := c.client.SendRequest("keepalive@golang.org", true, nil)
	return err == nil
}

// TestConnection mengetes koneksi ke OLT
func TestConnection(cfg Config) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	conn.Close()
	return nil
}
