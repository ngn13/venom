package user

import (
	"math"
	"os"
	"path"
	"server/global"
	"server/lib"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetFiles(name string, id string) []global.TFile {
  var files []global.TFile
  
  for _, a := range global.Agents {
    if id != "" && id != a.ID {
      continue
    }

    if name == "" {
      files = append(files, a.Files...)
      continue
    }

    for _, f := range a.Files {
      if strings.Contains(f.Path, name){
        files = append(files, f)
      }
    }
  }

  return files
}

func Files(c *fiber.Ctx) error {
  name := c.Query("name") 
  id := c.Query("id")
  cur, max, min := GetPage(c)

  files := GetFiles(name, id)
  var results []global.TFile

  for i, f := range files {
    if i >= max || i < min {
      continue
    } 
    results = append(results, f)
  }

  pages := int64(math.Ceil(float64(len(files))/float64(PAGE_SIZE)))
  return lib.Render(c, "files", fiber.Map{
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "files": results,
    "name": name,
    "disk": lib.DiskUsage("/"),
  }) 
} 

func FilesDownload(c *fiber.Ctx) error {
  pth := c.Query("path")

  for _, a := range global.Agents {
    for _, f := range a.Files {
      if f.Path == pth { 
        fp := path.Join("db", "files", f.Data)
        return c.SendFile(fp)
      }
    }
  }

  return lib.Render(c, "files", fiber.Map{
    "error": "File not found",
    "files": GetFiles("", ""),
    "disk": lib.DiskUsage("/"),
  })
}

func FilesRemove(c *fiber.Ctx) error {
  pth := c.Query("path")
  found := false

  for _, a := range global.Agents {
    for i, f := range a.Files {
      if f.Path == pth { 
        fp := path.Join("db", "files", f.Data)
        os.Remove(fp)
        a.Files = append(a.Files[:i], a.Files[i+1:]...)
        lib.SaveDB("agents", global.Agents) 
        found = true
      }
    }
  }

  name := c.Query("name") 
  id := c.Query("id")
  cur, max, min := GetPage(c)

  files := GetFiles(name, id)
  var results []global.TFile

  for i, f := range files {
    if i >= max || i < min {
      continue
    } 
    results = append(results, f)
  }

  pages := int64(math.Ceil(float64(len(files))/float64(PAGE_SIZE)))
  if !found {
    return lib.Render(c, "files", fiber.Map{
      "error": "File not found",
      "pages": pages, 
      "current": cur,
      "size": len(results),
      "files": results,
      "name": name,
      "disk": lib.DiskUsage("/"),
    })
  }
  
  return lib.Render(c, "files", fiber.Map{
    "success": "File has been deleted",
    "pages": pages, 
    "current": cur,
    "size": len(results),
    "files": results,
    "name": name,
    "disk": lib.DiskUsage("/"),
  })
}
