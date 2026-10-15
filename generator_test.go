package main

import "testing"

func TestGenerateRecords(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "age", Type: "integer", Min: 18, Max: 65},
		{Name: "price", Type: "float", Min: 10, Max: 20},
		{Name: "code", Type: "string", Min: 3, Max: 8},
		{Name: "active", Type: "boolean"},
	}

	records := generateRecords(fields, 100)
	if len(records) != 100 {
		t.Fatalf("got %d records, want 100", len(records))
	}

	for i, record := range records {
		age, ok := record["age"].(int)
		if !ok || age < 18 || age > 65 {
			t.Errorf("record %d: invalid age %v", i, record["age"])
		}
		price, ok := record["price"].(float64)
		if !ok || price < 10 || price > 20 {
			t.Errorf("record %d: invalid price %v", i, record["price"])
		}
		code, ok := record["code"].(string)
		if !ok || len(code) < 3 || len(code) > 8 {
			t.Errorf("record %d: invalid code %v", i, record["code"])
		}
		if _, ok := record["active"].(bool); !ok {
			t.Errorf("record %d: invalid boolean %v", i, record["active"])
		}
	}
}
