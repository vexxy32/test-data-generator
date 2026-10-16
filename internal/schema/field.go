package schema

import "strings"

type FieldDefinition struct {
	Name      string
	Type      string
	Min       float64
	Max       float64
	Precision int
}

func AutomaticStringType(name string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case
		// firstNames
		"name", "names", "firstname", "firstnames", "first_name",

		// surnames
		"surname", "surnames",

		// emailDomains
		"email", "emails",
		"emaildomain", "emaildomains", "email_domain", "email_domains",

		// cities
		"city", "cities",

		// countries
		"country", "countries",

		// companies
		"company", "companies",

		// departments
		"department", "departments",

		// jobTitles
		"jobtitle", "jobtitles", "job_title", "job_titles",

		// streetAddresses
		"streetaddress", "streetaddresses", "street_address", "street_addresses",

		// phoneNumbers
		"phonenumber", "phonenumbers", "phone_number", "phone_numbers",

		// userRoles
		"userrole", "userroles", "user_role", "user_roles",

		// statuses
		"status", "statuses",

		// colors
		"color", "colors",

		// currencies
		"currency", "currencies",

		// languages
		"language", "languages",

		// contentTitles
		"contenttitle", "contenttitles", "content_title", "content_titles":
		return "string", true

	default:
		return "", false
	}
}
