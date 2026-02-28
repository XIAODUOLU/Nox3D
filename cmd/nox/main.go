package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/XIAODUOLU/Nox3D/pkg/loader"
	"github.com/XIAODUOLU/Nox3D/pkg/render"
	"github.com/XIAODUOLU/Nox3D/pkg/scene"
	"github.com/XIAODUOLU/Nox3D/pkg/terminal"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

var (
	// Flags
	fps        int
	autoRotate bool
	noRotate   bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nox",
		Short: "Nox3D - Terminal 3D Model Viewer",
		Long: `Nox3D is a terminal-based 3D model viewer written in pure Go.
It renders 3D models directly in your terminal using ASCII art and ANSI colors.`,
	}

	// Test command
	testCmd := &cobra.Command{
		Use:   "test [model-file]",
		Short: "Load and display a 3D model",
		Long: `Load a 3D model file and display it in the terminal.
Supports OBJ, GLB, and FBX formats.`,
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// File completion for model files
			return []string{"obj", "glb", "fbx"}, cobra.ShellCompDirectiveFilterFileExt
		},
		Run: func(cmd *cobra.Command, args []string) {
			runTest(args[0])
		},
	}
	testCmd.Flags().IntVar(&fps, "fps", 30, "Target frames per second")
	testCmd.Flags().BoolVar(&autoRotate, "auto-rotate", false, "Enable auto-rotation on startup")
	testCmd.Flags().BoolVar(&noRotate, "no-rotate", false, "Disable auto-rotation on startup (default)")

	// Completion command
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for nox.

To load completions:

