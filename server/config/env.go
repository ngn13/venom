package config

import "os"

func InDebug() bool {
	return os.Getenv("VENOM_DEBUG") == "1" ||
		os.Getenv("VENOM_DEBUG") == "yes" ||
		os.Getenv("VENOM_DEBUG") == "true"
}

func UseAntiVM() bool {
	return !(os.Getenv("VENOM_DISABLE_ANTIVM") == "1" ||
		os.Getenv("VENOM_DISABLE_ANTIVM") == "yes" ||
		os.Getenv("VENOM_DISABLE_ANTIVM") == "true")
}

func UseAntiDebug() bool {
	return !(os.Getenv("VENOM_DISABLE_ANTIDEBUG") == "1" ||
		os.Getenv("VENOM_DISABLE_ANTIDEBUG") == "yes" ||
		os.Getenv("VENOM_DISABLE_ANTIDEBUG") == "true")
}

func UseAllInt() bool {
	return os.Getenv("VENOM_ALLINT") == "1" ||
		os.Getenv("VENOM_ALLINT") == "yes" ||
		os.Getenv("VENOM_ALLINT") == "true"
}

func GetURL() string {
	envurl := os.Getenv("VENOM_URL")
	if envurl == "" {
		envurl = "http://127.0.0.1:8082"
	} else if envurl[len(envurl)-1] == '/' {
		envurl = envurl[:len(envurl)-1]
	}

	return envurl
}
