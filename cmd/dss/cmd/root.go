package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set by goreleaser
	Version = "dev"

	// specDir is the directory containing the design system spec
	specDir string
)

var rootCmd = &cobra.Command{
	Use:   "dss",
	Short: "systemspec-designsystem CLI",
	Long: `dss is a CLI tool for working with Design System Specifications.

It can generate code artifacts (CSS, TypeScript types, LLM prompts) from
a declarative design system specification, and validate that implementations
comply with the spec.

Example usage:
  dss generate --css ./src/index.css --types ./src/lib/types.ts
  dss validate ./src/components
  dss info`,
	Version: Version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&specDir, "dir", "d", ".", "directory containing the design system spec")
}

// getSpecDir returns the spec directory, defaulting to current directory
func getSpecDir() string {
	if specDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}
		return dir
	}
	return specDir
}
