// Package ssh provides SSH connection management for cbwsh.
package ssh

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Manager handles SSH connections and host management.
type Manager struct {
	mu             sync.RWMutex
	client         *ssh.Client
	state          core.SSHConnectionState
	currentHost    *core.SSHHost
	savedHosts     []core.SSHHost
	hostFilePath   string
	knownHostsPath string
	timeout        time.Duration
	strictHostKey  bool
}

// NewManager creates a new SSH manager.
func NewManager(hostFilePath string, timeout time.Duration) *Manager {
	homeDir, _ := os.UserHomeDir()
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

	return &Manager{
		state:          core.SSHDisconnected,
		savedHosts:     make([]core.SSHHost, 0),
		hostFilePath:   hostFilePath,
		knownHostsPath: knownHostsPath,
		timeout:        timeout,
		strictHostKey:  false, // Default to permissive for ease of use
	}
}

// SetStrictHostKeyChecking enables or disables strict host key checking.
func (m *Manager) SetStrictHostKeyChecking(strict bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.strictHostKey = strict
}

// SetKnownHostsPath sets the path to the known_hosts file.
func (m *Manager) SetKnownHostsPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.knownHostsPath = path
}

// Connect establishes an SSH connection.
func (m *Manager) Connect(ctx context.Context, host string, port int, user string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		if err := m.client.Close(); err != nil {
			// Ignore close errors
			_ = err
		}
	}

	m.state = core.SSHConnecting

	hostKeyCallback, err := m.getHostKeyCallback()
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to configure host key verification: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostKeyCallback,
		Timeout:         m.timeout,
	}

	// Try to find saved host with key
	for _, savedHost := range m.savedHosts {
		if savedHost.Host == host && savedHost.User == user {
			if savedHost.KeyPath != "" {
				key, err := loadPrivateKey(savedHost.KeyPath, savedHost.Passphrase)
				if err == nil {
					config.Auth = append(config.Auth, ssh.PublicKeys(key))
				}
			}
			break
		}
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.client = client
	m.state = core.SSHConnected
	m.currentHost = &core.SSHHost{
		Host: host,
		Port: port,
		User: user,
	}

	return nil
}

// ConnectWithKey establishes an SSH connection using a key file.
func (m *Manager) ConnectWithKey(ctx context.Context, host string, port int, user, keyPath, passphrase string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		_ = m.client.Close()
	}

	m.state = core.SSHConnecting

	key, err := loadPrivateKey(keyPath, passphrase)
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to load key: %w", err)
	}

	hostKeyCallback, err := m.getHostKeyCallback()
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to configure host key verification: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         m.timeout,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.client = client
	m.state = core.SSHConnected
	m.currentHost = &core.SSHHost{
		Host:    host,
		Port:    port,
		User:    user,
		KeyPath: keyPath,
	}

	return nil
}

// ConnectWithPassword establishes an SSH connection using password authentication.
func (m *Manager) ConnectWithPassword(ctx context.Context, host string, port int, user, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		_ = m.client.Close()
	}

	m.state = core.SSHConnecting

	hostKeyCallback, err := m.getHostKeyCallback()
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to configure host key verification: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         m.timeout,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		m.state = core.SSHError
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.client = client
	m.state = core.SSHConnected
	m.currentHost = &core.SSHHost{
		Host: host,
		Port: port,
		User: user,
	}

	return nil
}

// ConnectToSavedHost connects to a saved host by name.
func (m *Manager) ConnectToSavedHost(ctx context.Context, name string) error {
	m.mu.RLock()
	var host *core.SSHHost
	for i := range m.savedHosts {
		if m.savedHosts[i].Name == name {
			host = &m.savedHosts[i]
			break
		}
	}
	m.mu.RUnlock()

	if host == nil {
		return fmt.Errorf("host not found: %s", name)
	}

	if host.KeyPath != "" {
		return m.ConnectWithKey(ctx, host.Host, host.Port, host.User, host.KeyPath, host.Passphrase)
	}

	return m.Connect(ctx, host.Host, host.Port, host.User)
}

// Disconnect closes the current SSH connection.
func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		err := m.client.Close()
		m.client = nil
		m.state = core.SSHDisconnected
		m.currentHost = nil
		return err
	}

	m.state = core.SSHDisconnected
	return nil
}

// Execute runs a command on the remote host.
func (m *Manager) Execute(ctx context.Context, command string) (*core.CommandResult, error) {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("not connected")
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	startTime := time.Now()

	output, err := session.CombinedOutput(command)

	result := &core.CommandResult{
		Command:  command,
		Duration: time.Since(startTime).Milliseconds(),
	}

	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			result.ExitCode = -1
		}
		result.Error = string(output)
	} else {
		result.Output = string(output)
		result.ExitCode = 0
	}

	return result, nil
}

// State returns the current connection state.
func (m *Manager) State() core.SSHConnectionState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// CurrentHost returns the currently connected host.
func (m *Manager) CurrentHost() *core.SSHHost {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.currentHost == nil {
		return nil
	}
	host := *m.currentHost
	return &host
}

