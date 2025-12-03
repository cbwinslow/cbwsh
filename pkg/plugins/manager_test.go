package plugins_test

import (
	"context"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/plugins"
)

func TestBasePlugin(t *testing.T) {
	plugin := plugins.NewBasePlugin("test", core.PluginTypeCommand, "1.0.0")

	if plugin.Name() != "test" {
		t.Errorf("expected name 'test', got '%s'", plugin.Name())
	}

	if plugin.Type() != core.PluginTypeCommand {
		t.Errorf("expected command type")
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", plugin.Version())
	}

	if !plugin.Enabled() {
		t.Error("expected plugin to be enabled by default")
	}

	if plugin.Initialized() {
		t.Error("expected plugin to not be initialized by default")
	}
}

func TestPluginEnableDisable(t *testing.T) {
	plugin := plugins.NewBasePlugin("test", core.PluginTypeCommand, "1.0.0")

	err := plugin.Disable()
	if err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	if plugin.Enabled() {
		t.Error("expected plugin to be disabled")
	}

	err = plugin.Enable()
	if err != nil {
		t.Fatalf("enable failed: %v", err)
	}

	if !plugin.Enabled() {
		t.Error("expected plugin to be enabled")
	}
}

func TestPluginInitializeShutdown(t *testing.T) {
	plugin := plugins.NewBasePlugin("test", core.PluginTypeCommand, "1.0.0")
	ctx := context.Background()

	err := plugin.Initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	if !plugin.Initialized() {
		t.Error("expected plugin to be initialized")
	}

	err = plugin.Shutdown(ctx)
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if plugin.Initialized() {
		t.Error("expected plugin to not be initialized after shutdown")
	}
}

func TestPluginManager(t *testing.T) {
	manager := plugins.NewManager()

	plugin1 := plugins.NewBasePlugin("plugin1", core.PluginTypeCommand, "1.0.0")
	plugin2 := plugins.NewBasePlugin("plugin2", core.PluginTypeUI, "1.0.0")

	err := manager.Register(plugin1)
	if err != nil {
		t.Fatalf("register plugin1 failed: %v", err)
	}

	err = manager.Register(plugin2)
	if err != nil {
		t.Fatalf("register plugin2 failed: %v", err)
	}

	// Test duplicate registration
	err = manager.Register(plugin1)
	if err == nil {
		t.Error("expected error for duplicate registration")
	}

	// Test Get
	p, ok := manager.Get("plugin1")
	if !ok {
		t.Error("expected to find plugin1")
	}
	if p.Name() != "plugin1" {
		t.Errorf("expected plugin1, got %s", p.Name())
	}

	// Test List
	list := manager.List()
	if len(list) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(list))
	}

	// Test ListByType
	cmdPlugins := manager.ListByType(core.PluginTypeCommand)
	if len(cmdPlugins) != 1 {
		t.Errorf("expected 1 command plugin, got %d", len(cmdPlugins))
	}

	// Test Unregister
	err = manager.Unregister("plugin1")
	if err != nil {
		t.Fatalf("unregister failed: %v", err)
	}

	_, ok = manager.Get("plugin1")
	if ok {
		t.Error("expected plugin1 to be removed")
	}

	// Test unregister non-existent
	err = manager.Unregister("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent plugin")
	}
}

func TestPluginManagerInitializeShutdown(t *testing.T) {
	manager := plugins.NewManager()
	ctx := context.Background()

	plugin := plugins.NewBasePlugin("test", core.PluginTypeCommand, "1.0.0")
	err := manager.Register(plugin)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	err = manager.Initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Verify plugin was initialized
	p, _ := manager.Get("test")
	basePlugin := p.(*plugins.BasePlugin)
	if !basePlugin.Initialized() {
		t.Error("expected plugin to be initialized")
	}

	err = manager.Shutdown(ctx)
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}

func TestCommandPlugin(t *testing.T) {
	handler := func(args []string) (string, error) {
		return "executed: " + args[0], nil
	}

	plugin := plugins.NewCommandPlugin("test", "1.0.0", handler)

	result, err := plugin.Execute([]string{"arg1"})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if result != "executed: arg1" {
		t.Errorf("expected 'executed: arg1', got '%s'", result)
	}

	// Test disabled plugin
	err = plugin.Disable()
	if err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	_, err = plugin.Execute([]string{"arg1"})
	if err == nil {
		t.Error("expected error for disabled plugin")
	}
}

func TestHookPlugin(t *testing.T) {
	plugin := plugins.NewHookPlugin("test", "1.0.0")

	preExecuteCalled := false
	postExecuteCalled := false

	plugin.SetPreExecuteHook(func(cmd string) (string, error) {
		preExecuteCalled = true
		return cmd + " modified", nil
	})

	plugin.SetPostExecuteHook(func(result *core.CommandResult) error {
		postExecuteCalled = true
		return nil
	})

	// Test PreExecute
	modified, err := plugin.PreExecute("test")
	if err != nil {
		t.Fatalf("pre-execute failed: %v", err)
	}

	if !preExecuteCalled {
		t.Error("expected pre-execute hook to be called")
	}

	if modified != "test modified" {
		t.Errorf("expected 'test modified', got '%s'", modified)
	}

	// Test PostExecute
	result := &core.CommandResult{Command: "test", ExitCode: 0}
	err = plugin.PostExecute(result)
	if err != nil {
		t.Fatalf("post-execute failed: %v", err)
	}

	if !postExecuteCalled {
		t.Error("expected post-execute hook to be called")
	}
}

func TestFormatterPlugin(t *testing.T) {
	formatter := func(output string) (string, error) {
		return "[formatted] " + output, nil
	}

	plugin := plugins.NewFormatterPlugin("test", "1.0.0", formatter)

	result, err := plugin.Format("output")
	if err != nil {
		t.Fatalf("format failed: %v", err)
	}

	if result != "[formatted] output" {
		t.Errorf("expected '[formatted] output', got '%s'", result)
	}

	// Test disabled plugin
	err = plugin.Disable()
	if err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	result, err = plugin.Format("output")
	if err != nil {
		t.Fatalf("format failed: %v", err)
	}

	// Disabled plugin should return unmodified output
	if result != "output" {
		t.Errorf("expected 'output', got '%s'", result)
	}
}

func TestGlobalRegistry(t *testing.T) {
	plugins.GlobalRegistry.RegisterCreator("testPlugin", func() core.Plugin {
		return plugins.NewBasePlugin("testPlugin", core.PluginTypeCommand, "1.0.0")
	})

	available := plugins.GlobalRegistry.ListAvailable()
	found := false
	for _, name := range available {
		if name == "testPlugin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected testPlugin to be in available list")
	}

	plugin, ok := plugins.GlobalRegistry.Create("testPlugin")
	if !ok {
		t.Error("expected to create testPlugin")
	}

	if plugin.Name() != "testPlugin" {
		t.Errorf("expected name 'testPlugin', got '%s'", plugin.Name())
	}

	// Test non-existent creator
	_, ok = plugins.GlobalRegistry.Create("nonexistent")
	if ok {
		t.Error("expected false for non-existent creator")
	}
}
