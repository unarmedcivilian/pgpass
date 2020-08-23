package pgpass

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

// OpenDefault opens default pgpass file, which is ~/.pgpass.
// Current homedir will be retrieved by calling user.Current
// or using $HOME on failure.
func OpenDefault() (f *os.File, err error) {
	var homedir = os.Getenv("HOME")
	usr, err := user.Current()
	if err == nil {
		homedir = usr.HomeDir
	} else if homedir == "" {
		return
	}
	fileInfo, err := os.Stat(path.Join(homedir, ".pgpass"))
	if err != nil {
		return nil, err
	}
	if fileInfo.Mode().Perm()&(1<<2) != 0 {
		return nil, fmt.Errorf("pgpass file too open. set correct permissions using `chmod 600 ~/.pgpass`")
	}

	return os.Open(path.Join(homedir, ".pgpass"))
}