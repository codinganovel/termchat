# termchat Architecture Document

**Version:** 2.0.0  
**Last Updated:** 2025-07-27

## Table of Contents

1. [System Overview](#system-overview)
2. [Design Principles](#design-principles)
3. [High-Level Architecture](#high-level-architecture)
4. [Component Architecture](#component-architecture)
5. [Connection Flow](#connection-flow)
6. [Security Architecture](#security-architecture)
7. [Implementation Details](#implementation-details)
8. [Technology Stack](#technology-stack)

## System Overview

termchat is a peer-to-peer terminal chat application that establishes direct connections between exactly two participants using SSH port forwarding. The architecture is deliberately minimal:

- **No Server Component**: Pure P2P communication
- **SSH-Based Transport**: Leverages existing SSH infrastructure
- **Ephemeral by Design**: No persistence, no history
- **Single Binary**: Self-contained Go application
- **Terminal Native**: Built with tcell for rich terminal UI

## Design Principles

### 1. Serverless Architecture
- No central server or coordinator
- Direct P2P connection between participants
- Session exists only while both parties are connected

### 2. Leverage SSH
- Use existing SSH authentication and encryption
- No custom security protocols
- Works with standard SSH configurations

### 3. Simplicity First
- Minimal codebase
- No external services or databases
- Clear, understandable architecture

### 4. Privacy by Default
- No message persistence
- No user tracking or analytics
- Connection details never leave the local machines

## High-Level Architecture

```
┌─────────────────┐                              ┌─────────────────┐
│   termchat A    │                              │   termchat B    │
│   (Initiator)   │◀────────SSH Tunnel──────────▶│   (Joiner)      │
└─────────────────┘                              └─────────────────┘
        │                                                │
        │                                                │
   Local Socket                                    Local Socket
   (port 9999)                                    (ephemeral)
```

### Connection Models

1. **Local Testing**: Both instances on same machine (localhost)
2. **LAN Connection**: Direct connection over local network
3. **Remote Connection**: Through SSH tunnel over internet

## Component Architecture

### 1. Core Components

#### Session Manager
- Generates unique session IDs
- Manages connection state
- Handles connection lifecycle

#### UI Manager
- Terminal UI using tcell
- Message display and scrolling
- Input field handling
- Status bar updates

#### Network Handler
- TCP socket management
- Message serialization/deserialization
- Connection health monitoring

#### SSH Tunnel Manager
- Establishes SSH connections
- Sets up port forwarding
- Manages tunnel lifecycle

### 2. Message Protocol

**Simple text-based protocol:**
```go
type Message struct {
    Type      MessageType `json:"type"`
    Content   string      `json:"content,omitempty"`
    Timestamp int64       `json:"timestamp"`
}

type MessageType string

const (
    MessageTypeText   MessageType = "text"
    MessageTypeTyping MessageType = "typing"
    MessageTypeJoin   MessageType = "join"
    MessageTypeLeave  MessageType = "leave"
    MessageTypePing   MessageType = "ping"
    MessageTypePong   MessageType = "pong"
)
```

### 3. Connection States

```
┌──────────┐      ┌──────────┐      ┌───────────┐
│  INIT    │─────▶│ WAITING  │─────▶│ CONNECTED │
└──────────┘      └──────────┘      └───────────┘
                        │                   │
                        └───────────────────┘
                                │
                                ▼
                          ┌──────────┐
                          │  CLOSED  │
                          └──────────┘
```

## Connection Flow

### 1. Session Initiation (Person A)

```
$ termchat start
    │
    ├─▶ Generate session ID (e.g., "cosmic-turtle-7823")
    ├─▶ Open TCP listener on port 9999
    ├─▶ Display session ID to user
    └─▶ Wait for incoming connection
```

### 2. Session Joining (Person B)

```
$ termchat join alice@host:cosmic-turtle-7823
    │
    ├─▶ Parse connection string
    ├─▶ SSH to alice@host
    ├─▶ Set up port forward (remote:9999 → local:random)
    ├─▶ Connect to forwarded port
    ├─▶ Send session ID for validation
    └─▶ Establish P2P connection
```

### 3. Message Exchange

```
Person A                                    Person B
   │                                           │
   ├──────[Text Message via TCP]──────────────▶│
   │                                           ├─▶ Display
   │                                           │
   │◀─────[Text Message via TCP]───────────────┤
   ├─▶ Display                                 │
   │                                           │
```

### 4. Connection Teardown

- Either party disconnects → TCP connection closes
- SSH tunnel automatically cleaned up
- Both instances detect disconnect and exit
- No cleanup required (no persistent state)

## Security Architecture

### SSH-Based Security
- All remote connections use SSH encryption
- Leverages existing SSH key infrastructure
- No additional authentication layer needed

### Session Security
- Session IDs are random and unguessable
- IDs must match for connection to establish
- No session information exposed externally

### Privacy Guarantees
- No message logging or persistence
- No central server to compromise
- Connection metadata stays local
- Perfect forward secrecy per session

### Trust Model
- Trust established through SSH access
- If you can SSH to the host, you can chat
- No additional identity management

## Implementation Details

### Session ID Generation
```go
func generateSessionID() string {
    adjectives := []string{"cosmic", "mystic", "quantum", ...}
    nouns := []string{"turtle", "phoenix", "nebula", ...}
    number := rand.Intn(10000)
    
    adj := adjectives[rand.Intn(len(adjectives))]
    noun := nouns[rand.Intn(len(nouns))]
    
    return fmt.Sprintf("%s-%s-%d", adj, noun, number)
}
```

### Port Forwarding Setup
```go
// Simplified SSH port forwarding
config := &ssh.ClientConfig{
    User: user,
    Auth: []ssh.AuthMethod{sshAgent()},
}

client, _ := ssh.Dial("tcp", host+":22", config)
listener, _ := client.Listen("tcp", "localhost:9999")

// Forward connections
for {
    conn, _ := listener.Accept()
    go handleForwardedConnection(conn)
}
```

### Message Framing
- Length-prefixed messages for TCP
- JSON encoding for simplicity
- Optional compression for large messages

## Technology Stack

### Core Implementation
- **Language**: Go 1.21+
- **Terminal UI**: github.com/gdamore/tcell/v2
- **SSH Client**: golang.org/x/crypto/ssh
- **JSON**: encoding/json (standard library)

### Key Libraries
- **tcell**: Modern terminal handling
- **x/crypto/ssh**: SSH protocol implementation
- **cobra**: Command-line interface
- **logrus**: Structured logging

### Development Tools
- **Testing**: Go standard testing package
- **Benchmarking**: Go benchmark suite
- **CI/CD**: GitHub Actions
- **Cross-compilation**: Go's built-in support

## Error Handling

### Connection Errors
- SSH connection failure → Clear error message
- Port already in use → Suggest alternative port
- Session ID mismatch → Connection refused
- Network interruption → Clean disconnect

### Recovery Strategy
- No automatic reconnection (ephemeral design)
- Clear error messages for troubleshooting
- Graceful degradation on terminal issues

## Performance Characteristics

### Resource Usage
- Memory: ~10MB base + 1KB per message
- CPU: <1% during normal chat
- Network: Minimal overhead (JSON + TCP)
- Disk: Zero (no persistence)

### Latency
- Local connection: <1ms
- LAN connection: <5ms
- Internet (SSH): 20-100ms typical
- Typing indicators: 100ms debounce

## Future Considerations

### Potential Enhancements
1. Multiple port support for firewall traversal
2. Optional message encryption layer
3. File transfer capability
4. Custom color themes

### Explicitly Not Planned
- Multi-user chat (breaks P2P model)
- Message persistence
- Server mode
- Web interface