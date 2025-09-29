/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jamlie/colors"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .ask.toml",
	Long:  "Initialize the configuration file (.ask.toml) in the user directory",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		configPath := filepath.Join(homeDir, ".ask.toml")
		redColor := colors.New(colors.RedFg, colors.WithBold)

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			content := `gemini_api = ""`
			if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
				fmt.Fprintln(os.Stderr, redColor.String(fmt.Sprintf("failed to create config: %v", err)))
				return
			}
			greenColor := colors.New(colors.GreenFg, colors.WithBold)
			fmt.Println(greenColor.String(fmt.Sprintf("Created config at %s", configPath)))
		} else if err != nil {
			fmt.Fprintln(os.Stderr, redColor.String(fmt.Sprintf("error checking config: %v", err)))
			return
		} else {
			fmt.Fprintln(os.Stderr, redColor.String(fmt.Sprintf("Config already exists at %s", configPath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
