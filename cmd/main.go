package main

import (
	"html/template"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var id = 1

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
	Id    int
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

func (d *Data) getIndexByID(id int) int {
	for idx, contact := range d.Contacts {
		if contact.Id == id {
			return idx
		}
	}
	return -1
}

func newStore() *Store {
	return &Store{
		Data: Data{
			Contacts: Contacts{
				{1, "AR", "ar@email.com"},
				{2, "AA", "aa@email.com"},
			},
		},
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.Renderer = newTemplate()

	store := newStore()

	e.GET("/", func(c echo.Context) error {
		store.Form = newFormData()
		return c.Render(http.StatusOK, "index", store)
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

		contact := Contact{id, name, email}
		id++

		store.Data.Contacts = append(store.Data.Contacts, contact)

		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(time.Second * 2)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}

		contactIdx := store.Data.getIndexByID(id)

		if contactIdx == -1 {
			return c.String(http.StatusNotFound, "Contact not found")
		}

		store.Data.Contacts = slices.Delete(store.Data.Contacts, contactIdx, contactIdx+1)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
