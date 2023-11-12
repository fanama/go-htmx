package handler

import (
	"fanama/go-htmx/domain/entity"
	"html"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Component = entity.Component

type SectionData = entity.SectionData

type HandlerResume struct {
	HandlerError
}

const (
	HeaderPath       = "components/molecules/headder.html"
	ExperiencePath   = "components/atoms/experience.html"
	ErrorPath        = "components/atoms/error.html"
	SectionPath      = "components/molecules/section.html"
	DefaultComponent = "body"
)

func BuildHandlerResume(errorPath string, getHTML func(path string, data interface{}) (string, error)) HandlerResume {
	return HandlerResume{HandlerError{ErrorPath: errorPath, GetHTML: getHTML}}
}

func (this *HandlerResume) GetResume(resume entity.Resume) (string, error) {

	headerHTML, err := this.GetHTML(HeaderPath, &resume)
	if err != nil {
		return err.Error(), err
	}

	var sectionsHTML strings.Builder
	for _, section := range resume.Sections {
		var renderedExperienceList strings.Builder

		for _, experience := range section.Experiences {
			experienceHTML, err := this.GetHTML(ExperiencePath, experience)
			if err != nil {
				return err.Error(), err
			}
			renderedExperienceList.WriteString(html.UnescapeString(experienceHTML))
		}

		sectionData := SectionData{
			Title:   section.Title,
			Content: renderedExperienceList.String(),
		}

		sectionHTML, err := this.GetHTML(SectionPath, sectionData)
		if err != nil {
			return err.Error(), err
		}
		sectionsHTML.WriteString(html.UnescapeString(sectionHTML))
	}

	return headerHTML + "<div id=\"resume\" class=\"grid grid-cols-2\">" + sectionsHTML.String() + "</div>", nil

}

func (this *HandlerResume) CreateResume(c *fiber.Ctx) error {

	var resume entity.Resume

	err := c.BodyParser(&resume)

	if err != nil {
		return this.HandleError(c, err)
	}

	result, err := this.GetResume(resume)

	if err != nil {
		return this.HandleError(c, err)
	}

	return c.Type("html").SendString(result)
}

func (this *HandlerResume) GetExample(c *fiber.Ctx) error {

	resume, err := entity.ReadResumeFromFile("./example.json")

	if err != nil {
		return this.HandleError(c, err)

	}

	result, err := this.GetResume(*resume)

	if err != nil {
		return this.HandleError(c, err)
	}

	return c.Type("html").SendString(result)
}
