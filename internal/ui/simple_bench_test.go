package ui

import (
	"testing"
)

func BenchmarkWrapText(b *testing.B) {
	tests := []struct {
		name  string
		text  string
		width int
	}{
		{
			name:  "ShortText",
			text:  "Hello world",
			width: 80,
		},
		{
			name:  "LongSingleWord",
			text:  "supercalifragilisticexpialidocious",
			width: 20,
		},
		{
			name:  "Paragraph",
			text:  "This is a longer paragraph of text that needs to be wrapped across multiple lines to fit within the specified width constraint.",
			width: 40,
		},
		{
			name:  "VeryLongText",
			text:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
			width: 50,
		},
	}
	
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = wrapText(tt.text, tt.width)
			}
		})
	}
}

func BenchmarkChatMsgCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ChatMsg{
			Content: "This is a test message",
			FromMe:  i%2 == 0,
		}
	}
}

func BenchmarkMessageFormatting(b *testing.B) {
	messages := []ChatMsg{
		{Content: "Short message", FromMe: true},
		{Content: "This is a medium length message for testing", FromMe: false},
		{Content: "This is a very long message that would need to be wrapped across multiple lines in a typical terminal window width", FromMe: true},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := messages[i%len(messages)]
		var displayText string
		if msg.FromMe {
			displayText = "you: " + msg.Content
		} else {
			displayText = "peer: " + msg.Content
		}
		_ = displayText
	}
}