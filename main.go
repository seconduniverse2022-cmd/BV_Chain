package main

import (
	_ "embed"

	"github.com/bitvault/command/root"
	"github.com/bitvault/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
