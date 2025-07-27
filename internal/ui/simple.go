package ui

import (
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/sam/termchat/pkg/protocol"
)

type SimpleUI struct {
	screen    tcell.Screen
	messages  []ChatMsg
	input     string
	cursorPos int
	sessionID string
	scrollPos int  // 0 = bottom (newest), increases as you scroll up
	mu        sync.Mutex
	
	onMessage func(string)
	onQuit    func()
}

type ChatMsg struct {
	Content string
	FromMe  bool
}

func NewSimple(sessionID string) (*SimpleUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	
	if err := screen.Init(); err != nil {
		return nil, err
	}
	
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.Clear()
	
	return &SimpleUI{
		screen:    screen,
		messages:  make([]ChatMsg, 0),
		sessionID: sessionID,
	}, nil
}

func (ui *SimpleUI) Close() {
	ui.screen.Fini()
}

func (ui *SimpleUI) SetCallbacks(onMessage func(string), onQuit func()) {
	ui.onMessage = onMessage
	ui.onQuit = onQuit
}

func (ui *SimpleUI) Run() {
	ui.draw()
	
	for {
		ev := ui.screen.PollEvent()
		if ev == nil {
			return
		}
		
		switch ev := ev.(type) {
		case *tcell.EventKey:
			ui.handleKey(ev)
		case *tcell.EventResize:
			ui.screen.Sync()
			ui.draw()
		}
	}
}

func (ui *SimpleUI) handleKey(ev *tcell.EventKey) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	
	switch ev.Key() {
	case tcell.KeyCtrlD:
		if ui.onQuit != nil {
			ui.onQuit()
		}
		return
		
	case tcell.KeyEnter:
		if ui.input != "" {
			if ui.input == "/quit" {
				if ui.onQuit != nil {
					ui.onQuit()
				}
				return
			}
			
			// Add message to display
			ui.messages = append(ui.messages, ChatMsg{
				Content: ui.input,
				FromMe:  true,
			})
			
			// Reset scroll to bottom when sending
			ui.scrollPos = 0
			
			// Send message
			if ui.onMessage != nil {
				ui.onMessage(ui.input)
			}
			
			ui.input = ""
			ui.cursorPos = 0
		}
		
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if ui.cursorPos > 0 {
			ui.input = ui.input[:ui.cursorPos-1] + ui.input[ui.cursorPos:]
			ui.cursorPos--
		}
		
	case tcell.KeyLeft:
		if ui.cursorPos > 0 {
			ui.cursorPos--
		}
		
	case tcell.KeyRight:
		if ui.cursorPos < len(ui.input) {
			ui.cursorPos++
		}
		
	case tcell.KeyUp, tcell.KeyCtrlK:
		// Scroll up (older messages)
		if ui.scrollPos < len(ui.messages)-1 {
			ui.scrollPos++
		}
		
	case tcell.KeyDown, tcell.KeyCtrlJ:
		// Scroll down (newer messages)
		if ui.scrollPos > 0 {
			ui.scrollPos--
		}
		
	default:
		if ev.Rune() != 0 {
			ui.input = ui.input[:ui.cursorPos] + string(ev.Rune()) + ui.input[ui.cursorPos:]
			ui.cursorPos++
		}
	}
	
	ui.draw()
}

func (ui *SimpleUI) draw() {
	ui.screen.Clear()
	width, height := ui.screen.Size()
	
	// Draw session ID at top
	sessionText := "Session: " + ui.sessionID
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)
	for i, r := range sessionText {
		if i < width {
			ui.screen.SetContent(i, 0, r, nil, style)
		}
	}
	
	// Draw messages with boxes
	y := 2
	
	// Calculate visible message range based on scroll
	totalMessages := len(ui.messages)
	if totalMessages == 0 {
		return
	}
	
	// Calculate how many messages can fit
	availableHeight := height - 4 // Leave room for header and input
	messagesPerScreen := availableHeight / 4 // Each message box takes ~4 lines
	
	// Calculate start index based on scroll position
	startIdx := totalMessages - messagesPerScreen - ui.scrollPos
	if startIdx < 0 {
		startIdx = 0
	}
	
	endIdx := startIdx + messagesPerScreen
	if endIdx > totalMessages {
		endIdx = totalMessages
	}
	
	// Draw visible messages
	for i := startIdx; i < endIdx; i++ {
		if y+3 >= height-2 {
			break
		}
		msg := ui.messages[i]
		displayText := ""
		if msg.FromMe {
			displayText = "you: " + msg.Content
		} else {
			displayText = "peer: " + msg.Content
		}
		ui.drawMessageBox(1, y, width-2, displayText)
		y += 4
	}
	
	// Draw input line at bottom
	inputY := height - 1
	for i, r := range ui.input {
		if i < width {
			ui.screen.SetContent(i, inputY, r, nil, tcell.StyleDefault)
		}
	}
	
	// Show cursor
	ui.screen.ShowCursor(ui.cursorPos, inputY)
	ui.screen.Show()
}

func (ui *SimpleUI) drawMessageBox(x, y, maxWidth int, msg string) {
	// Wrap text if needed
	lines := wrapText(msg, maxWidth-2)
	boxHeight := len(lines) + 2
	boxWidth := maxWidth
	
	// Find actual width needed
	actualWidth := 0
	for _, line := range lines {
		if len(line) > actualWidth {
			actualWidth = len(line)
		}
	}
	if actualWidth < maxWidth-2 {
		boxWidth = actualWidth + 2
	}
	
	// Top border
	ui.screen.SetContent(x, y, '┌', nil, tcell.StyleDefault)
	for i := 1; i < boxWidth-1; i++ {
		ui.screen.SetContent(x+i, y, '─', nil, tcell.StyleDefault)
	}
	ui.screen.SetContent(x+boxWidth-1, y, '┐', nil, tcell.StyleDefault)
	
	// Message lines
	for i, line := range lines {
		ui.screen.SetContent(x, y+i+1, '│', nil, tcell.StyleDefault)
		for j, r := range line {
			ui.screen.SetContent(x+j+1, y+i+1, r, nil, tcell.StyleDefault)
		}
		ui.screen.SetContent(x+boxWidth-1, y+i+1, '│', nil, tcell.StyleDefault)
	}
	
	// Bottom border
	ui.screen.SetContent(x, y+boxHeight-1, '└', nil, tcell.StyleDefault)
	for i := 1; i < boxWidth-1; i++ {
		ui.screen.SetContent(x+i, y+boxHeight-1, '─', nil, tcell.StyleDefault)
	}
	ui.screen.SetContent(x+boxWidth-1, y+boxHeight-1, '┘', nil, tcell.StyleDefault)
}

func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	
	var lines []string
	words := strings.Fields(text)
	currentLine := ""
	
	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}

func (ui *SimpleUI) AddMessage(content string) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	
	ui.messages = append(ui.messages, ChatMsg{
		Content: content,
		FromMe:  false,
	})
	// Reset scroll to see new message
	ui.scrollPos = 0
	ui.draw()
}

func (ui *SimpleUI) DisplayMessage(msg protocol.Message) {
	if msg.Type == protocol.MessageTypeText {
		ui.mu.Lock()
		defer ui.mu.Unlock()
		
		ui.messages = append(ui.messages, ChatMsg{
			Content: msg.Content,
			FromMe:  false,
		})
		// Reset scroll to see new message
		ui.scrollPos = 0
		ui.draw()
	}
}