// Package secrets provides secure secrets management for cbwsh.
package secrets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// EncryptionBackend represents the encryption method to use.
type EncryptionBackend string

const (
	// BackendAES uses the built-in AES-256-GCM encryption.
	BackendAES EncryptionBackend = "aes"
	// BackendAge uses age encryption.
	BackendAge EncryptionBackend = "age"
	// BackendGPG uses GPG encryption.
	BackendGPG EncryptionBackend = "gpg"
)

// GitBackend represents the git management tool to use.
type GitBackend string

const (
	// GitBackendNone uses no git integration.
	GitBackendNone GitBackend = "none"
	// GitBackendGit uses standard git.
	GitBackendGit GitBackend = "git"
	// GitBackendYadm uses yadm (Yet Another Dotfiles Manager).
	GitBackendYadm GitBackend = "yadm"
)

// ExtendedManager provides extended secrets management with multiple backends.
type ExtendedManager struct {
	mu                sync.RWMutex
	baseManager       *Manager
	encryptionBackend EncryptionBackend
	gitBackend        GitBackend
	storePath         string
	recipientKey      string // For age/GPG
	gitRepoPath       string
}

// NewExtendedManager creates a new extended secrets manager.
func NewExtendedManager(storePath string, backend EncryptionBackend, gitBackend GitBackend) *ExtendedManager {
	return &ExtendedManager{
		baseManager:       NewManager(storePath),
		encryptionBackend: backend,
		gitBackend:        gitBackend,
		storePath:         storePath,
	}
}

// SetRecipientKey sets the recipient key for age/GPG encryption.
func (m *ExtendedManager) SetRecipientKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recipientKey = key
}

// SetGitRepoPath sets the git repository path.
func (m *ExtendedManager) SetGitRepoPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gitRepoPath = path
}

// Initialize initializes the extended manager.
func (m *ExtendedManager) Initialize(masterPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Initialize(masterPassword)
	case BackendAge:
		return m.initializeAge()
	case BackendGPG:
		return m.initializeGPG()
	default:
		return m.baseManager.Initialize(masterPassword)
	}
}

func (m *ExtendedManager) initializeAge() error {
	// Check if age is available
	if _, err := exec.LookPath("age"); err != nil {
		return fmt.Errorf("age not found in PATH: %w", err)
	}

	// Create store directory
	dir := filepath.Dir(m.storePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}

	// Create empty secrets file
	return os.WriteFile(m.storePath, []byte("{}"), 0o600)
}

func (m *ExtendedManager) initializeGPG() error {
	// Check if gpg is available
	if _, err := exec.LookPath("gpg"); err != nil {
		return fmt.Errorf("gpg not found in PATH: %w", err)
	}

	// Create store directory
	dir := filepath.Dir(m.storePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}

	// Create empty secrets file
	return os.WriteFile(m.storePath, []byte("{}"), 0o600)
}

// Unlock unlocks the secrets store.
func (m *ExtendedManager) Unlock(masterPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Unlock(masterPassword)
	case BackendAge, BackendGPG:
		// Age and GPG handle their own authentication
		return nil
	default:
		return m.baseManager.Unlock(masterPassword)
	}
}

// Lock locks the secrets store.
func (m *ExtendedManager) Lock() error {
	return m.baseManager.Lock()
}

// IsUnlocked returns whether the store is unlocked.
func (m *ExtendedManager) IsUnlocked() bool {
	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.IsUnlocked()
	case BackendAge, BackendGPG:
		return true // External tools manage authentication
	default:
		return m.baseManager.IsUnlocked()
	}
}

// Store securely stores a secret.
func (m *ExtendedManager) Store(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Store(key, value)
	case BackendAge:
		return m.storeWithAge(key, value)
	case BackendGPG:
		return m.storeWithGPG(key, value)
	default:
		return m.baseManager.Store(key, value)
	}
}

