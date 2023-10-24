package main

import (
	"bytes"
	"fanama/go-htmx/domain/entity"
	"fmt"
	"html"
	"html/template"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Component = entity.Component
type Resume = entity.Resume

var resume = entity.Resume{
	Name:       "John Doe",
	Profession: "Software Engineer",
	Sections: []entity.Section{
		{
			Title: "Work Experience",
			Experiences: []entity.Experience{
				{
					Title:       "Software Engineer Intern",
					Description: "Worked on several key projects including...",
					Company:     "Google",
					Achievements: []string{
						"Developed a new algorithm for image classification.",
						"Improved the performance of the search engine by 10%.",
					},
					StartYear: 2022,
					EndYear:   2023,
				},
				{
					Title:       "Backend Developer",
					Description: "Focused on developing scalable microservices...",
					Company:     "Amazon",
					Achievements: []string{
						"Optimized data storage solutions.",
						"Led a team for a critical infrastructure project.",
					},
					StartYear: 2020,
					EndYear:   2022,
				},
			},
		},
		{
			Title: "Education",
			Experiences: []entity.Experience{
				{
					Title:       "Bachelor in Computer Science",
					Description: "Studied various aspects of CS...",
					Company:     "MIT",
					Achievements: []string{
						"Graduated with Honors",
						"Participated in several coding competitions.",
					},
					StartYear: 2016,
					EndYear:   2020,
				},
			},
		},
	},
}

func GetHTML(path string, data interface{}) (string, error) {
	if path == "" || data == nil {
		return "", fmt.Errorf("invalid input - path: %s, data: %v", path, data)
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

const (
	HeaderPath       = "components/molecules/headder.html"
	ExperiencePath   = "components/atoms/experience.html"
	SectionPath      = "components/molecules/section.html"
	DefaultComponent = "body"
)

type SectionData struct {
	Title   string
	Content string
}

// Handle error responses
func handleError(c *fiber.Ctx, err error) error {
	fmt.Println(err)
	return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("<div>ERROR: %s</div>", err.Error()))
}

// Build component path
func buildComponentPath(folder, name string) string {
	return "./components/" + folder + "/" + name + ".html"
}

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/cv", func(c *fiber.Ctx) error {
		headerHTML, err := GetHTML(HeaderPath, resume)
		if err != nil {
			return handleError(c, err)
		}

		var sectionsHTML strings.Builder
		for _, section := range resume.Sections {
			var renderedExperienceList strings.Builder

			for _, experience := range section.Experiences {
				experienceHTML, err := GetHTML(ExperiencePath, experience)
				if err != nil {
					return handleError(c, err)
				}
				renderedExperienceList.WriteString(html.UnescapeString(experienceHTML))
			}

			sectionData := SectionData{
				Title:   section.Title,
				Content: renderedExperienceList.String(),
			}

			sectionHTML, err := GetHTML(SectionPath, sectionData)
			if err != nil {
				return handleError(c, err)
			}
			sectionsHTML.WriteString(html.UnescapeString(sectionHTML))
		}

		return c.Type("html").SendString(headerHTML + sectionsHTML.String())
	})

	app.Get("/:folder/:name", func(c *fiber.Ctx) error {
		name := c.Params("name", DefaultComponent)
		folder := c.Params("folder", DefaultComponent)
		htmlElement, err := GetHTML(buildComponentPath(folder, name), &Component{Title: name, ID: name})
		if err != nil {
			return handleError(c, err)
		}
		return c.SendString(htmlElement)
	})

	log.Fatal(app.Listen(":3000"))
}
