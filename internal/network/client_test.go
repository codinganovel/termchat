package network

import (
	"testing"
)

func TestParseConnectionString(t *testing.T) {
	tests := []struct {
		input    string
		wantUser string
		wantHost string
		wantSess string
		wantPort int
		wantErr  bool
	}{
		{
			input:    "alice@localhost:cosmic-turtle-123",
			wantUser: "alice",
			wantHost: "localhost",
			wantSess: "cosmic-turtle-123",
			wantPort: 9999,
			wantErr:  false,
		},
		{
			input:    "bob@192.168.1.1:mystic-dragon-456:8080",
			wantUser: "bob",
			wantHost: "192.168.1.1",
			wantSess: "mystic-dragon-456",
			wantPort: 8080,
			wantErr:  false,
		},
		{
			input:    "user@host.com:session-id-789:22222",
			wantUser: "user",
			wantHost: "host.com",
			wantSess: "session-id-789",
			wantPort: 22222,
			wantErr:  false,
		},
		{
			input:   "invalid-format",
			wantErr: true,
		},
		{
			input:   "missing@host",
			wantErr: true,
		},
		{
			input:   "missing:session",
			wantErr: true,
		},
		{
			input:   "@host:session",
			wantErr: true,
		},
		{
			input:   "user@:session",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			info, err := ParseConnectionString(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseConnectionString(%q) expected error, got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("ParseConnectionString(%q) unexpected error: %v", tt.input, err)
				return
			}
			
			if info.User != tt.wantUser {
				t.Errorf("User = %q, want %q", info.User, tt.wantUser)
			}
			
			if info.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", info.Host, tt.wantHost)
			}
			
			if info.SessionID != tt.wantSess {
				t.Errorf("SessionID = %q, want %q", info.SessionID, tt.wantSess)
			}
			
			if info.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", info.Port, tt.wantPort)
			}
		})
	}
}

func TestConnectionInfoValidation(t *testing.T) {
	// Test empty user
	info, err := ParseConnectionString("@host:session")
	if err == nil {
		t.Error("Expected error for empty user")
	}
	
	// Test empty host
	info, err = ParseConnectionString("user@:session")
	if err == nil {
		t.Error("Expected error for empty host")
	}
	
	// Valid case
	info, err = ParseConnectionString("user@host:session")
	if err != nil {
		t.Errorf("Unexpected error for valid input: %v", err)
	}
	
	if info == nil {
		t.Error("Expected non-nil ConnectionInfo")
	}
}