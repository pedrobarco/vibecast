package player

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	vlcCmd   *exec.Cmd
	vlcMutex sync.Mutex
	vlcPid   int
)

// PlayWithVLC launches VLC with the given URL if not running, or reuses the same process by sending a new URL.
func PlayWithVLC(url string) error {
	vlcMutex.Lock()
	defer vlcMutex.Unlock()

	// If VLC is not running, start it with the given URL
	if vlcCmd == nil || vlcPid == 0 || !isProcessRunning(vlcPid) {
		cmd := exec.Command("vlc", "--no-video-title-show", "--quiet", url)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Start(); err != nil {
			return err
		}
		vlcCmd = cmd
		vlcPid = cmd.Process.Pid
		go func() {
			_ = cmd.Wait()
			vlcMutex.Lock()
			vlcCmd = nil
			vlcPid = 0
			vlcMutex.Unlock()
		}()
		return nil
	}

	// If VLC is running, send the new URL to the running process using osascript (macOS only)
	// This will use AppleScript to tell VLC to open the new URL in the same instance
	// On Linux, you could use dbus or xdotool, but here we focus on macOS
	if isMac() {
		script := fmt.Sprintf(`tell application "VLC" to OpenURL "%s"`, url)
		osascript := exec.Command("osascript", "-e", script)
		return osascript.Run()
	}

	// On other platforms, fallback: kill and restart VLC with the new URL
	_ = vlcCmd.Process.Kill()
	vlcCmd = nil
	vlcPid = 0
	cmd := exec.Command("vlc", "--no-video-title-show", "--quiet", url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return err
	}
	vlcCmd = cmd
	vlcPid = cmd.Process.Pid
	go func() {
		_ = cmd.Wait()
		vlcMutex.Lock()
		vlcCmd = nil
		vlcPid = 0
		vlcMutex.Unlock()
	}()
	return nil
}

// isProcessRunning checks if a process with the given pid is running.
func isProcessRunning(pid int) bool {
	if pid == 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, sending signal 0 checks for existence
	err = proc.Signal(os.Signal(syscall(0)))
	return err == nil
}

// isMac returns true if running on macOS.
func isMac() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "darwin") ||
		strings.Contains(strings.ToLower(os.Getenv("GOOS")), "darwin") ||
		(strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "apple") && os.Getenv("TERM_PROGRAM_VERSION") != "")
}

// StopVLC stops the VLC process if running.
func StopVLC() error {
	vlcMutex.Lock()
	defer vlcMutex.Unlock()
	if vlcCmd != nil {
		_ = vlcCmd.Process.Kill()
		vlcCmd = nil
		vlcPid = 0
	}
	return nil
}

// syscall is a helper to convert int to os.Signal for signal 0
func syscall(sig int) os.Signal {
	return os.Signal(syscallRaw(sig))
}

// syscallRaw is a platform-specific syscall number for signal 0
func syscallRaw(sig int) int {
	return sig
}
