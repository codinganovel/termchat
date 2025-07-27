package session

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sam/termchat/pkg/protocol"
)

type State int

const (
	StateCreated State = iota
	StateWaiting
	StateActive
	StateEnded
)

type Session struct {
	ID        string
	StartTime time.Time
	Messages  []protocol.Message
	State     State
	mu        sync.RWMutex
}

var (
	adjectives = []string{
		"cosmic", "mystic", "stellar", "quantum",
		"cyber", "neon", "turbo", "ultra",
		"alpha", "omega", "shadow", "ghost",
		"plasma", "lunar", "solar", "astral",
	}
	
	nouns = []string{
		"phoenix", "dragon", "tiger", "eagle",
		"shark", "wolf", "bear", "fox",
		"turtle", "falcon", "raven", "nebula",
		"comet", "meteor", "galaxy", "nova",
	}
)

// No need for init() - Go 1.20+ automatically seeds rand

func GenerateSessionID() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	num := rand.Intn(10000)
	return fmt.Sprintf("%s-%s-%d", adj, noun, num)
}

func New() *Session {
	return &Session{
		ID:        GenerateSessionID(),
		StartTime: time.Now(),
		Messages:  make([]protocol.Message, 0),
		State:     StateCreated,
	}
}

func (s *Session) AddMessage(msg protocol.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msg.Timestamp = time.Now().UnixMilli()
	s.Messages = append(s.Messages, msg)
}

func (s *Session) GetMessages() []protocol.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return append([]protocol.Message{}, s.Messages...)
}

func (s *Session) SetState(state State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State = state
}

func (s *Session) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.State
}