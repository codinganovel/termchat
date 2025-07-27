package session

import (
	"strings"
	"testing"
	
	"github.com/sam/termchat/pkg/protocol"
)

func TestGenerateSessionID(t *testing.T) {
	// Test format
	id := GenerateSessionID()
	parts := strings.Split(id, "-")
	
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts in session ID, got %d: %s", len(parts), id)
	}
	
	// Test uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateSessionID()
		if ids[id] {
			t.Errorf("Duplicate session ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestNewSession(t *testing.T) {
	s := New()
	
	if s.ID == "" {
		t.Error("Session ID should not be empty")
	}
	
	if s.StartTime.IsZero() {
		t.Error("StartTime should be set")
	}
	
	if s.State != StateCreated {
		t.Errorf("Initial state should be StateCreated, got %v", s.State)
	}
	
	if len(s.Messages) != 0 {
		t.Error("Messages should be empty initially")
	}
}

func TestSessionState(t *testing.T) {
	s := New()
	
	// Test state transitions
	s.SetState(StateWaiting)
	if s.GetState() != StateWaiting {
		t.Errorf("Expected StateWaiting, got %v", s.GetState())
	}
	
	s.SetState(StateActive)
	if s.GetState() != StateActive {
		t.Errorf("Expected StateActive, got %v", s.GetState())
	}
}

func TestAddMessage(t *testing.T) {
	s := New()
	
	msg := protocol.Message{
		Type:    protocol.MessageTypeText,
		Content: "Hello, world!",
	}
	
	s.AddMessage(msg)
	
	messages := s.GetMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
	
	if messages[0].Content != "Hello, world!" {
		t.Errorf("Expected message content 'Hello, world!', got %s", messages[0].Content)
	}
	
	// Check timestamp was set
	if messages[0].Timestamp == 0 {
		t.Error("Message timestamp should be set")
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	done := make(chan bool)
	
	// Concurrent writes
	go func() {
		for i := 0; i < 100; i++ {
			s.AddMessage(protocol.Message{
				Type:    protocol.MessageTypeText,
				Content: "msg1",
			})
		}
		done <- true
	}()
	
	go func() {
		for i := 0; i < 100; i++ {
			s.AddMessage(protocol.Message{
				Type:    protocol.MessageTypeText,
				Content: "msg2",
			})
		}
		done <- true
	}()
	
	// Concurrent reads
	go func() {
		for i := 0; i < 100; i++ {
			_ = s.GetMessages()
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Verify total messages
	messages := s.GetMessages()
	if len(messages) != 200 {
		t.Errorf("Expected 200 messages, got %d", len(messages))
	}
}