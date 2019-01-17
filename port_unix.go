// +build linux darwin

package ps

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func associatedPorts(pid int) (ports []uint16, err error) {
	//cmd := fmt.Sprintf("ss -l -p -n | grep \"pid=%d,\"", pid)
	cmd := fmt.Sprintf("lsof -p 18275 | grep LISTEN", pid)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return ports, fmt.Errorf("Failed to execute command: %s", err)
	}
	outputLines := strings.Split(string(out), "\n")
	var port string
	for _, line := range outputLines {
		for i := 0; i < len(line); i++ {
			r := line[i : i+1]
			if r == " " && len(port) > 1 {
				p, err := strconv.Atoi(port)
				if err != nil {
					return ports, err
				}
				ports = append(ports, uint16(p))
				port = ""
				continue
			}
			// skip everything between [ and ]
			if r == "[" {
				for line[i:i+1] != "]" {
					i++
				}
			}
			_, err := strconv.Atoi(r)
			if err != nil {
				continue
			}
			if string(r) == " " {
				continue
			}
			if i > 0 && string(out[i-1]) == ":" {
				port = string(r)
				continue
			}
			if port != "" {
				port += string(r)
			}
		}
	}

	return
}
