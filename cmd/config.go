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

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interact with the whoiam config file",
	Long:  ``,
}

func addEntrypoint(cmd *cobra.Command, args []string) {
	accountName, _ := cmd.Flags().GetString("name")
	accountNumber, _ := cmd.Flags().GetString("account")

	if accountName == "" || accountNumber == "" {
		fmt.Println("--name and --account are required")
		os.Exit(1)
	}

	cfgPath, err := internal.NewConfigPath()
	internal.HandleError(err)

	cfg, err := cfgPath.LoadConfig()
	internal.HandleError(err)

	if cfg.AccountExists(accountName) {
		fmt.Printf("Account %s already exists in the config\n", accountName)
		os.Exit(1)
	}

	account := cfg.GetAccountByNumber(accountNumber)
	if account != "" {
		fmt.Printf("Account Number %s already exists in the config as %s\n", accountNumber, account)
		os.Exit(1)
	}

	err = cfg.AddAccount(accountName, accountNumber)
	internal.HandleError(err)

	err = cfgPath.SaveConfig(cfg)
	internal.HandleError(err)

	fmt.Printf("Added -- Account Number: %s Account Name: %s\n", accountNumber, accountName)

	os.Exit(0)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add to the whoiam config file",
	Run:   addEntrypoint,
}

func viewEntrypoint(cmd *cobra.Command, args []string) {
	cfg, err := internal.LoadEffectiveConfig()
	internal.HandleError(err)

	cfg.PrintConfigTable()

	os.Exit(0)
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the effective whoiam config (global + project-local merged)",
	Run:   viewEntrypoint,
}

func deleteEntrypoint(cmd *cobra.Command, args []string) {
	accountName, _ := cmd.Flags().GetString("name")

	if accountName == "" {
		fmt.Println("--name is required")
		os.Exit(1)
	}

	cfgPath, err := internal.NewConfigPath()
	internal.HandleError(err)

	cfg, err := cfgPath.LoadConfig()
	internal.HandleError(err)

	if !cfg.AccountExists(accountName) {
		fmt.Printf("Account %s does not exist in config\n", accountName)
		os.Exit(1)
	}

	accountNumber := cfg.Accounts[accountName]
	cfg.DeleteAccount(accountName)

	err = cfgPath.SaveConfig(cfg)
	internal.HandleError(err)

	fmt.Printf("Deleted -- Account Name: %s Number: %s\n", accountName, accountNumber)

	os.Exit(0)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete from the whoiam config file",
	Run:   deleteEntrypoint,
}

func initEntrypoint(cmd *cobra.Command, args []string) {
	cfgPath, err := internal.NewConfigPath()
	internal.HandleError(err)

	if cfgPath.Exists() {
		fmt.Println("Config file already exists")
		os.Exit(1)
	}

	err = cfgPath.Create()
	internal.HandleError(err)

	cfg, err := internal.NewTemplateConfig()
	internal.HandleError(err)

	err = cfgPath.SaveConfig(cfg)
	internal.HandleError(err)

	os.Exit(0)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new whoiam config file",
	Run:   initEntrypoint,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(viewCmd)

	configCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("name", "n", "", "AWS Account Name")
	addCmd.Flags().StringP("account", "a", "", "AWS Account Number")

	configCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("name", "n", "", "AWS Account Name")

}
