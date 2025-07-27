package session

import (
	"testing"
	
	"github.com/sam/termchat/pkg/protocol"
)

func BenchmarkGenerateSessionID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateSessionID()
	}
}

func BenchmarkNewSession(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkAddMessage(b *testing.B) {
	s := New()
	msg := protocol.Message{
		Type:    protocol.MessageTypeText,
		Content: "This is a test message for benchmarking purposes",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.AddMessage(msg)
	}
}

func BenchmarkGetMessages(b *testing.B) {
	s := New()
	// Add some messages
	for i := 0; i < 100; i++ {
		s.AddMessage(protocol.Message{
			Type:    protocol.MessageTypeText,
			Content: "Test message",
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.GetMessages()
	}
}

func BenchmarkConcurrentAddMessage(b *testing.B) {
	s := New()
	msg := protocol.Message{
		Type:    protocol.MessageTypeText,
		Content: "Concurrent test message",
	}
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.AddMessage(msg)
		}
	})
}

func BenchmarkSessionMemoryUsage(b *testing.B) {
	b.Run("1000Messages", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := New()
			for j := 0; j < 1000; j++ {
				s.AddMessage(protocol.Message{
					Type:    protocol.MessageTypeText,
					Content: "This is a longer message to simulate real chat content with more bytes",
				})
			}
		}
	})
	
	b.Run("10000Messages", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := New()
			for j := 0; j < 10000; j++ {
				s.AddMessage(protocol.Message{
					Type:    protocol.MessageTypeText,
					Content: "Message content",
				})
			}
		}
	})
}