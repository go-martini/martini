package martini

import (
	"os"
	"path/filepath"
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
	if Env == Dev || Env == Test {
		var err error
		Root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		path, err := filepath.Abs(os.Args[0])
		if err != nil {
			panic(err)
		}
		Root = filepath.Dir(path)
	}
}
