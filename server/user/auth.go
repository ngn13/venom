package user

import (
	"log"
	"server/global"
	"server/lang"
	"server/lib"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
  lib.RemoveUserCookie(c.Cookies("auth"))
  c.ClearCookie()
  return c.Redirect("/login")
}

func Login(c *fiber.Ctx) error {
  if lang.CheckLang(c.Query("lang")) && !lib.CheckDB() {
    global.Settings.Language = strings.Clone(c.Query("lang"))
  }

  var selected lang.Lang
  for _, l := range lang.Langs{
    if l.Code == global.Settings.Language {
      selected = l
    }
  }

  if !lib.CheckDB() {
    return lib.Render(c, "setup", fiber.Map{
      "selected": selected,
      "langs": lang.Langs,
    })
  }

  return lib.Render(c, "login", fiber.Map{})
}

func LoginPost(c *fiber.Ctx) error {
  if !lib.CheckDB() {
    return lib.Render(c, "setup", fiber.Map{})
  }

  data := struct{
    Pwd string `form:"pwd"`
  }{}

  err := c.BodyParser(&data)
  if err != nil {
    log.Printf("Error on login: %s", err.Error())
    return lib.Render(c, "login", fiber.Map{
      "error": lang.GetLang("login.fail"),
    })
  }

  if global.Settings.Password != lib.GetHash(data.Pwd) {
    return lib.Render(c, "login", fiber.Map{
      "error": lang.GetLang("login.fail"),
    })
  }

  c.Cookie(&fiber.Cookie{
    Name: "auth",
    Value: lib.AddUserCookie(),
  })

  return c.Redirect("/home")
}

func Complete(c *fiber.Ctx) error {
  if lib.CheckDB() {
    return c.Redirect("/login")
  }
 
  sd := struct{
    Lang string `form:"lang"`
    Pwd  string `form:"pwd"`
  }{}

  err := c.BodyParser(&sd)
  if err != nil || sd.Lang == "" || 
      sd.Pwd == "" || !lang.CheckLang(sd.Lang) {
    return c.Redirect("/login")
  }

  encpwd := lib.GetHash(sd.Pwd)
  log.Print("Completing setup")
  log.Printf("Language: %s", sd.Lang)
  log.Printf("Password: %s", encpwd) 
  lib.MakeDB()

  global.Settings.Password = encpwd 
  global.Settings.Language = strings.Clone(sd.Lang) 
  lib.SaveDB("settings", global.Settings)
  return c.Redirect("/login")
}
