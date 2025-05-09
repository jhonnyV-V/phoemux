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
	"fmt"

	"github.com/jhonnyV-V/phoemux/tmux"
	"github.com/spf13/cobra"
)

var (
	target      string
	attach      string
	dumb_attach bool
)

// killCmd represents the list command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill current tmux session",
	Long: `kill command
Kills current tmux session using tmux kill-session:
phoemux kill`,
	Example: "phoemux run kill -t react-app -a server-app",
	Run: func(cmd *cobra.Command, args []string) {
		tmuxEnvExist := tmux.IsInsideTmux()
		if !tmuxEnvExist && target == "" {
			fmt.Printf("You are not in a tmux session\n")
			return
		}

		if target == "" {
			target = tmux.GetCurrentSessionName()
		}

		if dumb_attach {
			ash := tmux.Ash{}
			if tmuxEnvExist {
				ash.SessionName = tmux.GetOtherSession()
				if ash.SessionName == "" {
					fmt.Printf("can't find other tmux session\n")
					return
				}
			}
			tmux.ChangeSession(ash)
		}

		if attach != "" {
			tmux.ChangeSession(tmux.Ash{
				SessionName: attach,
			})
		}

		tmux.Kill(target)
	},
}

func init() {
	killCmd.Flags().StringVarP(&target, "target", "t", "", "target session name")
	killCmd.Flags().StringVarP(&attach, "attach", "a", "", "attach to session name")
	killCmd.Flags().BoolVarP(&dumb_attach, "dumb-attach", "d", false, "run attach without arguments")
	killCmd.MarkFlagsMutuallyExclusive("attach", "dumb-attach")
	killCmd.RegisterFlagCompletionFunc("target", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return tmux.GetListOfSessions(), cobra.ShellCompDirectiveNoFileComp
	})
	killCmd.RegisterFlagCompletionFunc("attach", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return tmux.GetOthersSessions(), cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.AddCommand(killCmd)
}
