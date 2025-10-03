package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"
)

func testMCPTool(t *testing.T, tool string, args map[string]interface{}) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      tool,
			"arguments": args,
		},
	}

	reqBytes, _ := json.Marshal(request)
	reqStr := string(reqBytes) + "\n"

	cmd := exec.Command("../bin/dispatch-mcp-server")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(100 * time.Millisecond)

	stdin.Write([]byte(reqStr))
	stdin.Close()

	responseChan := make(chan string, 1)
	go func() {
		output := make([]byte, 8192)
		if n, err := stdout.Read(output); err == nil {
			responseChan <- string(output[:n])
		}
	}()

	select {
	case response := <-responseChan:
		t.Logf("MCP Response: %s", response)
	case <-time.After(10 * time.Second):
		t.Log("MCP server timed out")
	}
}

func TestCreateEstimate(t *testing.T) {
	if os.Getenv("DISPATCH_AUTH_TOKEN") == "" {
		t.Skip("Skipping - Dispatch credentials not available")
	}

	pickupInfo := map[string]interface{}{
		"business_name": "Test Business",
		"location": map[string]interface{}{
			"address": map[string]interface{}{
				"street":   "123 Main St",
				"city":     "San Francisco",
				"state":    "CA",
				"zip_code": "94105",
				"country":  "US",
			},
		},
	}

	dropOffs := []map[string]interface{}{
		{
			"business_name": "Drop Off Business",
			"location": map[string]interface{}{
				"address": map[string]interface{}{
					"street":   "456 Oak Ave",
					"city":     "San Francisco",
					"state":    "CA",
					"zip_code": "94110",
					"country":  "US",
				},
			},
		},
	}

	args := map[string]interface{}{
		"pickup_info":  jsonString(pickupInfo),
		"drop_offs":    jsonString(dropOffs),
		"vehicle_type": "cargo_van",
	}

	testMCPTool(t, "create_estimate", args)
}

func jsonString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