Bash:
  $ source <(nox completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ nox completion bash > /etc/bash_completion.d/nox
  # macOS:
  $ nox completion bash > /usr/local/etc/bash_completion.d/nox

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ nox completion zsh > "${fpath[1]}/_nox"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ nox completion fish | source
  # To load completions for each session, execute once:
  $ nox completion fish > ~/.config/fish/completions/nox.fish

PowerShell:
  PS> nox completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> nox completion powershell > nox.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(completionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runTest(modelPath string) {
	// Initialize terminal screen
	screen, err := terminal.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Close()

	// Load model
	mesh, err := loadModel(modelPath)
	if err != nil {
		screen.Close()
		fmt.Fprintf(os.Stderr, "Failed to load model: %v\n", err)
		os.Exit(1)
	}

	// Calculate model bounds for auto-scaling
	center := mesh.GetCenter()
	maxDim := mesh.GetMaxDimension()

	// Calculate optimal camera distance to fit model in view
	width, height := screen.Size()
	camera := scene.NewCamera()
	camera.Aspect = float32(width) / float32(height)

	// Set camera target to model center
	camera.Target = center

	// Calculate distance based on model size, FOV, and aspect ratio
	// Account for terminal character aspect ratio (chars are ~2x taller than wide)

	// Calculate distance needed to fit model horizontally and vertically
	// Use the smaller FOV dimension to ensure model fits in both directions
	effectiveAspect := camera.Aspect * 0.5 // Account for character aspect ratio
	var effectiveFOV float32
	if effectiveAspect < 1.0 {
		// Portrait mode: vertical FOV is limiting
		effectiveFOV = camera.FOV
	} else {
		// Landscape mode: horizontal FOV is limiting
		effectiveFOV = 2.0 * float32(math.Atan(math.Tan(float64(camera.FOV*0.5))/float64(effectiveAspect)))
	}

	// Calculate optimal distance with very tight framing to fill screen
	// Use smaller multiplier to make model appear much larger
	optimalDistance := (maxDim * 1.1) / (2.0 * float32(math.Tan(float64(effectiveFOV*0.5))))
	camera.Distance = optimalDistance * 0.5 // Reduce distance by 50% to make model 2x larger

	// Position camera at a nice angle with lower elevation
	camera.Yaw = float32(math.Pi * 0.25)   // 45 degrees horizontal
	camera.Pitch = float32(math.Pi * 0.05) // ~9 degrees up (much lower elevation)
	camera.UpdateOrbit()

	// Initialize renderer
	rasterizer := render.NewRasterizer(width, height)

	// Main loop
	running := true
	lastTime := time.Now()
	frameDuration := time.Second / time.Duration(fps)

	// Mouse state
	var mouseDown bool
	var lastMouseX, lastMouseY int
	modelRotation := float32(0)

	// Determine initial auto-rotate state
	shouldAutoRotate := true // Auto-rotate by default
	if noRotate {
		shouldAutoRotate = false
	}
	if autoRotate {
		shouldAutoRotate = true
	}

	rotationSpeed := float32(0.01) // Rotation speed per frame

	for running {
		// Handle events
		for screen.HasPendingEvent() {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					running = false
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'q', 'Q':
						running = false
					case 'r', 'R':
						// Toggle auto-rotation
						shouldAutoRotate = !shouldAutoRotate
					case 'w', 'W':
						shouldAutoRotate = false // Stop auto-rotation on manual control
						camera.Target.Y += 0.1
						camera.UpdateOrbit()
					case 's', 'S':
						shouldAutoRotate = false
						camera.Target.Y -= 0.1
						camera.UpdateOrbit()
					case 'a', 'A':
						shouldAutoRotate = false
						camera.Target.X -= 0.1
						camera.UpdateOrbit()
					case 'd', 'D':
						shouldAutoRotate = false
						camera.Target.X += 0.1
						camera.UpdateOrbit()
					case 'e', 'E':
						shouldAutoRotate = false
						modelRotation += 0.1
					case ' ':
						// Space bar: toggle auto-rotation
						shouldAutoRotate = !shouldAutoRotate
					case '1':
						rasterizer.Mode = render.RenderModeSolid
					case '2':
						rasterizer.Mode = render.RenderModeWireframe
					}
				}

			case *tcell.EventMouse:
				x, y := ev.Position()
				buttons := ev.Buttons()

				if buttons&tcell.Button1 != 0 { // Left button
					shouldAutoRotate = false // Stop auto-rotation on mouse interaction
					if mouseDown {
						deltaX := float32(x - lastMouseX)
						deltaY := float32(y - lastMouseY)
						camera.Rotate(deltaX*0.01, -deltaY*0.01)
					}
					mouseDown = true
					lastMouseX = x
					lastMouseY = y
				} else {
					mouseDown = false
				}

				// Mouse wheel
				if buttons&tcell.WheelUp != 0 {
					shouldAutoRotate = false // Stop auto-rotation on zoom
					camera.Zoom(-0.5)
				} else if buttons&tcell.WheelDown != 0 {
					shouldAutoRotate = false
					camera.Zoom(0.5)
				}

			case *tcell.EventResize:
				width, height = screen.Size()
				camera.Aspect = float32(width) / float32(height)

				// Recalculate optimal distance for new aspect ratio
				effectiveAspect := camera.Aspect * 0.5
				var effectiveFOV float32
				if effectiveAspect < 1.0 {
					effectiveFOV = camera.FOV
				} else {
					effectiveFOV = 2.0 * float32(math.Atan(math.Tan(float64(camera.FOV*0.5))/float64(effectiveAspect)))
				}
				optimalDistance := (maxDim * 1.1) / (2.0 * float32(math.Tan(float64(effectiveFOV*0.5))))
				camera.Distance = optimalDistance * 0.5 // Same 50% reduction as initial setup
				camera.UpdateOrbit()

				// Preserve mode across resize
				oldMode := rasterizer.Mode
				rasterizer = render.NewRasterizer(width, height)
				rasterizer.Mode = oldMode
			}
		}

		// Render
		currentTime := time.Now()
		if currentTime.Sub(lastTime) >= frameDuration {
			lastTime = currentTime

			// Auto-rotate if enabled
			if shouldAutoRotate {
				modelRotation += rotationSpeed
			}

			screen.Clear()

			// Calculate MVP matrix
			// First, center the model at origin
			centerTransform := scene.Translation(-center.X, -center.Y, -center.Z)

			// Apply rotation
			rotation := scene.RotationY(modelRotation)

			// Combine transformations: Rotation -> Center
			model := rotation.Multiply(centerTransform)

			view := camera.GetViewMatrix()
			projection := camera.GetProjectionMatrix()
			mvp := projection.Multiply(view).Multiply(model)

			// Render mesh
			rasterizer.Render(mesh, mvp, screen)

			// Flush to screen
			screen.Flush()
		}

		// Small sleep to prevent CPU spinning
		time.Sleep(time.Millisecond * 10)
	}
}

func loadModel(filepath string) (*scene.Mesh, error) {
	// Get file extension
	ext := ""
	for i := len(filepath) - 1; i >= 0; i-- {
		if filepath[i] == '.' {
			ext = strings.ToLower(filepath[i:])
			break
		}
	}

	switch ext {
	case ".obj":
		l := loader.NewOBJLoader()
		return l.Load(filepath)
	case ".glb", ".gltf":
		l := loader.NewGLBLoader()
		return l.Load(filepath)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (supported: .obj, .glb, .gltf)", ext)
	}
}
