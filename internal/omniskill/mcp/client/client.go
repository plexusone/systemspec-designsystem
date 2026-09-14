// Package client provides MCP client for connecting to other MCP servers.
// This is a local stub that will be replaced with github.com/plexusone/omniskill/mcp/client
// when that package becomes available.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
)

// Client manages connections to remote MCP servers.
type Client struct {
	name    string
	version string
}

// New creates a new MCP client.
func New(name, version string, _ any) *Client {
	return &Client{
		name:    name,
		version: version,
	}
}

// Session represents an active connection to an MCP server.
type Session struct {
	client  *mcp.Client
	session *mcp.ClientSession
	tools   []*mcp.Tool
}

// ConnectCommand connects to an MCP server via subprocess.
func (c *Client) ConnectCommand(ctx context.Context, cmd *exec.Cmd, _ any) (*Session, error) {
	// Create MCP client
	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    c.name,
		Version: c.version,
	}, nil)

	// Create command transport for subprocess
	transport := &mcp.CommandTransport{Command: cmd}

	// Connect
	session, err := mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to server: %w", err)
	}

	// List tools
	toolsResp, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("listing tools: %w", err)
	}

	return &Session{
		client:  mcpClient,
		session: session,
		tools:   toolsResp.Tools,
	}, nil
}

// Close closes the session.
func (s *Session) Close() error {
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}

// SkillOption configures skill creation.
type SkillOption func(*proxySkillConfig)

type proxySkillConfig struct {
	name        string
	description string
}

// WithSkillName sets the skill name.
func WithSkillName(name string) SkillOption {
	return func(c *proxySkillConfig) {
		c.name = name
	}
}

// WithSkillDescription sets the skill description.
func WithSkillDescription(desc string) SkillOption {
	return func(c *proxySkillConfig) {
		c.description = desc
	}
}

// AsSkill wraps all discovered tools as a skill.
func (s *Session) AsSkill(opts ...SkillOption) skill.Skill {
	cfg := &proxySkillConfig{
		name:        "remote",
		description: "Remote MCP tools",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return &proxySkill{
		session: s,
		config:  cfg,
	}
}

// proxySkill wraps remote MCP tools as a skill.
type proxySkill struct {
	session *Session
	config  *proxySkillConfig
}

func (p *proxySkill) Name() string {
	return p.config.name
}

func (p *proxySkill) Description() string {
	return p.config.description
}

func (p *proxySkill) Init(_ context.Context) error {
	return nil
}

func (p *proxySkill) Close() error {
	return nil
}

func (p *proxySkill) Tools() []skill.Tool {
	tools := make([]skill.Tool, len(p.session.tools))
	for i, t := range p.session.tools {
		tools[i] = &proxyTool{
			session: p.session,
			mcpTool: t,
		}
	}
	return tools
}

// proxyTool wraps a remote MCP tool.
type proxyTool struct {
	session *Session
	mcpTool *mcp.Tool
}

func (t *proxyTool) Name() string {
	return t.mcpTool.Name
}

func (t *proxyTool) Description() string {
	return t.mcpTool.Description
}

func (t *proxyTool) Parameters() skill.Parameters {
	params := make(skill.Parameters)

	// InputSchema is `any`, typically map[string]any from JSON unmarshaling
	schemaMap, ok := t.mcpTool.InputSchema.(map[string]any)
	if !ok {
		return params
	}

	properties, ok := schemaMap["properties"].(map[string]any)
	if !ok {
		return params
	}

	// Get required array
	requiredSet := make(map[string]bool)
	if required, ok := schemaMap["required"].([]any); ok {
		for _, r := range required {
			if s, ok := r.(string); ok {
				requiredSet[s] = true
			}
		}
	}

	for name, propAny := range properties {
		prop, ok := propAny.(map[string]any)
		if !ok {
			continue
		}

		param := skill.Parameter{
			Required: requiredSet[name],
		}
		if typ, ok := prop["type"].(string); ok {
			param.Type = typ
		}
		if desc, ok := prop["description"].(string); ok {
			param.Description = desc
		}

		params[name] = param
	}

	return params
}

func (t *proxyTool) Execute(ctx context.Context, params map[string]any) (any, error) {
	// Call remote tool
	result, err := t.session.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      t.mcpTool.Name,
		Arguments: params,
	})
	if err != nil {
		return nil, fmt.Errorf("calling tool: %w", err)
	}

	// Extract text content
	for _, content := range result.Content {
		if tc, ok := content.(*mcp.TextContent); ok {
			// Try to parse as JSON
			var data any
			if err := json.Unmarshal([]byte(tc.Text), &data); err == nil {
				return data, nil
			}
			return tc.Text, nil
		}
	}

	return result.Content, nil
}
