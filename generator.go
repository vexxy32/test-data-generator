package main

import (
	"math"
	"math/rand/v2"
)

type FieldDefinition struct {
	Name string
	Type string
	Min  float64
	Max  float64
}

func generateValue(field FieldDefinition) any {
	switch field.Type {

	case "integer":
		min := int(field.Min)
		max := int(field.Max)

		return min + rand.IntN(max-min+1)

	case "float":

		value := field.Min + rand.Float64()*(field.Max-field.Min)
		return math.Round(value*100) / 100

	case "string":
		const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		minLength := int(field.Min)
		maxLength := int(field.Max)

		length := minLength + rand.IntN(maxLength-minLength+1)
		randomString := make([]byte, length)

		for i := range randomString {
			randomIndex := rand.IntN(len(alphabet))
			randomString[i] = alphabet[randomIndex]
		}

		return string(randomString)

	case "boolean":
		var randB bool
		randB = rand.IntN(2) == 1
		return randB
	}

	return nil
}

func generateRecord(fields []FieldDefinition) map[string]any {
	record := map[string]any{}

	for _, field := range fields {
		record[field.Name] = generateValue(field)
	}

	return record
}

func generateRecords(fields []FieldDefinition, count int) []map[string]any {
	records := []map[string]any{}

	for i := 0; i < count; i++ {
		record := generateRecord(fields)
		records = append(records, record)
	}

	return records
}

