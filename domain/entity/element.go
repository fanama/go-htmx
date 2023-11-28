package entity

type ResumePage struct {
	Header   string
	Sections string
}

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
