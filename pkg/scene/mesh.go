package scene

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}

// Vec2 represents a 2D vector (for texture coordinates)
type Vec2 struct {
	X, Y float32
}

// Triangle represents a single triangle face
type Triangle struct {
	V0, V1, V2    Vec3 // Vertex positions
	N0, N1, N2    Vec3 // Vertex normals
	UV0, UV1, UV2 Vec2 // Texture coordinates (optional)
	C0, C1, C2    Vec3 // Vertex colors (RGB, optional)
}

// Material represents material properties
type Material struct {
	BaseColor    Vec3    // Base color (RGB)
	Metallic     float32 // Metallic factor
	Roughness    float32 // Roughness factor
	HasBaseColor bool    // Whether base color is set
}

// Mesh represents a 3D mesh
type Mesh struct {
	Vertices  []Vec3
	Normals   []Vec3
	UVs       []Vec2
	Colors    []Vec3 // Vertex colors (optional)
	Indices   []uint32
	Triangles []Triangle
	Material  Material // Material properties
}

// NewMesh creates a new mesh
func NewMesh() *Mesh {
	return &Mesh{
		Vertices:  make([]Vec3, 0),
		Normals:   make([]Vec3, 0),
		UVs:       make([]Vec2, 0),
		Colors:    make([]Vec3, 0),
		Indices:   make([]uint32, 0),
		Triangles: make([]Triangle, 0),
		Material: Material{
			BaseColor:    Vec3{X: 0.8, Y: 0.8, Z: 0.8}, // Default gray
			Metallic:     0.0,
			Roughness:    0.5,
			HasBaseColor: false,
		},
	}
}

// BuildTriangles constructs triangle list from indexed data
func (m *Mesh) BuildTriangles() {
	m.Triangles = make([]Triangle, 0, len(m.Indices)/3)

	for i := 0; i < len(m.Indices); i += 3 {
		idx0 := m.Indices[i]
		idx1 := m.Indices[i+1]
		idx2 := m.Indices[i+2]

		tri := Triangle{
			V0: m.Vertices[idx0],
			V1: m.Vertices[idx1],
			V2: m.Vertices[idx2],
		}

		// Add normals if available
		if len(m.Normals) > 0 && int(idx0) < len(m.Normals) && int(idx1) < len(m.Normals) && int(idx2) < len(m.Normals) {
			tri.N0 = m.Normals[idx0]
			tri.N1 = m.Normals[idx1]
			tri.N2 = m.Normals[idx2]
		}

		// Add UVs if available
		if len(m.UVs) > 0 && int(idx0) < len(m.UVs) && int(idx1) < len(m.UVs) && int(idx2) < len(m.UVs) {
			tri.UV0 = m.UVs[idx0]
			tri.UV1 = m.UVs[idx1]
			tri.UV2 = m.UVs[idx2]
		}

		// Add vertex colors if available
		if len(m.Colors) > 0 && int(idx0) < len(m.Colors) && int(idx1) < len(m.Colors) && int(idx2) < len(m.Colors) {
			tri.C0 = m.Colors[idx0]
			tri.C1 = m.Colors[idx1]
			tri.C2 = m.Colors[idx2]
		}

		m.Triangles = append(m.Triangles, tri)
	}
}

// GetBounds calculates the bounding box of the mesh
func (m *Mesh) GetBounds() (min, max Vec3) {
	if len(m.Vertices) == 0 {
		return Vec3{}, Vec3{}
	}

	min = m.Vertices[0]
	max = m.Vertices[0]

	for _, v := range m.Vertices {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.Z < min.Z {
			min.Z = v.Z
		}
		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
		if v.Z > max.Z {
			max.Z = v.Z
		}
	}

	return min, max
}

// GetCenter calculates the center of the mesh
func (m *Mesh) GetCenter() Vec3 {
	min, max := m.GetBounds()
	return Vec3{
		X: (min.X + max.X) * 0.5,
		Y: (min.Y + max.Y) * 0.5,
		Z: (min.Z + max.Z) * 0.5,
	}
}

// GetSize calculates the size of the mesh
func (m *Mesh) GetSize() Vec3 {
	min, max := m.GetBounds()
	return Vec3{
		X: max.X - min.X,
		Y: max.Y - min.Y,
		Z: max.Z - min.Z,
	}
}

// GetMaxDimension returns the largest dimension of the mesh
func (m *Mesh) GetMaxDimension() float32 {
	size := m.GetSize()
	maxDim := size.X
	if size.Y > maxDim {
		maxDim = size.Y
	}
	if size.Z > maxDim {
		maxDim = size.Z
	}
	return maxDim
}
