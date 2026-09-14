// Package server provides MCP server runtime for skills.
// This is a local stub that will be replaced with github.com/plexusone/omniskill/mcp/server
// when that package becomes available.
package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
)

// Runtime manages MCP server operations and skill registration.
type Runtime struct {
	impl   *mcp.Implementation
	skills []skill.Skill
	server *mcp.Server
}

// New creates a new MCP server runtime.
func New(impl *mcp.Implementation, _ any) *Runtime {
	return &Runtime{
		impl:   impl,
		skills: []skill.Skill{},
	}
}

// RegisterSkill adds a skill to the runtime.
func (r *Runtime) RegisterSkill(s skill.Skill) {
	r.skills = append(r.skills, s)
}

// ServeStdio starts the MCP server over stdio.
func (r *Runtime) ServeStdio(ctx context.Context) error {
	// Create MCP server
	r.server = mcp.NewServer(r.impl, nil)

	// Register all tools from all skills
	for _, s := range r.skills {
		for _, t := range s.Tools() {
			r.registerTool(t)
		}
	}

	// Create stdio transport and run
	transport := &mcp.StdioTransport{}
	return r.server.Run(ctx, transport)
}

// registerTool registers a single tool with the MCP server.
func (r *Runtime) registerTool(t skill.Tool) {
	// Build JSON Schema for parameters as map[string]any
	schema := buildInputSchema(t.Parameters())

	// Create MCP tool definition
	mcpTool := &mcp.Tool{
		Name:        t.Name(),
		Description: t.Description(),
		InputSchema: schema,
	}

	// Register with raw handler
	r.server.AddTool(mcpTool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse arguments from json.RawMessage
		params := make(map[string]any)
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &params); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
					},
					IsError: true,
				}, nil
			}
		}

		// Execute tool
		result, err := t.Execute(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
				},
				IsError: true,
			}, nil
		}

		// Serialize result to JSON
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling result: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(resultJSON)},
			},
		}, nil
	})
}

// buildInputSchema creates a JSON Schema as map[string]any from skill parameters.
func buildInputSchema(params skill.Parameters) map[string]any {
	properties := make(map[string]any)
	var required []string

	for name, param := range params {
		prop := map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}
		if len(param.Enum) > 0 {
			prop["enum"] = param.Enum
		}
		if param.Default != nil {
			prop["default"] = param.Default
		}
		properties[name] = prop

		if param.Required {
			required = append(required, name)
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}
