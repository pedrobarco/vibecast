package player

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	vlcCmd   *exec.Cmd
	vlcMutex sync.Mutex
	vlcConn  net.Conn
)

// StartVLCWithRC starts VLC with the RC (remote control) interface on a local TCP port.
func StartVLCWithRC(port int) error {
	vlcMutex.Lock()
	defer vlcMutex.Unlock()
	if vlcCmd != nil {
		return nil // Already started
	}
	args := []string{
		"--intf", "rc",
		"--rc-host", fmt.Sprintf("127.0.0.1:%d", port),
		"--no-video-title-show",
		"--quiet",
	}
	vlcCmd = exec.Command("vlc", args...)
	return vlcCmd.Start()
}

// ConnectToVLC connects to the VLC RC interface.
func ConnectToVLC(port int) error {
	vlcMutex.Lock()
	defer vlcMutex.Unlock()
	if vlcConn != nil {
		return nil // Already connected
	}
	var err error
	for i := 0; i < 10; i++ {
		vlcConn, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return err
}

// PlayWithVLC switches the currently playing channel in a persistent VLC instance.
func PlayWithVLC(url string) error {
	const port = 4212
	vlcMutex.Lock()
	defer vlcMutex.Unlock()

	// Start VLC if not running
	if vlcCmd == nil {
		if err := StartVLCWithRC(port); err != nil {
			return err
		}
	}
	// Connect to RC interface if not connected
	if vlcConn == nil {
		if err := ConnectToVLC(port); err != nil {
			return err
		}
	}

	// Send "add" command to play the new URL
	cmd := fmt.Sprintf("add %s\n", url)
	_, err := vlcConn.Write([]byte(cmd))
	return err
}

// StopVLC stops the persistent VLC instance.
func StopVLC() error {
	vlcMutex.Lock()
	defer vlcMutex.Unlock()
	if vlcConn != nil {
		vlcConn.Write([]byte("quit\n"))
		vlcConn.Close()
		vlcConn = nil
	}
	if vlcCmd != nil {
		vlcCmd.Process.Kill()
		vlcCmd = nil
	}
	return nil
}
