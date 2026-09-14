package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	bindOutput   string
	bindFormat   string
	bindStrategy string
)

var bindCmd = &cobra.Command{
	Use:   "bind",
	Short: "Generate theme bindings from design system",
	Long: `Generate CSS, TypeScript, or SCSS theme bindings from design system tokens.

Uses the themeBindings configuration in the design system to map application
tokens to component theming contracts.

Output formats:
  css        CSS custom properties (default)
  typescript TypeScript constants with type safety
  scss       SCSS variables plus CSS custom properties

Mapping strategies:
  explicit   Only use defined mappings, skip unmapped tokens (default)
  semantic   Auto-map by semantic field, fall back to defaults
  inherit    Use component defaults for all unmapped tokens

Examples:
  dss bind --output ./theme.css
  dss bind --format typescript --output ./theme.ts
  dss bind --strategy semantic --output ./theme.css
  dss bind -d ./design-system -o ./bindings.css`,
	RunE: runBind,
}

func init() {
	rootCmd.AddCommand(bindCmd)

	bindCmd.Flags().StringVarP(&bindOutput, "output", "o", "", "output file (default: stdout)")
	bindCmd.Flags().StringVarP(&bindFormat, "format", "f", "css", "output format: css, typescript, scss")
	bindCmd.Flags().StringVar(&bindStrategy, "strategy", "explicit", "mapping strategy: explicit, semantic, inherit")
}

func runBind(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	if len(ds.ThemeBindings) == 0 {
		fmt.Fprintln(os.Stderr, "No theme bindings defined in design system.")
		fmt.Fprintln(os.Stderr, "Add a 'themeBindings' section to your design-system.json or create a themeBindings.json file.")
		return nil
	}

	// Validate format
	var format dss.BindingFormat
	switch bindFormat {
	case "css":
		format = dss.FormatCSS
	case "typescript", "ts":
		format = dss.FormatTypeScript
	case "scss":
		format = dss.FormatSCSS
	default:
		return fmt.Errorf("unsupported format: %s (use: css, typescript, scss)", bindFormat)
	}

	// Validate strategy
	switch bindStrategy {
	case "explicit", "semantic", "inherit":
		// Valid
	default:
		return fmt.Errorf("unsupported strategy: %s (use: explicit, semantic, inherit)", bindStrategy)
	}

	// Generate bindings
	opts := dss.BindingOptions{
		Format:          format,
		SpecDir:         dir,
		DefaultStrategy: bindStrategy,
	}

	bindings, err := dss.GenerateBindings(ds, opts)
	if err != nil {
		return fmt.Errorf("generating bindings: %w", err)
	}

	// Print warnings
	for _, b := range bindings {
		for _, w := range b.Warnings {
			fmt.Fprintf(os.Stderr, "warning: %s: %s\n", b.Component, w)
		}
	}

	// Output
	var out *os.File
	if bindOutput == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(bindOutput)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer out.Close()
	}

	if err := dss.WriteBindings(out, bindings); err != nil {
		return fmt.Errorf("writing bindings: %w", err)
	}

	if bindOutput != "" {
		fmt.Fprintf(os.Stderr, "Generated bindings written to %s\n", bindOutput)
	}

	return nil
}
