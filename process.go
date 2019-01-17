// ps provides an API for finding and listing processes in a platform-agnostic
// way.
//
// NOTE: If you're reading these docs online via GoDocs or some other system,
// you might only see the Unix docs. This project makes heavy use of
// platform-specific implementations. We recommend reading the source if you
// are interested.
package ps

import "errors"

// Process is the generic interface that is implemented on every platform
// and provides common operations for processes.
type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}

// Processes returns all processes.
//
// This of course will be a point-in-time snapshot of when this method was
// called. Some operating systems don't provide snapshot capability of the
// process table, in which case the process table returned might contain
// ephemeral entities that happened to be running when this was called.
func Processes() ([]Process, error) {
	return processes()
}

// FindProcess looks up a single process by pid.
//
// Process will be nil and error will be nil if a matching process is
// not found.
func FindProcess(pid int) (Process, error) {
	return findProcess(pid)
}

func ProcessByName(executable string) (p Process, err error) {
	processes, err := Processes()
	if err != nil {
		return
	}
	for _, p := range processes {
		if p.Executable() == executableName(executable) {
			return p, err
		}
	}
	err = errors.New("process not found")
	return
}

func AssociatedPorts(pid int) ([]uint16, error) {
	return associatedPorts(pid)
}

func executableName(name string) string {
	return name
}
