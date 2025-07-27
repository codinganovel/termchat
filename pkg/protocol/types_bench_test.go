package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func BenchmarkNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewMessage(MessageTypeText, "Hello, world!")
	}
}

func BenchmarkNewHandshakeMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewHandshakeMessage(MessageTypeHello, "test-session-123")
	}
}

func BenchmarkMessageMarshal(b *testing.B) {
	msg := &Message{
		Type:      MessageTypeText,
		Content:   "This is a test message for benchmarking JSON marshaling performance",
		SessionID: "cosmic-turtle-1234",
		Timestamp: time.Now().UnixMilli(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(msg)
	}
}

func BenchmarkMessageUnmarshal(b *testing.B) {
	msg := &Message{
		Type:      MessageTypeText,
		Content:   "This is a test message for benchmarking JSON unmarshaling performance",
		SessionID: "cosmic-turtle-1234",
		Timestamp: time.Now().UnixMilli(),
	}
	
	data, _ := json.Marshal(msg)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded Message
		_ = json.Unmarshal(data, &decoded)
	}
}

func BenchmarkMessageSizes(b *testing.B) {
	b.Run("SmallMessage", func(b *testing.B) {
		msg := &Message{
			Type:      MessageTypeText,
			Content:   "Hi",
			Timestamp: time.Now().UnixMilli(),
		}
		
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(msg)
		}
	})
	
	b.Run("MediumMessage", func(b *testing.B) {
		msg := &Message{
			Type:      MessageTypeText,
			Content:   "This is a medium-sized message that represents typical chat content",
			Timestamp: time.Now().UnixMilli(),
		}
		
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(msg)
		}
	})
	
	b.Run("LargeMessage", func(b *testing.B) {
		content := ""
		for j := 0; j < 100; j++ {
			content += "This is a large message to test performance with bigger payloads. "
		}
		
		msg := &Message{
			Type:      MessageTypeText,
			Content:   content,
			Timestamp: time.Now().UnixMilli(),
		}
		
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(msg)
		}
	})
}

func BenchmarkRoundTrip(b *testing.B) {
	msg := &Message{
		Type:      MessageTypeText,
		Content:   "Round trip benchmark message",
		SessionID: "test-123",
		Timestamp: time.Now().UnixMilli(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(msg)
		var decoded Message
		_ = json.Unmarshal(data, &decoded)
	}
}