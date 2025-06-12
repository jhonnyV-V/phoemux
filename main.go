/*
Copyright © 2024 Jhonny Varela jhonny_varela_visbal@hotmail.com
*/
package main

import (
	"github.com/jhonnyV-V/phoemux/cmd"
	"github.com/jhonnyV-V/phoemux/version"
)

func main() {
	cmd.SetVersion(version.Version)
	cmd.Execute()
}
