package main

import (
	_ "embed"
	"masbench/cmd"
)

//go:embed VERSION
var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
