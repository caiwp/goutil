package file

import (
	"os"
)

func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func FileOrSymlinkExists(filename string) (bool, error) {
	if _, err := os.Lstat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}

func ReadDirNoStat(dirname string) ([]string, error) {
	if dirname == "" {
		dirname = "."
	}

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}
