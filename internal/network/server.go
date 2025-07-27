package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sam/termchat/internal/session"
	"github.com/sam/termchat/pkg/protocol"
)

type Server struct {
	session    *session.Session
	listener   net.Listener
	conn       net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	mu         sync.Mutex
	
	onMessage  func(protocol.Message)
	onConnect  func()
	onDisconnect func()
}

func NewServer(sess *session.Session) *Server {
	return &Server{
		session: sess,
	}
}

func (s *Server) SetCallbacks(onMessage func(protocol.Message), onConnect, onDisconnect func()) {
	s.onMessage = onMessage
	s.onConnect = onConnect
	s.onDisconnect = onDisconnect
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	
	s.listener = listener
	s.session.SetState(session.StateWaiting)
	
	go s.acceptConnections()
	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		
		s.mu.Lock()
		if s.conn != nil {
			conn.Close()
			s.mu.Unlock()
			continue
		}
		
		s.conn = conn
		s.encoder = json.NewEncoder(conn)
		s.decoder = json.NewDecoder(conn)
		s.mu.Unlock()
		
		go s.handleConnection()
	}
}

func (s *Server) handleConnection() {
	defer func() {
		s.mu.Lock()
		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
		}
		s.mu.Unlock()
		
		if s.onDisconnect != nil {
			s.onDisconnect()
		}
	}()
	
	if err := s.performHandshake(); err != nil {
		return
	}
	
	s.session.SetState(session.StateActive)
	if s.onConnect != nil {
		s.onConnect()
	}
	
	for {
		var msg protocol.Message
		if err := s.decoder.Decode(&msg); err != nil {
			return
		}
		
		s.session.AddMessage(msg)
		
		if s.onMessage != nil {
			s.onMessage(msg)
		}
		
		switch msg.Type {
		case protocol.MessageTypePing:
			s.SendMessage(protocol.NewMessage(protocol.MessageTypePong, ""))
		case protocol.MessageTypeLeave:
			return
		}
	}
}

func (s *Server) performHandshake() error {
	var hello protocol.Message
	if err := s.decoder.Decode(&hello); err != nil {
		return err
	}
	
	if hello.Type != protocol.MessageTypeHello {
		s.sendError("Expected HELLO message")
		return fmt.Errorf("invalid handshake: expected HELLO, got %s", hello.Type)
	}
	
	if hello.SessionID != s.session.ID {
		s.sendError("Session ID mismatch")
		return fmt.Errorf("session ID mismatch: expected %s, got %s", s.session.ID, hello.SessionID)
	}
	
	welcome := protocol.NewHandshakeMessage(protocol.MessageTypeWelcome, s.session.ID)
	if err := s.encoder.Encode(welcome); err != nil {
		return err
	}
	
	var ready protocol.Message
	if err := s.decoder.Decode(&ready); err != nil {
		return err
	}
	
	if ready.Type != protocol.MessageTypeReady {
		return fmt.Errorf("invalid handshake: expected READY, got %s", ready.Type)
	}
	
	return nil
}

func (s *Server) SendMessage(msg *protocol.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.encoder == nil {
		return fmt.Errorf("not connected")
	}
	
	return s.encoder.Encode(msg)
}

func (s *Server) sendError(errMsg string) {
	msg := protocol.NewMessage(protocol.MessageTypeError, errMsg)
	s.encoder.Encode(msg)
}

func (s *Server) Stop() {
	s.mu.Lock()
	
	if s.conn != nil {
		// Send leave message before closing
		if s.encoder != nil {
			msg := protocol.NewMessage(protocol.MessageTypeLeave, "Host disconnected")
			s.encoder.Encode(msg)
		}
		s.conn.Close()
		s.conn = nil
	}
	
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
	
	s.mu.Unlock()
	
	s.session.SetState(session.StateEnded)
}