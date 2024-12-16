package user

import (
	"server/lib"

	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	ip := lib.GetIP(c)
	if ip != "127.0.0.1" {
		return lib.APIError(c, 404)
	}

	return c.Redirect("/")
}

func Middleware(c *fiber.Ctx) error {
	ip := lib.GetIP(c)
	if ip != "127.0.0.1" && (c.Path() != "/ntp" || c.Method() != "POST") {
		return lib.APIError(c, 404)
	}

	if ip != "127.0.0.1" && c.Path() == "/ntp" && c.Method() == "POST" {
		return c.Next()
	}

	auth := c.Cookies("auth")
	if !lib.CheckUserCookie(auth) &&
		(c.Path() != "/login" && c.Path() != "/complete") {
		return c.Redirect("/login")
	}

	if !lib.CheckUserCookie(auth) &&
		(c.Path() == "/login" || c.Path() == "/complete") {
		return c.Next()
	}

	if !lib.CheckUserCookie(auth) {
		return c.Redirect("/login")
	}

	if lib.CheckUserCookie(auth) &&
		(c.Path() == "/login" || c.Path() == "/complete") {
		return c.Redirect("/home")
	}

	if c.Path() == "/" {
		return c.Redirect("/home")
	}
	return c.Next()
}
