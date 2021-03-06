package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/vcs"
)

type Prefetcher interface {
	fetchHash(url string, revision string) (string, error)
}

func PrefetcherFor(typ vcs.Type) Prefetcher {
	switch typ {
	case vcs.Git:
		return &gitPrefetcher{}
	case vcs.Hg:
		return &hgPrefetcher{}
	default:
		return nil
	}
}

func cmdStdout(command string, arguments ...string) (string, error) {
	cmd := exec.Command(command, arguments...)
	// var out bytes.Buffer
	// cmd.Stdout = &out
	var out []byte

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Command output", string(out[:]))
		return "", err
	}

	return string(out[:]), nil
}

type gitPrefetcher struct{}

func (p *gitPrefetcher) fetchHash(url string, revision string) (string, error) {
	out, err := cmdStdout("nix-prefetch-git", "--url", url, "--rev", revision)
	fmt.Println("Command ", "nix-prefetch-git", "--url", url, "--rev", revision)
	if err != nil {
		return "", err
	}

	// extract hash from response
	res := &struct {
		SHA256 string `json:"sha256"`
	}{}
	if err := json.Unmarshal([]byte(out), res); err != nil {
		return "", err
	}

	return res.SHA256, nil
}

type hgPrefetcher struct{}

func (p *hgPrefetcher) fetchHash(url string, revision string) (string, error) {
	out, err := cmdStdout("nix-prefetch-hg", url, revision)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
