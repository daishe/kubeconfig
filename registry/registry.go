package registry

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"
)

func KubectlConfigPath() (string, error) {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		home, ok = os.LookupEnv("USERPROFILE") // windows
	}
	if !ok {
		return "", fmt.Errorf("$HOME environment vatiabble not set")
	}

	path, err := filepath.Abs(home)
	if err != nil {
		return "", fmt.Errorf("cannot obtain the absolute path to the user home directory: %w", err)
	}

	return filepath.Join(path, ".kube", "config"), nil
}

func Find(registry string, hash []byte) (string, bool, error) {
	ls, err := List(registry)
	if err != nil {
		return "", false, err
	}

	for _, f := range ls {
		h, err := Hash(f)
		if err != nil {
			return "", false, err
		}
		if bytes.Equal(hash, h) {
			return f, true, nil
		}
	}

	return "", false, nil
}

func Path() (string, error) {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("$HOME environment vatiabble not set")
	}

	path, err := filepath.Abs(home)
	if err != nil {
		return "", fmt.Errorf("cannot obtain the absolute path to the user home directory: %w", err)
	}

	return filepath.Join(path, ".kubeconfig"), nil
}

func OverrodePath(override string) (string, error) {
	path, err := filepath.Abs(override)
	if err != nil {
		return "", fmt.Errorf("cannot obtain the absolute path to the kubeconfig registry: %w", err)
	}
	return path, nil
}

func PathToName(registry string, path string) string {
	path = strings.TrimPrefix(path, registry)
	path = strings.TrimPrefix(path, string(os.PathSeparator))
	path = strings.ReplaceAll(path, string(os.PathSeparator), "/")
	return path
}

func NameToPath(registry string, name string) string {
	return filepath.Join(registry, strings.ReplaceAll(name, "/", string(os.PathSeparator)))
}

func Hash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", path, err)
	}
	defer f.Close()

	hash := sha3.New256()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("cannot while readint file %q: %w", path, err)
	}
	return hash.Sum(nil), nil
}

func Exist(registry string, name string) (bool, error) {
	path := filepath.Join(registry, strings.ReplaceAll(name, "/", string(os.PathSeparator)))
	stat, err := os.Stat(path)
	if err != nil {
		if err != os.ErrNotExist {
			return false, fmt.Errorf("cannot stat file %q: %w", path, err)
		}
		return false, nil
	}
	return stat.Mode().IsRegular(), nil
}

func List(registry string) ([]string, error) {
	var res []string
	if err := listDirReclusive(registry, &res); err != nil {
		return nil, err
	}
	sort.Strings(res)
	return res, nil
}

func listDirReclusive(path string, result *[]string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("connot read directory %q: %w", path, err)
	}

	for _, file := range files {
		switch {
		case file.Type().IsRegular():
			*result = append(*result, filepath.Join(path, file.Name()))
		case file.IsDir():
			if err := listDirReclusive(filepath.Join(path, file.Name()), result); err != nil {
				return err
			}
		}
	}

	return nil
}

func ListWithCmp(registry string, hash []byte) ([]string, []bool, error) {
	ls, err := List(registry)
	if err != nil {
		return nil, nil, err
	}

	cmp := make([]bool, 0, len(ls))

	for _, f := range ls {
		h, err := Hash(f)
		if err != nil {
			return nil, nil, err
		}
		cmp = append(cmp, bytes.Equal(hash, h))
	}

	return ls, cmp, nil
}

func Read(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", path, err)
	}
	return f, nil
}

func Write(path string, content io.Reader) error {
	return writeWithFlag(path, content, os.O_RDWR|os.O_CREATE|os.O_EXCL)
}

func ForceWrite(path string, content io.Reader) error {
	return writeWithFlag(path, content, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
}

func writeWithFlag(path string, content io.Reader, flag int) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create directory %q for file %q: %w", dir, path, err)
	}

	f, err := os.OpenFile(path, flag, 0640)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("writing to file %q filed: %w", path, err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("writing to file %q filed: %w", path, err)
	}

	return nil
}
