package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	cssOutput      string
	typesOutput    string
	llmOutput      string
	cssFormat      string
	packageOutput  string
	packageScope   string
	packageName    string
	packageTargets string
	packageDryRun  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code artifacts from design system spec",
	Long: `Generate CSS, TypeScript types, LLM prompts, and NPM packages from a design system specification.

By default, outputs to stdout. Use flags to write to files.

Examples:
  # Preview all outputs
  dss generate

  # Generate specific files
  dss generate --css ./src/index.css --types ./src/lib/types.ts --llm ./DESIGN_CONTEXT.md

  # Generate from a different directory
  dss generate -d ./design-system --css ./web/src/index.css

  # Generate NPM package with all targets
  dss generate --package ./dist --targets all

  # Generate NPM package with specific targets
  dss generate --package ./dist --targets css,tailwind,shadcn`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVar(&cssOutput, "css", "", "output path for CSS (default: stdout)")
	generateCmd.Flags().StringVar(&typesOutput, "types", "", "output path for TypeScript types (default: stdout)")
	generateCmd.Flags().StringVar(&llmOutput, "llm", "", "output path for LLM prompt (default: stdout)")
	generateCmd.Flags().StringVar(&cssFormat, "css-format", "tailwind4", "CSS format: tailwind4, css-vars, scss, mkdocs-material")

	// NPM package generation flags
	generateCmd.Flags().StringVarP(&packageOutput, "package", "p", "", "output directory for NPM package")
	generateCmd.Flags().StringVarP(&packageScope, "scope", "s", "", "NPM scope (e.g., @myorg)")
	generateCmd.Flags().StringVarP(&packageName, "name", "n", "", "package name (default: design-tokens)")
	generateCmd.Flags().StringVarP(&packageTargets, "targets", "t", "css,tailwind", "comma-separated targets: css,tailwind,shadcn,mkdocs-material,scss,json,w3c,all")
	generateCmd.Flags().BoolVar(&packageDryRun, "dry-run", false, "preview package generation without writing files")

	rootCmd.AddCommand(generateCmd)
}

func runGenerate(_ *cobra.Command, _ []string) error {
	dir := getSpecDir()
	ctx := context.Background()

	// Load the design system and create service
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}
	service := dss.NewService(ds)

	// Validate
	if err := ds.Validate(); err != nil {
		return fmt.Errorf("design system validation failed: %w", err)
	}

	meta := service.GetMeta(ctx)
	fmt.Fprintf(os.Stderr, "Loaded design system: %s v%s\n", meta.Name, meta.Version)

	// If --package flag is set, generate NPM package
	if packageOutput != "" {
		return runPackageGenerate(ds)
	}

	// Generate CSS
	cssOpts := dss.DefaultCSSOptions()
	switch cssFormat {
	case "tailwind4":
		cssOpts.Format = "tailwind4"
	case "vars", "css-vars":
		cssOpts.Format = "css-vars"
	case "scss":
		cssOpts.Format = "scss"
	case "mkdocs-material", "mkdocs":
		cssOpts.Format = "mkdocs-material"
	default:
		return fmt.Errorf("unknown CSS format: %s (valid: tailwind4, css-vars, scss, mkdocs-material)", cssFormat)
	}

	css, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return fmt.Errorf("generating CSS: %w", err)
	}

	if cssOutput != "" {
		if err := os.WriteFile(cssOutput, []byte(css), 0600); err != nil {
			return fmt.Errorf("writing CSS: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Generated CSS: %s\n", cssOutput)
	} else if typesOutput == "" && llmOutput == "" {
		// Only print to stdout if no files specified
		fmt.Println("=== CSS ===")
		fmt.Println(css)
	}

	// Generate LLM prompt via service
	llmOpts := &dss.PromptOptions{
		Format:               "markdown",
		IncludeFoundations:   true,
		IncludeComponents:    true,
		IncludePatterns:      true,
		IncludeAccessibility: true,
		IncludeAntiPatterns:  true,
		MaxExamples:          3,
	}
	llmPrompt, err := service.GenerateLLMPrompt(ctx, llmOpts)
	if err != nil {
		return fmt.Errorf("generating LLM prompt: %w", err)
	}

	if llmOutput != "" {
		if err := os.WriteFile(llmOutput, []byte(llmPrompt), 0600); err != nil {
			return fmt.Errorf("writing LLM prompt: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Generated LLM prompt: %s\n", llmOutput)
	} else if cssOutput == "" && typesOutput == "" {
		fmt.Println("\n=== LLM Prompt ===")
		fmt.Println(llmPrompt)
	}

	// Generate TypeScript types
	components := service.ListComponents(ctx)
	if len(components) > 0 {
		reactOpts := dss.DefaultReactOptions()
		types, err := ds.GenerateReactTypes(reactOpts)
		if err != nil {
			return fmt.Errorf("generating TypeScript types: %w", err)
		}

		if typesOutput != "" {
			if err := os.WriteFile(typesOutput, []byte(types), 0600); err != nil {
				return fmt.Errorf("writing TypeScript types: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Generated TypeScript types: %s\n", typesOutput)
		} else if cssOutput == "" && llmOutput == "" {
			fmt.Println("\n=== TypeScript Types ===")
			fmt.Println(types)
		}
	}

	fmt.Fprintf(os.Stderr, "Done!\n")
	return nil
}

// runPackageGenerate handles NPM package generation.
func runPackageGenerate(ds *dss.DesignSystem) error {
	opts := dss.DefaultPackageOptions()
	opts.OutputDir = packageOutput
	opts.DryRun = packageDryRun

	if packageScope != "" {
		opts.Scope = packageScope
	}

	if packageName != "" {
		opts.PackageName = packageName
	}

	opts.Targets = dss.ParseTargets(packageTargets)

	if packageDryRun {
		fmt.Fprintf(os.Stderr, "Dry run: would generate package to %s\n", packageOutput)
		fmt.Fprintf(os.Stderr, "Targets: %v\n", opts.Targets)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Generating NPM package to %s\n", packageOutput)
	fmt.Fprintf(os.Stderr, "Targets: %v\n", opts.Targets)

	if err := ds.GeneratePackage(opts); err != nil {
		return fmt.Errorf("generating package: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Package generated successfully!\n")
	fmt.Fprintf(os.Stderr, "\nTo publish:\n")
	fmt.Fprintf(os.Stderr, "  cd %s && npm publish\n", packageOutput)

	return nil
}
