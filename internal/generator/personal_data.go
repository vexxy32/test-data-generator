package generator

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var firstNames = []string{
	"Alex", "Alexander", "Alexey", "Andrey", "Anton", "Artem", "Artur", "Boris",
	"Danil", "David", "Denis", "Dmitry", "Egor", "Evgeny", "Fedor", "George",
	"Gleb", "Igor", "Ilya", "Ivan", "Kirill", "Konstantin", "Leonid", "Lev",
	"Maksim", "Mark", "Mikhail", "Nikita", "Nikolay", "Oleg", "Pavel", "Petr",
	"Roman", "Ruslan", "Sergey", "Stanislav", "Stepan", "Timofey", "Timur",
	"Vadim", "Valery", "Vasily", "Victor", "Vitaly", "Vladislav", "Vyacheslav",
	"Yaroslav", "Yury", "James", "John", "Robert", "Michael", "William",
}

var surnames = []string{
	"Alexeev", "Andreev", "Antonov", "Belov", "Bogdanov", "Borisov", "Brown",
	"Chernov", "Davydov", "Davis", "Dmitriev", "Fedorov", "Frolov", "Gerasimov",
	"Golubev", "Gorbachev", "Gromov", "Gusev", "Ilyin", "Ivanov", "Kalinin",
	"Karpov", "Kazakov", "Kiselev", "Komarov", "Konev", "Konovalov", "Korolev",
	"Koval", "Kozlov", "Krylov", "Kudryavtsev", "Kulikov", "Kuzmin", "Kuznetsov",
	"Lebedev", "Maksimov", "Markov", "Matveev", "Melnikov", "Mikhailov", "Miller",
	"Morozov", "Nikolaev", "Novikov", "Orlov", "Pavlov", "Petrov", "Poliakov",
	"Popov", "Romanov", "Semenov", "Sergeev", "Smirnov", "Smith", "Sobolev",
	"Sokolov", "Soloviev", "Stepanov", "Tarasov", "Taylor", "Tikhonov", "Titov",
	"Tkachev", "Trofimov", "Vasiliev", "Vinogradov", "Volkov", "Vorobiev",
	"Voronov", "Yakovlev", "Zaitsev", "Zakharov", "Zhukov",
}

var emailDomains = []string{
	"gmail.com",
	"yandex.ru",
	"mail.ru",
	"outlook.com",
	"icloud.com",
	"proton.me",
	"yahoo.com",
}

var cities = []string{
	"Amsterdam", "Berlin", "Dubai", "Istanbul", "Kazan", "Lipetsk", "London",
	"Madrid", "Minsk", "Moscow", "New York", "Novosibirsk", "Paris", "Prague",
	"Rome", "Saint Petersburg", "Samara", "Seoul", "Singapore", "Tokyo",
	"Toronto", "Vienna", "Voronezh", "Warsaw", "Yekaterinburg",
}

var countries = []string{
	"Armenia", "Belarus", "Brazil", "Canada", "China", "Cyprus", "France",
	"Georgia", "Germany", "Italy", "Japan", "Kazakhstan", "Netherlands",
	"Poland", "Russia", "Serbia", "Spain", "Turkey", "UAE", "United Kingdom",
	"United States", "Uzbekistan",
}

var companies = []string{
	"Apex Logistics", "Astra Tech", "ByteForge", "CloudScale", "CoreDynamics",
	"CyberPulse", "DeltaPoint", "FinFlow", "GreenWave", "Hyperion Soft",
	"IronClad Systems", "Kvantum Labs", "MegaCorp", "NeuraLink Systems",
	"NextGen Media", "Nova Digital", "OmniGroup", "PrimeHoldings",
	"Quantum Leap", "RedRock Media", "SkyNetics", "Starlight Industries",
	"Synthex", "TensorCraft", "Vanguard Solutions", "Vertex IT",
}

