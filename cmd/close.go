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

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "close your vault",
	RunE:  closeFile,
}

func init() {
	rootCmd.AddCommand(closeCmd)
}

func closeFile(cmd *cobra.Command, args []string) error {
	closedVaultName, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("unable to generate closed vault name: %w", err)
	}

	gzipVault := fmt.Sprintf("%s.tgz", closedVaultName.String())
	if err := tgz.CreateTarGz("./open", gzipVault); err != nil {
		return fmt.Errorf("unable to create tgz file: %w", err)
	}

	gzipVaultContent, err := os.ReadFile(gzipVault)
	if err != nil {
		return fmt.Errorf("unabe to read gizip vault: %w", err)
	}

	fmt.Println("Password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("unable to read password: %w", err)
	}

	encVaultContent, err := enc.Encrypt(gzipVaultContent, password)
	if err != nil {
		return fmt.Errorf("unable to encrypt data: %w", err)
	}

	encGzipVault := fmt.Sprintf("%s.enc", gzipVault)
	if err := os.WriteFile(encGzipVault, []byte(encVaultContent), 0666); err != nil {
		return fmt.Errorf("unable to write encrypted data: %w", err)
	}

	// clean-up
	os.Remove(gzipVault)
	os.RemoveAll("./open")

	return nil
}
