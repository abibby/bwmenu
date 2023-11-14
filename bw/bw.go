package bw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/abibby/nulls"
)

type Object struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}

type Item struct {
	Object
	FolderID *nulls.String `json:"folderId"`
	Login    *Login
}

type Login struct {
	// URIs []*URI `json:"uris"`
	Password string `json:"password"`
}

type Client struct {
	sessionCacheFile string
	session          string
}

var ErrVaultLocked = errors.New("vault locked")

func New(sessionCacheFile string) (*Client, error) {
	b, err := os.ReadFile(sessionCacheFile)
	if errors.Is(err, os.ErrNotExist) {
	} else if err != nil {
		return nil, err
	}
	return &Client{
		sessionCacheFile: sessionCacheFile,
		session:          string(b),
	}, nil
}

func (c *Client) bw(v any, args ...string) error {
	b, err := exec.Command("bw", append([]string{"--nointeraction", "--session", c.session}, args...)...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(b), "Vault is locked.") {
			err = ErrVaultLocked
		}
		return NewBWError(err, "bw "+strings.Join(args, " "))
	}
	// fmt.Printf("%s\n", b)
	if v == nil {
		return nil
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return NewBWError(err, "bw "+strings.Join(args, " "))
	}
	return nil
}

func (c *Client) ListItems() ([]*Item, error) {
	items := []*Item{}
	err := c.bw(&items, "list", "items")
	return items, err
}

type Folder struct {
	Object
}

func (c *Client) ListFolders() ([]*Folder, error) {
	folders := []*Folder{}
	err := c.bw(&folders, "list", "folders")
	return folders, err
}

func (c *Client) Unlock(password string) error {
	b, err := exec.Command("bw", "unlock", "--nointeraction", "--raw", password).CombinedOutput()
	if err != nil {
		return NewBWError(fmt.Errorf("%s: %w", firstLine(b), err), "bw unlock")
	}
	c.session = string(firstLine(b))
	err = os.WriteFile(c.sessionCacheFile, b, 0644)
	if err != nil {
		return NewBWError(err, "bw unlock")
	}
	return nil
}

func (c *Client) Sync() error {
	err := c.bw(nil, "sync")
	return err
}

func firstLine(s []byte) []byte {
	return bytes.SplitN(s, []byte("\n"), 2)[0]
}
