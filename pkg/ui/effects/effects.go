// Package effects provides visual effects for cbwsh using harmonica physics.
package effects

import (
	"math"
	"sync"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// EffectType represents the type of visual effect.
type EffectType int

const (
	// EffectWater creates a water wave effect.
	EffectWater EffectType = iota
	// EffectFire creates a fire/heat effect.
	EffectFire
	// EffectParticle creates a particle effect.
	EffectParticle
	// EffectWave creates a wave effect.
	EffectWave
	// EffectFluid creates a fluid dynamics effect.
	EffectFluid
)

// Physics constants for simulations.
const (
	// DefaultDecayFactor controls how quickly values decay in simulations.
	DefaultDecayFactor = 0.99
)

// ColorIntensity maps intensity values to colors.
type ColorIntensity struct {
	Low    lipgloss.Color
	Medium lipgloss.Color
	High   lipgloss.Color
}

// DefaultWaterColors provides water-themed colors.
var DefaultWaterColors = ColorIntensity{
	Low:    lipgloss.Color("27"),  // Dark blue
	Medium: lipgloss.Color("39"),  // Blue
	High:   lipgloss.Color("117"), // Light cyan
}

// DefaultFireColors provides fire-themed colors.
var DefaultFireColors = ColorIntensity{
	Low:    lipgloss.Color("52"),  // Dark red
	Medium: lipgloss.Color("208"), // Orange
	High:   lipgloss.Color("226"), // Yellow
}

// FluidCell represents a single cell in a fluid simulation.
type FluidCell struct {
	Density   float64
	VelocityX float64
	VelocityY float64
}

// FluidSimulation simulates fluid dynamics for visual effects.
type FluidSimulation struct {
	mu          sync.RWMutex
	width       int
	height      int
	cells       [][]FluidCell
	viscosity   float64
	diffusion   float64
	spring      harmonica.Spring
	colorScheme ColorIntensity
}

// NewFluidSimulation creates a new fluid simulation.
func NewFluidSimulation(width, height int, viscosity, diffusion float64) *FluidSimulation {
	cells := make([][]FluidCell, height)
	for i := range cells {
		cells[i] = make([]FluidCell, width)
	}

	return &FluidSimulation{
		width:       width,
		height:      height,
		cells:       cells,
		viscosity:   viscosity,
		diffusion:   diffusion,
		spring:      harmonica.NewSpring(harmonica.FPS(60), 5.0, 0.5),
		colorScheme: DefaultWaterColors,
	}
}

// SetColorScheme sets the color scheme for the simulation.
func (f *FluidSimulation) SetColorScheme(colors ColorIntensity) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.colorScheme = colors
}

// AddDensity adds density at a specific point.
func (f *FluidSimulation) AddDensity(x, y int, amount float64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if x >= 0 && x < f.width && y >= 0 && y < f.height {
		f.cells[y][x].Density += amount
	}
}

// AddVelocity adds velocity at a specific point.
func (f *FluidSimulation) AddVelocity(x, y int, vx, vy float64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if x >= 0 && x < f.width && y >= 0 && y < f.height {
		f.cells[y][x].VelocityX += vx
		f.cells[y][x].VelocityY += vy
	}
}

// Step advances the simulation by one time step.
func (f *FluidSimulation) Step() {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Diffuse density
	f.diffuseDensity()

	// Advect density
	f.advectDensity()

	// Diffuse velocity
	f.diffuseVelocity()

	// Apply decay
	f.applyDecay()
}

func (f *FluidSimulation) diffuseDensity() {
	diffRate := f.diffusion * 0.1
	for y := 1; y < f.height-1; y++ {
		for x := 1; x < f.width-1; x++ {
			avg := (f.cells[y-1][x].Density + f.cells[y+1][x].Density +
				f.cells[y][x-1].Density + f.cells[y][x+1].Density) / 4.0
			f.cells[y][x].Density += (avg - f.cells[y][x].Density) * diffRate
		}
	}
}

