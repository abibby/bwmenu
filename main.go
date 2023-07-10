package main

import (
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
	cacheFile := path.Join(cacheDir, "bwmenu_names_cache")

	namesChan := make(chan []string, 2)
	itemsChan := make(chan []*bw.Item, 1)

	cache, err := os.ReadFile(cacheFile)
	if os.IsNotExist(err) {
	} else if err != nil {
		log.Fatal(err)
	} else {
		namesChan <- strings.Split(string(cache), "\n")
	}
	go func() {
		items, err := bw.ListItems()
		if err != nil {
			log.Fatal(err)
		}
		itemsChan <- items
		names := make([]string, 0, len(items))
		for _, item := range items {
			if item.Login != nil {
				names = append(names, item.Name)
			}
		}
		err = os.WriteFile(cacheFile, []byte(strings.Join(names, "\n")), 0644)
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
	items := <-itemsChan
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
}
