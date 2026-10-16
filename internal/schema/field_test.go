package schema

import "testing"

func TestAutomaticStringType(t *testing.T) {
	for _, name := range []string{"name", "email", "city"} {
		fieldType, automatic := AutomaticStringType(name)
		if !automatic || fieldType != "string" {
			t.Errorf("AutomaticStringType(%q) = %q, %t; want string, true", name, fieldType, automatic)
		}
	}

	if _, automatic := AutomaticStringType("description"); automatic {
		t.Error("description should remain a manually configured field")
	}
}
