package lang

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"server/global"
	"strings"
)

var Langs []Lang
type Lang struct {
  JSON map[string]string
  Code string
  Name string
}

func CheckLang(name string) bool {
  for _, l := range Langs {
    if l.Code == name {
      return true
    }
  }
  return false
}

func GetLang(name string) string {
  for _, l := range Langs {
    if l.Code == global.Settings.Language {
      return l.JSON[name]
    }
  }
  return "[TRANSLATION NOT FOUND]"
}

func LoadLang(f string) (Lang, error) {
  var res Lang
  data, err := os.ReadFile(f)
  if err != nil {
    return res, err
  }

  err = json.Unmarshal(data, &res.JSON)
  if err != nil {
    return res, err
  }

  res.Name = res.JSON["lang.flag"]+" "+res.JSON["lang.name"]
  res.Code = res.JSON["lang.code"] 
  return res, nil
}

func LoadLangs() {
  entires, err := os.ReadDir("lang")
  if err != nil {
    log.Panic("Cannot access languages")
  }

  for _, e := range entires {
    name := e.Name()
    if !strings.HasSuffix(name, ".json") {
      continue
    }

    fp := path.Join("lang", name)
    newlang, err := LoadLang(fp)
    if err != nil {
      log.Printf("Failed to load %s: %s", 
        name, err.Error())
      continue
    }

    Langs = append(Langs, newlang)
  }
}
