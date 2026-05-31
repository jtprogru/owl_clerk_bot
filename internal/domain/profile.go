package domain

import "time"

type Category string

const (
	CategoryUnknown Category = "unknown"
	CategoryHR      Category = "hr"
	CategoryCoop    Category = "coop"
	CategoryFriend  Category = "friend"
	CategoryMisc    Category = "misc"
	CategorySpam    Category = "spam"
)

func (c Category) String() string { return string(c) }

func (c Category) Valid() bool {
	switch c {
	case CategoryUnknown, CategoryHR, CategoryCoop, CategoryFriend, CategoryMisc, CategorySpam:
		return true
	}
	return false
}

func (c Category) Title() string {
	switch c {
	case CategoryHR:
		return "HR"
	case CategoryCoop:
		return "Сотрудничество"
	case CategoryFriend:
		return "Дружба"
	case CategoryMisc:
		return "Разное"
	case CategorySpam:
		return "Спам"
	default:
		return "Неизвестно"
	}
}

func ParseCategoryFromTitle(s string) (Category, bool) {
	switch s {
	case "HR":
		return CategoryHR, true
	case "Сотрудничество":
		return CategoryCoop, true
	case "Дружба":
		return CategoryFriend, true
	case "Разное":
		return CategoryMisc, true
	}
	return CategoryUnknown, false
}

type Profile struct {
	UID         int64
	FirstName   string
	LastName    string
	Username    string
	Category    Category
	IsBlocked   bool
	Notes       string
	Contact     string
	FirstSeenAt time.Time
	LastSeenAt  time.Time
}

func (p Profile) DisplayName() string {
	name := p.FirstName
	if p.LastName != "" {
		if name != "" {
			name += " "
		}
		name += p.LastName
	}
	if name == "" && p.Username != "" {
		name = "@" + p.Username
	}
	if name == "" {
		name = "(без имени)"
	}
	return name
}
