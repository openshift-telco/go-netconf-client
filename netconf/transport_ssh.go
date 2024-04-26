// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found here
// https://github.com/Juniper/go-netconf/blob/master/LICENSE.

// The content has been modified from the original version, but the initial code
// remains from Juniper Networks, following above licence.

package netconf

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	// sshDefaultPort is the default SSH port used when communicating with
	// NETCONF
	sshDefaultPort = 830
	// sshNetconfSubsystem sets the SSH subsystem to NETCONF
	sshNetconfSubsystem = "netconf"
)

// TransportSSH maintains the information necessary to communicate with the
// remote device over SSH
type TransportSSH struct {
	transportBasicIO
	sshClient  *ssh.Client
	sshSession *ssh.Session
}

// Close closes an existing SSH session and socket if they exist.
func (t *TransportSSH) Close() error {
	// If TransportSSH is nil ignore closing ssh session
	if t == nil {
		return nil
	}

	// Close the SSH Session if we have one
	if t.sshSession != nil {
		if err := t.sshSession.Close(); err != nil {
			// If we receive an error when trying to close the session, then
			// lets try to close the socket, otherwise it will be left open
			err := t.sshClient.Close()
			if err != nil {
				return err
			}
			return err
		}
	}

	// Close the socket
	if t.sshClient != nil {
		return t.sshClient.Close()
	}
	return fmt.Errorf("no connection to close")
}

// Dial connects and establishes SSH sessions
//
// target can be an IP address (e.g.) 172.16.1.1 which utilizes the default
// NETCONF over SSH port of 830.  Target can also specify a port with the
// following format <host>:<port (e.g. 172.16.1.1:22)
//
// config takes a ssh.ClientConfig connection. See documentation for
// go.crypto/ssh for documentation.  There is a helper function SSHConfigPassword
// thar returns a ssh.ClientConfig for simple username/password authentication
func (t *TransportSSH) Dial(target string, config *ssh.ClientConfig) error {
	if !strings.Contains(target, ":") {
		target = fmt.Sprintf("%s:%d", target, sshDefaultPort)
	}

	var err error

	t.sshClient, err = ssh.Dial("tcp", target, config)
	if err != nil {
		return err
	}

	err = t.setupSession()
	return err
}

// DialSSH creates a new SSH Transport.
// See TransportSSH.Dial for arguments.
func DialSSH(target string, config *ssh.ClientConfig) (*TransportSSH, error) {
	t := new(TransportSSH)
	err := t.Dial(target, config)
	if err != nil {
		err := t.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return t, nil
}

// DialSSHTimeout creates a new SSH Transport with timeout.
// See TransportSSH.Dial for arguments.
// The timeout value is used for both connection establishment and Read/Write operations.
func DialSSHTimeout(target string, config *ssh.ClientConfig, timeout time.Duration) (*TransportSSH, error) {
	bareConn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return nil, err
	}

	conn := &deadlineConn{Conn: bareConn, timeout: timeout}
	t, err := connToTransport(conn, config)
	if err != nil {
		if t != nil {
			err := t.Close()
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(timeout / 2)
		defer ticker.Stop()
		for range ticker.C {
			_, _, err := t.sshClient.Conn.SendRequest("KEEP_ALIVE", true, nil)
			if err != nil {
				return
			}
		}
	}()

	return t, nil
}

// NoDialSSH - create a new TransportSSH from given ssh Client.
func NoDialSSH(sshClient *ssh.Client) (*TransportSSH, error) {
	t := new(TransportSSH)
	t.sshClient = sshClient
	err := t.setupSession()
	if err != nil {
		err := t.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return t, nil
}

// SSHConfigPubKeyFile is a convenience function that takes a username, private key
// and passphrase and returns a new ssh.ClientConfig setup to pass credentials
// to DialSSH
func SSHConfigPubKeyFile(user string, file string, passphrase string) (*ssh.ClientConfig, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(buf)
	if len(rest) > 0 {
		return nil, fmt.Errorf("pem: unable to decode file %s", file)
	}

	if x509.IsEncryptedPEMBlock(block) {
		b, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, err
		}
		buf = pem.EncodeToMemory(&pem.Block{
			Type:  block.Type,
			Bytes: b,
		})
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}, nil

}

func connToTransport(conn net.Conn, config *ssh.ClientConfig) (*TransportSSH, error) {
	c, channel, reqs, err := ssh.NewClientConn(conn, conn.RemoteAddr().String(), config)
	if err != nil {
		return nil, err
	}

	t := &TransportSSH{}
	t.sshClient = ssh.NewClient(c, channel, reqs)

	err = t.setupSession()
	if err != nil {
		return nil, err
	}

	return t, nil
}

type deadlineConn struct {
	net.Conn
	timeout time.Duration
}

func (c *deadlineConn) Read(b []byte) (n int, err error) {
	_ = c.SetReadDeadline(time.Now().Add(c.timeout))
	return c.Conn.Read(b)
}

func (c *deadlineConn) Write(b []byte) (n int, err error) {
	_ = c.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.Conn.Write(b)
}

func (t *TransportSSH) setupSession() error {
	var err error

	t.sshSession, err = t.sshClient.NewSession()
	if err != nil {
		return err
	}

	writer, err := t.sshSession.StdinPipe()
	if err != nil {
		return err
	}

	reader, err := t.sshSession.StdoutPipe()
	if err != nil {
		return err
	}

	t.ReadWriteCloser = NewReadWriteCloser(reader, writer)
	return t.sshSession.RequestSubsystem(sshNetconfSubsystem)
}
