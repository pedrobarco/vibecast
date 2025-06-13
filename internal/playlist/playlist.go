package playlist

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
)

type Channel struct {
	Name string
	URL  string
}

func LoadM3U(pathOrURL string) ([]Channel, error) {
	var r io.Reader
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		r = resp.Body
	} else {
		f, err := os.Open(pathOrURL)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	return parseM3U(r)
}

func parseM3U(r io.Reader) ([]Channel, error) {
	var channels []Channel
	scanner := bufio.NewScanner(r)
	var name string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#EXTINF:") {
			// Example: #EXTINF:-1 tvg-id="" tvg-name="Channel" ... ,Channel Name
			parts := strings.SplitN(line, ",", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[1])
			}
		} else if line != "" && !strings.HasPrefix(line, "#") {
			channels = append(channels, Channel{Name: name, URL: line})
			name = ""
		}
	}
	return channels, scanner.Err()
}
