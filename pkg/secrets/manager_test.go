package secrets_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/secrets"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	
	manager := secrets.NewManager("/tmp/test_secrets.enc")
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		t.Error("expected secrets file to be created")
	}
}

func TestSetAndGet(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Set a secret
	testKey := "api_key"
	testValue := []byte("super-secret-value")
	
	err = manager.Store(testKey, testValue)
	if err != nil {
		t.Fatalf("failed to store secret: %v", err)
	}
	
	// Get the secret
	value, err := manager.Retrieve(testKey)
	if err != nil {
		t.Fatalf("failed to retrieve secret: %v", err)
	}
	
	if string(value) != string(testValue) {
		t.Errorf("expected %s, got %s", testValue, value)
	}
}

func TestGetNonExistent(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Try to get non-existent secret
	_, err = manager.Retrieve("non-existent")
	if err == nil {
		t.Error("expected error for non-existent secret")
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Set and then delete
	testKey := "temp_key"
	err = manager.Store(testKey, []byte("temp-value"))
	if err != nil {
		t.Fatalf("failed to store secret: %v", err)
	}
	
	err = manager.Delete(testKey)
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}
	
	// Verify it's gone
	_, err = manager.Retrieve(testKey)
	if err == nil {
		t.Error("expected error after deleting secret")
	}
}

func TestList(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Add multiple secrets
	secrets := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}
	
	for k, v := range secrets {
		if err := manager.Store(k, v); err != nil {
			t.Fatalf("failed to store secret %s: %v", k, err)
		}
	}
	
	// List all keys
	keys, err := manager.List()
	if err != nil {
		t.Fatalf("failed to list keys: %v", err)
	}
	if len(keys) != len(secrets) {
		t.Errorf("expected %d keys, got %d", len(secrets), len(keys))
	}
	
	// Verify all keys are present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	
	for k := range secrets {
		if !keyMap[k] {
			t.Errorf("expected key %s in list", k)
		}
	}
}

func TestPersistence(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	password := "test-password"
	
	// Create manager and add secrets
	manager1 := secrets.NewManager(storePath)
	err := manager1.Initialize(password)
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	testKey := "persistent_key"
	testValue := []byte("persistent-value")
	
	err = manager1.Store(testKey, testValue)
	if err != nil {
		t.Fatalf("failed to store secret: %v", err)
	}
	
	// Store is automatically saved
	
	// Create new manager and unlock
	manager2 := secrets.NewManager(storePath)
	err = manager2.Unlock(password)
	if err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}
	
	// Verify secret persisted
	value, err := manager2.Retrieve(testKey)
	if err != nil {
		t.Fatalf("failed to retrieve secret: %v", err)
	}
	
	if string(value) != string(testValue) {
		t.Errorf("expected %s, got %s", testValue, value)
	}
}

func TestWrongPassword(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	// Create with one password
	manager1 := secrets.NewManager(storePath)
	err := manager1.Initialize("correct-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	err = manager1.Store("key", []byte("value"))
	if err != nil {
		t.Fatalf("failed to store secret: %v", err)
	}
	
	// Try to unlock with wrong password
	manager2 := secrets.NewManager(storePath)
	err = manager2.Unlock("wrong-password")
	if err == nil {
		t.Error("expected error when unlocking with wrong password")
	}
}

func TestDeleteAll(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Add secrets
	manager.Store("key1", []byte("value1"))
	manager.Store("key2", []byte("value2"))
	
	// Delete all
	keys, _ := manager.List()
	for _, k := range keys {
		err = manager.Delete(k)
		if err != nil {
			t.Fatalf("failed to delete %s: %v", k, err)
		}
	}
	
	// Verify all are gone
	keys, err = manager.List()
	if err != nil {
		t.Fatalf("failed to list keys: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys after delete all, got %d", len(keys))
	}
}

func TestLockUnlock(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	password := "test-password"
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize(password)
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Add a secret
	manager.Store("key", []byte("value"))
	
	// Lock the manager
	err = manager.Lock()
	if err != nil {
		t.Fatalf("failed to lock: %v", err)
	}
	
	// Try to access secret while locked - should fail
	_, err = manager.Retrieve("key")
	if err == nil {
		t.Error("expected error when accessing locked secrets")
	}
	
	// Unlock
	err = manager.Unlock(password)
	if err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}
	
	// Should be able to access now
	value, err := manager.Retrieve("key")
	if err != nil {
		t.Fatalf("failed to retrieve secret after unlock: %v", err)
	}
	
	if string(value) != "value" {
		t.Errorf("expected 'value', got %s", value)
	}
}

func TestMultipleSecrets(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Add many secrets
	count := 100
	for i := 0; i < count; i++ {
		key := string(rune('a' + (i % 26))) + string(rune('0' + (i / 26)))
		value := []byte("value-" + key)
		
		err = manager.Store(key, value)
		if err != nil {
			t.Fatalf("failed to store secret %s: %v", key, err)
		}
	}
	
	// Verify all can be retrieved
	keys, err := manager.List()
	if err != nil {
		t.Fatalf("failed to list keys: %v", err)
	}
	if len(keys) != count {
		t.Errorf("expected %d keys, got %d", count, len(keys))
	}
}

func TestEmptyValue(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Set empty value
	err = manager.Store("empty", []byte(""))
	if err != nil {
		t.Fatalf("failed to store empty secret: %v", err)
	}
	
	// Get empty value
	value, err := manager.Retrieve("empty")
	if err != nil {
		t.Fatalf("failed to retrieve empty secret: %v", err)
	}
	
	if len(value) != 0 {
		t.Errorf("expected empty value, got %s", value)
	}
}

func TestBinaryData(t *testing.T) {
	t.Parallel()
	
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "secrets.enc")
	
	manager := secrets.NewManager(storePath)
	err := manager.Initialize("test-password")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	
	// Test with binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	
	err = manager.Store("binary", binaryData)
	if err != nil {
		t.Fatalf("failed to store binary secret: %v", err)
	}
	
	value, err := manager.Retrieve("binary")
	if err != nil {
		t.Fatalf("failed to retrieve binary secret: %v", err)
	}
	
	if len(value) != len(binaryData) {
		t.Errorf("expected length %d, got %d", len(binaryData), len(value))
	}
	
	for i, b := range binaryData {
		if value[i] != b {
			t.Errorf("byte %d: expected %x, got %x", i, b, value[i])
		}
	}
}
