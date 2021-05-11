package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ninostephen/bookings/pkg/config"
	"github.com/ninostephen/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("couldn't not get template from template cache")
	}

	td = AddDefaultData(td)

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to browser:", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// pages hold the absolute path to all page.tmpl files
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// page will have the absolute path to the page template.
	for _, page := range pages {
		// name stores only the file name

		name := filepath.Base(page)

		// ts is the template set which holds the newly created template with custom functions (funcMap) attached
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// matches hold all the layout.tmpl files
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Println(myCache)
			return myCache, err
		}

		// matches will have a length greater than 0 if layout files exist
		if len(matches) > 0 {

			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {

				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
