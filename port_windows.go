// +build windows

package ps

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const anySize = 1

type mibTcprowOwnerPid struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
	DwOwningPid  uint32
}
type mibTcptableOwnerPid struct {
	DwNumEntries uint32
	Table        [anySize]mibTcprowOwnerPid
}
type PMIB_TCPTABLE_OWNER_PID_ALL *mibTcptableOwnerPid

var (
	// Library
	modiphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")

	// Function
	procGetExtendedTCPTable = modiphlpapi.NewProc("GetExtendedTcpTable")
)

// netStat holds a netstat record
type netStat struct {
	LocalAddr string
	LocalPort uint16
	OwningPid int
}

// associatedPorts list specfic pid netstats include: tcp tcp6 udp udp6
func associatedPorts(pid int) (ports []uint16, err error) {
	stats, err := getProcInet(pid)

	if err != nil {
		return
	}

	for _, stat := range stats {
		ports = append(ports, stat.LocalPort)
	}
	return
}

func getProcInet(pid int) ([]netStat, error) {
	stats := make([]netStat, 0)

	s, err := getTCP4Stat()
	if err != nil {
		return nil, err
	}

	for _, ns := range s {
		if ns.OwningPid != pid {
			continue
		}
		stats = append(stats, ns)
	}

	return stats, nil
}

func getTCP4Stat() ([]netStat, error) {
	var (
		pmibtable PMIB_TCPTABLE_OWNER_PID_ALL
		buf       []byte
		size      uint32
	)

	for {
		if len(buf) > 0 {
			pmibtable = (*mibTcptableOwnerPid)(unsafe.Pointer(&buf[0]))
		}
		err := getExtendedTcpTable(uintptr(unsafe.Pointer(pmibtable)),
			&size,
			true,
			syscall.AF_INET,
			0)
		if err == nil {
			break
		}
		if err != windows.ERROR_INSUFFICIENT_BUFFER {
			return nil, err
		}
		buf = make([]byte, size)
	}

	if int(pmibtable.DwNumEntries) == 0 {
		return nil, nil
	}

	stats := make([]netStat, 0)
	index := int(unsafe.Sizeof(pmibtable.DwNumEntries))
	step := int(unsafe.Sizeof(pmibtable.Table))

	for i := 0; i < int(pmibtable.DwNumEntries); i++ {
		mibs := (*mibTcprowOwnerPid)(unsafe.Pointer(&buf[index]))

		ns := netStat{
			LocalAddr: parseIPv4(mibs.DwLocalAddr),
			LocalPort: decodePort(mibs.DwLocalPort),
			OwningPid: int(mibs.DwOwningPid),
		}
		stats = append(stats, ns)

		index += step
	}
	return stats, nil
}

func getExtendedTcpTable(pTcpTable uintptr, pdwSize *uint32, bOrder bool, ulAf uint32, reserved uint32) (errcode error) {

	r1, _, _ := syscall.Syscall6(procGetExtendedTCPTable.Addr(), 6, pTcpTable, uintptr(unsafe.Pointer(pdwSize)),
		getUintptrFromBool(bOrder), uintptr(ulAf), 5, uintptr(reserved))
	if r1 != 0 {
		errcode = syscall.Errno(r1)
	}
	return
}

func decodePort(port uint32) uint16 {
	return syscall.Ntohs(uint16(port))
}

func parseIPv4(addr uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", addr&255, addr>>8&255, addr>>16&255, addr>>24&255)
}
func getUintptrFromBool(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}
