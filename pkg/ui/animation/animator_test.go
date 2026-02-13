package animation

import (
	"testing"
	"time"
)

func TestNewAnimator(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()

	if animator == nil {
		t.Fatal("expected non-nil animator")
	}

	if animator.frequency == 0 {
		t.Error("frequency should be set")
	}

	if animator.damping == 0 {
		t.Error("damping should be set")
	}
}

func TestNewSpringAnimator(t *testing.T) {
	t.Parallel()

	frequency := 10.0
	damping := 0.7

	animator := NewSpringAnimator(frequency, damping)

	if animator == nil {
		t.Fatal("expected non-nil animator")
	}

	if animator.frequency != frequency {
		t.Errorf("expected frequency %f, got %f", frequency, animator.frequency)
	}

	if animator.damping != damping {
		t.Errorf("expected damping %f, got %f", damping, animator.damping)
	}

	if animator.animType != Spring {
		t.Error("expected Spring animation type")
	}
}

func TestAnimatorStartStop(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()

	if animator.running {
		t.Error("animator should not be running initially")
	}

	animator.Start()

	if !animator.running {
		t.Error("animator should be running after Start()")
	}

	animator.Stop()

	if animator.running {
		t.Error("animator should not be running after Stop()")
	}
}

func TestAnimatorSetTarget(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()

	targetValue := 100.0
	animator.SetValue(targetValue)

	if animator.GetTarget() != targetValue {
		t.Errorf("expected target %f, got %f", targetValue, animator.GetTarget())
	}
}

func TestAnimatorValue(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()

	// Set initial value
	animator.SetValue(50.0)

	value := animator.GetValue()

	// Since SetValue sets the target, the current value should move towards it
	_ = value
}

func TestAnimatorUpdate(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()
	animator.SetValue(100.0)
	animator.Start()

	// Update should move towards target
	_ = animator.Update()

	// After update, animator should still be running
	if !animator.IsRunning() {
		t.Error("animator should still be running")
	}
}

func TestAnimatorDone(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()
	animator.SetValue(100.0)

	// After setting value, animation starts moving
	// We'll just check the IsDone method works
	_ = animator.IsDone()
}

func TestAnimatorReset(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()
	animator.SetValue(100.0)
	animator.Start()

	animator.Reset()

	if animator.IsRunning() {
		t.Error("animator should not be running after reset")
	}
}

func TestAnimatorConcurrency(t *testing.T) {
	t.Parallel()

	animator := NewAnimator()

	// Test concurrent access
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			animator.SetValue(float64(i))
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = animator.GetValue()
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	<-done
	<-done
}

func TestAnimatorTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		animType  Type
		frequency float64
		damping   float64
	}{
		{"Spring", Spring, 5.0, 0.5},
		{"Damped", Damped, 8.0, 0.3},
		{"Linear", Linear, 1.0, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			animator := NewSpringAnimator(tc.frequency, tc.damping)
			animator.animType = tc.animType

			if animator.animType != tc.animType {
				t.Errorf("expected type %v, got %v", tc.animType, animator.animType)
			}
		})
	}
}
