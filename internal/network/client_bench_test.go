package network

import (
	"testing"
)

func BenchmarkParseConnectionString(b *testing.B) {
	testStrings := []string{
		"alice@localhost:cosmic-turtle-123",
		"bob@192.168.1.100:mystic-dragon-456",
		"charlie@remote-server.example.com:stellar-phoenix-789",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseConnectionString(testStrings[i%len(testStrings)])
	}
}

func BenchmarkParseConnectionStringShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseConnectionString("u@h:s")
	}
}

func BenchmarkParseConnectionStringLong(b *testing.B) {
	longString := "verylongusername@very-long-hostname.with.multiple.subdomains.example.com:extremely-verbose-session-identifier-12345"
	
	for i := 0; i < b.N; i++ {
		_, _ = ParseConnectionString(longString)
	}
}