package lib

import (
	"agent/vars"
	"encoding/json"
	"errors"
)

var Cfg vars.TConfig = vars.TConfig{
	Quiet: false,
}

func ReadConfig() error {
	dec := Decode(vars.Config)
	if dec == "" {
		return errors.New("failed to decode config")
	}

	err := ParseConfig(string(dec))
	if err != nil {
		return err
	}

	if len(Cfg.Token) < 15 {
		return errors.New("bad token length")
	}

	return nil
}

func ParseConfig(c string) error {
	err := json.Unmarshal([]byte(c), &Cfg)
	if err != nil {
		return err
	}
	return nil
}
