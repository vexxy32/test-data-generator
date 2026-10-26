package schema

import "strings"

// FieldDefinition описывает поле будущей записи. Min и Max задают диапазон
// чисел или длину обычной строки; Precision нужен только для float.
// AllowEmpty разрешает пустую строку для string или null для остальных типов.
type FieldDefinition struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Precision  int     `json:"precision"`
	AllowEmpty bool    `json:"allowEmpty"`
}

// automaticFields хранит все допустимые названия специальных строковых полей.
// Значение карты — каноническое имя, которое использует генератор.
var automaticFields = map[string]string{
	"name": "name", "names": "name", "firstname": "name", "firstnames": "name", "first_name": "name",
	"surname": "surname", "surnames": "surname",
	"email": "email", "emails": "email",
	"emaildomain": "emaildomain", "emaildomains": "emaildomain", "email_domain": "emaildomain", "email_domains": "emaildomain",
	"city": "city", "cities": "city",
	"country": "country", "countries": "country",
	"company": "company", "companies": "company",
	"department": "department", "departments": "department",
	"jobtitle": "jobtitle", "jobtitles": "jobtitle", "job_title": "jobtitle", "job_titles": "jobtitle",
	"streetaddress": "streetaddress", "streetaddresses": "streetaddress", "street_address": "streetaddress", "street_addresses": "streetaddress",
	"phonenumber": "phonenumber", "phonenumbers": "phonenumber", "phone_number": "phonenumber", "phone_numbers": "phonenumber",
	"userrole": "userrole", "userroles": "userrole", "user_role": "userrole", "user_roles": "userrole",
	"status": "status", "statuses": "status",
	"color": "color", "colors": "color",
	"currency": "currency", "currencies": "currency",
	"language": "language", "languages": "language",
	"contenttitle": "contenttitle", "contenttitles": "contenttitle", "content_title": "contenttitle", "content_titles": "contenttitle",
}

// AutomaticStringKind возвращает каноническое имя специального поля.
func AutomaticStringKind(name string) (string, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	kind, found := automaticFields[name]
	return kind, found
}

// AutomaticStringType сообщает TUI, для каких имён тип string выбирается сам.
func AutomaticStringType(name string) (string, bool) {
	_, found := AutomaticStringKind(name)
	return "string", found
}
