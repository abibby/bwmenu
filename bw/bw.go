package bw

import (
	"encoding/json"
	"os/exec"
)

type Item struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Login  *Login
}

type Login struct {
	// URIs []*URI `json:"uris"`
	Password string `json:"password"`
}

func bw(v any, args ...string) error {
	b, err := exec.Command("bw", args...).Output()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func ListItems() ([]*Item, error) {
	items := []*Item{}
	err := bw(&items, "list", "items")
	return items, err
}