var departments = []string{
	"Backend Engineering", "Customer Support", "Data Science", "DevOps & Infrastructure",
	"Finance & Accounting", "Frontend Engineering", "Human Resources", "Legal",
	"Marketing & Growth", "Mobile Development", "Operations", "Product Management",
	"Quality Assurance", "Sales & Business Development", "Security & Compliance",
	"UI/UX Design",
}

var jobTitles = []string{
	"Accountant", "Backend Developer", "Business Analyst", "CTO", "Customer Support Lead",
	"Data Engineer", "DevOps Engineer", "Frontend Developer", "Fullstack Developer",
	"HR Manager", "Head of Sales", "Marketing Specialist", "Operations Manager",
	"Product Designer", "Product Manager", "Project Manager", "QA Automation Engineer",
	"Security Engineer", "System Administrator", "Team Lead",
}

var streetAddresses = []string{
	"10 Downing Street", "1600 Pennsylvania Ave", "221B Baker Street",
	"42 Wallaby Way", "742 Evergreen Terrace", "Abbey Road 3", "Baker Street 12",
	"Broadway 101", "Lenina 15", "Main Street 24", "Market Street 5",
	"Mira Avenue 42", "Park Avenue 88", "Sovetskaya 10", "Wall Street 1",
}

var phoneNumbers = []string{
	"+1-202-555-0143", "+1-202-555-0188", "+44-20-7946-0912", "+44-20-7946-0995",
	"+49-30-1234567", "+7-900-123-45-67", "+7-903-555-35-35", "+7-916-234-56-78",
	"+7-920-345-67-89", "+7-999-888-77-66",
}

var userRoles = []string{
	"admin", "analyst", "auditor", "creator", "editor", "guest",
	"manager", "member", "moderator", "owner", "subscriber", "superadmin",
	"support", "tester", "viewer",
}

var statuses = []string{
	"active", "approved", "archived", "blocked", "canceled", "completed",
	"declined", "draft", "failed", "in_progress", "inactive", "new",
	"paused", "pending", "processing", "rejected", "resolved", "suspended",
}

var colors = []string{
	"Amber", "Aqua", "Beige", "Black", "Blue", "Bronze", "Brown", "Coral",
	"Cyan", "Emerald", "Gold", "Gray", "Green", "Indigo", "Magenta",
	"Navy", "Olive", "Orange", "Pink", "Purple", "Red", "Silver",
	"Teal", "Violet", "White", "Yellow",
}

var currencies = []string{
	"AED", "AMD", "BYN", "CAD", "CHF", "CNY", "EUR", "GBP", "GEL",
	"JPY", "KZT", "PLN", "RUB", "TRY", "USD", "USDT",
}

var languages = []string{
	"Arabic", "Chinese", "English", "French", "German", "Hindi", "Italian",
	"Japanese", "Kazakh", "Korean", "Polish", "Portuguese", "Russian",
	"Spanish", "Turkish",
}

var contentTitles = []string{
	"Annual Financial Review", "Breaking News Digest", "Bug Report & Reproduction",
	"Customer Feedback Summary", "Design System Overview", "Feature Implementation Plan",
	"Getting Started Guide", "Incident Postmortem", "Infrastructure Migration Notes",
	"Performance Benchmark Results", "Project Kickoff Agenda", "Quarterly Roadmap",
	"Release Notes v2.0", "Security Audit Findings", "Weekly Team Sync",
}

// randomItem выбирает одно значение из непустого списка готовых вариантов.
func randomItem(items []string, randomGenerator *rand.Rand) string {
	return items[randomGenerator.IntN(len(items))]
}

// randomEmail строит адрес из имени и фамилии текущей записи.
func randomEmail(firstName, surname string, randomGenerator *rand.Rand) string {
	return fmt.Sprintf(
		"%s.%s%d@%s",
		strings.ToLower(firstName),
		strings.ToLower(surname),
		randomGenerator.IntN(10000),
		randomItem(emailDomains, randomGenerator),
	)
}
