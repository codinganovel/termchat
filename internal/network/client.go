package network

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/sam/termchat/internal/session"
	"github.com/sam/termchat/pkg/protocol"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	session     *session.Session
	conn        net.Conn
	encoder     *json.Encoder
	decoder     *json.Decoder
	sshClient   *ssh.Client
	mu          sync.Mutex
	
	onMessage    func(protocol.Message)
	onConnect    func()
	onDisconnect func()
}

type ConnectionInfo struct {
	User      string
	Host      string
	SessionID string
	Port      int
}

func NewClient(sess *session.Session) *Client {
	return &Client{
		session: sess,
	}
}

func (c *Client) SetCallbacks(onMessage func(protocol.Message), onConnect, onDisconnect func()) {
	c.onMessage = onMessage
	c.onConnect = onConnect
	c.onDisconnect = onDisconnect
}

func ParseConnectionString(connStr string) (*ConnectionInfo, error) {
	parts := strings.Split(connStr, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("invalid format, expected user@host:session-id or user@host:session-id:port")
	}
	
	userHost := parts[0]
	sessionID := parts[1]
	port := 9999 // default port
	
	// Optional port
	if len(parts) == 3 {
		customPort, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", parts[2])
		}
		if customPort < 1 || customPort > 65535 {
			return nil, fmt.Errorf("port must be between 1 and 65535")
		}
		port = customPort
	}
	
	uhParts := strings.Split(userHost, "@")
	if len(uhParts) != 2 {
		return nil, fmt.Errorf("invalid user@host format")
	}
	
	user := uhParts[0]
	host := uhParts[1]
	
	if user == "" {
		return nil, fmt.Errorf("user cannot be empty")
	}
	
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	
	return &ConnectionInfo{
		User:      user,
		Host:      host,
		SessionID: sessionID,
		Port:      port,
	}, nil
}

func (c *Client) ConnectLocal(addr string, sessionID string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	
	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)
	c.session.ID = sessionID
	
	if err := c.performHandshake(); err != nil {
		c.conn.Close()
		return err
	}
	
	c.session.SetState(session.StateActive)
	if c.onConnect != nil {
		c.onConnect()
	}
	
	go c.handleConnection()
	return nil
}

func (c *Client) ConnectViaSSH(connInfo *ConnectionInfo) error {
	sshConfig := &ssh.ClientConfig{
		User: connInfo.User,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	
	if authMethod := getSSHAuthMethod(); authMethod != nil {
		sshConfig.Auth = append(sshConfig.Auth, authMethod)
	}
	
	sshAddr := fmt.Sprintf("%s:22", connInfo.Host)
	sshClient, err := ssh.Dial("tcp", sshAddr, sshConfig)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}
	
	c.sshClient = sshClient
	
	conn, err := sshClient.Dial("tcp", fmt.Sprintf("localhost:%d", connInfo.Port))
	if err != nil {
		sshClient.Close()
		return fmt.Errorf("failed to connect through SSH tunnel: %w", err)
	}
	
	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)
	c.session.ID = connInfo.SessionID
	
	if err := c.performHandshake(); err != nil {
		c.conn.Close()
		sshClient.Close()
		return err
	}
	
	c.session.SetState(session.StateActive)
	if c.onConnect != nil {
		c.onConnect()
	}
	
	go c.handleConnection()
	return nil
}

func (c *Client) handleConnection() {
	defer func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		if c.sshClient != nil {
			c.sshClient.Close()
			c.sshClient = nil
		}
		c.mu.Unlock()
		
		if c.onDisconnect != nil {
			c.onDisconnect()
		}
	}()
	
	for {
		var msg protocol.Message
		if err := c.decoder.Decode(&msg); err != nil {
			return
		}
		
		c.session.AddMessage(msg)
		
		if c.onMessage != nil {
			c.onMessage(msg)
		}
		
		switch msg.Type {
		case protocol.MessageTypePing:
			c.SendMessage(protocol.NewMessage(protocol.MessageTypePong, ""))
		case protocol.MessageTypeLeave:
			return
		}
	}
}

func (c *Client) performHandshake() error {
	hello := protocol.NewHandshakeMessage(protocol.MessageTypeHello, c.session.ID)
	if err := c.encoder.Encode(hello); err != nil {
		return err
	}
	
	var welcome protocol.Message
	if err := c.decoder.Decode(&welcome); err != nil {
		return err
	}
	
	if welcome.Type == protocol.MessageTypeError {
		return fmt.Errorf("server error: %s", welcome.Content)
	}
	
	if welcome.Type != protocol.MessageTypeWelcome {
		return fmt.Errorf("invalid handshake: expected WELCOME, got %s", welcome.Type)
	}
	
	ready := protocol.NewMessage(protocol.MessageTypeReady, "")
	if err := c.encoder.Encode(ready); err != nil {
		return err
	}
	
	return nil
}

func (c *Client) SendMessage(msg *protocol.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.encoder == nil {
		return fmt.Errorf("not connected")
	}
	
	return c.encoder.Encode(msg)
}

func (c *Client) Stop() {
	c.mu.Lock()
	
	if c.conn != nil {
		// Send leave message before closing
		if c.encoder != nil {
			msg := protocol.NewMessage(protocol.MessageTypeLeave, "Peer disconnected")
			c.encoder.Encode(msg)
		}
		c.conn.Close()
		c.conn = nil
	}
	
	if c.sshClient != nil {
		c.sshClient.Close()
		c.sshClient = nil
	}
	
	c.mu.Unlock()
	
	c.session.SetState(session.StateEnded)
}