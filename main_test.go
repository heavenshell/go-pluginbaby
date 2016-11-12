package main

import (
	"os/user"
	"path"
	"reflect"
	"regexp"
	"testing"
)

func TestShouldGetVimHome(t *testing.T) {
	value := getVimHome()
	usr, _ := user.Current()
	expected := path.Join(usr.HomeDir, ".vim")
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
	t.Logf("%s", "test pass")
}

func TestShouldListRepositories(t *testing.T) {
	home := getVimHome()
	repositories := listRepositories(home, "pack")
	if len(repositories) == 0 {
		t.Fatalf("Fail")
	}
}

func TestListRepositoriesShouldContainsRepositoryObject(t *testing.T) {
	home := getVimHome()
	repositories := listRepositories(home, "pack")
	for i := range repositories {
		if reflect.TypeOf(repositories[i]).String() != "*main.Repository" {
			t.Fatalf("Fail")
		}
	}
}

func TestRepositoryObjectShouldContainsPath(t *testing.T) {
	home := getVimHome()
	repositories := listRepositories(home, "pack")
	for i := range repositories {
		if m, _ := regexp.MatchString("^/*[.vim]*[bundle|ftbundle]*", repositories[i].path); !m {
			t.Fatalf("Fail")
		}
	}
}
