# termchat Implementation Plan

**Version:** 2.0.0  
**Last Updated:** 2025-07-27

## Executive Summary

This document outlines the implementation plan for termchat, a serverless P2P terminal chat application. The simplified architecture allows for rapid development with a focus on core functionality.

**Total Timeline**: 4-6 weeks for full implementation  
**MVP Timeline**: 2-3 weeks

## Development Phases

### Phase 1: Core Foundation (Week 1)

**Goal**: Establish Go project structure and basic terminal UI

#### Milestones

1. **Project Setup**
   - Initialize Go module structure
   - Setup GitHub repository
   - Configure CI with GitHub Actions
   - Create Makefile for builds

2. **Terminal UI Framework**
   - Integrate tcell library
   - Create basic UI layout (header, chat area, input)
   - Implement keyboard input handling
   - Build message display system

3. **Core Data Structures**
   ```go
   type Session struct {
       ID        string
       StartTime time.Time
       Messages  []Message
       State     SessionState
   }
   
   type Message struct {
       Type      MessageType
       Content   string
       Timestamp time.Time
       FromMe    bool
   }
   
   type MessageType string
   const (
       MessageTypeText   MessageType = "text"
       MessageTypeTyping MessageType = "typing"
       MessageTypeSystem MessageType = "system"
   )
   ```

#### Deliverables
- Working terminal UI that accepts input
- Basic project structure with modules
- Makefile with build/test targets
- Initial README

#### Success Criteria
- Can run `termchat` and see UI
- Can type and see text in input field
- Clean separation of UI and logic

### Phase 2: P2P Networking (Week 2)

**Goal**: Implement TCP-based P2P connection

#### Milestones

1. **Session Management**
   - Session ID generation (adjective-noun-number)
   - TCP listener on port 9999
   - Connection acceptance logic
   - Session state management

2. **Message Protocol**
   - JSON message encoding/decoding
   - Message framing over TCP
   - Handshake protocol implementation
   - Keep-alive mechanism

3. **Connection Handler**
   ```go
   func handleConnection(conn net.Conn, session *Session) {
       decoder := json.NewDecoder(conn)
       encoder := json.NewEncoder(conn)
       
       // Handshake
       if err := performHandshake(decoder, encoder, session); err != nil {
           return
       }
       
       // Message loop
       for {
           var msg Message
           if err := decoder.Decode(&msg); err != nil {
               break
           }
           handleMessage(msg, session)
       }
   }
   ```

#### Deliverables
- TCP server/client implementation
- Working P2P connection locally
- Message exchange between instances
- Protocol documentation

#### Success Criteria
- Two local instances can connect
- Messages flow bidirectionally
- Clean disconnect handling
- Robust error handling

### Phase 3: SSH Integration (Week 3)

**Goal**: Enable connections through SSH tunneling

#### Milestones

1. **SSH Client Integration**
   - Integrate golang.org/x/crypto/ssh
   - Parse connection strings (user@host:session)
   - Establish SSH connections
   - SSH agent authentication support

2. **Port Forwarding**
   - Remote port forwarding setup
   - Local port management
   - Tunnel lifecycle management
   - Error handling for SSH failures

3. **Connection Flow**
   ```go
   func joinViaSSH(connStr string) error {
       user, host, sessionID := parseConnectionString(connStr)
       
       // SSH connection
       client, err := connectSSH(user, host)
       if err != nil {
           return err
       }
       
       // Port forwarding
       listener, err := client.Listen("tcp", "localhost:9999")
       if err != nil {
           return err
       }
       
       // Connect through tunnel
       conn, err := net.Dial("tcp", listener.Addr().String())
       return handleJoinConnection(conn, sessionID)
   }
   ```

#### Deliverables
- SSH tunneling implementation
- Connection string parsing
- SSH error handling
- Cross-machine testing

#### Success Criteria
- Can connect between different machines
- SSH key authentication works
- Clear error messages for SSH issues
- Tunnel cleanup on disconnect

### Phase 4: Features & Polish (Week 4)

**Goal**: Add typing indicators and UI improvements

#### Milestones

1. **Typing Indicators**
   - Detect typing state changes
   - Send typing notifications
   - Display "X is typing..." in UI
   - Typing timeout handling

2. **UI Enhancements**
   - Message timestamps
   - Connection status indicator
   - Smooth scrolling for messages
   - Better error display
   - Color themes support

3. **Command System**
   - `/quit` - Exit gracefully
   - `/clear` - Clear screen
   - `/help` - Show commands
   - `Ctrl+L` - Redraw screen

#### Deliverables
- Typing indicators working
- Polished UI experience
- Command system
- User documentation

#### Success Criteria
- Typing indicators work reliably
- UI feels responsive
- All commands implemented
- No visual glitches

### Phase 5: Testing & Documentation (Week 5)

**Goal**: Comprehensive testing and documentation

#### Milestones

1. **Test Suite**
   - Unit tests for all components
   - Integration tests for P2P flow
   - Mock SSH for testing
   - Benchmarks for performance

2. **Documentation**
   - User guide with examples
   - Technical architecture docs
   - Protocol specification
   - Troubleshooting guide

3. **Error Handling**
   - Comprehensive error types
   - User-friendly error messages
   - Connection failure recovery
   - Graceful degradation

