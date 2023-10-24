package entity

type Resume struct {
	Name       string    `json:"name"`
	Profession string    `json:"profession"`
	Sections   []Section `json:"sections"`
}

type Section struct {
	Title       string       `json:"title"`
	Experiences []Experience `json:"experiences"`
}

type Experience struct {
	Title        string   `json:"title"`
	Description  string   `json:"desctiption"`
	Company      string   `json:"company"`
	Achievements []string `json:"Achievements"`
	StartYear    int      `json:"start_year"`
	EndYear      int      `json:"end_year"`
}
