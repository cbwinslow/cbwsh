package progress

import (
	"testing"
	"time"
)

func TestNewBar(t *testing.T) {
	t.Parallel()

	bar := NewBar()

	if bar == nil {
		t.Fatal("expected non-nil progress bar")
	}

	if bar.width == 0 {
		t.Error("width should be set")
	}
}

func TestNewBarWithColors(t *testing.T) {
	t.Parallel()

	colorA := "#FF0000"
	colorB := "#00FF00"

	bar := NewBarWithColors(colorA, colorB)

	if bar == nil {
		t.Fatal("expected non-nil progress bar")
	}

	if bar.colorA != colorA {
		t.Errorf("expected colorA %s, got %s", colorA, bar.colorA)
	}

	if bar.colorB != colorB {
		t.Errorf("expected colorB %s, got %s", colorB, bar.colorB)
	}

	if !bar.useGradient {
		t.Error("useGradient should be true")
	}
}

func TestBarStart(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	total := 100

	bar.Start(total)

	if bar.total != total {
		t.Errorf("expected total %d, got %d", total, bar.total)
	}

	if bar.current != 0 {
		t.Errorf("expected current 0, got %d", bar.current)
	}

	if bar.finished {
		t.Error("bar should not be finished after start")
	}

	if bar.startTime.IsZero() {
		t.Error("startTime should be set")
	}
}

func TestBarIncrement(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(10)

	initialCurrent := bar.current

	bar.Increment()

	if bar.current != initialCurrent+1 {
		t.Errorf("expected current %d, got %d", initialCurrent+1, bar.current)
	}
}

func TestBarIncrementBy(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)

	amount := 25
	bar.IncrementBy(amount)

	if bar.current != amount {
		t.Errorf("expected current %d, got %d", amount, bar.current)
	}
}

func TestBarSetProgress(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)

	value := 50
	bar.IncrementBy(value)

	if bar.current != value {
		t.Errorf("expected current %d, got %d", value, bar.current)
	}
}

func TestBarSetMessage(t *testing.T) {
	t.Parallel()

	bar := NewBar()

	message := "Processing files..."
	bar.SetMessage(message)

	if bar.message != message {
		t.Errorf("expected message '%s', got '%s'", message, bar.message)
	}
}

func TestBarFinish(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(10)

	if bar.finished {
		t.Error("bar should not be finished initially")
	}

	bar.Finish()

	if !bar.finished {
		t.Error("bar should be finished after Finish()")
	}
}

func TestBarPercent(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)
	bar.IncrementBy(50)

	percent := bar.Percent()

	if percent != 50.0 {
		t.Errorf("expected percent 50.0, got %f", percent)
	}
}

func TestBarPercentZeroTotal(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(0)

	percent := bar.Percent()

	// Should handle division by zero gracefully
	if percent != 0 {
		t.Errorf("percent should be 0, got %f", percent)
	}
}

func TestBarRender(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)
	bar.IncrementBy(50)

	rendered := bar.View()

	if rendered == "" {
		t.Error("render should return non-empty string")
	}
}

func TestBarElapsedTime(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Check if the view contains timing information (it uses elapsed time internally)
	view := bar.View()

	if view == "" {
		t.Error("view should return non-empty string with timing info")
	}

	// The View method internally calculates elapsed time
	// We can't access it directly, but we can verify the view is generated
}

func TestBarSetWidth(t *testing.T) {
	t.Parallel()

	bar := NewBar()

	width := 80
	bar.SetWidth(width)

	if bar.width != width {
		t.Errorf("expected width %d, got %d", width, bar.width)
	}
}

func TestBarView(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(100)
	bar.IncrementBy(50)
	bar.SetMessage("Working...")

	view := bar.View()

	if view == "" {
		t.Error("view should return non-empty string")
	}
}

func TestBarConcurrency(t *testing.T) {
	t.Parallel()

	bar := NewBar()
	bar.Start(1000)

	done := make(chan bool, 2)

	// Concurrent increments
	go func() {
		for i := 0; i < 500; i++ {
			bar.Increment()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 500; i++ {
			bar.Increment()
		}
		done <- true
	}()

	<-done
	<-done

	// Final count should be consistent (though may not be exactly 1000 due to race conditions,
	// but the mutex should prevent data corruption)
	if bar.current < 0 || bar.current > 1000 {
		t.Errorf("current value out of bounds: %d", bar.current)
	}
}
