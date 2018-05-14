package goforjj

import (
	"testing"
)

func TestNewYamlPlugin(t *testing.T) {
	t.Log("Expect NewYamlPlugin to initialize YamlPlugin object.")

	// --- Setting test context ---
	// --- Run the test ---
	ret := NewYamlPlugin()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected YamlPlugin to exit. Got nil.")
		return
	}
	if ret.instancesDetails == nil {
		t.Error("Expected YamlPlugin to initialized properly. instancesDetails is nil.")
		return
	}
}
