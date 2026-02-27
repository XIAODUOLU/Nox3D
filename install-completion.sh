#!/bin/bash

# Shell completion installation script for Nox3D

set -e

BINARY="./bin/nox"
SHELL_TYPE="${1:-bash}"

if [ ! -f "$BINARY" ]; then
    echo "Error: nox binary not found. Please run 'make build' first."
    exit 1
fi

echo "Installing shell completion for $SHELL_TYPE..."

case "$SHELL_TYPE" in
    bash)
        # Check if running on Linux or macOS
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            COMPLETION_DIR="/etc/bash_completion.d"
            if [ ! -d "$COMPLETION_DIR" ]; then
                COMPLETION_DIR="$HOME/.bash_completion.d"
                mkdir -p "$COMPLETION_DIR"
            fi
            
            if [ -w "$COMPLETION_DIR" ]; then
                $BINARY completion bash > "$COMPLETION_DIR/nox"
                echo "✓ Bash completion installed to $COMPLETION_DIR/nox"
            else
                echo "Installing to user directory (requires sudo for system-wide)..."
                mkdir -p "$HOME/.bash_completion.d"
                $BINARY completion bash > "$HOME/.bash_completion.d/nox"
                echo "✓ Bash completion installed to $HOME/.bash_completion.d/nox"
                echo ""
                echo "Add this to your ~/.bashrc:"
                echo "  source $HOME/.bash_completion.d/nox"
            fi
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            COMPLETION_DIR="/usr/local/etc/bash_completion.d"
            if [ ! -d "$COMPLETION_DIR" ]; then
                COMPLETION_DIR="$HOME/.bash_completion.d"
                mkdir -p "$COMPLETION_DIR"
            fi
            
            $BINARY completion bash > "$COMPLETION_DIR/nox"
            echo "✓ Bash completion installed to $COMPLETION_DIR/nox"
        fi
        
        echo ""
        echo "To activate completion, run:"
        echo "  source <($BINARY completion bash)"
        echo ""
        echo "Or restart your shell."
        ;;
        
    zsh)
        # Find zsh completion directory
        if [ -n "$ZSH_VERSION" ]; then
            COMPLETION_DIR="${fpath[1]}"
        else
            COMPLETION_DIR="$HOME/.zsh/completion"
        fi
        
        mkdir -p "$COMPLETION_DIR"
        $BINARY completion zsh > "$COMPLETION_DIR/_nox"
        echo "✓ Zsh completion installed to $COMPLETION_DIR/_nox"
        
        echo ""
        echo "Make sure you have the following in your ~/.zshrc:"
        echo "  autoload -U compinit; compinit"
        echo ""
        echo "Then restart your shell or run:"
        echo "  source ~/.zshrc"
        ;;
        
    fish)
        COMPLETION_DIR="$HOME/.config/fish/completions"
        mkdir -p "$COMPLETION_DIR"
        $BINARY completion fish > "$COMPLETION_DIR/nox.fish"
        echo "✓ Fish completion installed to $COMPLETION_DIR/nox.fish"
        
        echo ""
        echo "Completion will be available in new fish sessions."
        echo "Or run: source $COMPLETION_DIR/nox.fish"
        ;;
        
    powershell)
        echo "Generating PowerShell completion script..."
        $BINARY completion powershell > nox.ps1
        echo "✓ PowerShell completion script generated: nox.ps1"
        
        echo ""
        echo "To use it, add this to your PowerShell profile:"
        echo "  . $(pwd)/nox.ps1"
        echo ""
        echo "Or run once in your session:"
        echo "  . ./nox.ps1"
        ;;
        
    *)
        echo "Error: Unsupported shell type: $SHELL_TYPE"
        echo "Supported shells: bash, zsh, fish, powershell"
        exit 1
        ;;
esac

echo ""
echo "Installation complete!"
