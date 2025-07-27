package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(MessageTypeText, "Hello, world!")
	
	if msg.Type != MessageTypeText {
		t.Errorf("Expected type %s, got %s", MessageTypeText, msg.Type)
	}
	
	if msg.Content != "Hello, world!" {
		t.Errorf("Expected content 'Hello, world!', got %s", msg.Content)
	}
	
	if msg.Timestamp == 0 {
		t.Error("Timestamp should be set")
	}
	
	// Verify timestamp is recent (within last second)
	now := time.Now().UnixMilli()
	if msg.Timestamp > now || msg.Timestamp < now-1000 {
		t.Error("Timestamp should be recent")
	}
}

func TestNewHandshakeMessage(t *testing.T) {
	sessionID := "test-session-123"
	msg := NewHandshakeMessage(MessageTypeHello, sessionID)
	
	if msg.Type != MessageTypeHello {
		t.Errorf("Expected type %s, got %s", MessageTypeHello, msg.Type)
	}
	
	if msg.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, msg.SessionID)
	}
	
	if msg.Content != "" {
		t.Errorf("Content should be empty for handshake, got %s", msg.Content)
	}
}

func TestMessageSerialization(t *testing.T) {
	original := &Message{
		Type:      MessageTypeText,
		Content:   "Test message",
		SessionID: "test-123",
		Timestamp: time.Now().UnixMilli(),
	}
	
	// Serialize
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	// Deserialize
	var decoded Message
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	
	// Compare
	if decoded.Type != original.Type {
		t.Errorf("Type mismatch: expected %s, got %s", original.Type, decoded.Type)
	}
	
	if decoded.Content != original.Content {
		t.Errorf("Content mismatch: expected %s, got %s", original.Content, decoded.Content)
	}
	
	if decoded.SessionID != original.SessionID {
		t.Errorf("SessionID mismatch: expected %s, got %s", original.SessionID, decoded.SessionID)
	}
	
	if decoded.Timestamp != original.Timestamp {
		t.Errorf("Timestamp mismatch: expected %d, got %d", original.Timestamp, decoded.Timestamp)
	}
}

func TestMessageTypes(t *testing.T) {
	// Verify all message types are distinct
	types := []MessageType{
		MessageTypeHello,
		MessageTypeWelcome,
		MessageTypeReady,
		MessageTypeText,
		MessageTypePing,
		MessageTypePong,
		MessageTypeLeave,
		MessageTypeError,
	}
	
	seen := make(map[MessageType]bool)
	for _, mt := range types {
		if seen[mt] {
			t.Errorf("Duplicate message type: %s", mt)
		}
		seen[mt] = true
	}
	
	// Verify expected values
	if MessageTypeHello != "hello" {
		t.Errorf("MessageTypeHello should be 'hello', got %s", MessageTypeHello)
	}
	
	if MessageTypeText != "text" {
		t.Errorf("MessageTypeText should be 'text', got %s", MessageTypeText)
	}
}

func TestOmitEmptyFields(t *testing.T) {
	msg := &Message{
		Type:      MessageTypeText,
		Timestamp: time.Now().UnixMilli(),
		// Content and SessionID are empty
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	jsonStr := string(data)
	
	// Verify omitempty works
	if contains(jsonStr, "content") {
		t.Error("Empty content field should be omitted")
	}
	
	if contains(jsonStr, "session_id") {
		t.Error("Empty session_id field should be omitted")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}