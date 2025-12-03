// Package plugins provides plugin management for cbwsh.
package plugins

import (
	"context"
	"fmt"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// BasePlugin provides a base implementation for plugins.
type BasePlugin struct {
	mu          sync.RWMutex
	name        string
	pluginType  core.PluginType
	version     string
	enabled     bool
	initialized bool
}

// NewBasePlugin creates a new base plugin.
func NewBasePlugin(name string, pluginType core.PluginType, version string) *BasePlugin {
	return &BasePlugin{
		name:       name,
		pluginType: pluginType,
		version:    version,
		enabled:    true,
	}
}

// Name returns the plugin name.
func (p *BasePlugin) Name() string {
	return p.name
}

// Type returns the plugin type.
func (p *BasePlugin) Type() core.PluginType {
	return p.pluginType
}

// Version returns the plugin version.
func (p *BasePlugin) Version() string {
	return p.version
}

// Initialize sets up the plugin.
func (p *BasePlugin) Initialize(_ context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.initialized = true
	return nil
}

// Shutdown cleans up the plugin.
func (p *BasePlugin) Shutdown(_ context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.initialized = false
	return nil
}

// Enabled returns whether the plugin is enabled.
func (p *BasePlugin) Enabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// Enable enables the plugin.
func (p *BasePlugin) Enable() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = true
	return nil
}

// Disable disables the plugin.
func (p *BasePlugin) Disable() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = false
	return nil
}

// Initialized returns whether the plugin is initialized.
func (p *BasePlugin) Initialized() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.initialized
}

// Manager manages plugins.
type Manager struct {
	mu      sync.RWMutex
	plugins map[string]core.Plugin
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]core.Plugin),
	}
}

// Register adds a plugin to the manager.
func (m *Manager) Register(plugin core.Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[plugin.Name()]; exists {
		return fmt.Errorf("plugin already registered: %s", plugin.Name())
	}

	m.plugins[plugin.Name()] = plugin
	return nil
}

// Unregister removes a plugin from the manager.
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	delete(m.plugins, name)
	return nil
}

// Get returns a plugin by name.
func (m *Manager) Get(name string) (core.Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, exists := m.plugins[name]
	return plugin, exists
}

// List returns all registered plugins.
func (m *Manager) List() []core.Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]core.Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		result = append(result, plugin)
	}
	return result
}

// ListByType returns plugins of a specific type.
func (m *Manager) ListByType(pluginType core.PluginType) []core.Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []core.Plugin
	for _, plugin := range m.plugins {
		if plugin.Type() == pluginType {
			result = append(result, plugin)
		}
	}
	return result
}

// Initialize initializes all registered plugins.
func (m *Manager) Initialize(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, plugin := range m.plugins {
		if plugin.Enabled() {
			if err := plugin.Initialize(ctx); err != nil {
				return fmt.Errorf("failed to initialize plugin %s: %w", plugin.Name(), err)
			}
		}
	}
	return nil
}

// Shutdown shuts down all registered plugins.
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for _, plugin := range m.plugins {
		if err := plugin.Shutdown(ctx); err != nil {
			lastErr = fmt.Errorf("failed to shutdown plugin %s: %w", plugin.Name(), err)
		}
	}
	return lastErr
}

// CommandPlugin is a plugin that provides commands.
type CommandPlugin struct {
	*BasePlugin
	handler func(args []string) (string, error)
}

// NewCommandPlugin creates a new command plugin.
func NewCommandPlugin(name, version string, handler func(args []string) (string, error)) *CommandPlugin {
	return &CommandPlugin{
		BasePlugin: NewBasePlugin(name, core.PluginTypeCommand, version),
		handler:    handler,
	}
}

// Execute executes the command.
func (p *CommandPlugin) Execute(args []string) (string, error) {
	if !p.Enabled() {
		return "", fmt.Errorf("plugin %s is disabled", p.Name())
	}
	return p.handler(args)
}

// HookPlugin is a plugin that provides lifecycle hooks.
type HookPlugin struct {
	*BasePlugin
	onPreExecute  func(command string) (string, error)
	onPostExecute func(result *core.CommandResult) error
}

// NewHookPlugin creates a new hook plugin.
func NewHookPlugin(name, version string) *HookPlugin {
	return &HookPlugin{
		BasePlugin: NewBasePlugin(name, core.PluginTypeHook, version),
	}
}

// SetPreExecuteHook sets the pre-execute hook.
func (p *HookPlugin) SetPreExecuteHook(hook func(command string) (string, error)) {
	p.onPreExecute = hook
}

// SetPostExecuteHook sets the post-execute hook.
func (p *HookPlugin) SetPostExecuteHook(hook func(result *core.CommandResult) error) {
	p.onPostExecute = hook
}

// PreExecute runs before command execution.
func (p *HookPlugin) PreExecute(command string) (string, error) {
	if !p.Enabled() || p.onPreExecute == nil {
		return command, nil
	}
	return p.onPreExecute(command)
}

// PostExecute runs after command execution.
func (p *HookPlugin) PostExecute(result *core.CommandResult) error {
	if !p.Enabled() || p.onPostExecute == nil {
		return nil
	}
	return p.onPostExecute(result)
}

// FormatterPlugin is a plugin that formats output.
type FormatterPlugin struct {
	*BasePlugin
	formatter func(output string) (string, error)
}

// NewFormatterPlugin creates a new formatter plugin.
func NewFormatterPlugin(name, version string, formatter func(output string) (string, error)) *FormatterPlugin {
	return &FormatterPlugin{
		BasePlugin: NewBasePlugin(name, core.PluginTypeFormatter, version),
		formatter:  formatter,
	}
}

// Format formats the output.
func (p *FormatterPlugin) Format(output string) (string, error) {
	if !p.Enabled() {
		return output, nil
	}
	return p.formatter(output)
}

// Registry provides global plugin registration.
type Registry struct {
	mu       sync.RWMutex
	creators map[string]func() core.Plugin
}

// GlobalRegistry is the global plugin registry.
var GlobalRegistry = &Registry{
	creators: make(map[string]func() core.Plugin),
}

// RegisterCreator registers a plugin creator function.
func (r *Registry) RegisterCreator(name string, creator func() core.Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creators[name] = creator
}

// Create creates a plugin by name.
func (r *Registry) Create(name string) (core.Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	creator, exists := r.creators[name]
	if !exists {
		return nil, false
	}
	return creator(), true
}

// ListAvailable returns all available plugin names.
func (r *Registry) ListAvailable() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]string, 0, len(r.creators))
	for name := range r.creators {
		result = append(result, name)
	}
	return result
}
