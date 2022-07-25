// Package internal contains functions internal to the renew command.
package internal

import (
	"debug/buildinfo"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// installPath is the path where "go install" would install a binary. It is
// lazy-initialized.
var installPath string

// goEnv holds the results of "go env" run once.
var goEnv struct {
	sync.Once
	m   map[string]string
	err error
}

// getGoEnv returns the value of a key named in "go env".
func getGoEnv(key string) (string, error) {
	goEnv.Once.Do(func() {
		var out []byte
		out, goEnv.err = exec.Command("go", "env", "-json").Output()
		if goEnv.err != nil {
			return
		}
		goEnv.m = make(map[string]string)
		goEnv.err = json.Unmarshal(out, &goEnv.m)
	})

	if goEnv.err != nil {
		return "", goEnv.err
	}

	v, ok := goEnv.m[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in go env", key)
	}
	return v, nil
}

func init() {
	var err error
	installPath, err = findInstallPath()
	if err != nil {
		log.Fatalf("init failed: %v", err)
	}
}

// findInstallPath returns the path where "go install" would install a binary.
func findInstallPath() (string, error) {
	// See "go help install" for more info on Go's installation strategy.
	s, err := getGoEnv("GOBIN")
	if s != "" && err == nil {
		return s, nil
	}

	s, err = getGoEnv("GOPATH")
	if s != "" && err == nil {
		return s + "/bin", nil
	}

	s, err = os.UserHomeDir()
	if s != "" && err == nil {
		return s + "/go/bin", nil
	}

	return "", fmt.Errorf("unable to determine go install path: %v", err)
}

// Binary contains information about a binary that is installed.
type Binary struct {
	Name       string // The binary name, like "foo".
	LocalPath  string // The binary's local path, like "/Users/calvin/go/bin/foo".
	ImportPath string // The binary's import path, like "github.com/bar/foo".
}

// InstalledBinaries returns all binaries under the path where "go install"
// would install a binary.
func InstalledBinaries() ([]Binary, error) {
	entries, err := os.ReadDir(installPath)
	if err != nil {
		return nil, err
	}

	var bins []Binary
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		bin, err := BinaryFor(e.Name())
		if err != nil {
			return nil, err
		}
		bins = append(bins, bin)
	}

	return bins, nil
}

// BinaryFor returns binary metadata for the given binary name.
func BinaryFor(name string) (Binary, error) {
	localPath := filepath.Join(installPath, name)

	buildInfo, err := buildinfo.ReadFile(localPath)
	if err != nil {
		return Binary{}, err
	}

	return Binary{
		Name:       name,
		LocalPath:  localPath,
		ImportPath: buildInfo.Path,
	}, nil
}

// Updater is a client for updating binaries.
type Updater struct {
	mtx    sync.RWMutex
	Stdout io.Writer
	Stderr io.Writer
}

// NewUpdater returns a new Updater with Stdout and Stderr set to os.Stdout and
// os.Stderr respectively.
func NewUpdater() *Updater {
	return &Updater{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Update updates the given binary to the latest version. It runs "go install"
// and writes any errors to w.
func (u *Updater) Update(bin Binary) error {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	cmd := exec.Command("go", "install", bin.ImportPath+"@latest")
	cmd.Stdout = u.Stdout
	cmd.Stderr = u.Stderr
	return cmd.Run()
}
