package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunExecutable launches the application at the given path with arguments.
// It starts the process and detaches it so it doesn't block the server.
func RunExecutable(path string, args string) error {
	// Debug logging
	fmt.Printf("[Launcher] Launching: %s Args: %s\n", path, args)

	// Validate path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("executable not found: %s", path)
	}

	// Prepare command
	// On Windows, complex argument parsing can be tricky usually.
	// But os/exec handles most of it.
	// If args is a raw string, we might need to split it if we use Command(name, arg1, arg2...)
	// However, windows often takes the whole command line.
	// Let's try to split fields simply for now, similar to python's shlex or split.
	// For better windows compatibility with spaces in args, might need care.
	// Python `subprocess.Popen` with string args on Windows works well.
	// Go `exec.Command` expects a slice of args.
	
	// Simple split by space (imperfect for quoted args, but a start)
	// TODO: meaningful arg parsing if users use quotes.
	var cmdArgs []string
	if args != "" {
		cmdArgs = strings.Fields(args) // This is naive, doesn't handle "quoted args"
	}

	cmd := exec.Command(path, cmdArgs...)
	
	// Determine working directory (usually the folder containing the exe)
	cmd.Dir = filepath.Dir(path)
	
	// Detach process logic?
	// exec.Start() starts it. If we don't Wait(), it runs in background.
	// But if the parent (us) dies, what happens? 
	// On Windows, sending it to background usually works fine with Start unless we attach pipes.
	
	if err := cmd.Start(); err != nil {
		return err
	}
	
	// We do NOT wait. 
	// To avoid zombies in unix we'd need cleanup, but in Windows it's different.
	// Actually, we should probably Release the process so we don't keep a handle.
	if cmd.Process != nil {
		cmd.Process.Release()
	}

	return nil
}

