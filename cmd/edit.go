/*
Copyright © 2024 Jhonny Varela jhonny_varela_visbal@hotmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/jhonnyV-V/phoemux/core"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit existing ash",
	Long: `edit command.
opens existing Ash with the editor set in $EDITOR or nano as a fallback:

phoemux edit <project_name>`,
	Args:    cobra.MinimumNArgs(1),
	Example: "phoemux edit <project_name>",
	Run: func(cmd *cobra.Command, args []string) {
		phoemuxConfigPath := core.GetConfigPath()
		core.Edit(phoemuxConfigPath, args[0])
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		phoemuxConfigPath := core.GetConfigPath()

		ashes, err := core.GetSimpleList(phoemuxConfigPath)

		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return ashes, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
