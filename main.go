package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mutemaniac/crypfile/encryption"
	"github.com/mutemaniac/crypfile/util"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	// If not enough args, return help text
	if len(os.Args) < 3 {
		printHelp()
		os.Exit(0)
	}

	function := os.Args[1]

	switch function {
	case "help":
		printHelp()
	case "encrypt":
		encryptHandle()
	case "decrypt":
		decryptHandle()
	default:
		fmt.Println("format error.")
		os.Exit(1)
	}

}

func printHelp() {
	fmt.Println("O(∩_∩)O")
}

func encryptHandle() {
	fmt.Println("Start to encryption.")

	//Read all files
	fmt.Println("os.Args", os.Args)
	files, err := findFile(os.Args[2:])
	if err != nil {
		fmt.Println("Cannot find any file.")
		panic(err)
	}
	if len(files) < 1 {
		panic("File not found")
	}
	password := getPassword()

	fmt.Println("\nEncrypting...")
	for _, file := range files {
		err = encryption.Encrypt(file, password)
		if err != nil {
			fmt.Println("Encrypt "+file+" failure. ", err)
		}
	}

	fmt.Println("\nFile successfully protected")
}

func decryptHandle() {
	fmt.Println("O(∩_∩)O")
	if len(os.Args) < 3 {
		println("Missing the path to the file. For more information run CryptoGo help")
		os.Exit(0)
	}

	//Read all files
	files, err := findFile(os.Args[2:])
	if err != nil {
		panic(err)
	}
	files = validateEncryptedFiles(files)
	if len(files) < 1 {
		panic("File not found")
	}

	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(0)

	fmt.Println("\nDecrypting...")
	for _, file := range files {
		err := encryption.Decrypt(file, password)
		if err != nil {
			fmt.Println("Decrypt "+file+" failure. ", err)
		}
	}
	fmt.Println("\nFile successfully decrypted.")
}

func getPassword() []byte {
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(0)
	fmt.Print("\nConfirm password: ")
	password2, _ := terminal.ReadPassword(0)
	if !validatePassword(password, password2) {
		fmt.Print("\nPasswords do not match. Please try again.\n")
		return getPassword()
	}
	return password
}

func validatePassword(password1 []byte, password2 []byte) bool {
	if !bytes.Equal(password1, password2) {
		return false
	}

	return true
}

// findFile find all files, return the absolut path of files
func findFile(filestrs []string) ([]string, error) {
	var files []string
	for _, filestr := range filestrs {
		subfiles, err := filepath.Glob(filestr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		files = append(files, subfiles...)
	}

	for i, file := range files {
		if fileinfo, err := os.Stat(file); os.IsNotExist(err) && fileinfo.IsDir() {
			continue
		}
		adspath, err := filepath.Abs(file)
		if err != nil {
			continue
		}
		files[i] = adspath
	}
	return files, nil
}

func validateEncryptedFiles(files []string) []string {
	var encryptedFiles []string
	for _, file := range files {
		if filepath.Ext(file) == util.EncryptedSuffix {
			encryptedFiles = append(encryptedFiles, file)
		}
	}
	return encryptedFiles
}
