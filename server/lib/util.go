package lib

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"server/global"
	"server/lang"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DiskStatus struct {
  All     int64
  Used    int64
  Free    int64
  Percent float64 
}

func DiskUsage(path string) (disk DiskStatus) {
  var st syscall.Statfs_t
  err := syscall.Statfs(path, &st)
  if err != nil {
    return
  }
  
  disk.All = int64(st.Blocks) * int64(st.Bsize)
  disk.Free = int64(st.Bfree) * int64(st.Bsize)
  disk.Used = disk.All - disk.Free
  disk.Percent = float64(disk.Used)/float64(disk.All)*100
  return disk
}

func Reverse(s []byte) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

func XOR(input, key string) (string, error) {
	var output string
	if len(input) > len(key) {
		return output, errors.New("Key and data length does not match")
	}

	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i])
	}
	return output, nil
}

func GetIP(c *fiber.Ctx) string {
  if c.Get("X-Real-IP")==""{
    return c.IP()
  }
  return c.Get("X-Real-IP")
}

func APIError(c *fiber.Ctx, err int) error {
  return c.JSON(fiber.Map{
    "error": err,
  })
}

func APIReturn(c *fiber.Ctx, m fiber.Map) error {
  m["error"] = 0
  return c.JSON(m)
}


func Render(c *fiber.Ctx, p string, m fiber.Map) error {
  m["version"] = VERSION
  m["lang"] = global.Settings.Language 
  m["get"]  = lang.GetLang
  m["wt"]  = GetWKTime 
  m["h"]  = BytesToHuman 
  m["u"]  = URLEncode 
  return c.Render(p, m)
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func MakeRandom(n int) string {
  b := make([]rune, n)
  for i := range b {
    b[i] = chars[rand.Intn(len(chars))]
  }
  return string(b)
}

func GetHash(str string) string {
  hasher := sha512.New()
  hasher.Write([]byte(str))
  return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func GetBuild(agent global.TAgent) *global.TBuild {
  for _, b := range global.Builds {
    if b.Token == agent.Token {
      return &b
    }
  }
  return nil
}

func LFISafe(s string) bool {
  chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  for _, c1 := range s {
    fail := true
    for _, c2 := range chars {
      if c1 == c2 {
        fail = false
        break
      }
    }

    if fail {
      return false
    }
  }

  return true 
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

func GetWKTime(t string) string {
  ts, err := strconv.ParseInt(t, 10, 64)
  if err != nil {
    return ""
  }

  unix := (ts/1000000)-11644473600
  res := time.Unix(unix, 0)
  return res.Format("15:04:05 02/01/06")
}

func URLEncode(u string) string {
  return url.QueryEscape(u)
}