func (f *FluidSimulation) advectDensity() {
	newCells := make([][]FluidCell, f.height)
	for i := range newCells {
		newCells[i] = make([]FluidCell, f.width)
	}

	dt := 0.1
	for y := 1; y < f.height-1; y++ {
		for x := 1; x < f.width-1; x++ {
			// Back-trace position
			srcX := float64(x) - f.cells[y][x].VelocityX*dt
			srcY := float64(y) - f.cells[y][x].VelocityY*dt

			// Clamp to bounds
			srcX = math.Max(0.5, math.Min(float64(f.width)-1.5, srcX))
			srcY = math.Max(0.5, math.Min(float64(f.height)-1.5, srcY))

			// Bilinear interpolation
			x0 := int(srcX)
			y0 := int(srcY)
			x1 := x0 + 1
			y1 := y0 + 1

			if x1 >= f.width {
				x1 = f.width - 1
			}
			if y1 >= f.height {
				y1 = f.height - 1
			}

			sx := srcX - float64(x0)
			sy := srcY - float64(y0)

			d00 := f.cells[y0][x0].Density
			d10 := f.cells[y0][x1].Density
			d01 := f.cells[y1][x0].Density
			d11 := f.cells[y1][x1].Density

			density := (1-sx)*(1-sy)*d00 + sx*(1-sy)*d10 + (1-sx)*sy*d01 + sx*sy*d11

			newCells[y][x].Density = density
			newCells[y][x].VelocityX = f.cells[y][x].VelocityX
			newCells[y][x].VelocityY = f.cells[y][x].VelocityY
		}
	}

	f.cells = newCells
}

func (f *FluidSimulation) diffuseVelocity() {
	viscRate := f.viscosity * 0.1
	for y := 1; y < f.height-1; y++ {
		for x := 1; x < f.width-1; x++ {
			avgVX := (f.cells[y-1][x].VelocityX + f.cells[y+1][x].VelocityX +
				f.cells[y][x-1].VelocityX + f.cells[y][x+1].VelocityX) / 4.0
			avgVY := (f.cells[y-1][x].VelocityY + f.cells[y+1][x].VelocityY +
				f.cells[y][x-1].VelocityY + f.cells[y][x+1].VelocityY) / 4.0

			f.cells[y][x].VelocityX += (avgVX - f.cells[y][x].VelocityX) * viscRate
			f.cells[y][x].VelocityY += (avgVY - f.cells[y][x].VelocityY) * viscRate
		}
	}
}

func (f *FluidSimulation) applyDecay() {
	for y := range f.cells {
		for x := range f.cells[y] {
			f.cells[y][x].Density *= DefaultDecayFactor
			f.cells[y][x].VelocityX *= DefaultDecayFactor
			f.cells[y][x].VelocityY *= DefaultDecayFactor
		}
	}
}

// GetDensity returns the density at a specific point.
func (f *FluidSimulation) GetDensity(x, y int) float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if x >= 0 && x < f.width && y >= 0 && y < f.height {
		return f.cells[y][x].Density
	}
	return 0
}

// Render renders the fluid simulation as a string.
func (f *FluidSimulation) Render() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	chars := []rune{' ', '░', '▒', '▓', '█'}
	result := make([]rune, 0, (f.width+1)*f.height)

	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			density := f.cells[y][x].Density
			charIdx := int(math.Min(density*float64(len(chars)-1), float64(len(chars)-1)))
			if charIdx < 0 {
				charIdx = 0
			}
			result = append(result, chars[charIdx])
		}
		if y < f.height-1 {
			result = append(result, '\n')
		}
	}

	return string(result)
}

// RenderColored renders the fluid simulation with colors based on intensity.
func (f *FluidSimulation) RenderColored() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var result string
	chars := []rune{'░', '▒', '▓', '█'}

	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			density := f.cells[y][x].Density
			if density < 0.1 {
				result += " "
				continue
			}

			charIdx := int(math.Min(density*float64(len(chars)), float64(len(chars)-1)))
			if charIdx < 0 {
				charIdx = 0
			}

			var color lipgloss.Color
			if density < 0.33 {
				color = f.colorScheme.Low
			} else if density < 0.66 {
				color = f.colorScheme.Medium
			} else {
				color = f.colorScheme.High
			}

			style := lipgloss.NewStyle().Foreground(color)
			result += style.Render(string(chars[charIdx]))
		}
		if y < f.height-1 {
			result += "\n"
		}
	}

	return result
}

// WaterEffect creates an animated water wave effect.
type WaterEffect struct {
	mu      sync.RWMutex
	width   int
	height  int
	phase   float64
	springs []harmonica.Spring
	values  []float64
	colors  ColorIntensity
}

