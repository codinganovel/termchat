package network

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func getSSHAuthMethod() ssh.AuthMethod {
	if authMethod := sshAgentAuth(); authMethod != nil {
		return authMethod
	}
	
	return publicKeyAuth()
}

func sshAgentAuth() ssh.AuthMethod {
	if sshAuthSock := os.Getenv("SSH_AUTH_SOCK"); sshAuthSock != "" {
		if conn, err := net.Dial("unix", sshAuthSock); err == nil {
			agentClient := agent.NewClient(conn)
			return ssh.PublicKeysCallback(agentClient.Signers)
		}
	}
	return nil
}

func publicKeyAuth() ssh.AuthMethod {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	
	keyPaths := []string{
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
		filepath.Join(home, ".ssh", "id_dsa"),
	}
	
	for _, keyPath := range keyPaths {
		if key, err := readPrivateKey(keyPath); err == nil {
			return ssh.PublicKeys(key)
		}
	}
	
	return nil
}

func readPrivateKey(path string) (ssh.Signer, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			return nil, fmt.Errorf("key %s requires passphrase", path)
		}
		return nil, err
	}
	
	return signer, nil
}