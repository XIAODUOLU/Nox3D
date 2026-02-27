# Shell Completion Guide for Nox3D

Nox3D supports shell completion for Bash, Zsh, Fish, and PowerShell. This allows you to use Tab key to auto-complete commands, subcommands, and file paths.

## Quick Start

### Automatic Installation

Use the provided installation script:

```bash
# For Bash
./install-completion.sh bash

# For Zsh
./install-completion.sh zsh

# For Fish
./install-completion.sh fish

# For PowerShell
./install-completion.sh powershell
```

### Manual Installation

#### Bash

**Linux:**
```bash
# System-wide (requires sudo)
sudo ./bin/nox completion bash > /etc/bash_completion.d/nox

# User-specific
mkdir -p ~/.bash_completion.d
./bin/nox completion bash > ~/.bash_completion.d/nox
echo 'source ~/.bash_completion.d/nox' >> ~/.bashrc
source ~/.bashrc
```

**macOS:**
```bash
# Install bash-completion first if needed
brew install bash-completion

# Then install nox completion
./bin/nox completion bash > /usr/local/etc/bash_completion.d/nox
source ~/.bash_profile
```

**Temporary (current session only):**
```bash
source <(./bin/nox completion bash)
```

#### Zsh

```bash
# Make sure completion is enabled
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
mkdir -p ~/.zsh/completion
./bin/nox completion zsh > ~/.zsh/completion/_nox

# Add to fpath in ~/.zshrc
echo 'fpath=(~/.zsh/completion $fpath)' >> ~/.zshrc

# Reload
source ~/.zshrc
```

**Temporary (current session only):**
```bash
source <(./bin/nox completion zsh)
```

#### Fish

```bash
./bin/nox completion fish > ~/.config/fish/completions/nox.fish
```

Completion will be available immediately in new fish sessions.

**Temporary (current session only):**
```bash
./bin/nox completion fish | source
```

#### PowerShell

```powershell
# Generate completion script
./bin/nox completion powershell > nox.ps1

# Add to your PowerShell profile
# Find your profile location with: $PROFILE
Add-Content $PROFILE ". $(pwd)/nox.ps1"

# Or load for current session only
. ./nox.ps1
```

## What Gets Completed

### Commands
```bash
nox <Tab>
# Shows: test, completion, help
```

### Subcommands
```bash
nox test <Tab>
# Shows available .obj, .glb, .fbx files in current directory
```

### Flags
```bash
nox test --<Tab>
# Shows: --fps, --help
```

### Completion Command
```bash
nox completion <Tab>
# Shows: bash, zsh, fish, powershell
```

## Testing Completion

After installation, test that completion works:

```bash
# Type this and press Tab
nox <Tab>

# Should show:
# completion  help  test

# Type this and press Tab
nox test assets/obj/<Tab>

# Should show available .obj files
```

## Troubleshooting

### Bash: Completion not working

1. Check if bash-completion is installed:
   ```bash
   # Linux
   apt-get install bash-completion  # Debian/Ubuntu
   yum install bash-completion      # RHEL/CentOS
   
   # macOS
   brew install bash-completion
   ```

2. Make sure completion is sourced in ~/.bashrc:
   ```bash
   if [ -f /etc/bash_completion ]; then
       . /etc/bash_completion
   fi
   ```

3. Reload your shell:
   ```bash
   source ~/.bashrc
   ```

### Zsh: Completion not working

1. Make sure compinit is called in ~/.zshrc:
   ```bash
   autoload -U compinit
   compinit
   ```

2. Check if the completion file is in fpath:
   ```bash
   echo $fpath
   ```

3. Rebuild completion cache:
   ```bash
   rm -f ~/.zcompdump
   compinit
   ```

### Fish: Completion not working

1. Check if the completion file exists:
   ```bash
   ls ~/.config/fish/completions/nox.fish
   ```

2. Restart fish or run:
   ```bash
   source ~/.config/fish/completions/nox.fish
   ```

### PowerShell: Completion not working

1. Check execution policy:
   ```powershell
   Get-ExecutionPolicy
   # If Restricted, run:
   Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

2. Make sure the script is sourced in your profile:
   ```powershell
   notepad $PROFILE
   # Should contain: . path/to/nox.ps1
   ```

## Uninstalling Completion

### Bash
```bash
# Linux
sudo rm /etc/bash_completion.d/nox
# or
rm ~/.bash_completion.d/nox

# macOS
rm /usr/local/etc/bash_completion.d/nox
```

### Zsh
```bash
rm ~/.zsh/completion/_nox
# Remove fpath line from ~/.zshrc
```

### Fish
```bash
rm ~/.config/fish/completions/nox.fish
```

### PowerShell
```powershell
# Remove the source line from your profile
notepad $PROFILE
# Delete: . path/to/nox.ps1
```

## Advanced Usage

### Custom Completion Location

You can generate the completion script to any location:

```bash
# Bash
./bin/nox completion bash > /path/to/custom/location

# Then source it
source /path/to/custom/location
```

### Multiple Shell Support

If you use multiple shells, install completion for each:

```bash
./install-completion.sh bash
./install-completion.sh zsh
./install-completion.sh fish
```

## How It Works

Nox3D uses the Cobra library's built-in completion system:

1. **Command Completion**: Cobra automatically provides completion for all registered commands and subcommands
2. **Flag Completion**: All flags (like `--fps`) are automatically completed
3. **File Completion**: The `test` command uses `ValidArgsFunction` to filter files by extension (.obj, .glb, .fbx)
4. **Dynamic Completion**: Completion adapts based on context (e.g., only shows valid file extensions)

## Examples

```bash
# Complete command
nox t<Tab>
# Expands to: nox test

# Complete file with extension filter
nox test assets/<Tab>
# Shows only .obj, .glb, .fbx files

# Complete flags
nox test model.obj --f<Tab>
# Expands to: nox test model.obj --fps

# Complete completion shell type
nox completion z<Tab>
# Expands to: nox completion zsh
```

## See Also

- [Cobra Completion Documentation](https://github.com/spf13/cobra/blob/master/shell_completions.md)
- Main README: [README.md](../README.md)
- Code Documentation: [CODE_EXPLANATION.md](../CODE_EXPLANATION.md)
