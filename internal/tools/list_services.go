package tools

import (
	"context"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-services/internal/svc"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// ListServicesSchema returns the JSON Schema for the list_services tool.
func ListServicesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"filter": map[string]any{
				"type":        "string",
				"description": "Optional filter string to match against service names",
			},
		},
	})
	return s
}

// ListServices returns a tool handler that lists OS services.
func ListServices() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		filter := helpers.GetString(req.Arguments, "filter")

		var output string
		var err error

		if svc.IsDarwin() {
			output, err = svc.RunLaunchctl(ctx, "list")
		} else if svc.IsLinux() {
			output, err = svc.RunSystemctl(ctx, "list-units", "--type=service", "--no-pager")
		} else {
			return helpers.ErrorResult("unsupported_os", "Only macOS (launchctl) and Linux (systemctl) are supported"), nil
		}

		if err != nil {
			return helpers.ErrorResult("exec_error", err.Error()), nil
		}

		if filter != "" {
			lines := strings.Split(output, "\n")
			var filtered []string
			for _, line := range lines {
				if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
					filtered = append(filtered, line)
				}
			}
			// Keep header line if present
			if len(lines) > 0 && len(filtered) > 0 {
				header := lines[0]
				if !strings.Contains(strings.ToLower(header), strings.ToLower(filter)) {
					filtered = append([]string{header}, filtered...)
				}
			}
			output = strings.Join(filtered, "\n")
		}

		return helpers.TextResult(output), nil
	}
}
