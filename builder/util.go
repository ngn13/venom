package builder

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func MakeRandom(s int) []byte {
	ret := make([]byte, s)
	rand.Read(ret)
	return ret
}

func LoadMapping(file string) (map[string]interface{}, error) {
	var mapping map[string]interface{}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func Reverse(s []byte) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

func XOR(input, key []byte) ([]byte, error) {
	var output []byte
	if len(input) > len(key) {
    return output, fmt.Errorf("key and data length does not match (%d > %d)",
      len(input), len(key))
	}

	for i := 0; i < len(input); i++ {
		output = append(output, input[i]^key[i])
	}

	return output, nil
}

func GetMD5(v []byte) string {
	sum := md5.Sum(v)
	return fmt.Sprintf("%x", sum)
}

func Encode(key []byte, str []byte) (string, error) {
	xored, err := XOR(str, key)
	if err != nil {
		return "", err
	}

	var n []byte = []byte{}
	for i := range xored {
		n = append(n, (xored[i] + 3))
	}

	enc := base64.StdEncoding.EncodeToString([]byte(n))
	return Reverse([]byte(enc)), nil
}

func CopyFile(src string, dst string) error {
	all, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, all, 0644)
	if err != nil {
		return err
	}

	return nil
}
