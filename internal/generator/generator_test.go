package generator

import (
	"math"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

func collectRecords(t *testing.T, fields []schema.FieldDefinition, count int, seed uint64) []map[string]any {
	t.Helper()
	var records []map[string]any
	if err := GenerateRecordsStream(fields, count, seed, func(record map[string]any) error {
		records = append(records, record)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return records
}

func TestGenerateRecordsStream(t *testing.T) {
	fields := []schema.FieldDefinition{
		{Name: "age", Type: "integer", Min: 18, Max: 65},
		{Name: "price", Type: "float", Min: 10, Max: 20, Precision: 2},
		{Name: "code", Type: "string", Min: 3, Max: 8},
		{Name: "name", Type: "string"},
		{Name: "active", Type: "boolean"},
	}

	records := collectRecords(t, fields, 100, 42)
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
		if math.Abs(price*100-math.Round(price*100)) > 1e-9 {
			t.Errorf("record %d: price has more than two decimal places: %v", i, price)
		}
		code, ok := record["code"].(string)
		if !ok || len(code) < 3 || len(code) > 8 {
			t.Errorf("record %d: invalid code %v", i, record["code"])
		}
		name, ok := record["name"].(string)
		if !ok || !slices.Contains(firstNames, name) {
			t.Errorf("record %d: invalid name %v", i, record["name"])
		}
		if _, ok := record["active"].(bool); !ok {
			t.Errorf("record %d: invalid boolean %v", i, record["active"])
		}
	}
}

func TestNameSurnameAndEmailBelongToSameRecord(t *testing.T) {
	fields := []schema.FieldDefinition{
		{Name: "EMAILS", Type: "string"},
		{Name: "first_name", Type: "string"},
		{Name: "surnames", Type: "string"},
	}

	for _, record := range collectRecords(t, fields, 100, 42) {
		name := strings.ToLower(record["first_name"].(string))
		surname := strings.ToLower(record["surnames"].(string))
		email := record["EMAILS"].(string)
		if !strings.HasPrefix(email, name+"."+surname) {
			t.Fatalf("email %q does not match name %q and surname %q", email, name, surname)
		}
	}
}

func TestGenerateRecordsStreamWithSameSeedReturnsSameData(t *testing.T) {
	fields := []schema.FieldDefinition{
		{Name: "name", Type: "string"},
		{Name: "email", Type: "string"},
		{Name: "age", Type: "integer", Min: 18, Max: 65},
	}

	first := collectRecords(t, fields, 3, 2026)
	second := collectRecords(t, fields, 3, 2026)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("same seed produced different records")
	}
}
