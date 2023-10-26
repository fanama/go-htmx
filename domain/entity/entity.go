package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Resume struct {
	Name       string    `json:"name"`
	Profession string    `json:"profession"`
	Contact    Contact   `json:"contact"`
	Sections   []Section `json:"sections"`
	Skills     Skills    `json:"skills"`
}

type Contact struct {
	Email     string `json:"email"`
	Location  string `json:"location"`
	Linkedin  string `json:"linkedin"`
	Phone     string `json:"phone"`
	Portfolio string `json:"portfolio"`
}

type Section struct {
	Title       string       `json:"title"`
	Experiences []Experience `json:"experiences"`
}

type Experience struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Company      string   `json:"company"`
	Achievements []string `json:"Achievements"`
	StartYear    string   `json:"start_date"`
	EndYear      string   `json:"end_date"`
}

type Skills struct {
	Technologies []string `json:"technologies"`
}

func (resume *Resume) AddExperience(sectionTitle string, experience Experience) {

	for _, section := range resume.Sections {
		if section.Title == sectionTitle {
			section.Experiences = append(section.Experiences, experience)
		}
	}
}

func (resume *Resume) AddSection(section Section) {
	resume.Sections = append(resume.Sections, section)
}

func ReadResumeFromFile(filename string) (*Resume, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Unmarshal the content into the Resume struct
	var resume Resume
	err = json.Unmarshal(content, &resume)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &resume, nil
}
