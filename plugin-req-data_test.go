package goforjj

import "testing"

func TestNewReqData(t *testing.T) {
	t.Log("Expect NewReqData to return a valid request.")

	// --- Setting test context ---

	// --- Run the test ---
	// Creating the request object
	ret := NewReqData()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected NewReqData to return a request Object. Got nil.")
	} else if ret.Forj == nil {
		t.Error("Expected NewReqData.Forj to be initialized. Got nil.")
	} else if ret.Objects == nil {
		t.Error("Expected NewReqData.Objects to be initialized. Got nil.")
	}

}

func TestAddObjectActions(t *testing.T) {
	t.Log("Expect AddObjectActions to add properly data in the request.")

	// --- Setting test context ---
	req := NewReqData()

	const (
		test      = "test1"
		testValue = "test"
		obj1      = "obj1"
		instance1 = "instance1"
	)

	if req == nil {
		return
	}

	keys := make(InstanceKeys)
	keys[test] = &ValueStruct{value: testValue}
	extent := make(InstanceExtentKeys)
	creds := make(map[string]string)

	// --- Run the test ---
	// Check if we can load an instance without 'master' loader. We should not.
	req.AddObjectActions(obj1, instance1, keys, extent, creds)

	// --- Start testing ---
	if req.Objects == nil {
		t.Error("Expected req.Objects to be initialized. Got nil.")
	} else if v1 := len(req.Objects); v1 != 1 {
		t.Errorf("Expected req.Objects to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.Objects[obj1]; !f1 {
		t.Errorf("Expected req.Objects to contain '%s'. Not found.", obj1)
	} else if v3, f2 := v2[instance1]; !f2 {
		t.Errorf("Expected '%s' to have '%s'. Not found.", obj1, instance1)
	} else if v4, f3 := v3[test]; !f3 {
		t.Errorf("Expected '%s/%s' to have '%s'. Not found.", obj1, instance1, test)
	} else if v5, ok1 := v4.(*ValueStruct); !ok1 {
		t.Errorf("Expected '%s/%s/%s' to have '%s'. Is not.", obj1, instance1, test, "*ValueStruct")
	} else if v5.value != testValue {
		t.Errorf("Expected '%s/%s/%s=%s' to have '%s'. Got '%s'.", obj1, instance1, test, "*ValueStruct", v5.value, testValue)
	}

	// --- Setting test context ---
	const (
		test2      = "test2"
		testValue2 = "testValue2"
	)

	if req == nil {
		return
	}

	keys = make(InstanceKeys)
	keys[test] = &ValueStruct{value: testValue}
	extent = make(InstanceExtentKeys)
	extent[test2] = &ValueStruct{value: testValue2}

	// --- Run the test ---
	// Check if we can load an instance without 'master' loader. We should not.
	req.AddObjectActions(obj1, instance1, keys, extent, creds)

	// --- Start testing ---
	if req.Objects == nil {
		t.Error("Expected req.Objects to be initialized. Got nil.")
	} else if v1 := len(req.Objects); v1 != 1 {
		t.Errorf("Expected req.Objects to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.Objects[obj1]; !f1 {
		t.Errorf("Expected req.Objects to contain '%s'. Not found.", obj1)
	} else if v3, f2 := v2[instance1]; !f2 {
		t.Errorf("Expected '%s' to have '%s'. Not found.", obj1, instance1)
	} else if v4, f3 := v3[test]; !f3 {
		t.Errorf("Expected '%s/%s' to have '%s'. Not found.", obj1, instance1, test)
	} else if v5, ok1 := v4.(*ValueStruct); !ok1 {
		t.Errorf("Expected '%s/%s/%s' to have '%s'. Is not.", obj1, instance1, test, "*ValueStruct")
	} else if v5.value != testValue {
		t.Errorf("Expected '%s/%s/%s=%s' to have '%s'. Got '%s'.", obj1, instance1, test, "*ValueStruct", testValue, v5.value)
	} else if v6, f4 := v3["extent"]; !f4 {
		t.Errorf("Expected '%s/%s' to have '%s'. Not found.", obj1, instance1, test)
	} else if v7, ok2 := v6.(InstanceExtentKeys); !ok2 {
		t.Errorf("Expected '%s/%s/%s' to have '%s'. Is not.", obj1, instance1, test, "*ValueStruct")
	} else if v8, f5 := v7[test2]; !f5 {
		t.Errorf("Expected '%s/%s/extent' to have '%s'. Not found.", obj1, instance1, test2)
	} else if v8.value != testValue2 {
		t.Errorf("Expected '%s/%s/extent/%s=%s' to have '%s'. Got '%s'.", obj1, instance1, test, "*ValueStruct", testValue2, v5.value)
	}

	// --- Setting test context ---
	const (
		credKey   = "key"
		credValue = "value"
	)

	if req == nil {
		return
	}

	creds[credKey] = credValue

	// --- Run the test ---
	// Check if we can load an instance without 'master' loader. We should not.
	req.AddObjectActions(obj1, instance1, keys, extent, creds)

	// --- Start testing ---
	if req.Objects == nil {
		t.Error("Expected req.Objects to be initialized. Got nil.")
	} else if v1 := len(req.Objects); v1 != 1 {
		t.Errorf("Expected req.Objects to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.Objects[obj1]; !f1 {
		t.Errorf("Expected req.Objects to contain '%s'. Not found.", obj1)
	} else if v3, f2 := v2[instance1]; !f2 {
		t.Errorf("Expected '%s' to have '%s'. Not found.", obj1, instance1)
	} else if v4, f3 := v3[test]; !f3 {
		t.Errorf("Expected '%s/%s' to have '%s'. Not found.", obj1, instance1, test)
	} else if v5, ok1 := v4.(*ValueStruct); !ok1 {
		t.Errorf("Expected '%s/%s/%s' to have '%s'. Is not.", obj1, instance1, test, "*ValueStruct")
	} else if v5.value != testValue {
		t.Errorf("Expected '%s/%s/%s=%s' to have '%s'. Got '%s'.", obj1, instance1, test, "*ValueStruct", testValue, v5.value)
	} else if v6, f4 := v3["extent"]; !f4 {
		t.Errorf("Expected '%s/%s' to have '%s'. Not found.", obj1, instance1, test)
	} else if v7, ok2 := v6.(InstanceExtentKeys); !ok2 {
		t.Errorf("Expected '%s/%s/%s' to have '%s'. Is not.", obj1, instance1, test, "*ValueStruct")
	} else if v8, f5 := v7[test2]; !f5 {
		t.Errorf("Expected '%s/%s/extent' to have '%s'. Not found.", obj1, instance1, test2)
	} else if v8.value != testValue2 {
		t.Errorf("Expected '%s/%s/extent/%s=%s' to have '%s'. Got '%s'.", obj1, instance1, test, "*ValueStruct", testValue2, v5.value)
	} else if req.Creds == nil {
		t.Error("Expected req.Creds to be initialized. Got nil.")
	} else if v9, f6 := req.Creds[obj1+"-"+instance1+"-"+credKey]; !f6 {
		t.Errorf("Expected req.Creds to contain '%s'. Not found.", obj1+"-"+instance1+"-"+credKey)
	} else if v9 != credValue {
		t.Errorf("Expected req.Creds[%s] to be '%s'. Got '%s'", obj1+"-"+instance1+"-"+credKey, credValue, v9)
	}

}

func TestSetForjFlag(t *testing.T) {
	t.Log("Expect SetForjFlag to add properly data in the request.")

	// --- Setting test context ---
	req := NewReqData()

	const (
		test      = "test1"
		testValue = "test"
		obj1      = "obj1"
		instance1 = "instance1"
	)

	if req == nil {
		return
	}

	// --- Run the test ---
	// Set basic forj value
	req.SetForjFlag(test, testValue, false, false)

	// --- Start testing ---
	if req.Forj == nil {
		t.Error("Expected req.Forj to be initialized. Got nil.")
	} else if v1 := len(req.Forj); v1 != 1 {
		t.Errorf("Expected req.Forj to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.Forj[test]; !f1 {
		t.Errorf("Expected req.Forj to contain '%s'. Not found.", test)
	} else if v2 != testValue {
		t.Errorf("Expected '%s' to have '%s'. Got '%s'.", test, testValue, v2)
	}

	// --- Setting test context ---
	const (
		test2      = "test2"
		testValue2 = "testValue2"
	)

	if req == nil {
		return
	}

	// --- Run the test ---
	// Check we can add a new one but as extent
	req.SetForjFlag(test2, testValue2, false, true)

	// --- Start testing ---
	if req.ForjExtent == nil {
		t.Error("Expected req.Forj to be initialized. Got nil.")
	} else if v1 := len(req.ForjExtent); v1 != 1 {
		t.Errorf("Expected req.Forj to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.ForjExtent[test2]; !f1 {
		t.Errorf("Expected req.Forj to contain '%s'. Not found.", test2)
	} else if v2 != testValue2 {
		t.Errorf("Expected '%s' to have '%s'. Got '%s'.", test2, testValue2, v2)
	}

	// --- Setting test context ---
	const ()

	req = NewReqData()

	if req == nil {
		return
	}

	// --- Run the test ---
	// Check if we can set a cred value. The Forj value is kept for compatibility.
	req.SetForjFlag(test, testValue, true, false)

	// --- Start testing ---
	if req.Forj == nil {
		t.Error("Expected req.Forj to be initialized. Got nil.")
	} else if v1 := len(req.Forj); v1 != 1 {
		t.Errorf("Expected req.Forj to contains only one element. Got %d.", v1)
	} else if v2, f1 := req.Forj[test]; !f1 { // The Forj value is kept for compatibility.
		t.Errorf("Expected req.Forj to contain '%s'. Not found.", test)
	} else if v2 != testValue {
		t.Errorf("Expected '%s' to have '%s'. Got '%s'.", test, testValue, v2)
	} else if v3, f2 := req.Creds[test]; !f2 {
		t.Errorf("Expected req.Forj to contain '%s'. Not found.", test)
	} else if v3 != testValue {
		t.Errorf("Expected '%s' to have '%s'. Got '%s'.", test, testValue, v3)
	}

	// --- Setting test context ---
	const ()

	if req == nil {
		return
	}

	// --- Run the test ---
	// Check if we can set an extent cred.
	req.SetForjFlag(test2, testValue2, true, true)

	// --- Start testing ---
	if req.ForjExtent != nil {
		t.Error("Expected req.Extent to be nil. Got initialized.")
	} else if req.Creds == nil {
		t.Error("Expected req.Creds to be initialized. Got nil.")
	} else if v3, f2 := req.Creds[test]; !f2 {
		t.Errorf("Expected req.Forj to contain '%s'. Not found.", test)
	} else if v3 != testValue {
		t.Errorf("Expected '%s' to have '%s'. Got '%s'.", test, testValue, v3)
	}

}
