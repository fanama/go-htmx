package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Component struct {
	Title     string
	ClassName string
	ID        string
	Content   string
}

func GetHTML(path string, data any) (string, error) {
	tmpl, _ := template.ParseFiles(path)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {

		fmt.Println(err)
		return "", err
	}
	return buf.String(), nil
}

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Delete("/", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/list/:repeat/:folder/:name", func(c *fiber.Ctx) error {
		name := c.Params("name", "body")
		folder := c.Params("folder", "body")
		repeat, _ := strconv.Atoi(c.Params("repeat", "0"))
		path := "./components/" + folder + "/" + name + ".html"

		list := []string{}

		for i := 0; i < repeat; i++ {
			page := &Component{Title: name, ID: name + strconv.Itoa(i)}
			htmlElement, _ := GetHTML(path, page)
			list = append(list, htmlElement)
		}

		return c.SendString(strings.Join(list, ""))
	})

	app.Get("/:folder/:name", func(c *fiber.Ctx) error {
		name := c.Params("name", "body")
		folder := c.Params("folder", "body")
		path := "./components/" + folder + "/" + name + ".html"

		page := &Component{Title: name, ID: name}

		htmlElement, err := GetHTML(path, page)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.SendString(htmlElement)
	})
	log.Fatal(app.Listen(":3000"))
}
