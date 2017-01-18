package goforjj

import "testing"

func TestActionKeys_AddKey(t *testing.T) {
	t.Log("Expect ActionKeys_AddKey() to Add a key in a map.")

	// --- Setting test context ---
	var a ActionKeys

	// --- Run the test ---
	b := a.AddKey("key", "value")

	// --- Start testing ---
	if a != nil {
		t.Error("Expected original Action to NOT be initialized. Got it.")
	}
	if b == nil {
		t.Error("Expected Action to be initialized. Got Nil.")
		return
	}
	a = b
	if v, found := b["key"]; !found {
		t.Errorf("Expected Action to be have a key '%s'. not found", "action")
		return
	} else {
		if v != "value" {
			t.Errorf("Expected action key '%s' to be have '%s'. not found", "key", "value")
		}
	}

	// --- start another test ---
	a.AddKey("key2", "value2")

	// --- Start testing ---
	if v, found := a["key2"]; !found {
		t.Errorf("Expected Action to be have a key '%s'. not found", "key2")
		return
	} else {
		if v != "value2" {
			t.Errorf("Expected action key '%s' to be have '%s'. not found", "key2", "value2")
		}
	}
	if v, found := b["key2"]; !found {
		t.Errorf("Expected Instance to be have action '%s'. not found", "action2")
		return
	} else {
		if v != "value2" {
			t.Errorf("Expected action key '%s' to be have '%s'. not found", "key2", "value2")
		}
	}

}

func TestInstanceActions_AddAction(t *testing.T) {
	t.Log("Expect AddActions() to Add an action in the map.")

	// --- Setting test context ---
	var a InstanceActions
	// --- Run the test ---
	b := a.AddAction("action", nil)
	// --- Start testing ---
	if a != nil {
		t.Error("Expected original Instance Action to NOT be initialized. Got it.")
	}
	if b == nil {
		t.Error("Expected Instance Action to be initialized. Got Nil.")
		return
	}

	// --- Update test context ---
	a = make(InstanceActions)

	// --- Run the test ---
	b = a.AddAction("action", nil)
	// --- Start testing ---
	if _, found := a["action"]; !found {
		t.Errorf("Expected Instance to be have action '%s'. not found", "action")
		return
	}
	if _, found := b["action"]; !found {
		t.Errorf("Expected Instance to be have action '%s'. not found", "action")
		return
	}

	// --- Updating test context ---
	a.AddAction("action2", nil)
	// --- Start testing ---
	if _, found := a["action2"]; !found {
		t.Errorf("Expected Instance to be have action '%s'. not found", "action2")
		return
	}
	if _, found := b["action2"]; !found {
		t.Errorf("Expected Instance to be have action '%s'. not found", "action2")
		return
	}
}

func TestPluginReqData_AddObjectActions(t *testing.T) {
	t.Log("Expect PluginReqData_AddObjectActions() to add an Object actions.")

	// --- Setting test context ---
	var d *PluginReqData

	// --- Run the test ---
	d.AddObjectActions("type", "instance", nil)
	// No exception should occur.

	d = new(PluginReqData)
	// --- Run the test ---
	d.AddObjectActions("type", "instance", nil)
	// --- Start testing ---
	if d.Objects == nil {
		t.Error("Expected Objects to be initialized. Got nil")
		return
	}
	if object, found_object := d.Objects["type"]; !found_object {
		t.Errorf("Expected Objects to have '%s'. Not found", "type")
		return
	} else {
		if _, found_action := object["instance"]; !found_action {
			t.Errorf("Expected Object instance '%s' to have '%s'. Not found", "type", "instance")
		}
	}
	// --- Updating test context ---
	a := make(InstanceActions)
	b := a.AddAction("action", nil)
	// --- Run the test ---
	d.AddObjectActions("type", "instance", b)
	// --- Start testing ---
	if _, found := d.Objects["type"]["instance"]["action"]; !found {
		t.Errorf("Expected Object instance '%s-%s' to have action '%s'. Not found", "type", "instance", "action")
	}
}
