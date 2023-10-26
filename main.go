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
	ErrorPath        = "components/atoms/error.html"
	SectionPath      = "components/molecules/section.html"
	DefaultComponent = "body"
)

type SectionData struct {
	Title   string
	Content string
}

type ErrorData struct {
	Message string
}

// Handle error responses
func handleError(c *fiber.Ctx, err error) error {
	fmt.Println(err)
	errorHTML, _ := GetHTML(ErrorPath, ErrorData{Message: err.Error()})
	return c.SendString(errorHTML)
}

// Build component path
func buildComponentPath(folder, name string) string {
	return "./components/" + folder + "/" + name + ".html"
}

func main() {

	resume, err := entity.ReadResumeFromFile("./example.json")

	if err != nil {
		fmt.Println(err)
		return

	}
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/cv", func(c *fiber.Ctx) error {

		headerHTML, err := GetHTML(HeaderPath, &resume)
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

		return c.Type("html").SendString(headerHTML + "<div id=\"resume\" class=\"grid grid-cols-2\">" + sectionsHTML.String() + "</div>")
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
