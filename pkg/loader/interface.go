package loader

import "github.com/XIAODUOLU/Nox3D/pkg/scene"

// Loader is the interface for loading 3D models
type Loader interface {
	Load(filepath string) (*scene.Mesh, error)
}
