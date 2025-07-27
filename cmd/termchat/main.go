package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sam/termchat/internal/network"
	"github.com/sam/termchat/internal/session"
	"github.com/sam/termchat/internal/ui"
	"github.com/sam/termchat/pkg/protocol"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	port    int
	
	rootCmd = &cobra.Command{
		Use:   "termchat",
		Short: "Serverless P2P terminal chat over SSH",
		Long: `termchat is a serverless, peer-to-peer terminal chat application that works over SSH.
It enables secure, ephemeral one-on-one conversations between two developers without any infrastructure requirements.`,
	}
	
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a new P2P session and wait for connection",
		Run:   startSession,
	}
	
	joinCmd = &cobra.Command{
		Use:   "join user@host:session-id",
		Short: "Join a session via SSH tunnel",
		Args:  cobra.ExactArgs(1),
		Run:   joinSession,
	}
	
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("termchat version %s\n", version)
		},
	}
)

func init() {
	startCmd.Flags().IntVar(&port, "port", 9999, "Port to listen on")
	
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(joinCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func startSession(cmd *cobra.Command, args []string) {
	sess := session.New()
	
	fmt.Printf("Session started: %s\n", sess.ID)
	fmt.Printf("Listening on port %d\n", port)
	fmt.Println()
	fmt.Println("Share this with your chat partner:")
	fmt.Printf("  termchat join user@host:%s", sess.ID)
	if port != 9999 {
		fmt.Printf(":%d", port)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("Waiting for connection...")
	
	server := network.NewServer(sess)
	
	if err := server.Start(port); err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			fmt.Fprintf(os.Stderr, "Error: Port %d is already in use.\n", port)
			fmt.Fprintf(os.Stderr, "Try: lsof -ti:%d | xargs kill -9\n", port)
			fmt.Fprintf(os.Stderr, "Or use a different port: termchat start --port 9998\n")
		} else {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		}
		os.Exit(1)
	}
	
	ui, err := ui.NewSimple(sess.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize UI: %v\n", err)
		os.Exit(1)
	}
	defer ui.Close()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	stopChan := make(chan struct{})
	
	server.SetCallbacks(
		func(msg protocol.Message) {
			ui.DisplayMessage(msg)
		},
		func() {
			ui.AddMessage("[Connected]")
		},
		func() {
			ui.AddMessage("[Disconnected]")
			close(stopChan)
		},
	)
	
	ui.SetCallbacks(
		func(text string) {
			msg := protocol.NewMessage(protocol.MessageTypeText, text)
			server.SendMessage(msg)
		},
		func() {
			close(stopChan)
		},
	)
	
	go ui.Run()
	
	select {
	case <-sigChan:
		fmt.Println("\nShutting down...")
	case <-stopChan:
	}
	
	server.Stop()
}

func joinSession(cmd *cobra.Command, args []string) {
	connInfo, err := network.ParseConnectionString(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid connection string: %v\n", err)
		os.Exit(1)
	}
	
	sess := session.New()
	sess.ID = connInfo.SessionID
	
	fmt.Printf("Connecting via SSH to %s@%s...\n", connInfo.User, connInfo.Host)
	
	client := network.NewClient(sess)
	
	isLocal := connInfo.Host == "localhost" || connInfo.Host == "127.0.0.1"
	
	if isLocal {
		err = client.ConnectLocal(fmt.Sprintf("localhost:%d", connInfo.Port), connInfo.SessionID)
	} else {
		err = client.ConnectViaSSH(connInfo)
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Connected! Type your messages below.")
	fmt.Println()
	
	ui, err := ui.NewSimple(sess.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize UI: %v\n", err)
		os.Exit(1)
	}
	defer ui.Close()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	stopChan := make(chan struct{})
	
	client.SetCallbacks(
		func(msg protocol.Message) {
			ui.DisplayMessage(msg)
		},
		func() {
			ui.AddMessage("[Connected to session]")
		},
		func() {
			ui.AddMessage("[Disconnected]")
			close(stopChan)
		},
	)
	
	ui.SetCallbacks(
		func(text string) {
			msg := protocol.NewMessage(protocol.MessageTypeText, text)
			client.SendMessage(msg)
		},
		func() {
			close(stopChan)
		},
	)
	
	go ui.Run()
	
	select {
	case <-sigChan:
		fmt.Println("\nShutting down...")
	case <-stopChan:
	}
	
	client.Stop()
}