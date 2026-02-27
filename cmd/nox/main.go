package main

import (
	"fmt"
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
	fps int
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

	// Initialize camera
	width, height := screen.Size()
	camera := scene.NewCamera()
	camera.Aspect = float32(width) / float32(height)
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
					case 'w', 'W':
						camera.Target.Y += 0.1
						camera.UpdateOrbit()
					case 's', 'S':
						camera.Target.Y -= 0.1
						camera.UpdateOrbit()
					case 'a', 'A':
						camera.Target.X -= 0.1
						camera.UpdateOrbit()
					case 'd', 'D':
						camera.Target.X += 0.1
						camera.UpdateOrbit()
					case 'e', 'E':
						modelRotation += 0.1
					}
				}

			case *tcell.EventMouse:
				x, y := ev.Position()
				buttons := ev.Buttons()

				if buttons&tcell.Button1 != 0 { // Left button
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
					camera.Zoom(-0.5)
				} else if buttons&tcell.WheelDown != 0 {
					camera.Zoom(0.5)
				}

			case *tcell.EventResize:
				width, height = screen.Size()
				camera.Aspect = float32(width) / float32(height)
				rasterizer = render.NewRasterizer(width, height)
			}
		}

		// Render
		currentTime := time.Now()
		if currentTime.Sub(lastTime) >= frameDuration {
			lastTime = currentTime

			screen.Clear()

			// Calculate MVP matrix
			model := scene.RotationY(modelRotation)
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
	ext := strings.ToLower(filepath[len(filepath)-4:])

	var l loader.Loader

	switch ext {
	case ".obj":
		l = loader.NewOBJLoader()
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	return l.Load(filepath)
}
