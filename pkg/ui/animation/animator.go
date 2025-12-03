// Package animation provides animation utilities using harmonica for cbwsh.
package animation

import (
	"math"
	"sync"
	"time"

	"github.com/charmbracelet/harmonica"
)

// Type represents the type of animation.
type Type int

const (
	// Spring uses spring physics for animation.
	Spring Type = iota
	// Damped uses damped oscillation.
	Damped
	// Linear uses linear interpolation.
	Linear
)

// Animator provides smooth animations using harmonica.
type Animator struct {
	mu        sync.RWMutex
	spring    harmonica.Spring
	animType  Type
	current   float64
	target    float64
	velocity  float64
	running   bool
	done      bool
	lastTick  time.Time
	frequency float64
	damping   float64
}

// NewSpringAnimator creates a new spring-based animator.
func NewSpringAnimator(frequency, damping float64) *Animator {
	return &Animator{
		spring:    harmonica.NewSpring(harmonica.FPS(60), frequency, damping),
		animType:  Spring,
		frequency: frequency,
		damping:   damping,
		lastTick:  time.Now(),
	}
}

// NewAnimator creates a new animator with default spring settings.
func NewAnimator() *Animator {
	return NewSpringAnimator(5.0, 0.5)
}

// Start starts the animation.
func (a *Animator) Start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.running = true
	a.done = false
	a.lastTick = time.Now()
}

// Stop stops the animation.
func (a *Animator) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.running = false
}

// Update updates the animation state and returns the current value as a string.
func (a *Animator) Update() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return ""
	}

	now := time.Now()
	dt := now.Sub(a.lastTick).Seconds()
	a.lastTick = now

	switch a.animType {
	case Spring:
		a.current, a.velocity = a.spring.Update(a.current, a.velocity, a.target)
	case Damped:
		a.current, a.velocity = a.spring.Update(a.current, a.velocity, a.target)
	case Linear:
		diff := a.target - a.current
		step := diff * dt * 5
		if abs(diff) < 0.001 {
			a.current = a.target
		} else {
			a.current += step
		}
	}

	// Check if animation is done
	if abs(a.current-a.target) < 0.001 && abs(a.velocity) < 0.001 {
		a.done = true
		a.current = a.target
	}

	return ""
}

// SetValue sets the target value for the animation.
func (a *Animator) SetValue(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.target = value
	a.done = false
}

// GetValue returns the current animated value.
func (a *Animator) GetValue() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.current
}

// GetTarget returns the target value.
func (a *Animator) GetTarget() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.target
}

// IsDone returns whether the animation has completed.
func (a *Animator) IsDone() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.done
}

// IsRunning returns whether the animation is running.
func (a *Animator) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// Reset resets the animator to initial state.
func (a *Animator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.current = 0
	a.target = 0
	a.velocity = 0
	a.done = false
	a.running = false
}

// SetSpringParams updates the spring parameters.
func (a *Animator) SetSpringParams(frequency, damping float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.frequency = frequency
	a.damping = damping
	a.spring = harmonica.NewSpring(harmonica.FPS(60), frequency, damping)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Group manages multiple animators.
type Group struct {
	mu        sync.RWMutex
	animators map[string]*Animator
}

// NewGroup creates a new animation group.
func NewGroup() *Group {
	return &Group{
		animators: make(map[string]*Animator),
	}
}

// Add adds an animator to the group.
func (g *Group) Add(name string, animator *Animator) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.animators[name] = animator
}

// Get returns an animator by name.
func (g *Group) Get(name string) (*Animator, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	animator, exists := g.animators[name]
	return animator, exists
}

// Remove removes an animator from the group.
func (g *Group) Remove(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.animators, name)
}

// StartAll starts all animators.
func (g *Group) StartAll() {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, animator := range g.animators {
		animator.Start()
	}
}

// StopAll stops all animators.
func (g *Group) StopAll() {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, animator := range g.animators {
		animator.Stop()
	}
}

// UpdateAll updates all animators.
func (g *Group) UpdateAll() {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, animator := range g.animators {
		animator.Update()
	}
}

// AllDone returns whether all animations are done.
func (g *Group) AllDone() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, animator := range g.animators {
		if !animator.IsDone() {
			return false
		}
	}
	return true
}

// Easing provides easing functions.
type Easing struct{}

// EaseInOut applies ease-in-out to a value (0-1).
func (Easing) EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - (-2*t+2)*(-2*t+2)/2
}

// EaseIn applies ease-in to a value (0-1).
func (Easing) EaseIn(t float64) float64 {
	return t * t
}

// EaseOut applies ease-out to a value (0-1).
func (Easing) EaseOut(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// Bounce applies a bounce effect to a value (0-1).
func (Easing) Bounce(t float64) float64 {
	n1 := 7.5625
	d1 := 2.75

	if t < 1/d1 {
		return n1 * t * t
	} else if t < 2/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	}
	t -= 2.625 / d1
	return n1*t*t + 0.984375
}

// Elastic applies an elastic effect to a value (0-1).
func (Easing) Elastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	return -math.Pow(2, 10*t-10) * math.Sin((t*10-10.75)*(2*math.Pi)/3)
}

// Transition represents an animated transition between states.
type Transition struct {
	mu        sync.RWMutex
	animator  *Animator
	fromValue float64
	toValue   float64
	duration  time.Duration
	startTime time.Time
	easing    func(float64) float64
}

// NewTransition creates a new transition.
func NewTransition(from, to float64, duration time.Duration) *Transition {
	animator := NewAnimator()
	animator.current = from
	animator.target = to
	return &Transition{
		animator:  animator,
		fromValue: from,
		toValue:   to,
		duration:  duration,
		easing:    Easing{}.EaseInOut,
	}
}

// Start starts the transition.
func (t *Transition) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.startTime = time.Now()
	t.animator.Start()
}

// Update updates the transition.
func (t *Transition) Update() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	elapsed := time.Since(t.startTime)
	progress := float64(elapsed) / float64(t.duration)
	if progress > 1 {
		progress = 1
	}

	easedProgress := t.easing(progress)
	return t.fromValue + (t.toValue-t.fromValue)*easedProgress
}

// IsDone returns whether the transition is complete.
func (t *Transition) IsDone() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return time.Since(t.startTime) >= t.duration
}

// SetEasing sets the easing function.
func (t *Transition) SetEasing(easing func(float64) float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.easing = easing
}
