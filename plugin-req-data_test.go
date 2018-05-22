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

	// --- Run the test ---
	// Check if we can load an instance without 'master' loader. We should not.
	req.AddObjectActions(obj1, instance1, keys, extent)

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
	req.AddObjectActions(obj1, instance1, keys, extent)

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
	// Check if we can load an instance without 'master' loader. We should not.
	req.SetForjFlag(test, testValue, false)

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
	// Check if we can load an instance without 'master' loader. We should not.
	req.SetForjFlag(test2, testValue2, true)

	// --- Start testing ---
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

}
