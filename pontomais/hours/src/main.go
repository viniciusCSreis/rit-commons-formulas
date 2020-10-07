// This is the main class.
// Where you will extract the inputs asked on the config.json file and call the formula's method(s).

package main

import (
	"formula/pkg/formula"
	"os"
)

func main() {
	username := os.Getenv("PONTOMAIS_LOGIN")
	password := os.Getenv("PONTOMAIS_PASSWORD")

	formula.Formula{
		Username:     username,
		Password:     password,
	}.Run()
}
