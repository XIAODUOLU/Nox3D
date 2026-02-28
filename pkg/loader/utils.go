package loader

import "github.com/XIAODUOLU/Nox3D/pkg/scene"

// generateNormals generates vertex normals from face normals
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
