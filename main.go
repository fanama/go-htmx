package main

import (
	"bytes"
	"fanama/go-htmx/domain/entity"
	"fanama/go-htmx/infra/handler"
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Component = entity.Component

type SectionData = entity.SectionData

type ErrorData = entity.ErrorData

const (
	HeaderPath       = "components/molecules/headder.html"
	ExperiencePath   = "components/atoms/experience.html"
	ErrorPath        = "components/atoms/error.html"
	SectionPath      = "components/molecules/section.html"
	DefaultComponent = "body"
)

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

// Build component path
func buildComponentPath(folder, name string) string {
	return "./components/" + folder + "/" + name + ".html"
}

func main() {
	resumeHandler := handler.BuildHandlerResume(ErrorPath, GetHTML)

	app := fiber.New()

	app.Static("/", "./public")

	app.Post("/cv", resumeHandler.CreateResume)

	app.Get("/cv", resumeHandler.GetExample)

	app.Get("/:folder/:name", func(c *fiber.Ctx) error {
		name := c.Params("name", DefaultComponent)
		folder := c.Params("folder", DefaultComponent)
		htmlElement, err := GetHTML(buildComponentPath(folder, name), &Component{Title: name, ID: name})
		if err != nil {
			return resumeHandler.HandleError(c, err)
		}
		return c.SendString(htmlElement)
	})

	log.Fatal(app.Listen(":3000"))
}
