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
}

// Mesh represents a 3D mesh
type Mesh struct {
	Vertices  []Vec3
	Normals   []Vec3
	UVs       []Vec2
	Indices   []uint32
	Triangles []Triangle
}

// NewMesh creates a new mesh
func NewMesh() *Mesh {
	return &Mesh{
		Vertices:  make([]Vec3, 0),
		Normals:   make([]Vec3, 0),
		UVs:       make([]Vec2, 0),
		Indices:   make([]uint32, 0),
		Triangles: make([]Triangle, 0),
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
		if len(m.Normals) > 0 {
			tri.N0 = m.Normals[idx0]
			tri.N1 = m.Normals[idx1]
			tri.N2 = m.Normals[idx2]
		}

		// Add UVs if available
		if len(m.UVs) > 0 {
			tri.UV0 = m.UVs[idx0]
			tri.UV1 = m.UVs[idx1]
			tri.UV2 = m.UVs[idx2]
		}

		m.Triangles = append(m.Triangles, tri)
	}
}
