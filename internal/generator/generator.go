package generator

import (
	"math"
	"math/rand/v2"
	"strings"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

// generateFieldValue выбирает способ генерации по типу одного поля.
func generateFieldValue(field schema.FieldDefinition, firstName, surname string, randomGenerator *rand.Rand) any {
	switch field.Type {

	case "integer":
		min := int(field.Min)
		max := int(field.Max)

		return min + randomGenerator.IntN(max-min+1)

	case "float":
		value := field.Min + randomGenerator.Float64()*(field.Max-field.Min)
		scale := math.Pow10(field.Precision)
		return math.Round(value*scale) / scale

	case "string":
		return generateStringValue(field, firstName, surname, randomGenerator)

	case "boolean":
		return randomGenerator.IntN(2) == 1

	}

	return nil
}

// generateStringValue использует готовые значения для специальных имён полей.
// Для обычного строкового поля создаёт строку заданной длины.
func generateStringValue(field schema.FieldDefinition, firstName, surname string, randomGenerator *rand.Rand) string {
	switch strings.ToLower(strings.TrimSpace(field.Name)) {
	case "name", "names", "firstname", "firstnames", "first_name":
		return firstName

	case "surname", "surnames":
		return surname

	case "email", "emails":
		return randomEmail(firstName, surname, randomGenerator)

	case "emaildomain", "emaildomains", "email_domain", "email_domains":
		return randomItem(emailDomains, randomGenerator)

	case "city", "cities":
		return randomItem(cities, randomGenerator)

	case "country", "countries":
		return randomItem(countries, randomGenerator)

	case "company", "companies":
		return randomItem(companies, randomGenerator)

	case "department", "departments":
		return randomItem(departments, randomGenerator)

	case "jobtitle", "jobtitles", "job_title", "job_titles":
		return randomItem(jobTitles, randomGenerator)

	case "streetaddress", "streetaddresses", "street_address", "street_addresses":
		return randomItem(streetAddresses, randomGenerator)

	case "phonenumber", "phonenumbers", "phone_number", "phone_numbers":
		return randomItem(phoneNumbers, randomGenerator)

	case "userrole", "userroles", "user_role", "user_roles":
		return randomItem(userRoles, randomGenerator)

	case "status", "statuses":
		return randomItem(statuses, randomGenerator)

	case "color", "colors":
		return randomItem(colors, randomGenerator)

	case "currency", "currencies":
		return randomItem(currencies, randomGenerator)

	case "language", "languages":
		return randomItem(languages, randomGenerator)

	case "contenttitle", "contenttitles", "content_title", "content_titles":
		return randomItem(contentTitles, randomGenerator)

	default:
		return generateRandomString(field, randomGenerator)
	}
}

// generateRandomString создаёт строку из букв и цифр в пределах заданной длины.
func generateRandomString(field schema.FieldDefinition, randomGenerator *rand.Rand) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	minLength := int(field.Min)
	maxLength := int(field.Max)
	length := minLength + randomGenerator.IntN(maxLength-minLength+1)
	randomString := make([]byte, length)

	for i := range randomString {
		randomIndex := randomGenerator.IntN(len(alphabet))
		randomString[i] = alphabet[randomIndex]
	}

	return string(randomString)
}

// generateRecord собирает одну запись; имя, фамилия и email в ней согласованы.
func generateRecord(fields []schema.FieldDefinition,
	randomGenerator *rand.Rand) map[string]any {
	record := make(map[string]any)

	firstName := randomItem(firstNames, randomGenerator)
	surname := randomItem(surnames, randomGenerator)

	for _, field := range fields {
		record[field.Name] = generateFieldValue(field, firstName, surname, randomGenerator)
	}
	return record
}

// GenerateRecords создаёт указанное число записей и добавляет каждой номер _test_id.
func GenerateRecords(
	fields []schema.FieldDefinition,
	count int,
	seed uint64,
) []map[string]any {
	pcg := rand.NewPCG(seed, 1)
	randomGenerator := rand.New(pcg)

	records := make([]map[string]any, 0, count)

	for i := 0; i < count; i++ {
		record := generateRecord(fields, randomGenerator)
		record["_test_id"] = i + 1
		records = append(records, record)
	}

	return records
}
