package lib

import (
	"agent/vars"
	"net"
	"strings"
	"syscall"

	"github.com/StackExchange/wmi"
	. "github.com/klauspost/cpuid/v2"
	"github.com/mitchellh/go-ps"
)

type Win32_DiskDrive struct {
	PNPDeviceID string
	Size        uint64
}

func CheckVM_Disk() bool {
	var disks []Win32_DiskDrive
	q := wmi.CreateQuery(&disks, "")
	err := wmi.Query(q, &disks)
	if err != nil {
		return false
	}

	for _, d := range disks {
		for _, b := range vars.Bad_disknames {
			if strings.Contains(strings.ToLower(d.PNPDeviceID), Decode(b)) {
				return true
			}
		}
	}

	return false
}

func CheckVM_MAC() bool {
	infs, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, i := range infs {
		mac := i.HardwareAddr.String()
		for _, b := range vars.Bad_macs {
			if strings.HasPrefix(strings.ToUpper(mac), Decode(b)) {
				return true
			}
		}
	}

	return false
}

func CheckVM() bool {
	if CPU.VM() {
		return true
	}

	if CheckVM_Disk() {
		return true
	}

	if CheckVM_MAC() {
		return true
	}

	return false
}

func CheckDebug_Proc() bool {
	procs, err := ps.Processes()
	if err != nil {
		return false
	}

	for _, p := range procs {
		for _, b := range vars.Bad_procs {
			if strings.HasPrefix(
				strings.ToLower(p.Executable()), strings.ToLower(Decode(b))) {
				return true
			}
		}
	}

	return false
}

func CheckDebug_Present() bool {
	var (
		ModKernel32           = syscall.NewLazyDLL(Decode(vars.Path_kerneldll))
		ProcIsDebuggerPresent = ModKernel32.NewProc(Decode(vars.Key_isdebugger))
	)

	flag, _, _ := ProcIsDebuggerPresent.Call()
	return flag == 1
}

func CheckDebug() bool {
	if CheckDebug_Present() {
		return true
	}

	if CheckDebug_Proc() {
		return true
	}

	return false
}