func (m *ExtendedManager) storeWithAge(key string, value []byte) error {
	if m.recipientKey == "" {
		return fmt.Errorf("age recipient key not set")
	}

	// Create a temporary file for input
	tmpFile, err := os.CreateTemp("", "cbwsh-secret-*")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(value); err != nil {
		return err
	}
	tmpFile.Close()

	// Encrypt with age
	secretPath := filepath.Join(filepath.Dir(m.storePath), key+".age")
	cmd := exec.Command("age", "-r", m.recipientKey, "-o", secretPath, tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("age encryption failed: %s: %w", string(output), err)
	}

	return m.commitToGit(secretPath, "Add secret: "+key)
}

func (m *ExtendedManager) storeWithGPG(key string, value []byte) error {
	if m.recipientKey == "" {
		return fmt.Errorf("gpg recipient key not set")
	}

	// Create a temporary file for input
	tmpFile, err := os.CreateTemp("", "cbwsh-secret-*")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(value); err != nil {
		return err
	}
	tmpFile.Close()

	// Encrypt with GPG
	secretPath := filepath.Join(filepath.Dir(m.storePath), key+".gpg")
	cmd := exec.Command("gpg", "--encrypt", "--recipient", m.recipientKey, "--output", secretPath, tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gpg encryption failed: %s: %w", string(output), err)
	}

	return m.commitToGit(secretPath, "Add secret: "+key)
}

// Retrieve gets a stored secret.
func (m *ExtendedManager) Retrieve(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Retrieve(key)
	case BackendAge:
		return m.retrieveWithAge(key)
	case BackendGPG:
		return m.retrieveWithGPG(key)
	default:
		return m.baseManager.Retrieve(key)
	}
}

func (m *ExtendedManager) retrieveWithAge(key string) ([]byte, error) {
	secretPath := filepath.Join(filepath.Dir(m.storePath), key+".age")

	cmd := exec.Command("age", "--decrypt", secretPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("age decryption failed: %w", err)
	}

	return output, nil
}

func (m *ExtendedManager) retrieveWithGPG(key string) ([]byte, error) {
	secretPath := filepath.Join(filepath.Dir(m.storePath), key+".gpg")

	cmd := exec.Command("gpg", "--decrypt", secretPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gpg decryption failed: %w", err)
	}

	return output, nil
}

// Delete removes a stored secret.
func (m *ExtendedManager) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Delete(key)
	case BackendAge:
		secretPath := filepath.Join(filepath.Dir(m.storePath), key+".age")
		if err := os.Remove(secretPath); err != nil {
			return err
		}
		return m.commitToGit(secretPath, "Remove secret: "+key)
	case BackendGPG:
		secretPath := filepath.Join(filepath.Dir(m.storePath), key+".gpg")
		if err := os.Remove(secretPath); err != nil {
			return err
		}
		return m.commitToGit(secretPath, "Remove secret: "+key)
	default:
		return m.baseManager.Delete(key)
	}
}

// List returns all stored secret keys.
func (m *ExtendedManager) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.List()
	case BackendAge:
		return m.listFiles(".age")
	case BackendGPG:
		return m.listFiles(".gpg")
	default:
		return m.baseManager.List()
	}
}

func (m *ExtendedManager) listFiles(extension string) ([]string, error) {
	dir := filepath.Dir(m.storePath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, extension) {
			keys = append(keys, strings.TrimSuffix(name, extension))
		}
	}

	return keys, nil
}

// Exists checks if a secret exists.
func (m *ExtendedManager) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch m.encryptionBackend {
	case BackendAES:
		return m.baseManager.Exists(key)
	case BackendAge:
		secretPath := filepath.Join(filepath.Dir(m.storePath), key+".age")
		_, err := os.Stat(secretPath)
		return err == nil
	case BackendGPG:
		secretPath := filepath.Join(filepath.Dir(m.storePath), key+".gpg")
		_, err := os.Stat(secretPath)
		return err == nil
	default:
		return m.baseManager.Exists(key)
	}
}

