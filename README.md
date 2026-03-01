# Orchestra Plugin: devtools-services

A tools plugin for the [Orchestra MCP](https://github.com/orchestra-mcp/framework) framework.

## Install

```bash
go install github.com/orchestra-mcp/plugin-devtools-services/cmd@latest
```

## Usage

Add to your `plugins.yaml`:

```yaml
- id: tools.devtools-services
  binary: ./bin/devtools-services
  enabled: true
```

## Tools

| Tool | Description |
|------|-------------|
| `hello` | Say hello to someone |

## Related Packages

- [sdk-go](https://github.com/orchestra-mcp/sdk-go) — Plugin SDK
- [gen-go](https://github.com/orchestra-mcp/gen-go) — Generated Protobuf types
