package builder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Ctx struct {
	Key    []byte
	Config []byte
  Debug  bool
  Dir    string
  Out    string
}

type Entry struct {
  Path string
  Dir  bool
}

func InstallSrc(src string, holder string, val string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	data = bytes.ReplaceAll(data, []byte(holder), []byte(val))
	err = os.WriteFile(src, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func InstallSrcsPlain(ctx *Ctx, srcs []string, holder string, val string) error {
	for _, s := range srcs {
		if strings.HasSuffix(s, ".go") {
			err := InstallSrc(s, holder, val)
			if err != nil {
				return fmt.Errorf("failed to install to %s: %s",
					s, err.Error())
			}
		}
	}

	ctx.Config = bytes.ReplaceAll(ctx.Config, []byte(holder), []byte(val))
	return nil
}

func InstallSrcs(ctx *Ctx, srcs []string, holder string, val string) error {
	enc, err := Encode(ctx.Key, []byte(val))
	if err != nil {
		return err
	}

	mapfmt := fmt.Sprintf("MAP$%s$", holder)
	return InstallSrcsPlain(ctx, srcs, mapfmt, enc)
}

func FindSrcs(dir string) ([]Entry, error) {
	srcs := []Entry{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return srcs, fmt.Errorf("failed to read directory: %s", err.Error())
	}

	for _, e := range entries {
		if e.Name() == "tmp" {
			continue
		}

		fp := path.Join(dir, e.Name())
		st, err := os.Stat(fp)
		if err != nil {
  		return srcs, fmt.Errorf("failed to stat %s: %s", fp, err.Error())
		}

		if st.IsDir() {
      srcs = append(srcs, Entry{
        Path: fp,
        Dir: true,
      })
      newsrcs, err := FindSrcs(fp)
      if err != nil {
        return srcs, err
      }
			srcs = append(srcs, newsrcs...)
			continue
		}

		if strings.HasSuffix(fp, ".go") ||
			strings.HasSuffix(fp, ".mod") ||
			strings.HasSuffix(fp, ".sum") {
			srcs = append(srcs, Entry{
        Path: fp,
        Dir: false,
      })
		}
	}

	return srcs, nil
}

func Run(ctx *Ctx) error {
  map_path := path.Join(ctx.Dir, "map.json")
  tmp_path := path.Join(ctx.Dir, "tmp")
  out_path, err := filepath.Abs(ctx.Out)
  if err != nil {
    return fmt.Errorf("failed to resolve output path: %s", err.Error())
  }

	maps, err := LoadMapping(map_path)
	if err != nil {
    return fmt.Errorf("failed to load mapping: %s", err.Error())
	}

	_, err = os.Stat(tmp_path)
	if err != nil && !os.IsNotExist(err) {
    return fmt.Errorf("cannot stat %s: %s", 
      tmp_path, err.Error())
	}

	os.RemoveAll(tmp_path)
	err = os.Mkdir(tmp_path, os.ModePerm)
	if err != nil {
    return fmt.Errorf("failed to create %s directory: %s",
      tmp_path, err.Error())
	}

  if len(ctx.Key)==0 {
	  ctx.Key = MakeRandom(120 + len(ctx.Config))
  }

	found, err := FindSrcs(ctx.Dir)
  if err != nil {
    return fmt.Errorf("failed to find sources: %s", err.Error())
  }

  for i, e := range found {
    found[i].Path = strings.Replace(e.Path, ctx.Dir, "", 1)
  }

	srcs := []string{}
	for _, e := range found {
		dst := path.Join(tmp_path, e.Path)
    src := path.Join(ctx.Dir, e.Path)

    if e.Dir {
      os.Mkdir(dst, os.ModePerm)
      continue
    }

		err = CopyFile(src, dst)
		if err != nil {
			return fmt.Errorf("failed to copy source: %s", 
        err.Error())
		}
    srcs = append(srcs, dst)
  }


	success := 0
	for h, v := range maps {
		switch v.(type) {
		case string:
			val := v.(string)
			err = InstallSrcs(ctx, srcs, h, val)
			if err != nil {
				return fmt.Errorf("failed to install mapping %s: %s", h, err.Error())
			}
			success++
		case []interface{}:
			list := v.([]interface{})
			for i, val := range list {
				subh := fmt.Sprintf("%s.%d", h, i)
				err = InstallSrcs(ctx, srcs, subh, val.(string))
				if err != nil {
					return fmt.Errorf("failed to install mapping %s: %s", h, err.Error())
				}
				success++
			}
		}
	}

	var keyholder string
	for i := 0; i < 5420; i++ {
		keyholder += "ABC"
	}

	err = InstallSrcsPlain(ctx, srcs, keyholder, base64.StdEncoding.EncodeToString(ctx.Key))
	if err != nil {
    return fmt.Errorf("failed to install key: %s", err.Error())
	}

	err = InstallSrcs(ctx, srcs, "CONFIG", string(ctx.Config))
	if err != nil {
    return fmt.Errorf("failed to install config: %s", err.Error())
	}

  key_path := path.Join(tmp_path, "key.dump")
	err = os.WriteFile(key_path, ctx.Key, 0644)
	if err != nil {
    return fmt.Errorf("failed to dump key: %s", err.Error())
	}


	var cmd *exec.Cmd
	if !ctx.Debug {
		cmd = exec.Command("garble", "build", "-ldflags", "-s -w -H=windowsgui", "-o", out_path)
	} else {
		cmd = exec.Command("go", "build", "-o", out_path)
	}

	cmd.Env = append(cmd.Env, os.Environ()...)
	if runtime.GOOS != "windows" {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "CC=x86_64-w64-mingw32-gcc")
		cmd.Env = append(cmd.Env, "CXX=x86_64-w64-mingw32-g++")
	}

	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	cmd.Env = append(cmd.Env, "GOOS=windows")
	cmd.Dir = tmp_path 

	out, err := cmd.CombinedOutput()
	if err != nil {
    return fmt.Errorf("failed to run build command: %s\noutput: %s", 
      err.Error(), string(out))
	}

	if ctx.Debug {
		goto END
	}

	cmd = exec.Command("upx", out_path)
	cmd.Dir = tmp_path 

	out, err = cmd.CombinedOutput()
	if err != nil {
    return fmt.Errorf("failed to run upx command: %s\noutput: %s", 
      err.Error(), string(out))
	}

END:
  return nil
}

