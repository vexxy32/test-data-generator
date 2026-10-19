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

func TestAutomaticStringKind(t *testing.T) {
	tests := map[string]string{
		"First_Name":    "name",
		"job_titles":    "jobtitle",
		"email_domains": "emaildomain",
	}

	for name, wantKind := range tests {
		kind, automatic := AutomaticStringKind(name)
		if !automatic || kind != wantKind {
			t.Errorf("AutomaticStringKind(%q) = %q, %t; want %q, true", name, kind, automatic, wantKind)
		}
	}
}
