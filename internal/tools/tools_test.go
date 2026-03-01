package tools

import (
	"context"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// callTool invokes a tool handler with the given argument map and returns the response.
// The test fails immediately if the handler returns a Go-level error.
func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var argStruct *structpb.Struct
	if args != nil {
		s, err := structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("structpb.NewStruct: %v", err)
		}
		argStruct = s
	}
	req := &pluginv1.ToolRequest{Arguments: argStruct}
	resp, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

// isError returns true when the response carries a failure (Success == false).
func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

// errorCode extracts the error_code field from a response.
func errorCode(resp *pluginv1.ToolResponse) string {
	return resp.GetErrorCode()
}

// --- list_services ---

// TestListServices_NoArgs calls list_services with no arguments. It must return
// a non-nil response without a Go-level error. The tool runs the real OS command
// so either a success or an OS-level error response is acceptable.
func TestListServices_NoArgs(t *testing.T) {
	handler := ListServices()
	resp := callTool(t, handler, nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	// Accept both success and any error code — the real command is invoked.
}

// --- start_service ---

// TestStartService_MissingName verifies that omitting the required "name" field
// returns a validation_error.
func TestStartService_MissingName(t *testing.T) {
	handler := StartService()
	resp := callTool(t, handler, nil)
	if !isError(resp) {
		t.Fatal("expected error response for missing name")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("error_code: got %q, want %q", errorCode(resp), "validation_error")
	}
}

// TestStartService_UnknownService verifies that starting a service that does
// not exist returns an error response (exec_error on macOS/Linux, unsupported_os
// on other platforms).
func TestStartService_UnknownService(t *testing.T) {
	handler := StartService()
	resp := callTool(t, handler, map[string]any{"name": "orchestra-nonexistent-service-xyz"})
	if !isError(resp) {
		t.Fatal("expected error response for unknown service")
	}
}

// --- stop_service ---

// TestStopService_MissingName verifies that omitting the required "name" field
// returns a validation_error.
func TestStopService_MissingName(t *testing.T) {
	handler := StopService()
	resp := callTool(t, handler, nil)
	if !isError(resp) {
		t.Fatal("expected error response for missing name")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("error_code: got %q, want %q", errorCode(resp), "validation_error")
	}
}

// TestStopService_UnknownService verifies that stopping a service that does not
// exist returns an error response.
func TestStopService_UnknownService(t *testing.T) {
	handler := StopService()
	resp := callTool(t, handler, map[string]any{"name": "orchestra-nonexistent-service-xyz"})
	if !isError(resp) {
		t.Fatal("expected error response for unknown service")
	}
}

// --- restart_service ---

// TestRestartService_MissingName verifies that omitting the required "name" field
// returns a validation_error.
func TestRestartService_MissingName(t *testing.T) {
	handler := RestartService()
	resp := callTool(t, handler, nil)
	if !isError(resp) {
		t.Fatal("expected error response for missing name")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("error_code: got %q, want %q", errorCode(resp), "validation_error")
	}
}

// TestRestartService_UnknownService verifies that restarting a service that does
// not exist returns an error response.
func TestRestartService_UnknownService(t *testing.T) {
	handler := RestartService()
	resp := callTool(t, handler, map[string]any{"name": "orchestra-nonexistent-service-xyz"})
	if !isError(resp) {
		t.Fatal("expected error response for unknown service")
	}
}

// --- service_logs ---

// TestServiceLogs_MissingName verifies that omitting the required "name" field
// returns a validation_error.
func TestServiceLogs_MissingName(t *testing.T) {
	handler := ServiceLogs()
	resp := callTool(t, handler, nil)
	if !isError(resp) {
		t.Fatal("expected error response for missing name")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("error_code: got %q, want %q", errorCode(resp), "validation_error")
	}
}

// TestServiceLogs_UnknownService retrieves logs for a service that does not
// exist. On macOS the "log show --predicate process==" command may succeed with
// empty output (no match), so we accept both success and error responses —
// either is a valid outcome for a nonexistent service name.
func TestServiceLogs_UnknownService(t *testing.T) {
	handler := ServiceLogs()
	resp := callTool(t, handler, map[string]any{"name": "orchestra-nonexistent-service-xyz"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	// Both success (empty log) and exec_error / unsupported_os are acceptable.
}

// --- service_info ---

// TestServiceInfo_MissingName verifies that omitting the required "name" field
// returns a validation_error.
func TestServiceInfo_MissingName(t *testing.T) {
	handler := ServiceInfo()
	resp := callTool(t, handler, nil)
	if !isError(resp) {
		t.Fatal("expected error response for missing name")
	}
	if errorCode(resp) != "validation_error" {
		t.Errorf("error_code: got %q, want %q", errorCode(resp), "validation_error")
	}
}

// TestServiceInfo_UnknownService verifies that querying info for a service that
// does not exist returns an error response (exec_error, not_found, or
// unsupported_os depending on the platform and fallback path).
func TestServiceInfo_UnknownService(t *testing.T) {
	handler := ServiceInfo()
	resp := callTool(t, handler, map[string]any{"name": "orchestra-nonexistent-service-xyz"})
	if !isError(resp) {
		t.Fatal("expected error response for unknown service")
	}
}
