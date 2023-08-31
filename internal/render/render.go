package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ddytert/bookings/internal/config"
	"github.com/ddytert/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer Sets the config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	// Get the template cache from app config
	var tc map[string]*template.Template
	if app.UseCache {
		//log.Println("Rendering template from cache", tmpl)
		tc = app.TemplateCache
	} else {
		log.Println("New rendering of template", tmpl)
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println("can't get template from cache")
		return errors.New("can't get template from cache")
	}
	// render the template
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser")
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from .templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
