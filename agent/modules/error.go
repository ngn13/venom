package modules

import (
	"agent/lib"
	"agent/vars"
	"sync"
)

func DisplayError(wg *sync.WaitGroup) {
	lib.MessageBox(0, lib.Cfg.Error.Message,
		lib.Cfg.Error.Title, 0x00000010|0x00010000)

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_error))
	wg.Done()
}
