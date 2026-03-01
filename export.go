package devtoolsservices

import (
	"github.com/orchestra-mcp/plugin-devtools-services/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all service management tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}