func (m *ExtendedManager) commitToGit(filePath, message string) error {
	if m.gitBackend == GitBackendNone {
		return nil
	}

	var gitCmd string
	switch m.gitBackend {
	case GitBackendGit:
		gitCmd = "git"
	case GitBackendYadm:
		gitCmd = "yadm"
	default:
		return nil
	}

	ctx := context.Background()

	// Add the file
	addCmd := exec.CommandContext(ctx, gitCmd, "add", filePath)
	addCmd.Dir = m.gitRepoPath
	if output, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s add failed: %s: %w", gitCmd, string(output), err)
	}

	// Commit
	commitCmd := exec.CommandContext(ctx, gitCmd, "commit", "-m", message)
	commitCmd.Dir = m.gitRepoPath
	if output, err := commitCmd.CombinedOutput(); err != nil {
		// Git commit returns exit code 1 when there's nothing to commit,
		// which is not an error in our use case. We check the exit code
		// to distinguish between this and actual errors.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// Check if it's truly "nothing to commit" by looking at output
			if strings.Contains(string(output), "nothing to commit") {
				return nil
			}
		}
		return fmt.Errorf("%s commit failed: %s: %w", gitCmd, string(output), err)
	}

	return nil
}

// SyncWithGit syncs secrets with the git repository.
func (m *ExtendedManager) SyncWithGit() error {
	if m.gitBackend == GitBackendNone {
		return nil
	}

	var gitCmd string
	switch m.gitBackend {
	case GitBackendGit:
		gitCmd = "git"
	case GitBackendYadm:
		gitCmd = "yadm"
	default:
		return nil
	}

	ctx := context.Background()

	// Pull latest changes
	pullCmd := exec.CommandContext(ctx, gitCmd, "pull", "--rebase")
	pullCmd.Dir = m.gitRepoPath
	if output, err := pullCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s pull failed: %s: %w", gitCmd, string(output), err)
	}

	return nil
}

// PushToGit pushes secrets to the git repository.
func (m *ExtendedManager) PushToGit() error {
	if m.gitBackend == GitBackendNone {
		return nil
	}

	var gitCmd string
	switch m.gitBackend {
	case GitBackendGit:
		gitCmd = "git"
	case GitBackendYadm:
		gitCmd = "yadm"
	default:
		return nil
	}

	ctx := context.Background()

	// Push changes
	pushCmd := exec.CommandContext(ctx, gitCmd, "push")
	pushCmd.Dir = m.gitRepoPath
	if output, err := pushCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s push failed: %s: %w", gitCmd, string(output), err)
	}

	return nil
}

// GetEncryptionBackend returns the current encryption backend.
func (m *ExtendedManager) GetEncryptionBackend() EncryptionBackend {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.encryptionBackend
}

// GetGitBackend returns the current git backend.
func (m *ExtendedManager) GetGitBackend() GitBackend {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gitBackend
}

// APIKeyManager provides specialized management for API keys.
type APIKeyManager struct {
	*ExtendedManager
	keyPrefix string
}

// NewAPIKeyManager creates a new API key manager.
func NewAPIKeyManager(manager *ExtendedManager) *APIKeyManager {
	return &APIKeyManager{
		ExtendedManager: manager,
		keyPrefix:       "apikey_",
	}
}

// StoreAPIKey stores an API key.
func (m *APIKeyManager) StoreAPIKey(name string, key []byte) error {
	return m.Store(m.keyPrefix+name, key)
}

// GetAPIKey retrieves an API key.
func (m *APIKeyManager) GetAPIKey(name string) ([]byte, error) {
	return m.Retrieve(m.keyPrefix + name)
}

// DeleteAPIKey deletes an API key.
func (m *APIKeyManager) DeleteAPIKey(name string) error {
	return m.Delete(m.keyPrefix + name)
}

// ListAPIKeys returns all stored API key names.
func (m *APIKeyManager) ListAPIKeys() ([]string, error) {
	keys, err := m.List()
	if err != nil {
		return nil, err
	}

	var apiKeys []string
	for _, key := range keys {
		if strings.HasPrefix(key, m.keyPrefix) {
			apiKeys = append(apiKeys, strings.TrimPrefix(key, m.keyPrefix))
		}
	}

	return apiKeys, nil
}

// APIKeyExists checks if an API key exists.
func (m *APIKeyManager) APIKeyExists(name string) bool {
	return m.Exists(m.keyPrefix + name)
}
