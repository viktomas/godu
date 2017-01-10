package main

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getIgnoredFolders() map[string]struct{} {
	usr, err := user.Current()
	if err != nil {
		log.Println("Wasn't able to retrieve current user at runtime")
		return map[string]struct{}{}
	}
	ignoreFileName := filepath.Join(usr.HomeDir, ".goduignore")
	if _, err := os.Stat(ignoreFileName); os.IsNotExist(err) {
		return map[string]struct{}{}
	}
	ignoreFile, err := os.Open(ignoreFileName)
	if err != nil {
		log.Printf("Failed to read ingorefile because %s\n", err.Error())
		return map[string]struct{}{}
	}
	defer ignoreFile.Close()
	scanner := bufio.NewScanner(ignoreFile)
	ignoredFolders := map[string]struct{}{}
	for scanner.Scan() {
		ignoredFolders[scanner.Text()] = struct{}{}
	}
	return ignoredFolders
}
