package filter

import "testing"

func TestUtil(t *testing.T)  {
	if IsFalse(0) != true {
		t.Fail()
	}

	if IsTrue(1) != true {
		t.Fail()
	}

	if IsTrue("1") != true {
		t.Fail()
	}
}