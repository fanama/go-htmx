package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Headder struct {
	Title string
}

type User struct {
	Username string
	Password string
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

	app.Post("/api/login", func(c *fiber.Ctx) error {

		user := User{Username: c.FormValue("username"), Password: c.FormValue("password")}

		htmlElement, _ := GetHTML("./components/forms/user.html", user)

		return c.SendString(htmlElement)
	})

	app.Static("/", "./public")
	app.Static("/components", "./components")

	app.Delete("/", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/components/:name", func(c *fiber.Ctx) error {
		name := c.Params("name", "body")
		path := "./components/" + name + ".html"

		page := &Headder{Title: name}

		htmlElement, err := GetHTML(path, page)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.SendString(htmlElement)
	})
	log.Fatal(app.Listen(":3000"))
}
