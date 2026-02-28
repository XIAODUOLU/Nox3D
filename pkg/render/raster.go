package render

import (
	"math"

	"github.com/XIAODUOLU/Nox3D/pkg/scene"
	"github.com/gdamore/tcell/v2"
)

// Rasterizer performs software rasterization
type Rasterizer struct {
	width   int
	height  int
	zBuffer [][]float32
}

// NewRasterizer creates a new rasterizer
func NewRasterizer(width, height int) *Rasterizer {
	r := &Rasterizer{
		width:   width,
		height:  height,
		zBuffer: make([][]float32, height),
	}

	for i := range r.zBuffer {
		r.zBuffer[i] = make([]float32, width)
	}

	return r
}

// ClearDepth clears the depth buffer
func (r *Rasterizer) ClearDepth() {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.zBuffer[y][x] = math.MaxFloat32
		}
	}
}

// ScreenBuffer represents the rendering output
type ScreenBuffer interface {
	SetPixel(x, y int, char rune, color tcell.Color)
}

// Render renders a mesh to the screen buffer
func (r *Rasterizer) Render(mesh *scene.Mesh, mvp scene.Mat4, buffer ScreenBuffer) {
	r.ClearDepth()

	// Enhanced lighting setup with multiple light sources
	// Main light (key light) - from upper right
	keyLight := scene.Normalize(scene.Vec3{X: 0.6, Y: 0.8, Z: 1.0})
	// Fill light - from left
	fillLight := scene.Normalize(scene.Vec3{X: -0.5, Y: 0.3, Z: 0.5})

	for _, tri := range mesh.Triangles {
		// Transform vertices
		v0 := mvp.TransformVec3(tri.V0)
		v1 := mvp.TransformVec3(tri.V1)
		v2 := mvp.TransformVec3(tri.V2)

		// Convert to screen space (accounting for character aspect ratio)
		s0 := r.toScreen(v0)
		s1 := r.toScreen(v1)
		s2 := r.toScreen(v2)

		// Backface culling
		edge1 := scene.Vec3{X: s1.X - s0.X, Y: s1.Y - s0.Y, Z: 0}
		edge2 := scene.Vec3{X: s2.X - s0.X, Y: s2.Y - s0.Y, Z: 0}
		cross := edge1.X*edge2.Y - edge1.Y*edge2.X
		if cross <= 0 {
			continue
		}

		// Calculate face normal for lighting
		worldEdge1 := scene.Sub(tri.V1, tri.V0)
		worldEdge2 := scene.Sub(tri.V2, tri.V0)
		faceNormal := scene.Normalize(scene.Cross(worldEdge1, worldEdge2))

		// Enhanced lighting calculation with multiple lights
		// Key light contribution (main directional light)
		keyIntensity := scene.Dot(faceNormal, keyLight)
		if keyIntensity < 0 {
			keyIntensity = 0
		}

		// Fill light contribution (softer secondary light)
		fillIntensity := scene.Dot(faceNormal, fillLight)
		if fillIntensity < 0 {
			fillIntensity = 0
		}

		// Specular highlight simulation (view-dependent)
		// Calculate reflection vector for specular
		viewDir := scene.Normalize(scene.Vec3{X: 0, Y: 0, Z: 1})
		reflectDir := scene.Sub(
			scene.Vec3{
				X: 2 * faceNormal.X * keyIntensity,
				Y: 2 * faceNormal.Y * keyIntensity,
				Z: 2 * faceNormal.Z * keyIntensity,
			},
			keyLight,
		)
		specular := scene.Dot(reflectDir, viewDir)
		if specular < 0 {
			specular = 0
		}
		specular = float32(math.Pow(float64(specular), 16)) // Shininess factor

		// Combine lighting: ambient + diffuse (key + fill) + specular
		ambient := float32(0.15)
		diffuse := keyIntensity*0.7 + fillIntensity*0.3
		intensity := ambient + diffuse + specular*0.4

		// Clamp intensity
		if intensity > 1.0 {
			intensity = 1.0
		}

		// Get color and character based on intensity
		char, color := r.intensityToChar(intensity)

		// Rasterize triangle
		r.rasterizeTriangle(s0, s1, s2, char, color, buffer)
	}
}

