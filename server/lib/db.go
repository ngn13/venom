package lib

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

func CheckDB() bool {
	st, err := os.Stat("db")
	if os.IsNotExist(err) {
		return false
	}

	if err != nil || !st.IsDir() {
		log.Panic("Cannot access to database directory")
	}

	st, err = os.Stat("db/settings.json")
	if os.IsNotExist(err) {
		return false
	}

	if err != nil || st.IsDir() {
		log.Panic("Cannot access to settings file")
	}

	return true
}

func MakeDB() error {
	err := os.Mkdir("db", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadDB(name string, obj interface{}) error {
	file := path.Join("db", name+".json")
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}

func SaveDB(name string, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	file := path.Join("db", name+".json")
	err = os.WriteFile(file, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
