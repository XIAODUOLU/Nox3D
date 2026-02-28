package loader

import (
	"fmt"

	"github.com/XIAODUOLU/Nox3D/pkg/scene"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// GLBLoader loads GLB/GLTF files
type GLBLoader struct{}

// NewGLBLoader creates a new GLB loader
func NewGLBLoader() *GLBLoader {
	return &GLBLoader{}
}

// Load loads a GLB/GLTF file
func (l *GLBLoader) Load(filepath string) (*scene.Mesh, error) {
	// Load GLTF document
	doc, err := gltf.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GLB file: %w", err)
	}

	mesh := scene.NewMesh()

	// Process all meshes in the document
	for _, gltfMesh := range doc.Meshes {
		for _, primitive := range gltfMesh.Primitives {
			// Extract material information if available
			if primitive.Material != nil {
				mat := doc.Materials[*primitive.Material]
				if mat.PBRMetallicRoughness != nil {
					pbr := mat.PBRMetallicRoughness
					if pbr.BaseColorFactor != nil {
						mesh.Material.BaseColor = scene.Vec3{
							X: float32(pbr.BaseColorFactor[0]),
							Y: float32(pbr.BaseColorFactor[1]),
							Z: float32(pbr.BaseColorFactor[2]),
						}
						mesh.Material.HasBaseColor = true
					}
					if pbr.MetallicFactor != nil {
						mesh.Material.Metallic = float32(*pbr.MetallicFactor)
					}
					if pbr.RoughnessFactor != nil {
						mesh.Material.Roughness = float32(*pbr.RoughnessFactor)
					}
				}
			}

			// Get position accessor
			posAccessor, ok := primitive.Attributes[gltf.POSITION]
			if !ok {
				continue
			}

			// Read positions
			// GLTF uses right-handed Y-up coordinate system
			// We need to convert to match OBJ behavior (which seems to work correctly)
			positions, err := modeler.ReadPosition(doc, doc.Accessors[posAccessor], nil)
			if err != nil {
				return nil, fmt.Errorf("failed to read positions: %w", err)
			}
			for _, pos := range positions {
				// GLTF: X right, Y up, Z out (towards viewer)
				// Convert Y-up to our coordinate system
				// Swap Y and Z, negate Y to flip vertical direction
				mesh.Vertices = append(mesh.Vertices, scene.Vec3{
					X: pos[0],
					Y: -pos[2], // Z becomes -Y (flip vertical)
					Z: pos[1],  // Y becomes Z (depth)
				})
			}

			// Read normals if available
			if normalAccessor, ok := primitive.Attributes[gltf.NORMAL]; ok {
				normals, err := modeler.ReadNormal(doc, doc.Accessors[normalAccessor], nil)
				if err == nil {
					for _, norm := range normals {
						// Apply same transformation to normals
						mesh.Normals = append(mesh.Normals, scene.Vec3{
							X: norm[0],
							Y: -norm[2],
							Z: norm[1],
						})
					}
				}
			}

			// Read texture coordinates if available
			if texCoordAccessor, ok := primitive.Attributes[gltf.TEXCOORD_0]; ok {
				texCoords, err := modeler.ReadTextureCoord(doc, doc.Accessors[texCoordAccessor], nil)
				if err == nil {
					for _, tc := range texCoords {
						mesh.UVs = append(mesh.UVs, scene.Vec2{
							X: tc[0],
							Y: tc[1],
						})
					}
				}
			}

			// Read vertex colors if available (COLOR_0 attribute)
			if colorAccessor, ok := primitive.Attributes[gltf.COLOR_0]; ok {
				accessor := doc.Accessors[colorAccessor]
				colors, err := modeler.ReadColor(doc, accessor, nil)
				if err == nil {
					for _, col := range colors {
						// Colors can be RGB or RGBA, we only use RGB
						// Normalize from 0-255 to 0-1 range
						mesh.Colors = append(mesh.Colors, scene.Vec3{
							X: float32(col[0]) / 255.0,
							Y: float32(col[1]) / 255.0,
							Z: float32(col[2]) / 255.0,
						})
					}
				}
			}

			// Read indices
			if primitive.Indices != nil {
				indices, err := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], nil)
				if err != nil {
					return nil, fmt.Errorf("failed to read indices: %w", err)
				}
				for _, idx := range indices {
					mesh.Indices = append(mesh.Indices, idx)
				}
			} else {
				// If no indices, create them sequentially
				for i := uint32(0); i < uint32(len(positions)); i++ {
					mesh.Indices = append(mesh.Indices, i)
				}
			}
		}
	}

	// Generate normals if not present
	if len(mesh.Normals) == 0 {
		mesh.Normals = generateNormals(mesh)
	}

	// Build triangles
	mesh.BuildTriangles()

	return mesh, nil
}
