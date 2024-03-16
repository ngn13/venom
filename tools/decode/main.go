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

func Decode(s string) (string, error) {
	rev := Reverse([]byte(s))
	dec, err := base64.StdEncoding.DecodeString(rev)
	if err != nil {
		return "", err
	}

	var n []byte = []byte{}
	for i := range dec {
		n = append(n, (dec[i] - 3))
	}

	xored, err := XOR([]byte(n), Key)
	if err != nil {
		return "", err
	}

	return string(xored), nil
}

func main() {
	if len(os.Args) != 2 {
		Print("Venom decode tool")
		Print("=======================================================")
		Print("Reads ./tmpbuild/key.dump and decodes provided string\n")
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

	res, err := Decode(os.Args[1])
	if err != nil {
		Print("Failed to decode string: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s", res)
}
