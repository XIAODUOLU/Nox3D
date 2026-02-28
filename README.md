# Nox3D

> Render 3D models in the dark. A terminal-based 3D viewer with ASCII/ANSI art.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Nox3D brings 3D graphics to your terminal using pure Go. View and interact with 3D models using ASCII art and cyberpunk-style ANSI colors—no GUI required.

## Features

- 🎨 **ASCII Art Rendering** - Software rasterization with depth buffering
- 🖱️ **Interactive Controls** - Mouse and keyboard navigation
- 🚀 **Pure Go** - No CGO dependencies, static binary compilation
- 🎯 **Shell Completion** - Tab completion for Bash, Zsh, Fish, PowerShell
- 🌈 **Cyberpunk Aesthetic** - Neon colors on black background
- 📦 **Cross-platform** - Linux, macOS, Windows

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/XIAODUOLU/Nox3D.git
cd Nox3D

# Build
make build

# Or build static binary
make build-static
```

### Usage

```bash
# View a 3D model
./bin/nox test assets/obj/nox.obj

# With custom FPS
./bin/nox test model.obj --fps 60

# Show help
./bin/nox --help
```

### Controls

| Action | Control |
|--------|---------|
| Rotate camera | Mouse drag |
| Zoom | Mouse wheel |
| Move target | W/A/S/D |
| Rotate model | E |
| Toggle auto-rotate | R or Space |
| Switch to solid mode | 1 |
| Switch to wireframe mode | 2 |
| Exit | ESC or Ctrl+C or Q |

## Shell Completion

Enable tab completion for your shell:

```bash
# Bash
./install-completion.sh bash

# Zsh
./install-completion.sh zsh

# Fish
./install-completion.sh fish
```

## Supported Formats

- ✅ OBJ (.obj)
- ✅ GLTF/GLB (.glb)
- 🚧 FBX (.fbx) - Coming soon

## Requirements

- Go 1.24+
- Terminal with 24-bit true color support
- Mouse-enabled terminal (most modern terminals)

### Recommended Terminal Settings

For the best viewing experience:

- **Use a monospaced (fixed-width) font** - This ensures proper character alignment and aspect ratio
- **Reduce font size** - Smaller fonts provide higher "resolution" and more detail in the rendered models
- **Enable true color support** - Most modern terminals support 24-bit color by default
- **Recommended fonts**:
  - JetBrains Mono
  - Fira Code
  - Consolas
  - Monaco
  - Courier New

**Tip**: Try reducing your terminal font size to 8-10pt for the best visual quality. The smaller the characters, the more detailed the 3D rendering will appear.

## Building

```bash
# Development build
make build

# Static binary (no dependencies)
make build-static

# Cross-compile for all platforms
make build-all

# Run tests
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Built with:
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [cobra](https://github.com/spf13/cobra) - CLI framework

---

**Note:** Nox3D is designed for headless servers and SSH environments where traditional GUI tools aren't available. Perfect for DevOps workflows and terminal enthusiasts.
