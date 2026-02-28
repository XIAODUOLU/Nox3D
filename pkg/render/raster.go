package render

import (
	"math"

	"github.com/XIAODUOLU/Nox3D/pkg/scene"
	"github.com/gdamore/tcell/v2"
)

// RenderMode represents the current rendering mode
type RenderMode int

const (
	RenderModeSolid RenderMode = iota
	RenderModeWireframe
)

// Rasterizer performs software rasterization
type Rasterizer struct {
	width   int
	height  int
	zBuffer [][]float32
	Mode    RenderMode
}

// NewRasterizer creates a new rasterizer
func NewRasterizer(width, height int) *Rasterizer {
	r := &Rasterizer{
		width:   width,
		height:  height,
		zBuffer: make([][]float32, height),
		Mode:    RenderModeSolid,
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
	// Main light (key light) - from upper right front
	keyLight := scene.Normalize(scene.Vec3{X: 0.5, Y: 0.7, Z: 1.0})
	// Fill light - from left
	fillLight := scene.Normalize(scene.Vec3{X: -0.6, Y: 0.2, Z: 0.4})
	// Rim light - from behind (for edge highlighting)
	rimLight := scene.Normalize(scene.Vec3{X: 0.0, Y: 0.3, Z: -1.0})

	// Get base color from material
	baseColor := mesh.Material.BaseColor
	hasColor := mesh.Material.HasBaseColor || len(mesh.Colors) > 0

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

		// Rim light contribution (edge highlighting for 3D depth)
		rimIntensity := scene.Dot(faceNormal, rimLight)
		if rimIntensity < 0 {
			rimIntensity = 0
		}
		// Rim light is stronger on edges (inverse of view angle)
		viewDir := scene.Normalize(scene.Vec3{X: 0, Y: 0, Z: 1})
		viewDot := scene.Dot(faceNormal, viewDir)
		if viewDot < 0 {
			viewDot = 0
		}
		rimIntensity *= (1.0 - viewDot) // Stronger on edges

		// Specular highlight simulation (view-dependent)
		// Calculate reflection vector for specular
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
		specular = float32(math.Pow(float64(specular), 32)) // Higher shininess for sharper highlights

		// Ambient occlusion approximation (darker in concave areas)
		// Use depth as a proxy - deeper areas are slightly darker
		ao := float32(1.0)
		avgDepth := (v0.Z + v1.Z + v2.Z) / 3.0
		if avgDepth > 0.5 {
			ao = 0.85 // Slightly darker for distant surfaces
		}

		// Combine lighting: ambient + diffuse (key + fill) + rim + specular
		ambient := float32(0.12) * ao // Lower ambient for more contrast
		diffuse := keyIntensity*0.65 + fillIntensity*0.25
		rim := rimIntensity * 0.3 // Rim light contribution
		intensity := ambient + diffuse + rim + specular*0.5

		// Clamp intensity
		if intensity > 1.0 {
			intensity = 1.0
		}

		// Get vertex color if available (average of triangle vertices)
		triColor := baseColor
		if len(mesh.Colors) > 0 {
			// Average vertex colors
			triColor = scene.Vec3{
				X: (tri.C0.X + tri.C1.X + tri.C2.X) / 3.0,
				Y: (tri.C0.Y + tri.C1.Y + tri.C2.Y) / 3.0,
				Z: (tri.C0.Z + tri.C1.Z + tri.C2.Z) / 3.0,
			}
		}

		if r.Mode == RenderModeWireframe {
			// Wireframe rendering - use thin characters for edges
			// Use a bright color for wireframe
			var wireColor tcell.Color
			if hasColor {
				// Boost brightness for wireframe visibility
				r := int32(float32(triColor.X) * 255 * 1.5)
				g := int32(float32(triColor.Y) * 255 * 1.5)
				b := int32(float32(triColor.Z) * 255 * 1.5)
				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}
				wireColor = tcell.NewRGBColor(r, g, b)
			} else {
				// Use bright cyan for default wireframe
				wireColor = tcell.NewRGBColor(0, 255, 255)
			}

			// Use thin line character for wireframe
			wireChar := rune('-')
			r.drawLine(s0, s1, wireChar, wireColor, buffer)
			r.drawLine(s1, s2, wireChar, wireColor, buffer)
			r.drawLine(s2, s0, wireChar, wireColor, buffer)
		} else {
			// Solid rasterization
			char, color := r.intensityToChar(intensity, triColor, hasColor)
			r.rasterizeTriangle(s0, s1, s2, char, color, buffer)
		}
	}
}