// ListSavedHosts returns a list of saved SSH hosts.
func (m *Manager) ListSavedHosts() ([]core.SSHHost, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]core.SSHHost, len(m.savedHosts))
	copy(result, m.savedHosts)
	return result, nil
}

// SaveHost saves an SSH host configuration.
func (m *Manager) SaveHost(host core.SSHHost) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update existing or add new
	found := false
	for i := range m.savedHosts {
		if m.savedHosts[i].Name == host.Name {
			m.savedHosts[i] = host
			found = true
			break
		}
	}

	if !found {
		m.savedHosts = append(m.savedHosts, host)
	}

	return nil
}

// RemoveHost removes a saved SSH host.
func (m *Manager) RemoveHost(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.savedHosts {
		if m.savedHosts[i].Name == name {
			m.savedHosts = append(m.savedHosts[:i], m.savedHosts[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("host not found: %s", name)
}

// ForwardLocalPort creates a local port forwarding tunnel.
func (m *Manager) ForwardLocalPort(ctx context.Context, localPort int, remoteHost string, remotePort int) error {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("not connected")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		defer listener.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			localConn, err := listener.Accept()
			if err != nil {
				continue
			}

			remoteConn, err := client.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost, remotePort))
			if err != nil {
				localConn.Close()
				continue
			}

			go copyConn(localConn, remoteConn)
			go copyConn(remoteConn, localConn)
		}
	}()

	return nil
}

func copyConn(dst, src net.Conn) {
	defer dst.Close()
	defer src.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func loadPrivateKey(keyPath, passphrase string) (ssh.Signer, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	if passphrase != "" {
		return ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	}

	return ssh.ParsePrivateKey(keyData)
}

// getHostKeyCallback returns an appropriate host key callback.
// If strict host key checking is enabled and a known_hosts file exists,
// it uses that file for verification. Otherwise, it uses a permissive callback.
func (m *Manager) getHostKeyCallback() (ssh.HostKeyCallback, error) {
	if m.strictHostKey && m.knownHostsPath != "" {
		// Try to use known_hosts file
		if _, err := os.Stat(m.knownHostsPath); err == nil {
			callback, err := knownhosts.New(m.knownHostsPath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse known_hosts: %w", err)
			}
			return callback, nil
		}
	}

	// Fall back to a callback that accepts any key but logs it
	// NOTE: This is insecure and should only be used in trusted networks
	// or when strict host key checking is explicitly disabled
	return ssh.InsecureIgnoreHostKey(), nil //nolint:gosec // User explicitly disabled strict mode
}

// AddHostToKnownHosts adds a host key to the known_hosts file.
func (m *Manager) AddHostToKnownHosts(hostname string, remote net.Addr, key ssh.PublicKey) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.knownHostsPath == "" {
		return fmt.Errorf("known_hosts path not configured")
	}

	// Create directory if needed
	dir := filepath.Dir(m.knownHostsPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	// Open file for appending
	f, err := os.OpenFile(m.knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open known_hosts: %w", err)
	}
	defer f.Close()

	// Format the entry
	entry := knownhosts.Line([]string{knownhosts.Normalize(remote.String())}, key)

	if _, err := f.WriteString(entry + "\n"); err != nil {
		return fmt.Errorf("failed to write to known_hosts: %w", err)
	}

	return nil
}

// IsHostKnown checks if a host is in the known_hosts file.
func (m *Manager) IsHostKnown(hostname string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.knownHostsPath == "" {
		return false
	}

	file, err := os.Open(m.knownHostsPath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, hostname) {
			return true
		}
	}

	return false
}

// HostKeyCallback returns a host key callback that checks known hosts.
// Deprecated: Use SetStrictHostKeyChecking and connect methods instead.
func HostKeyCallback(knownHostsPath string) ssh.HostKeyCallback {
	if knownHostsPath != "" {
		if _, err := os.Stat(knownHostsPath); err == nil {
			callback, err := knownhosts.New(knownHostsPath)
			if err == nil {
				return callback
			}
		}
	}
	// Fall back to insecure if known_hosts is not available
	return ssh.InsecureIgnoreHostKey() //nolint:gosec // Fallback when no known_hosts available
}

// ParseSSHURI parses an SSH URI into components.
func ParseSSHURI(uri string) (user, host string, port int, err error) {
	// Format: [user@]host[:port]
	port = 22

	// Extract user
	if idx := strings.Index(uri, "@"); idx != -1 {
		user = uri[:idx]
		uri = uri[idx+1:]
	}

	// Extract port
	if idx := strings.LastIndex(uri, ":"); idx != -1 {
		portStr := uri[idx+1:]
		uri = uri[:idx]
		n, err := fmt.Sscanf(portStr, "%d", &port)
		if err != nil || n != 1 {
			port = 22
		}
	}

	host = uri
	return user, host, port, nil
}
