package snmp

import (
	"context"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Client wraps gosnmp client
type Client struct {
	client *gosnmp.GoSNMP
}

// Config represents SNMP client configuration
type Config struct {
	Host      string
	Port      uint16
	Community string
	Timeout   time.Duration
	Retries   int
	MaxOids   int
}

// NewClient creates a new SNMP client
func NewClient(cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.Retries == 0 {
		cfg.Retries = 2
	}
	if cfg.MaxOids == 0 {
		cfg.MaxOids = 60
	}

	client := &gosnmp.GoSNMP{
		Target:    cfg.Host,
		Port:      cfg.Port,
		Community: cfg.Community,
		Version:   gosnmp.Version2c,
		Timeout:   cfg.Timeout,
		Retries:   cfg.Retries,
		MaxOids:   cfg.MaxOids,
	}

	return &Client{client: client}, nil
}

// Connect establishes connection
func (c *Client) Connect() error {
	return c.client.Connect()
}

// Close closes connection
func (c *Client) Close() error {
	if c.client.Conn != nil {
		return c.client.Conn.Close()
	}
	return nil
}

// Get performs SNMP GET
func (c *Client) Get(oids []string) (*gosnmp.SnmpPacket, error) {
	return c.client.Get(oids)
}

// Walk performs SNMP WALK
func (c *Client) Walk(oid string, fn func(gosnmp.SnmpPDU) error) error {
	return c.client.Walk(oid, fn)
}

// GetWithContext performs SNMP GET with context
func (c *Client) GetWithContext(ctx context.Context, oids []string) (*gosnmp.SnmpPacket, error) {
	// gosnmp doesn't support context natively, so we use a goroutine
	type result struct {
		pkt *gosnmp.SnmpPacket
		err error
	}

	resultCh := make(chan result, 1)
	go func() {
		pkt, err := c.client.Get(oids)
		resultCh <- result{pkt, err}
	}()

	select {
	case res := <-resultCh:
		return res.pkt, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// WalkWithContext performs SNMP WALK with context
func (c *Client) WalkWithContext(ctx context.Context, oid string, fn func(gosnmp.SnmpPDU) error) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- c.client.Walk(oid, fn)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
