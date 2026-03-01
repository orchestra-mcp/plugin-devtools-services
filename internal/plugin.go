package internal

import (
	"github.com/orchestra-mcp/plugin-devtools-services/internal/tools"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// ToolsPlugin registers all service management tools.
type ToolsPlugin struct{}

// RegisterTools registers all 6 service management tools with the plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("list_services",
		"List OS services (launchctl on macOS, systemctl on Linux)",
		tools.ListServicesSchema(), tools.ListServices())

	builder.RegisterTool("start_service",
		"Start an OS service by name",
		tools.StartServiceSchema(), tools.StartService())

	builder.RegisterTool("stop_service",
		"Stop an OS service by name",
		tools.StopServiceSchema(), tools.StopService())

	builder.RegisterTool("restart_service",
		"Restart an OS service by name",
		tools.RestartServiceSchema(), tools.RestartService())

	builder.RegisterTool("service_logs",
		"Get logs for an OS service",
		tools.ServiceLogsSchema(), tools.ServiceLogs())

	builder.RegisterTool("service_info",
		"Get status and detailed info for an OS service",
		tools.ServiceInfoSchema(), tools.ServiceInfo())
}
