package encryption

import (
	"fmt"
	"os"
	"path/filepath"

	"io"

	"github.com/mutemaniac/encryptfile/util"
)

var coverExistFile = false

const bufSize = 1024 * 1024 * 4

func Encrypt(file string, password []byte) error {
	fmt.Println("Encrypt " + file)
	interval := len(password)
	//源文件
	rf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := rf.Close(); err != nil {
			panic(err)
		}
	}()
	//加密后文件
	wf, err := os.Create(file + util.EncryptedSuffix)
	if err != nil {
		return err
	}
	defer func() {
		if err := wf.Close(); err != nil {
			panic(err)
		}
	}()

	//开始加密
	buf := make([]byte, bufSize)
	for {
		n, err := rf.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			buf[i] = buf[i] + password[i%interval]
		}
		if _, err = wf.Write(buf[0:n]); err != nil {
			return err
		}
		wf.Sync()
	}
	return nil
}

func Decrypt(file string, password []byte) error {
	interval := len(password)
	//源文件
	rf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := rf.Close(); err != nil {
			panic(err)
		}
	}()
	//解密后文件
	var extension = filepath.Ext(file)
	decryptedFilename := file[0 : len(file)-len(extension)]

	if !coverExistFile {
		if _, err = os.Stat(file); err == nil {
			fmt.Println("File " + decryptedFilename + " exist. Cover it? All[A], yes[y], no[n]:")
			var input string
			fmt.Scanln(&input)
			if input == "y" {
				//do nothing
			} else if input == "A" {
				coverExistFile = true
			} else {
				return nil
			}
		}
	}

	wf, err := os.Create(decryptedFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := wf.Close(); err != nil {
			panic(err)
		}
	}()

	//开始解密
	buffer := make([]byte, bufSize)
	for {
		n, err := rf.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			buffer[i] = buffer[i] - password[i%interval]
		}
		if _, err = wf.Write(buffer[:n]); err != nil {
			return err
		}
		wf.Sync()
	}
	return nil
}
