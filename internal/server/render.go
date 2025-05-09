package server

import (
  "html/template"
  "io"

  "github.com/labstack/echo/v4"
)

type Template struct {
  templates *template.Template
}

func NewTemplate(dir string) *Template {
  return &Template{
    templates: template.Must(template.ParseGlob(dir + "/*.html")),
  }
}

// Реализуем интерфейс echo.Renderer
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}