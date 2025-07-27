# termchat Protocol Specification

**Version:** 2.0.0  
**Last Updated:** 2025-07-27

## Table of Contents

1. [Protocol Overview](#protocol-overview)
2. [Message Format](#message-format)
3. [Connection Protocol](#connection-protocol)
4. [Message Types](#message-types)
5. [Error Handling](#error-handling)
6. [Session Management](#session-management)

## Protocol Overview

termchat uses a simple, text-based protocol over TCP:

- **Transport**: Direct TCP socket or SSH-tunneled TCP
- **Encoding**: JSON messages with newline delimiters
- **Connection**: One-to-one P2P only
- **Session**: Exists only while both parties connected

### Protocol Characteristics

- **Text-Based**: Human-readable JSON
- **Stateless**: No server-side state
- **Synchronous**: Direct message exchange
- **Ephemeral**: No persistence layer

## Message Format

### Wire Format

Each message is a JSON object followed by a newline character (`\n`):

```
{"type":"text","content":"Hello","timestamp":1234567890}\n
```

### Message Structure

```go
type Message struct {
    Type      string `json:"type"`
    Content   string `json:"content,omitempty"`
    SessionID string `json:"session_id,omitempty"`
    Timestamp int64  `json:"timestamp"`
}
```

### Core Fields

- **type**: Message type identifier (required)
- **content**: Message payload (optional, depends on type)
- **session_id**: Session identifier (used during handshake)
- **timestamp**: Unix timestamp in milliseconds

## Connection Protocol

### 1. Initiator (Person A)

When starting a session:

```bash
$ termchat start
```

The initiator:
1. Generates a unique session ID
2. Opens TCP listener on port 9999
3. Waits for incoming connection
4. Validates session ID on connect

### 2. Joiner (Person B)

When joining a session:

```bash
$ termchat join alice@host:session-id
```

The joiner:
1. Establishes SSH connection to host
2. Sets up port forwarding (remote:9999 → local:ephemeral)
3. Connects to forwarded port
4. Sends handshake with session ID

### 3. Handshake Sequence

```
Joiner                          Initiator
  │                                │
  ├─────[TCP Connect]─────────────▶│
  │                                │
  ├───[HELLO + Session ID]────────▶│
  │                                ├─▶ Validate
  │◀───────[WELCOME]───────────────┤
  │                                │
  ├────────[READY]───────────────▶│
  │                                │
  │◀========[Chat Active]=========▶│
```

## Message Types

### 1. Handshake Messages

#### HELLO (Joiner → Initiator)
```json
{
  "type": "hello",
  "session_id": "cosmic-turtle-7823",
  "timestamp": 1234567890
}
```

#### WELCOME (Initiator → Joiner)
```json
{
  "type": "welcome",
  "session_id": "cosmic-turtle-7823",
  "timestamp": 1234567890
}
```

#### READY (Joiner → Initiator)
```json
{
  "type": "ready",
  "timestamp": 1234567890
}
```

### 2. Chat Messages

#### TEXT
```json
{
  "type": "text",
  "content": "Hello, how are you?",
  "timestamp": 1234567890
}
```

#### TYPING
```json
{
  "type": "typing",
  "content": "true",
  "timestamp": 1234567890
}
```

### 3. Control Messages

#### PING
```json
{
  "type": "ping",
  "timestamp": 1234567890
}
```

#### PONG
```json
{
  "type": "pong",
  "timestamp": 1234567890
}
```

#### LEAVE
```json
{
  "type": "leave",
  "content": "User has disconnected",
  "timestamp": 1234567890
}
```

## Error Handling

### Connection Errors

#### SESSION_MISMATCH
```json
{
  "type": "error",
  "content": "Session ID mismatch",
  "timestamp": 1234567890
}
```

#### INVALID_MESSAGE
```json
{
  "type": "error",
  "content": "Invalid message format",
  "timestamp": 1234567890
}
```

### Error Recovery

- **Connection Lost**: No automatic reconnection
- **Invalid Session**: Connection refused
- **Protocol Error**: Connection terminated
- **SSH Failure**: User-friendly error message

### Common Error Scenarios

1. **Wrong Session ID**: Joiner provides incorrect ID
   - Initiator sends SESSION_MISMATCH error
   - Connection closed

2. **Port Already in Use**: Another termchat running
   - Error message to user
   - Suggest using different port

3. **SSH Connection Failed**: Can't reach host
   - Clear error about SSH failure
   - Check SSH access separately

4. **Network Interruption**: Connection drops
   - Both sides detect TCP close
   - Clean shutdown, no reconnect

## Session Management

### Session ID Format

Session IDs follow the pattern: `adjective-noun-number`

- **Examples**: `cosmic-turtle-7823`, `mystic-phoenix-1492`
- **Components**: 
  - Adjective from curated list
  - Noun from curated list  
  - Random 4-digit number
- **Uniqueness**: Random generation, collision unlikely

### Session Lifecycle

```
┌──────────┐      ┌──────────┐      ┌─────────┐
│ CREATED  │─────▶│ WAITING  │─────▶│ ACTIVE  │
└──────────┘      └──────────┘      └─────────┘
                                           │
                                           ▼
                                     ┌─────────┐
                                     │ ENDED   │
                                     └─────────┘
```

### State Transitions

1. **CREATED**: Session ID generated, listener started
2. **WAITING**: Listening for incoming connection
3. **ACTIVE**: Both parties connected, chat active
4. **ENDED**: Either party disconnected

### Implementation Notes

- No session persistence between runs
- Session ID only used for handshake validation
- No session state beyond the TCP connection
- Clean shutdown on any disconnect

## Protocol Design Rationale

### Why JSON?
- Human-readable for debugging
- Simple parsing in Go
- No binary protocol complexity
- Good enough performance for chat

### Why SSH?
- Existing authentication infrastructure
- Built-in encryption
- Firewall-friendly (port 22)
- No custom security to implement

### Why P2P?
- No server costs or maintenance
- Perfect privacy (no third party)
- Minimal attack surface
- Simple architecture