package player

import (
	"fmt"
	"os/exec"
	"runtime"
)

// PlayWithVLC opens the given URL in VLC using the OS default handler.
// This will reuse the running VLC instance if possible and show the video window.
func PlayWithVLC(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS: use 'open -a VLC <url>'
		cmd = exec.Command("open", "-a", "VLC", url)
	case "linux":
		// Linux: use 'xdg-open <url>'
		cmd = exec.Command("xdg-open", url)
	case "windows":
		// Windows: use 'start <url>'
		cmd = exec.Command("cmd", "/C", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
