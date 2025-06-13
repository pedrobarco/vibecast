package player

import (
	"os/exec"
)

func PlayWithVLC(url string) error {
	cmd := exec.Command("vlc", "--play-and-exit", url)
	return cmd.Start()
}
