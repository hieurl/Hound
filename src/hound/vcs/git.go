package vcs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	RegisterVCS("git", &GitDriver{})
}

type GitDriver struct{}

func (g *GitDriver) HeadHash(dir string, branch string) (string, error) {
	cmd := exec.Command(
		"git",
		"rev-parse",
		branch)
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), cmd.Wait()
}

func (g *GitDriver) Pull(dir string, branch string) (string, error) {
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return g.HeadHash(dir, branch)
}

func (g *GitDriver) Clone(dir, url string, branch string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"git",
		"clone",
		"-b",
		branch,
		url,
		rep)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return g.HeadHash(dir, branch)
}
