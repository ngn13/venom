package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

var Key []byte

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
		return output, errors.New("key and data length does not match")
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

func Encode(str string) (string, error) {
	xored, err := XOR([]byte(str), Key)
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

func main() {
	if len(os.Args) != 2 {
		Print("Venom encode tool")
		Print("=======================================================")
		Print("Reads ./tmpbuild/key.dump and encodes provided string\n")
		Print("Usage: %s <string>", os.Args[0])
		return
	}

	var err error
	Key, err = os.ReadFile("tmpbuild/key.dump")
	if err != nil {
		Print("Failed to load tmpbuild/key.dump")
		os.Exit(1)
	}
	Print("Loaded key: %s", GetMD5(Key))

	res, err := Encode(os.Args[1])
	if err != nil {
		Print("Failed to encode string: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s", res)
}
