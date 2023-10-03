package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/abibby/bwmenu/bw"
	"github.com/abibby/bwmenu/dmenu"
	"github.com/atotto/clipboard"
)

func main() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	nameCacheFile := path.Join(cacheDir, "bwmenu_names_cache")
	sessionCacheFile := path.Join(cacheDir, "bwmenu_session_cache")

	c, err := bw.New(sessionCacheFile)
	if err != nil {
		log.Fatal(err)
	}

	namesChan := make(chan []string, 2)
	itemsChan := make(chan []*bw.Item, 1)
	password := make(chan string, 1)
	passwordRequest := make(chan struct{}, 1)

	cache, err := os.ReadFile(nameCacheFile)
	if os.IsNotExist(err) {
	} else if err != nil {
		log.Fatal(err)
	} else {
		namesChan <- strings.Split(string(cache), "\n")
	}
	go func() {
		err := c.Sync()
		if err != nil {
			log.Fatal(err)
		}
		items, err := c.ListItems()
		if errors.Is(err, bw.ErrVaultLocked) {
			passwordRequest <- struct{}{}

			err = c.Unlock(<-password)
			if err != nil {
				log.Fatal(err)
			}
			items, err = c.ListItems()
		}
		if err != nil {
			log.Fatal(err)
		}

		itemsChan <- items

		folders, err := c.ListFolders()
		if err != nil {
			log.Fatal(err)
		}
		foldersMap := map[string]string{}
		for _, f := range folders {
			foldersMap[f.ID] = f.Name
		}
		names := make([]string, 0, len(items))
		for _, item := range items {
			if item.Login != nil {
				name := item.Name
				if fID, ok := item.FolderID.Ok(); ok {
					name = foldersMap[fID] + "/" + name
				}
				names = append(names, name)
			}
		}
		err = os.WriteFile(nameCacheFile, []byte(strings.Join(names, "\n")), 0644)
		if err != nil {
			log.Fatal(err)
		}
		namesChan <- names
	}()

	names := <-namesChan

	name, err := dmenu.Open(names)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case items := <-itemsChan:
			for _, item := range items {
				if item.Name == name {
					err = clipboard.WriteAll(item.Login.Password)
					if err != nil {
						log.Fatal(err)
					}
					return
				}
			}
			os.Exit(1)
		case <-passwordRequest:
			pass, err := dmenu.Open([]string{"."}, dmenu.ReturnNonMatches())
			if err != nil {
				log.Fatal(err)
			}
			password <- pass
		}
	}
}
