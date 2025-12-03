package effects_test

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/ui/effects"
)

func TestFluidSimulation(t *testing.T) {
	t.Parallel()
	sim := effects.NewFluidSimulation(10, 10, 0.1, 0.2)
	if sim == nil {
		t.Fatal("expected non-nil fluid simulation")
	}

	// Test adding density
	sim.AddDensity(5, 5, 1.0)
	density := sim.GetDensity(5, 5)
	if density != 1.0 {
		t.Errorf("expected density 1.0, got %f", density)
	}

	// Test step
	sim.Step()
	// After step, density should still be present but may have diffused
	density = sim.GetDensity(5, 5)
	if density < 0 {
		t.Error("expected non-negative density after step")
	}

	// Test render
	output := sim.Render()
	if output == "" {
		t.Error("expected non-empty render output")
	}

	// Test colored render
	coloredOutput := sim.RenderColored()
	if coloredOutput == "" {
		t.Error("expected non-empty colored render output")
	}
}

func TestFluidSimulationColorScheme(t *testing.T) {
	t.Parallel()
	sim := effects.NewFluidSimulation(10, 10, 0.1, 0.2)

	// Set fire colors
	sim.SetColorScheme(effects.DefaultFireColors)

	// Add density and render
	sim.AddDensity(5, 5, 0.5)
	output := sim.RenderColored()
	if output == "" {
		t.Error("expected non-empty colored render output")
	}
}

func TestWaterEffect(t *testing.T) {
	t.Parallel()
	water := effects.NewWaterEffect(20, 10)
	if water == nil {
		t.Fatal("expected non-nil water effect")
	}

	// Test update
	water.Update()

	// Test render
	output := water.Render()
	if output == "" {
		t.Error("expected non-empty render output")
	}

	// Test colored render
	coloredOutput := water.RenderColored()
	if coloredOutput == "" {
		t.Error("expected non-empty colored render output")
	}

	// Test disturb
	water.Disturb(10, 0.5)
	water.Update()
}

func TestWaterEffectColors(t *testing.T) {
	t.Parallel()
	water := effects.NewWaterEffect(20, 10)

	// Set custom colors
	water.SetColors(effects.ColorIntensity{
		Low:    "20",
		Medium: "40",
		High:   "60",
	})

	water.Update()
	output := water.RenderColored()
	if output == "" {
		t.Error("expected non-empty colored render output")
	}
}

func TestParticleSystem(t *testing.T) {
	t.Parallel()
	ps := effects.NewParticleSystem(40, 20, 100)
	if ps == nil {
		t.Fatal("expected non-nil particle system")
	}

	// Test emit
	ps.Emit(20, 10, 10, 0.5, 2.0)
	count := ps.ParticleCount()
	if count != 10 {
		t.Errorf("expected 10 particles, got %d", count)
	}

	// Test update
	ps.Update()
	// Particles should still exist after one update
	if ps.ParticleCount() == 0 {
		t.Error("expected particles after update")
	}

	// Test render
	output := ps.Render()
	if output == "" {
		t.Error("expected non-empty render output")
	}

	// Test colored render
	coloredOutput := ps.RenderColored()
	if coloredOutput == "" {
		t.Error("expected non-empty colored render output")
	}
}

func TestParticleSystemColors(t *testing.T) {
	t.Parallel()
	ps := effects.NewParticleSystem(40, 20, 100)

	// Set custom colors
	ps.SetColors(effects.ColorIntensity{
		Low:    "196",
		Medium: "208",
		High:   "226",
	})

	ps.Emit(20, 10, 5, 0.3, 1.5)
	output := ps.RenderColored()
	if output == "" {
		t.Error("expected non-empty colored render output")
	}
}

func TestParticleSystemGravity(t *testing.T) {
	t.Parallel()
	ps := effects.NewParticleSystem(40, 20, 100)

	// Set gravity
	ps.SetGravity(0.5)

	ps.Emit(20, 10, 5, 0.3, 1.5)

	// Update multiple times
	for i := 0; i < 10; i++ {
		ps.Update()
	}

	// Particles should have moved
	output := ps.Render()
	if output == "" {
		t.Error("expected non-empty render output")
	}
}

func TestDefaultColors(t *testing.T) {
	t.Parallel()
	// Test that default colors are defined
	if effects.DefaultWaterColors.Low == "" {
		t.Error("expected non-empty water low color")
	}
	if effects.DefaultWaterColors.Medium == "" {
		t.Error("expected non-empty water medium color")
	}
	if effects.DefaultWaterColors.High == "" {
		t.Error("expected non-empty water high color")
	}

	if effects.DefaultFireColors.Low == "" {
		t.Error("expected non-empty fire low color")
	}
	if effects.DefaultFireColors.Medium == "" {
		t.Error("expected non-empty fire medium color")
	}
	if effects.DefaultFireColors.High == "" {
		t.Error("expected non-empty fire high color")
	}
}

func TestFluidSimulationVelocity(t *testing.T) {
	t.Parallel()
	sim := effects.NewFluidSimulation(10, 10, 0.1, 0.2)

	// Add density and velocity
	sim.AddDensity(5, 5, 1.0)
	sim.AddVelocity(5, 5, 1.0, 0.5)

	// Step multiple times
	for i := 0; i < 5; i++ {
		sim.Step()
	}

	// Render should work
	output := sim.RenderColored()
	if output == "" {
		t.Error("expected non-empty render output")
	}
}