// drawLine draws a line using Bresenham's algorithm with adaptive character selection
func (r *Rasterizer) drawLine(v0, v1 scene.Vec3, char rune, color tcell.Color, buffer ScreenBuffer) {
	x0 := int(v0.X)
	y0 := int(v0.Y)
	x1 := int(v1.X)
	y1 := int(v1.Y)

	dx := math.Abs(float64(x1 - x0))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -math.Abs(float64(y1 - y0))
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	// Determine line orientation for better character selection
	isHorizontal := dx > math.Abs(dy)

	for {
		// Boundary check
		if x0 >= 0 && x0 < r.width && y0 >= 0 && y0 < r.height {
			// Calculate depth for the current pixel using linear interpolation
			t := float32(0.0)
			if dx > math.Abs(dy) {
				if int(v1.X) != int(v0.X) {
					t = float32(x0-int(v0.X)) / float32(int(v1.X)-int(v0.X))
				}
			} else {
				if int(v1.Y) != int(v0.Y) {
					t = float32(y0-int(v0.Y)) / float32(int(v1.Y)-int(v0.Y))
				}
			}
			z := v0.Z + t*(v1.Z-v0.Z)

			// Depth test - wireframe should be on top
			if z-0.01 < r.zBuffer[y0][x0] {
				r.zBuffer[y0][x0] = z - 0.01 // Bias to ensure wireframe is visible

				// Use different characters based on line orientation for thinner appearance
				lineChar := char
				if isHorizontal {
					lineChar = '-'
				} else {
					lineChar = '|'
				}

				buffer.SetPixel(x0, y0, lineChar, color)
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
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
func (r *Rasterizer) intensityToChar(intensity float32, baseColor scene.Vec3, hasColor bool) (rune, tcell.Color) {
	// Enhanced ASCII art mapping with better depth perception
	// Using pure ASCII characters for classic terminal aesthetic
	chars := []rune{' ', '.', ':', '-', '=', '+', '*', '#', '%', '@'}

	// Apply non-linear intensity curve for better contrast
	// This enhances the 3D depth perception
	intensity = float32(math.Pow(float64(intensity), 0.85))

	// Calculate color based on intensity and base color
	var finalColor tcell.Color

	if hasColor {
		// Use material/vertex color with intensity modulation
		r := int32(baseColor.X * intensity * 255)
		g := int32(baseColor.Y * intensity * 255)
		b := int32(baseColor.Z * intensity * 255)

		// Clamp values
		if r > 255 {
			r = 255
		}
		if g > 255 {
			g = 255
		}
		if b > 255 {
			b = 255
		}
		if r < 0 {
			r = 0
		}
		if g < 0 {
			g = 0
		}
		if b < 0 {
			b = 0
		}

		// Add slight color boost for better visibility
		r = int32(float32(r) * 1.15)
		g = int32(float32(g) * 1.15)
		b = int32(float32(b) * 1.15)
		if r > 255 {
			r = 255
		}
		if g > 255 {
			g = 255
		}
		if b > 255 {
			b = 255
		}

		finalColor = tcell.NewRGBColor(r, g, b)
	} else {
		// Default elegant color palette with stronger contrast
		colors := []tcell.Color{
			tcell.NewRGBColor(10, 5, 15),     // Very dark (deep shadows)
			tcell.NewRGBColor(30, 15, 45),    // Dark purple
			tcell.NewRGBColor(60, 30, 90),    // Medium-dark purple
			tcell.NewRGBColor(90, 50, 135),   // Medium purple
			tcell.NewRGBColor(120, 70, 180),  // Bright purple
			tcell.NewRGBColor(150, 95, 220),  // Brighter purple
			tcell.NewRGBColor(180, 125, 240), // Very bright purple
			tcell.NewRGBColor(210, 160, 255), // Light purple
			tcell.NewRGBColor(235, 195, 255), // Very light purple
			tcell.NewRGBColor(250, 220, 255), // Brightest (highlights)
		}

		idx := int(intensity * float32(len(colors)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(colors) {
			idx = len(colors) - 1
		}
		finalColor = colors[idx]
	}

	// Select character based on intensity
	charIdx := int(intensity * float32(len(chars)-1))
	if charIdx < 0 {
		charIdx = 0
	}
	if charIdx >= len(chars) {
		charIdx = len(chars) - 1
	}

	return chars[charIdx], finalColor
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
