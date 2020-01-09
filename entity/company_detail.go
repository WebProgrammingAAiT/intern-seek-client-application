package entity

type CompanyDetail struct {
	ID          int
	UserId      int
	Country     string
	City        string
	FocusArea   string
	Description string
	User        User
}
