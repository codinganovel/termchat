# termchat

**Version:** 2.0.0  
**Last Updated:** 2025-07-27

## Overview

termchat is a serverless, peer-to-peer terminal chat application that works over SSH. It enables secure, ephemeral one-on-one conversations between two developers without any infrastructure requirements. Simply start a session, share the ID, and chat directly through an SSH tunnel.

## Key Features

- **Zero Infrastructure**: No server needed - pure P2P over SSH tunneling
- **One-on-One Chat**: Designed for focused conversations between exactly 2 people
- **SSH-Based**: Leverages existing SSH access for secure connections
- **Ephemeral Sessions**: Chat exists only while both users are connected
- **Full Terminal UI**: Clean, responsive interface using tcell
- **No Persistence**: No logs, no history - completely private
- **Written in Go**: Fast, efficient, single binary distribution

## Quick Start

### Installation

```bash
# Install from binary
$ go install github.com/yourusername/termchat@latest
# or download pre-built binary
$ curl -L https://github.com/yourusername/termchat/releases/latest/download/termchat-$(uname -s)-$(uname -m) -o termchat
$ chmod +x termchat
```

### Basic Usage

#### Person A: Start a session
```bash
$ termchat start
Session started: cosmic-turtle-7823
Listening on local port 9999
Share this session ID with your chat partner
Waiting for connection...
```

#### Person B: Join the session
```bash
$ termchat join alice@workstation:cosmic-turtle-7823
Connecting via SSH to alice@workstation...
Connected! Type your messages below.
```

### Interface Overview

```
┌─────────────────────────────────────────────────────┐
│ termchat - cosmic-turtle-7823 (P2P with alice)     │
├─────────────────────────────────────────────────────┤
│ alice: Hey, can you help with this bug?            │
│ you: Sure, what's the issue?                       │
│ alice: The async handler is hanging...              │
│ you: Let me check the goroutine dump                │
│                                                     │
│ alice is typing...                                  │
├─────────────────────────────────────────────────────┤
│ > type your message here...                         │
└─────────────────────────────────────────────────────┘
```

## Commands

### Starting and Joining
- `termchat start` - Start a new P2P session and wait for connection
- `termchat join user@host:session-id` - Join a session via SSH tunnel

### In-Session Commands
- `/quit` or `Ctrl+D` - End the chat session
- `/clear` - Clear your local display
- `/help` - Show available commands
- `Ctrl+L` - Redraw the screen

## Use Cases

1. **Pair Programming**: Coordinate without leaving the terminal
2. **Quick Debugging**: Discuss issues with a colleague instantly
3. **Code Reviews**: Walk through changes in real-time
4. **Remote Assistance**: Help someone debug on their machine
5. **Secure Discussions**: Private conversations with no logs or history

## Security & Privacy

- **End-to-End Encryption**: All traffic goes through SSH
- **No Server**: Direct P2P connection means no third-party access
- **Zero Persistence**: Messages exist only in memory during the session
- **No Authentication**: If you can SSH to the host, you can chat
- **Session Isolation**: Each session ID creates a unique connection

## Requirements

- Go 1.21 or later (for building from source)
- Terminal with UTF-8 support
- SSH client and server (OpenSSH recommended)
- Unix-like OS (Linux, macOS, BSD)
- Open SSH port (usually 22) between participants

## Technical Details

### How It Works

1. Person A runs `termchat start`, which:
   - Generates a unique session ID
   - Opens a local TCP listener (default port 9999)
   - Waits for incoming connection

2. Person B runs `termchat join alice@host:session-id`, which:
   - Establishes SSH connection to alice@host
   - Sets up port forwarding from remote to local
   - Connects to the forwarded port
   - Establishes direct P2P chat connection

3. Messages flow directly between the two instances through the SSH tunnel

### Architecture Benefits

- **No Server Costs**: Zero infrastructure to maintain
- **Perfect Forward Secrecy**: Each session is completely isolated
- **Minimal Attack Surface**: Only SSH exposed, no custom protocols
- **Easy Deployment**: Single Go binary, no dependencies

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - See [LICENSE](LICENSE) for details