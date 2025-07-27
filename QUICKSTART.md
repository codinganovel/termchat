# termchat Quick Start

## Build
```bash
cd /Users/sam/Documents/coding/termchat
make build
```

## Test Locally (Two Terminals)

### Terminal 1 - Start a session:
```bash
./build/termchat start
```

Note the session ID in the header (e.g., "cosmic-turtle-7823")

### Terminal 2 - Join the session:
```bash
./build/termchat join localhost:SESSION-ID
```
Replace SESSION-ID with the actual ID from Terminal 1.

## Troubleshooting

### Port already in use
If you see "address already in use", kill the old process:
```bash
lsof -ti:9999 | xargs kill -9
```

### Can't build
Make sure you're in the termchat directory and have Go installed:
```bash
cd /Users/sam/Documents/coding/termchat
which go  # Should show /usr/local/go/bin/go or similar
```

## Commands
- Type messages and press Enter to send
- `/quit` or `Ctrl+D` - Exit
- `/clear` - Clear screen
- `Ctrl+L` - Redraw screen