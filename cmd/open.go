/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/pzolo85/git_vault/internal/enc"
	"github.com/pzolo85/git_vault/internal/tgz"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open a *.tgz.enc file and place the content inside a folder named 'open'",
	RunE:  openFile,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func openFile(cmd *cobra.Command, args []string) error {
	closedVault, err := tgz.GetLastTgzEncFile()
	if err != nil {
		return fmt.Errorf("unable to find closed vault: %w", err)
	}

	closedVaultContent, err := os.ReadFile(closedVault)
	if err != nil {
		return fmt.Errorf("unable to read closed vault: %w", err)
	}

	fmt.Println("Password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("unable to read password: %w", err)
	}

	decryptedVault, err := enc.Decrypt(string(closedVaultContent), string(password))
	if err != nil {
		return fmt.Errorf("unable to decrypt vault: %w", err)
	}

	if err := os.WriteFile("tmp.tgz", decryptedVault, 0666); err != nil {
		return fmt.Errorf("unable to write tmp encrypted file: %w", err)
	}

	if err := tgz.ExtractTarGz("tmp.tgz", "."); err != nil {
		return fmt.Errorf("unable to extract tmp encrypted file: %w", err)
	}

	// cleanup
	os.Remove(closedVault)
	os.Remove("tmp.tgz")

	return nil
}
