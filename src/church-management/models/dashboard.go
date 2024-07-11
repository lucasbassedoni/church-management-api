package models

type BirthdayUser struct {
	Name string `json:"name"`
	Day  int    `json:"day"`
	Age  int    `json:"age"`
}

type DashboardData struct {
	TotalUsers       int            `json:"total_users"`
	TotalMembers     int            `json:"total_members"`
	TotalNonBaptized int            `json:"total_non_baptized"`
	TotalBaptized    int            `json:"total_baptized"`
	RecentUsers      []string       `json:"recent_users"`
	BirthdayUsers    []BirthdayUser `json:"birthday_users"`
	UpcomingEvents   []Event        `json:"upcoming_events"`
}
