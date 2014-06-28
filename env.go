package martini

import (
	"os"
	"path/filepath"
	"strings"
)

// Envs
const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

// Env is the environment that Martini is executing in. The MARTINI_ENV is read on initialization to set this variable.
var Env = Dev
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("MARTINI_ENV"))

	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	// In the dev mode when used commands `go test`, `go run` determining Root
	// by a first argument is incorrect, because binary is putted into a temp directory.
	// In these cases used a current directory as martini.Root
	if strings.Contains(path, "go-build") && (strings.Contains(path, "command-line-arguments") || strings.Contains(path, "_test")) {
		Root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		Root = filepath.Dir(path)
	}
}
