package generator

import (
	"math"
	"math/rand/v2"
	"strings"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

func generateValue(field schema.FieldDefinition) any {
	switch field.Type {

	case "integer":
		min := int(field.Min)
		max := int(field.Max)

		return min + rand.IntN(max-min+1)

	case "float":
		value := field.Min + rand.Float64()*(field.Max-field.Min)
		scale := math.Pow10(field.Precision)
		return math.Round(value*scale) / scale

	case "string":
		return generateString(field)

	case "boolean":
		return rand.IntN(2) == 1

	}

	return nil
}

func generateString(field schema.FieldDefinition) string {
	switch strings.ToLower(strings.TrimSpace(field.Name)) {
	case "name", "names", "firstname", "firstnames", "first_name":
		return randomItem(firstNames)

	case "surname", "surnames":
		return randomItem(surnames)

	case "email", "emails":
		return randomEmail()

	case "emaildomain", "emaildomains", "email_domain", "email_domains":
		return randomItem(emailDomains)

	case "city", "cities":
		return randomItem(cities)

	case "country", "countries":
		return randomItem(countries)

	case "company", "companies":
		return randomItem(companies)

	case "department", "departments":
		return randomItem(departments)

	case "jobtitle", "jobtitles", "job_title", "job_titles":
		return randomItem(jobTitles)

	case "streetaddress", "streetaddresses", "street_address", "street_addresses":
		return randomItem(streetAddresses)

	case "phonenumber", "phonenumbers", "phone_number", "phone_numbers":
		return randomItem(phoneNumbers)

	case "userrole", "userroles", "user_role", "user_roles":
		return randomItem(userRoles)

	case "status", "statuses":
		return randomItem(statuses)

	case "color", "colors":
		return randomItem(colors)

	case "currency", "currencies":
		return randomItem(currencies)

	case "language", "languages":
		return randomItem(languages)

	case "contenttitle", "contenttitles", "content_title", "content_titles":
		return randomItem(contentTitles)

	default:
		return generateRandomString(field)
	}
}

func generateRandomString(field schema.FieldDefinition) string {
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
}

func generateRecord(fields []schema.FieldDefinition) map[string]any {
	record := map[string]any{}

	for _, field := range fields {
		record[field.Name] = generateValue(field)
	}

	return record
}

func GenerateRecords(fields []schema.FieldDefinition, count int) []map[string]any {
	records := []map[string]any{}

	for i := 0; i < count; i++ {
		record := generateRecord(fields)

		record["_test_id"] = i + 1

		records = append(records, record)
	}

	return records
}
