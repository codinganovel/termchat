package main

import (
	"testing"
	"time"
	
	"github.com/sam/termchat/internal/session"
	"github.com/sam/termchat/pkg/protocol"
)

// BenchmarkChatSession simulates a complete chat session
func BenchmarkChatSession(b *testing.B) {
	b.Run("ShortConversation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sess := session.New()
			
			// Simulate 20 message exchanges
			for j := 0; j < 20; j++ {
				msg := protocol.NewMessage(protocol.MessageTypeText, "Hello, how are you?")
				sess.AddMessage(*msg)
				
				reply := protocol.NewMessage(protocol.MessageTypeText, "I'm good, thanks!")
				sess.AddMessage(*reply)
			}
			
			// Get all messages
			_ = sess.GetMessages()
		}
	})
	
	b.Run("LongConversation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sess := session.New()
			
			// Simulate 500 message exchanges
			for j := 0; j < 500; j++ {
				msg := protocol.NewMessage(protocol.MessageTypeText, "This is message number "+string(rune(j)))
				sess.AddMessage(*msg)
			}
			
			// Get messages multiple times (simulating scrolling)
			for k := 0; k < 10; k++ {
				_ = sess.GetMessages()
			}
		}
	})
}

// BenchmarkMessageThroughput tests how many messages per second can be processed
func BenchmarkMessageThroughput(b *testing.B) {
	sess := session.New()
	msg := protocol.NewMessage(protocol.MessageTypeText, "Throughput test message")
	
	b.ResetTimer()
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		sess.AddMessage(*msg)
	}
	
	elapsed := time.Since(start)
	messagesPerSecond := float64(b.N) / elapsed.Seconds()
	b.ReportMetric(messagesPerSecond, "msgs/sec")
}

// BenchmarkMemoryPerMessage measures memory allocation per message
func BenchmarkMemoryPerMessage(b *testing.B) {
	b.ReportAllocs()
	
	sess := session.New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := protocol.NewMessage(protocol.MessageTypeText, "Memory benchmark message with some content")
		sess.AddMessage(*msg)
	}
}