package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-services/internal/svc"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// ServiceLogsSchema returns the JSON Schema for the service_logs tool.
func ServiceLogsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Name of the service to get logs for",
			},
			"lines": map[string]any{
				"type":        "number",
				"description": "Number of log lines to retrieve (default: 50)",
			},
		},
		"required": []any{"name"},
	})
	return s
}

// ServiceLogs returns a tool handler that retrieves service logs.
func ServiceLogs() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "name"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		name := helpers.GetString(req.Arguments, "name")
		lines := helpers.GetInt(req.Arguments, "lines")
		if lines <= 0 {
			lines = 50
		}

		var output string
		var err error

		if svc.IsDarwin() {
			// On macOS, use log show to get service logs filtered by subsystem/process.
			output, err = svc.RunCommand(ctx, "log", "show",
				"--predicate", fmt.Sprintf("process == %q", name),
				"--last", fmt.Sprintf("%dm", lines),
				"--style", "compact",
			)
		} else if svc.IsLinux() {
			output, err = svc.RunSystemctl(ctx, "status", name)
			if err != nil {
				// Try journalctl as fallback for more detailed logs.
				output, err = svc.RunCommand(ctx, "journalctl",
					"-u", name,
					"-n", fmt.Sprintf("%d", lines),
					"--no-pager",
				)
			}
		} else {
			return helpers.ErrorResult("unsupported_os", "Only macOS and Linux are supported"), nil
		}

		if err != nil {
			return helpers.ErrorResult("exec_error", err.Error()), nil
		}

		return helpers.TextResult(output), nil
	}
}