// toScreen converts NDC coordinates to screen space
func (r *Rasterizer) toScreen(v scene.Vec3) scene.Vec3 {
	// NDC is in [-1, 1], convert to screen coordinates
	// Note: Y coordinate needs adjustment for character aspect ratio
	// Terminal characters are ~2x taller than wide, so we compress Y by 0.5
	x := (v.X + 1.0) * float32(r.width) * 0.5
	y := (1.0 - v.Y) * float32(r.height) * 0.5
	return scene.Vec3{X: x, Y: y, Z: v.Z}
}

// intensityToChar maps light intensity to character and color with enhanced 3D effect
func (r *Rasterizer) intensityToChar(intensity float32) (rune, tcell.Color) {
	// Simple ASCII art mapping with good contrast
	chars := []rune{' ', '.', ':', '-', '=', '+', '*', '#', '@'}

	// Elegant purple color palette (from dark to bright, no white)
	colors := []tcell.Color{
		tcell.NewRGBColor(20, 10, 30),    // Very dark purple (shadows)
		tcell.NewRGBColor(50, 25, 75),    // Dark purple
		tcell.NewRGBColor(80, 40, 120),   // Medium-dark purple
		tcell.NewRGBColor(110, 60, 160),  // Medium purple
		tcell.NewRGBColor(140, 80, 200),  // Bright purple
		tcell.NewRGBColor(170, 110, 230), // Brighter purple
		tcell.NewRGBColor(190, 140, 240), // Very bright purple
		tcell.NewRGBColor(210, 170, 250), // Light purple
		tcell.NewRGBColor(230, 200, 255), // Very light purple (max)
	}

	idx := int(intensity * float32(len(chars)-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(chars) {
		idx = len(chars) - 1
	}

	return chars[idx], colors[idx]
}

// rasterizeTriangle rasterizes a triangle using scanline algorithm
func (r *Rasterizer) rasterizeTriangle(v0, v1, v2 scene.Vec3, char rune, color tcell.Color, buffer ScreenBuffer) {
	// Get bounding box
	minX := int(min(v0.X, min(v1.X, v2.X)))
	maxX := int(max(v0.X, max(v1.X, v2.X)))
	minY := int(min(v0.Y, min(v1.Y, v2.Y)))
	maxY := int(max(v0.Y, max(v1.Y, v2.Y)))

	// Clip to screen bounds
	if minX < 0 {
		minX = 0
	}
	if maxX >= r.width {
		maxX = r.width - 1
	}
	if minY < 0 {
		minY = 0
	}
	if maxY >= r.height {
		maxY = r.height - 1
	}

	// Rasterize
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := scene.Vec3{X: float32(x) + 0.5, Y: float32(y) + 0.5, Z: 0}

			// Calculate barycentric coordinates
			bary := barycentric(v0, v1, v2, p)

			// Check if point is inside triangle
			if bary.X >= 0 && bary.Y >= 0 && bary.Z >= 0 {
				// Interpolate depth
				z := v0.Z*bary.X + v1.Z*bary.Y + v2.Z*bary.Z

				// Depth test
				if z < r.zBuffer[y][x] {
					r.zBuffer[y][x] = z
					buffer.SetPixel(x, y, char, color)
				}
			}
		}
	}
}

// barycentric calculates barycentric coordinates
func barycentric(v0, v1, v2, p scene.Vec3) scene.Vec3 {
	v0v1 := scene.Vec3{X: v1.X - v0.X, Y: v1.Y - v0.Y, Z: 0}
	v0v2 := scene.Vec3{X: v2.X - v0.X, Y: v2.Y - v0.Y, Z: 0}
	v0p := scene.Vec3{X: p.X - v0.X, Y: p.Y - v0.Y, Z: 0}

	d00 := v0v1.X*v0v1.X + v0v1.Y*v0v1.Y
	d01 := v0v1.X*v0v2.X + v0v1.Y*v0v2.Y
	d11 := v0v2.X*v0v2.X + v0v2.Y*v0v2.Y
	d20 := v0p.X*v0v1.X + v0p.Y*v0v1.Y
	d21 := v0p.X*v0v2.X + v0p.Y*v0v2.Y

	denom := d00*d11 - d01*d01
	if denom == 0 {
		return scene.Vec3{X: -1, Y: -1, Z: -1}
	}

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1.0 - v - w

	return scene.Vec3{X: u, Y: v, Z: w}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
