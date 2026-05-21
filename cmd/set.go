/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
	"os"
)

func setEntrypoint(cmd *cobra.Command, args []string) {
	global, err := cmd.Flags().GetBool("global")
	internal.HandleError(err)

	if len(args) == 0 {
		if global {
			err := internal.ClearGlobalCurrentEnv()
			internal.HandleError(err)
			fmt.Println("Cleared global expected account")
		} else {
			err := internal.ClearCurrentEnv()
			internal.HandleError(err)
			fmt.Println("Cleared expected account")
		}
		return
	}

	accountName := args[0]

	cfg, err := internal.LoadEffectiveConfig()
	internal.HandleError(err)

	if !cfg.AccountExists(accountName) {
		fmt.Printf("Account %q does not exist in config\n", accountName)
		os.Exit(1)
	}

	if global {
		err = internal.WriteGlobalCurrentEnv(accountName)
		internal.HandleError(err)
		fmt.Printf("Global expected account set to %q\n", accountName)
	} else {
		err = internal.WriteCurrentEnv(accountName)
		internal.HandleError(err)
		fmt.Printf("Expected account set to %q\n", accountName)
	}
}

var setCmd = &cobra.Command{
	Use:   "set [env]",
	Short: "Set (or clear) the expected environment for this session",
	Long: `Writes the expected environment name to .whoiam/expected-env so that
subsequent 'whoiam exec' and 'whoiam validate' calls do not need an explicit --env flag.

Pass --global to write to ~/.whoiam/expected-env instead (applies across all projects).

Running 'whoiam set' with no argument clears the expectation.

Examples:
  whoiam set production           # set local expected environment
  whoiam set --global production  # set global expected environment
  whoiam set                      # clear local expected environment
  whoiam set --global             # clear global expected environment`,
	Args: cobra.MaximumNArgs(1),
	Run:  setEntrypoint,
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("global", false, "Write to ~/.whoiam/expected-env instead of the project-local file")
}
