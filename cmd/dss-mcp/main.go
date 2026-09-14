// Package main provides the MCP server for design system spec.
// It exposes design system operations as MCP tools, optionally including
// w3pilot browser tools for visual validation.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plexusone/systemspec-designsystem/internal/omniskill/mcp/client"
	"github.com/plexusone/systemspec-designsystem/internal/omniskill/mcp/server"
	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
	"github.com/plexusone/systemspec-designsystem/skills/designsystem"
	"github.com/spf13/cobra"
)

var (
	serverName    = "dss-mcp"
	serverVersion = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "dss-mcp",
	Short: "systemspec-designsystem MCP Server",
	Long: `MCP server for design system specification operations.

Exposes tools for:
  - Reading component specs, tokens, and patterns
  - Generating LLM context prompts
  - Validating implementations against the spec

Examples:
  # Start server with a spec
  dss-mcp --spec ./design-system/

  # Start with browser validation tools
  dss-mcp --spec ./design-system/ --browser

  # Use with Claude Desktop (add to config)
  # {
  #   "mcpServers": {
  #     "design-system": {
  #       "command": "dss-mcp",
  #       "args": ["--spec", "/path/to/spec"]
  #     }
  #   }
  # }
`,
	RunE: runServer,
}

var (
	specPath      string
	enableBrowser bool
)

func init() {
	rootCmd.Flags().StringVarP(&specPath, "spec", "s", "", "Path to design system spec (directory or file)")
	rootCmd.Flags().BoolVar(&enableBrowser, "browser", false, "Enable w3pilot browser tools for visual validation")
	_ = rootCmd.MarkFlagRequired("spec")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// 1. Load design system
	ds, err := dss.LoadDesignSystem(specPath)
	if err != nil {
		return fmt.Errorf("loading design system from %s: %w", specPath, err)
	}

	// 2. Create service and skill
	service := dss.NewService(ds)
	dsSkill := designsystem.New(service)

	if err := dsSkill.Init(ctx); err != nil {
		return fmt.Errorf("initializing designsystem skill: %w", err)
	}
	defer dsSkill.Close()

	// 3. Create MCP server runtime
	rt := server.New(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// 4. Register designsystem skill
	rt.RegisterSkill(dsSkill)

	// 5. Optionally add w3pilot skill (all 169 tools)
	if enableBrowser {
		w3pilotSkill, cleanup, err := connectW3Pilot(ctx)
		if err != nil {
			return fmt.Errorf("connecting to w3pilot: %w", err)
		}
		defer cleanup()
		rt.RegisterSkill(w3pilotSkill)
	}

	// 6. Serve via stdio
	return rt.ServeStdio(ctx)
}

// connectW3Pilot starts w3pilot MCP server and wraps all its tools as a skill.
// This enables AI agents to use browser automation for visual validation.
func connectW3Pilot(ctx context.Context) (skill.Skill, func(), error) {
	// Create MCP client
	c := client.New("dss-mcp", serverVersion, nil)

	// Start w3pilot as subprocess running its MCP server
	w3pilotCmd := exec.Command("w3pilot", "mcp", "serve")

	session, err := c.ConnectCommand(ctx, w3pilotCmd, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to w3pilot: %w", err)
	}

	// Auto-wrap all w3pilot tools as a skill
	// This discovers all 169 tools via MCP tools/list and creates proxy tools
	w3pilotSkill := session.AsSkill(
		client.WithSkillName("w3pilot"),
		client.WithSkillDescription("Browser automation for testing UI implementations against the design system"),
	)

	cleanup := func() {
		session.Close()
	}

	return w3pilotSkill, cleanup, nil
}
