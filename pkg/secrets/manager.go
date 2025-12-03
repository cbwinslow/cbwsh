// Package secrets provides secure secrets management for cbwsh.
package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/argon2"
)

// Manager provides encrypted secrets storage and retrieval.
type Manager struct {
	mu            sync.RWMutex
	storePath     string
	masterKeyHash []byte
	encryptionKey []byte
	secrets       map[string][]byte
	unlocked      bool
}

// Argon2 parameters for key derivation.
const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
)

// NewManager creates a new secrets manager.
func NewManager(storePath string) *Manager {
	return &Manager{
		storePath: storePath,
		secrets:   make(map[string][]byte),
	}
}

// Initialize sets up the secrets store with a master password.
func (m *Manager) Initialize(masterPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate salt for key derivation
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key from master password
	m.encryptionKey = argon2.IDKey(
		[]byte(masterPassword),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Store hash for verification
	hash := sha256.Sum256(m.encryptionKey)
	m.masterKeyHash = hash[:]

	m.secrets = make(map[string][]byte)
	m.unlocked = true

	// Create store directory if needed
	dir := filepath.Dir(m.storePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}

	// Save initial state with salt
	return m.saveStore(salt)
}

// Unlock unlocks the secrets store with the master password.
func (m *Manager) Unlock(masterPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Read store file
	data, err := os.ReadFile(m.storePath)
	if err != nil {
		return fmt.Errorf("failed to read store: %w", err)
	}

	// Parse store structure
	var store storeData
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("failed to parse store: %w", err)
	}

	// Decode salt
	salt, err := base64.StdEncoding.DecodeString(store.Salt)
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	// Derive key from password
	m.encryptionKey = argon2.IDKey(
		[]byte(masterPassword),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Verify key hash
	hash := sha256.Sum256(m.encryptionKey)
	m.masterKeyHash = hash[:]

	storedHash, err := base64.StdEncoding.DecodeString(store.KeyHash)
	if err != nil {
		return fmt.Errorf("failed to decode key hash: %w", err)
	}

	if !equalBytes(m.masterKeyHash, storedHash) {
		m.encryptionKey = nil
		m.masterKeyHash = nil
		return errors.New("invalid master password")
	}

	// Decrypt secrets
	m.secrets = make(map[string][]byte)
	for key, encValue := range store.Secrets {
		encData, err := base64.StdEncoding.DecodeString(encValue)
		if err != nil {
			continue
		}
		plaintext, err := m.decrypt(encData)
		if err != nil {
			continue
		}
		m.secrets[key] = plaintext
	}

	m.unlocked = true
	return nil
}

// Lock locks the secrets store.
func (m *Manager) Lock() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear sensitive data
	if m.encryptionKey != nil {
		for i := range m.encryptionKey {
			m.encryptionKey[i] = 0
		}
		m.encryptionKey = nil
	}

	m.secrets = make(map[string][]byte)
	m.unlocked = false
	return nil
}

// IsUnlocked returns whether the store is unlocked.
func (m *Manager) IsUnlocked() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.unlocked
}

// Store securely stores a secret.
func (m *Manager) Store(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.unlocked {
		return errors.New("secrets store is locked")
	}

	m.secrets[key] = value

	// Read existing store to get salt
	data, err := os.ReadFile(m.storePath)
	if err != nil {
		return fmt.Errorf("failed to read store: %w", err)
	}

	var store storeData
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("failed to parse store: %w", err)
	}

	salt, _ := base64.StdEncoding.DecodeString(store.Salt)
	return m.saveStore(salt)
}

// Retrieve gets a stored secret.
func (m *Manager) Retrieve(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.unlocked {
		return nil, errors.New("secrets store is locked")
	}

	value, exists := m.secrets[key]
	if !exists {
		return nil, fmt.Errorf("secret not found: %s", key)
	}

	// Return a copy to prevent modification
	result := make([]byte, len(value))
	copy(result, value)
	return result, nil
}

// Delete removes a stored secret.
func (m *Manager) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.unlocked {
		return errors.New("secrets store is locked")
	}

	delete(m.secrets, key)

	// Read existing store to get salt
	data, err := os.ReadFile(m.storePath)
	if err != nil {
		return fmt.Errorf("failed to read store: %w", err)
	}

	var store storeData
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("failed to parse store: %w", err)
	}

	salt, _ := base64.StdEncoding.DecodeString(store.Salt)
	return m.saveStore(salt)
}

// List returns all stored secret keys.
func (m *Manager) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.unlocked {
		return nil, errors.New("secrets store is locked")
	}

	result := make([]string, 0, len(m.secrets))
	for key := range m.secrets {
		result = append(result, key)
	}
	return result, nil
}

// Exists checks if a secret exists.
func (m *Manager) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.unlocked {
		return false
	}

	_, exists := m.secrets[key]
	return exists
}

// storeData is the on-disk format for the secrets store.
type storeData struct {
	Salt    string            `json:"salt"`
	KeyHash string            `json:"key_hash"`
	Secrets map[string]string `json:"secrets"`
}

func (m *Manager) saveStore(salt []byte) error {
	store := storeData{
		Salt:    base64.StdEncoding.EncodeToString(salt),
		KeyHash: base64.StdEncoding.EncodeToString(m.masterKeyHash),
		Secrets: make(map[string]string),
	}

	for key, value := range m.secrets {
		encrypted, err := m.encrypt(value)
		if err != nil {
			return fmt.Errorf("failed to encrypt secret %s: %w", key, err)
		}
		store.Secrets[key] = base64.StdEncoding.EncodeToString(encrypted)
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal store: %w", err)
	}

	return os.WriteFile(m.storePath, data, 0o600)
}

func (m *Manager) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (m *Manager) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// ChangePassword changes the master password.
func (m *Manager) ChangePassword(oldPassword, newPassword string) error {
	// First verify old password by unlocking
	if err := m.Unlock(oldPassword); err != nil {
		return fmt.Errorf("invalid old password: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate new salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive new encryption key
	m.encryptionKey = argon2.IDKey(
		[]byte(newPassword),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Update key hash
	hash := sha256.Sum256(m.encryptionKey)
	m.masterKeyHash = hash[:]

	// Save with new encryption
	return m.saveStore(salt)
}