// NewWaterEffect creates a new water wave effect.
func NewWaterEffect(width, height int) *WaterEffect {
	springs := make([]harmonica.Spring, width)
	values := make([]float64, width)

	for i := range springs {
		springs[i] = harmonica.NewSpring(harmonica.FPS(60), 3.0, 0.3)
	}

	return &WaterEffect{
		width:   width,
		height:  height,
		springs: springs,
		values:  values,
		colors:  DefaultWaterColors,
	}
}

// SetColors sets the color scheme.
func (w *WaterEffect) SetColors(colors ColorIntensity) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.colors = colors
}

// Update advances the water effect animation.
func (w *WaterEffect) Update() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.phase += 0.1

	for i := range w.values {
		target := math.Sin(w.phase+float64(i)*0.3) * 0.5
		var velocity float64
		w.values[i], velocity = w.springs[i].Update(w.values[i], velocity, target)
		_ = velocity // velocity is not accumulated across frames in this simple model
	}
}

// Render renders the water effect.
func (w *WaterEffect) Render() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	waveChars := []rune{'~', '≈', '∿', '≋'}
	result := make([]rune, 0, (w.width+1)*w.height)

	for y := 0; y < w.height; y++ {
		normalizedY := float64(y) / float64(w.height)

		for x := 0; x < w.width; x++ {
			waveHeight := (w.values[x] + 1) / 2 // Normalize to 0-1

			if normalizedY > 1-waveHeight {
				charIdx := int(waveHeight * float64(len(waveChars)-1))
				if charIdx >= len(waveChars) {
					charIdx = len(waveChars) - 1
				}
				result = append(result, waveChars[charIdx])
			} else {
				result = append(result, ' ')
			}
		}
		if y < w.height-1 {
			result = append(result, '\n')
		}
	}

	return string(result)
}

// RenderColored renders the water effect with colors.
func (w *WaterEffect) RenderColored() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	waveChars := []rune{'~', '≈', '∿', '≋'}
	var result string

	for y := 0; y < w.height; y++ {
		normalizedY := float64(y) / float64(w.height)

		for x := 0; x < w.width; x++ {
			waveHeight := (w.values[x] + 1) / 2

			if normalizedY > 1-waveHeight {
				charIdx := int(waveHeight * float64(len(waveChars)-1))
				if charIdx >= len(waveChars) {
					charIdx = len(waveChars) - 1
				}

				var color lipgloss.Color
				if normalizedY > 0.9 {
					color = w.colors.High
				} else if normalizedY > 0.7 {
					color = w.colors.Medium
				} else {
					color = w.colors.Low
				}

				style := lipgloss.NewStyle().Foreground(color)
				result += style.Render(string(waveChars[charIdx]))
			} else {
				result += " "
			}
		}
		if y < w.height-1 {
			result += "\n"
		}
	}

	return result
}

// Disturb creates a disturbance at a specific position.
func (w *WaterEffect) Disturb(x int, intensity float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if x >= 0 && x < w.width {
		w.values[x] += intensity
	}
}

// ParticleSystem simulates particles for visual effects.
type ParticleSystem struct {
	mu        sync.RWMutex
	width     int
	height    int
	particles []Particle
	gravity   float64
	springs   []harmonica.Spring
	colors    ColorIntensity
	randState uint64 // Per-instance random state for thread safety
}

// Particle represents a single particle.
type Particle struct {
	X, Y           float64
	VX, VY         float64
	Life           float64
	Size           float64
	ColorIntensity float64
}

// NewParticleSystem creates a new particle system.
func NewParticleSystem(width, height int, maxParticles int) *ParticleSystem {
	springs := make([]harmonica.Spring, maxParticles)
	for i := range springs {
		springs[i] = harmonica.NewSpring(harmonica.FPS(60), 4.0, 0.4)
	}

	return &ParticleSystem{
		width:     width,
		height:    height,
		particles: make([]Particle, 0, maxParticles),
		gravity:   0.1,
		springs:   springs,
		colors:    DefaultFireColors,
		randState: 12345, // Initial seed for deterministic effects
	}
}

// SetColors sets the color scheme.
func (ps *ParticleSystem) SetColors(colors ColorIntensity) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.colors = colors
}

// SetGravity sets the gravity value.
func (ps *ParticleSystem) SetGravity(gravity float64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.gravity = gravity
}