#### Deliverables
- >80% test coverage
- Complete documentation set
- Error handling guide
- Performance benchmarks

#### Success Criteria
- All tests passing
- Documentation complete
- No panic conditions
- Clear error messages

### Phase 6: Release Preparation (Week 6)

**Goal**: Prepare for public release

#### Milestones

1. **Build & Distribution**
   - Cross-platform builds (Linux, macOS, BSD)
   - GitHub releases automation
   - Homebrew formula
   - AUR package

2. **Configuration**
   - Custom port support (`--port`)
   - Alternative SSH options
   - Debug mode flag
   - Version command

3. **Final Polish**
   - Code cleanup and refactoring
   - Performance optimizations
   - Security review
   - License headers

#### Deliverables
- Release binaries for all platforms
- Installation instructions
- Contributing guidelines
- Security policy

#### Success Criteria
- Binaries work on target platforms
- Easy installation process
- No security vulnerabilities
- Ready for v1.0.0 release

## Technical Implementation Details

### Core Algorithms

#### Session ID Generation
```go
var (
    adjectives = []string{
        "cosmic", "mystic", "stellar", "quantum",
        "cyber", "neon", "turbo", "ultra",
    }
    nouns = []string{
        "phoenix", "dragon", "tiger", "eagle",
        "shark", "wolf", "bear", "fox",
    }
)

func generateSessionID() string {
    adj := adjectives[rand.Intn(len(adjectives))]
    noun := nouns[rand.Intn(len(nouns))]
    num := rand.Intn(10000)
    return fmt.Sprintf("%s-%s-%d", adj, noun, num)
}
```

#### Message Protocol Handler
```go
func handleMessage(msg Message, session *Session) {
    switch msg.Type {
    case MessageTypeText:
        session.Messages = append(session.Messages, msg)
        ui.DisplayMessage(msg)
        
    case MessageTypeTyping:
        isTyping := msg.Content == "true"
        ui.UpdateTypingIndicator(isTyping)
        
    case MessageTypePing:
        sendMessage(Message{
            Type:      MessageTypePong,
            Timestamp: time.Now(),
        })
    }
}
```

#### SSH Connection Setup
```go
func setupSSHTunnel(user, host string) (*ssh.Client, error) {
    // Try SSH agent first
    agentAuth, err := sshAgentAuth()
    if err != nil {
        // Fall back to default key
        agentAuth = publicKeyAuth()
    }
    
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{agentAuth},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: known_hosts
    }
    
    return ssh.Dial("tcp", host+":22", config)
}
```

## Risk Mitigation

### Technical Risks

1. **SSH Connection Issues**
   - Risk: Various SSH configurations and keys
   - Mitigation: Support multiple auth methods
   - Fallback: Clear troubleshooting guide

2. **Terminal Compatibility**
   - Risk: Different terminal emulators
   - Mitigation: Use tcell's compatibility layer
   - Fallback: Basic ASCII-only mode

3. **Port Conflicts**
   - Risk: Port 9999 already in use
   - Mitigation: Make port configurable
   - Fallback: Try alternative ports

### Design Risks

1. **Two-User Limitation**
   - Risk: Users want group chat
   - Mitigation: Clear documentation
   - Decision: Stay focused on P2P simplicity

2. **No Persistence**
   - Risk: Users lose conversation
   - Mitigation: Document ephemeral nature
   - Decision: Privacy over convenience

## Testing Strategy

### Unit Tests
- Session ID generation
- Message encoding/decoding
- UI component rendering
- Connection state management

### Integration Tests
- Local P2P connection
- SSH tunnel establishment
- Message flow end-to-end
- Disconnect handling

### Manual Testing Checklist
- [ ] Local connection (same machine)
- [ ] LAN connection (different machines)
- [ ] Internet connection (via SSH)
- [ ] Various SSH key types
- [ ] Terminal resize handling
- [ ] Typing indicators
- [ ] Error conditions
- [ ] Clean shutdown

### Performance Tests
- Message latency measurement
- Memory usage over time
- CPU usage during chat
- Large message handling

## Success Metrics

### MVP Success (Week 3)
- [ ] P2P chat working locally
- [ ] SSH connections functional
- [ ] <50ms local latency
- [ ] Clean UI with no glitches
- [ ] Basic error handling

### Final Success (Week 6)
- [ ] All features implemented
- [ ] Cross-platform binaries
- [ ] Typing indicators working
- [ ] Comprehensive documentation
- [ ] >80% test coverage
- [ ] <10MB memory usage
- [ ] <1% CPU usage idle
- [ ] Zero panics

## Resource Requirements

### Development
- 1 Go developer (full-time)
- GitHub account for hosting
- Test machines (Linux, macOS)
- Various terminal emulators

### Minimal Infrastructure
- GitHub Actions (free tier)
- No servers needed
- No databases
- No external services

### Development Tools
- Go 1.21+ toolchain
- tcell library
- SSH test environment
- Terminal emulators:
  - iTerm2 (macOS)
  - GNOME Terminal
  - Windows Terminal
  - Alacritty

## Conclusion

The serverless P2P architecture of termchat allows for rapid development with minimal complexity. By leveraging SSH for security and transport, we avoid implementing complex protocols while providing a secure, private communication channel.

The 6-week timeline is realistic given the simplified scope, and the resulting application will provide real value to developers who need quick, secure, ephemeral communication.