package templates

import (
	"html/template"
	"log"
	"net/http"
	"path"

	packr "github.com/gobuffalo/packr/v2"
)

type GalleryTemplateContext struct {
	RelativePath string
}

type Template struct {
	*template.Template
}

func GetTemplate(name string) (*Template, error) {
	box := packr.New("web", "../../web")
	templateName := path.Base(name)
	templateContent, err := box.FindString(name)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return nil, err
	}
	return &Template{tmpl}, nil
}

func (t *Template) ServeTemplate(data interface{}, w http.ResponseWriter, r *http.Request) {
	err := t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
