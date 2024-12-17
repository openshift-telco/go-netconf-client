package netconf

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// NewSessionFromSSHConfig established a NETCONF session connecting to the target using ssh client configuration.
func NewSessionFromSSHConfig(target string, config *ssh.ClientConfig, options ...SessionOption) (*Session, error) {
	t, err := DialSSH(target, config)
	if err != nil {
		return nil, fmt.Errorf("DialSSHTimeout: %w", err)
	}

	s, err := NewSession(t, options...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// NewSessionFromSSHConfigTimeout established a NETCONF session connecting to the target using ssh client configuration with timeout.
func NewSessionFromSSHConfigTimeout(ctx context.Context, target string, config *ssh.ClientConfig, timeout time.Duration, options ...SessionOption) (*Session, error) {
	t, err := DialSSHTimeout(target, config, timeout)
	if err != nil {
		return nil, fmt.Errorf("DialSSHTimeout: %w", err)
	}

	s, err := NewSession(t, options...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// NewSessionFromSSHClient established a NETCONF session over a given ssh client.
func NewSessionFromSSHClient(ctx context.Context, client *ssh.Client, options ...SessionOption) (*Session, error) {
	t, err := NoDialSSH(client)
	if err != nil {
		return nil, fmt.Errorf("NoDialSSH: %w", err)
	}

	s, err := NewSession(t, options...)
	if err != nil {
		return nil, err
	}

	return s, nil
}
