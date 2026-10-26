package generator

import (
	"errors"
	"math"
	"math/rand/v2"

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
	kind, automatic := schema.AutomaticStringKind(field.Name)

	if !automatic {
		return generateRandomString(field, randomGenerator)
	}

	switch kind {
	case "name":
		return firstName

	case "surname":
		return surname

	case "email":
		return randomEmail(firstName, surname, randomGenerator)

	case "emaildomain":
		return randomItem(emailDomains, randomGenerator)

	case "city":
		return randomItem(cities, randomGenerator)

	case "country":
		return randomItem(countries, randomGenerator)

	case "company":
		return randomItem(companies, randomGenerator)

	case "department":
		return randomItem(departments, randomGenerator)

	case "jobtitle":
		return randomItem(jobTitles, randomGenerator)

	case "streetaddress":
		return randomItem(streetAddresses, randomGenerator)

	case "phonenumber":
		return randomItem(phoneNumbers, randomGenerator)

	case "userrole":
		return randomItem(userRoles, randomGenerator)

	case "status":
		return randomItem(statuses, randomGenerator)

	case "color":
		return randomItem(colors, randomGenerator)

	case "currency":
		return randomItem(currencies, randomGenerator)

	case "language":
		return randomItem(languages, randomGenerator)

	case "contenttitle":
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

// GenerateRecordsStream генерирует записи по одной и сразу передаёт каждую в обработчик onRecord.
func GenerateRecordsStream(
	fields []schema.FieldDefinition,
	count int,
	seed uint64,
	onRecord func(record map[string]any) error,
) error {

	if count <= 0 {
		return errors.New("количество записей должно быть больше нуля")
	}
	if onRecord == nil {
		return errors.New("обработчик записи не задан")
	}

	pcg := rand.NewPCG(seed, 1)
	randomGenerator := rand.New(pcg)
	specialRandom := rand.New(rand.NewPCG(seed, 2))
	emptyRecordIndex := specialRandom.IntN(count)

	for i := 0; i < count; i++ {
		record := generateRecord(fields, randomGenerator)
		if i == emptyRecordIndex {
			for _, field := range fields {
				if field.AllowEmpty {
					if field.Type == "string" {
						record[field.Name] = ""
					} else {
						record[field.Name] = nil
					}
				}
			}
		}
		record["_test_id"] = i + 1
		if err := onRecord(record); err != nil {
			return err
		}
	}

	return nil
}
