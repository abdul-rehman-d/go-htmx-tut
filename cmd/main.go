package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Data struct {
	Contacts Contacts
}

type Store struct {
	Data Data
	Form FormData
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newStore() *Store {
	return &Store{
		Data: Data{
			Contacts: Contacts{
				{"AR", "ar@email.com"},
				{"AA", "aa@email.com"},
			},
		},
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	store := newStore()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", store)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if store.Data.hasEmail(email) {
			store.Form = newFormData()
			store.Form.Values["name"] = name
			store.Form.Values["email"] = email
			store.Form.Errors["email"] = "Email already exists."

			return c.Render(http.StatusUnprocessableEntity, "form", store.Form)
		}

		contact := Contact{name, email}

		store.Data.Contacts = append(store.Data.Contacts, contact)

		c.Render(200, "form", newFormData())
		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
