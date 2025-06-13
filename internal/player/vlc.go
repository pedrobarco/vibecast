package player

import (
	"fmt"
	"net"
	"time"
)

// PlayWithVLC connects to a running VLC instance with RC interface enabled and tells it to play the given URL.
// VLC must be started by the user with: vlc --intf rc --rc-host 127.0.0.1:4212
func PlayWithVLC(url string) error {
	const addr = "127.0.0.1:4212"
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return fmt.Errorf("could not connect to VLC RC interface at %s: %w", addr, err)
	}
	defer conn.Close()

	// Wait for VLC RC prompt (may not be necessary, but helps with timing)
	time.Sleep(200 * time.Millisecond)

	// Send "add" command to play the new URL
	cmd := fmt.Sprintf("add %s\n", url)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command to VLC: %w", err)
	}
	return nil
}
