package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"sync"
)

// Plugin paths.
var paths = []string{
	"bundle",
	"ftbundle",
	"pack/bundle/ever",
	"pack/bundle/opt",
}

// Repository object.
type Repository struct {
	path string
	vcs string
}

// Constructer.
func NewRepository(path, vcs string) *Repository {
	obj := new(Repository)
	obj.path = path
	obj.vcs = vcs

	return obj
}

var repos = []*Repository{}

// Run git pull.
func runGit(dir string) {
	os.Chdir(dir)
	fmt.Printf("git update %s\n", dir)
	cmd := exec.Command("git", "pull", "origin", "master")
	runCmd(cmd)
}

// Run hg update && hg pull.
func runHg(dir string) {
	os.Chdir(dir)
	fmt.Printf("hg pull && update %s\n", dir)
	cmdPull := exec.Command("hg", "pull")
	runCmd(cmdPull)
	cmdUpdate := exec.Command("hg", "update")
	runCmd(cmdUpdate)
}

// Execute command.
func runCmd(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command Run Error: %v\n", err)
		return
	}
}

// Get current user's $HOME.
func getVimHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(usr.HomeDir, ".vim")
}

// List all version controlled repositories.
func listRepositories(root, targetName string) []*Repository {
	targetPath := path.Join(root, targetName)
	dirs, err := ioutil.ReadDir(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dirs {
		if d.Name() == ".git" {
			r := NewRepository(targetPath, "git")
			repos = append(repos, r)
			continue
		} else if d.Name() == ".hg" {
			r := NewRepository(targetPath, "hg")
			repos = append(repos, r)
			continue
		} else if d.Name() == ".svn" {
			continue
		}

		if d.IsDir() {
			listRepositories(path.Join(root, targetName), d.Name())
		}
	}

	return repos
}

func run(v []*Repository) {
	var wg sync.WaitGroup
	for _, r := range v {
		wg.Add(1)
		go func(repo *Repository) {
			defer wg.Done()
			if repo.vcs == "git" {
				runGit(repo.path)
			} else if repo.vcs == "hg" {
				runHg(repo.path)
			}
		}(r)
	}
	wg.Wait()
}

func main() {
	home := getVimHome()
	repositories := map[string][]*Repository{}
	for _, v := range paths {
		if _, err := os.Stat(path.Join(home, v)); os.IsNotExist(err) {
			fmt.Printf("%s not exists.\n", v)
			continue
		}
		repos = make([]*Repository, 0)
		ret := listRepositories(home, v)
		repositories[v] = ret
	}

	for _, v := range repositories {
		run(v)
	}
}
