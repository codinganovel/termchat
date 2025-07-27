# Testing termchat

## Local Testing (Same Machine)

1. **Build the application:**
   ```bash
   make build
   ```

2. **Terminal 1 - Start a session:**
   ```bash
   ./build/termchat start
   ```
   - Note the session ID displayed in the header (e.g., "cosmic-turtle-7823")
   - The UI will show immediately

3. **Terminal 2 - Join the session:**
   ```bash
   ./build/termchat join localhost:SESSION-ID
   ```
   Replace SESSION-ID with the actual ID from Terminal 1

4. **Test the chat:**
   - Type messages and press Enter to send
   - Messages should appear in both terminals
   - Try these commands:
     - `/quit` or `Ctrl+D` - Exit the chat
     - `/clear` - Clear the screen
     - `Ctrl+L` - Redraw the screen

## Debugging Tips

If the UI seems frozen:
- Make sure you're running the latest build
- Check that port 9999 is not already in use
- Try `Ctrl+C` to exit and restart

## Known Issues
- The UI should now respond to keyboard input properly
- Backspace, arrow keys, and typing should all work
- If you still experience issues, rebuild with `make clean && make build`