// Emit emits new particles at a specific position.
func (ps *ParticleSystem) Emit(x, y float64, count int, spread, speed float64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i := 0; i < count; i++ {
		angle := (math.Pi / 2) + (math.Pi * (0.5 - ps.randFloat()) * spread)
		velocity := speed * (0.5 + ps.randFloat()*0.5)

		particle := Particle{
			X:              x,
			Y:              y,
			VX:             math.Cos(angle) * velocity,
			VY:             -math.Sin(angle) * velocity, // Negative for upward
			Life:           1.0,
			Size:           1.0,
			ColorIntensity: 1.0,
		}
		ps.particles = append(ps.particles, particle)
	}
}

// randFloat returns a pseudo-random float64 in [0, 1).
// Uses a simple LCG algorithm for deterministic, reproducible effects.
// This method is not cryptographically secure.
func (ps *ParticleSystem) randFloat() float64 {
	ps.randState = ps.randState*6364136223846793005 + 1442695040888963407
	return float64(ps.randState>>33) / float64(1<<31)
}

// Update advances the particle system.
func (ps *ParticleSystem) Update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	dt := 0.1
	aliveParticles := ps.particles[:0]

	for i := range ps.particles {
		p := &ps.particles[i]

		// Update position
		p.X += p.VX * dt
		p.Y += p.VY * dt

		// Apply gravity
		p.VY += ps.gravity

		// Decay life
		p.Life -= 0.02
		p.ColorIntensity = p.Life

		// Keep alive particles
		if p.Life > 0 && p.X >= 0 && p.X < float64(ps.width) &&
			p.Y >= 0 && p.Y < float64(ps.height) {
			aliveParticles = append(aliveParticles, *p)
		}
	}

	ps.particles = aliveParticles
}

// Render renders the particle system.
func (ps *ParticleSystem) Render() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	grid := make([][]rune, ps.height)
	for i := range grid {
		grid[i] = make([]rune, ps.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	chars := []rune{'·', '•', '●', '★'}

	for _, p := range ps.particles {
		x := int(p.X)
		y := int(p.Y)
		if x >= 0 && x < ps.width && y >= 0 && y < ps.height {
			charIdx := int(p.ColorIntensity * float64(len(chars)-1))
			if charIdx >= len(chars) {
				charIdx = len(chars) - 1
			}
			if charIdx < 0 {
				charIdx = 0
			}
			grid[y][x] = chars[charIdx]
		}
	}

	result := make([]rune, 0, (ps.width+1)*ps.height)
	for y := 0; y < ps.height; y++ {
		result = append(result, grid[y]...)
		if y < ps.height-1 {
			result = append(result, '\n')
		}
	}

	return string(result)
}

// RenderColored renders the particle system with colors.
func (ps *ParticleSystem) RenderColored() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// Build a map of particle positions and intensities
	type cell struct {
		char      rune
		intensity float64
	}
	grid := make(map[int]map[int]cell)
	for y := 0; y < ps.height; y++ {
		grid[y] = make(map[int]cell)
	}

	chars := []rune{'·', '•', '●', '★'}

	for _, p := range ps.particles {
		x := int(p.X)
		y := int(p.Y)
		if x >= 0 && x < ps.width && y >= 0 && y < ps.height {
			charIdx := int(p.ColorIntensity * float64(len(chars)-1))
			if charIdx >= len(chars) {
				charIdx = len(chars) - 1
			}
			if charIdx < 0 {
				charIdx = 0
			}
			existing, exists := grid[y][x]
			if !exists || p.ColorIntensity > existing.intensity {
				grid[y][x] = cell{char: chars[charIdx], intensity: p.ColorIntensity}
			}
		}
	}

	var result string
	for y := 0; y < ps.height; y++ {
		for x := 0; x < ps.width; x++ {
			c, exists := grid[y][x]
			if !exists {
				result += " "
				continue
			}

			var color lipgloss.Color
			if c.intensity < 0.33 {
				color = ps.colors.Low
			} else if c.intensity < 0.66 {
				color = ps.colors.Medium
			} else {
				color = ps.colors.High
			}

			style := lipgloss.NewStyle().Foreground(color)
			result += style.Render(string(c.char))
		}
		if y < ps.height-1 {
			result += "\n"
		}
	}

	return result
}

// ParticleCount returns the number of active particles.
func (ps *ParticleSystem) ParticleCount() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return len(ps.particles)
}
