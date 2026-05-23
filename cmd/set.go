/*
Copyright © 2024 Jesse Maitland jesse@cliwright.com
*/
package cmd

import (
	"fmt"
	"github.com/cliwright/whoiam/internal"
	"github.com/spf13/cobra"
)

func setEntrypoint(cmd *cobra.Command, args []string) error {
	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if global {
			if err := internal.ClearGlobalCurrentEnv(); err != nil {
				return err
			}
			cmd.Println("Cleared global expected account")
		} else {
			if err := internal.ClearCurrentEnv(); err != nil {
				return err
			}
			cmd.Println("Cleared expected account")
		}
		return nil
	}

	accountName := args[0]

	cfg, err := internal.LoadEffectiveConfig()
	if err != nil {
		return err
	}

	if !cfg.AccountExists(accountName) {
		return fmt.Errorf("account %q does not exist in config", accountName)
	}

	if global {
		if err := internal.WriteGlobalCurrentEnv(accountName); err != nil {
			return err
		}
		cmd.Printf("Global expected account set to %q\n", accountName)
	} else {
		if err := internal.WriteCurrentEnv(accountName); err != nil {
			return err
		}
		cmd.Printf("Expected account set to %q\n", accountName)
	}
	return nil
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
	RunE: setEntrypoint,
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("global", false, "Write to ~/.whoiam/expected-env instead of the project-local file")
}
