package entity

type Component struct {
	Title     string
	ClassName string
	ID        string
	Content   string
}

type SectionData struct {
	Title   string
	Content string
}

type ErrorData struct {
	Message string
}
