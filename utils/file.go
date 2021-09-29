package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

// FileExist ...
func FileExist(f string) bool {

	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// DirExist ...
func DirExist(dir string, autoCreate ...bool) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		if len(autoCreate) > 0 && autoCreate[0] {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return false
			}

			return true
		}

		return false
	}

	return true
}

// CopyFile ...
func CopyFile(src, dst string, backupDir ...string) (written int64, err error) {

	suffix := time.Now().Format("20060102150405")

	if !FileExist(src) {
		err = fmt.Errorf("file %s is not exists", src)
		return
	}

	if FileExist(dst) && len(backupDir) > 0 && len(backupDir[0]) > 0 {

		if !DirExist(backupDir[0], true) {
			err = fmt.Errorf("create backup dir %s failed, message is %s", backupDir[0], err.Error())
			return
		}

		name := filepath.Base(dst)
		backup := path.Join(backupDir[0], name) + "." + suffix

		var sf *os.File
		sf, err = os.Open(src)
		if err != nil {
			err = fmt.Errorf("open file %s failed, message is %s", src, err.Error())
			return
		}
		defer sf.Close()

		var bf *os.File
		bf, err = os.Create(backup)
		if err != nil {
			err = fmt.Errorf("open file %s failed, message is %s", backup, err.Error())
			return
		}
		defer bf.Close()

		written, err = io.Copy(sf, bf)
		if err != nil {
			err = fmt.Errorf("copy file failed, message is %s", err.Error())
			return
		}

		err = bf.Sync()
		if err != nil {
			err = fmt.Errorf("sync file %s failed, message is %s", backup, err.Error())
			return
		}
	}

	err = os.Rename(src, dst)
	if err != nil {
		err = fmt.Errorf("rename file %s to %s failed, message is %s", src, dst, err.Error())
	}

	return
}
