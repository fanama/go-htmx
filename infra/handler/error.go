package handler

import (
	"fanama/go-htmx/domain/entity"

	"github.com/gofiber/fiber/v2"
)

type ErrorData = entity.ErrorData

type HandlerError struct {
	ErrorPath string
	GetHTML   func(path string, data any) (string, error)
}

func BuildHandlerError(errorPath string, getHTML func(path string, data any) (string, error)) HandlerError {
	return HandlerError{ErrorPath: errorPath, GetHTML: getHTML}
}

// Handle error responses
func (this *HandlerError) HandleError(c *fiber.Ctx, err error) error {

	errorPath := this.ErrorPath
	data := ErrorData{Message: err.Error()}

	errorHTML, _ := this.GetHTML(errorPath, data)
	return c.SendString(errorHTML)
}
