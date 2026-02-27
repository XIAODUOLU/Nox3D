package loader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/XIAODUOLU/Nox3D/pkg/scene"
)

// OBJLoader loads OBJ files
type OBJLoader struct{}

// NewOBJLoader creates a new OBJ loader
func NewOBJLoader() *OBJLoader {
	return &OBJLoader{}
}

// Load loads an OBJ file
func (l *OBJLoader) Load(filepath string) (*scene.Mesh, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	mesh := scene.NewMesh()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "v": // Vertex
			if len(parts) >= 4 {
				x, _ := strconv.ParseFloat(parts[1], 32)
				y, _ := strconv.ParseFloat(parts[2], 32)
				z, _ := strconv.ParseFloat(parts[3], 32)
				mesh.Vertices = append(mesh.Vertices, scene.Vec3{
					X: float32(x),
					Y: float32(y),
					Z: float32(z),
				})
			}

		case "vn": // Normal
			if len(parts) >= 4 {
				x, _ := strconv.ParseFloat(parts[1], 32)
				y, _ := strconv.ParseFloat(parts[2], 32)
				z, _ := strconv.ParseFloat(parts[3], 32)
				mesh.Normals = append(mesh.Normals, scene.Vec3{
					X: float32(x),
					Y: float32(y),
					Z: float32(z),
				})
			}

		case "vt": // Texture coordinate
			if len(parts) >= 3 {
				u, _ := strconv.ParseFloat(parts[1], 32)
				v, _ := strconv.ParseFloat(parts[2], 32)
				mesh.UVs = append(mesh.UVs, scene.Vec2{
					X: float32(u),
					Y: float32(v),
				})
			}

		case "f": // Face
			if len(parts) >= 4 {
				// Parse face indices (supports v, v/vt, v/vt/vn, v//vn formats)
				indices := make([]uint32, 0, len(parts)-1)
				for i := 1; i < len(parts); i++ {
					vertexData := strings.Split(parts[i], "/")
					if len(vertexData) > 0 {
						idx, _ := strconv.ParseInt(vertexData[0], 10, 32)
						// OBJ indices are 1-based, convert to 0-based
						if idx > 0 {
							indices = append(indices, uint32(idx-1))
						} else if idx < 0 {
							// Negative indices count from the end
							indices = append(indices, uint32(len(mesh.Vertices)+int(idx)))
						}
					}
				}

				// Triangulate if needed (for quads or n-gons)
				if len(indices) >= 3 {
					// Simple fan triangulation
					for i := 1; i < len(indices)-1; i++ {
						mesh.Indices = append(mesh.Indices, indices[0], indices[i], indices[i+1])
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Generate normals if not present
	if len(mesh.Normals) == 0 {
		mesh.Normals = generateNormals(mesh)
	}

	// Build triangle list
	mesh.BuildTriangles()

	return mesh, nil
}

// generateNormals generates flat normals for a mesh
func generateNormals(mesh *scene.Mesh) []scene.Vec3 {
	normals := make([]scene.Vec3, len(mesh.Vertices))
	counts := make([]int, len(mesh.Vertices))

	// Calculate face normals and accumulate
	for i := 0; i < len(mesh.Indices); i += 3 {
		idx0 := mesh.Indices[i]
		idx1 := mesh.Indices[i+1]
		idx2 := mesh.Indices[i+2]

		v0 := mesh.Vertices[idx0]
		v1 := mesh.Vertices[idx1]
		v2 := mesh.Vertices[idx2]

		// Calculate face normal
		edge1 := scene.Sub(v1, v0)
		edge2 := scene.Sub(v2, v0)
		normal := scene.Cross(edge1, edge2)

		// Accumulate to vertex normals
		normals[idx0] = scene.Add(normals[idx0], normal)
		normals[idx1] = scene.Add(normals[idx1], normal)
		normals[idx2] = scene.Add(normals[idx2], normal)

		counts[idx0]++
		counts[idx1]++
		counts[idx2]++
	}

	// Normalize
	for i := range normals {
		if counts[i] > 0 {
			normals[i] = scene.Normalize(normals[i])
		}
	}

	return normals
}
