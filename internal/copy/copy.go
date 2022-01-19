package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// File copies a single file from src to dst
func FileCopy(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func DirCopy(src string, dst string) error {
	var err error
	var srcinfo os.FileInfo
	var fds []os.FileInfo

	fmt.Println("Copying directory " + src + " to " + dst)
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := src + string(os.PathSeparator) + fd.Name()
		dstfp := dst + string(os.PathSeparator) + fd.Name()

		if fd.IsDir() {
			if err = DirCopy(srcfp, dstfp); err != nil {
				fmt.Println(err) // continues if can't copy part
			}
		} else {
			if err = FileCopy(srcfp, dstfp); err != nil {
				fmt.Println(err) // continues if can't copy part
			}
		}
	}

	return nil
}
