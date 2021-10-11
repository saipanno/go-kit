package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	shutil "github.com/termie/go-shutil"
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
func CopyFile(src, dst string, backupDir ...string) (backup string, err error) {

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
		backup = path.Join(backupDir[0], name) + "." + suffix

		err = shutil.CopyFile(dst, backup, false)
		if err != nil {
			err = fmt.Errorf("backup file failed, message is %s", err.Error())
		}

		return
	}

	err = shutil.CopyFile(src, dst, false)
	if err != nil {
		err = fmt.Errorf("copy file %s to %s failed, message is %s", src, dst, err.Error())
	}

	return
}
