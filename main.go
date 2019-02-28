package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var defaultHexValue string
var iconDirectory string
var allIcons icondict
var customConfig config

const xmlHeading = `<?xml version="1.0" standalone="no"?><!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">`

type config struct {
	DefaultColor   string `yaml:"default_color"`
	ServiceName    string `yaml:"service_name"`
	ServiceTagline string `yaml:"service_tagline"`
	IconDirectory  string `yaml:"icon_directory"`
}

type icon struct {
	Name   string
	Source string
}

type icondict struct {
	Name    string
	Tagline string
	Icons   map[string]icon
}

func handler(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path[1:]) > 0 {
		// Serve an icon
		params, _ := url.ParseQuery(r.URL.RawQuery)

		// Setup defaults
		iconHex := customConfig.DefaultColor

		if len(params["hex"]) > 0 {
			iconHex = params["hex"][0]
		}

		requestedIconName := r.URL.Path[1:]

		if _, ok := allIcons.Icons[requestedIconName]; ok {
			svgContent := allIcons.Icons[requestedIconName].Source
			svgContent = strings.Replace(svgContent, "<path", "<path style='fill:#"+iconHex+"'", 1)

			print("--> Serving: " + r.URL.Path[1:] + "\n")
			w.Header().Set("Content-Type", "image/svg+xml")
			fmt.Fprint(w, xmlHeading+svgContent)
		} else {
			// File doesn't exist
			print("--> 404: " + r.URL.Path[1:] + "\n")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Icon doesn't exist")
		}
	} else {
		// Render an index page
		dat, err := ioutil.ReadFile("./docs/index.html")
		check(err)
		indexContent := string(dat)

		// Write the headers out
		w.Header().Set("Content-Type", "text/html")

		// Write the full page content out
		fmt.Fprint(w, indexContent)
	}
}

func main() {
	loadConfig()
	port := ":" + os.Getenv("PORT")

	if port == ":" {
		port = "localhost:8080"
	}

	readIconsIntoMemory()
	buildSite()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	print("--> Listening on http://" + port + "\n")
	log.Fatal(http.ListenAndServe(port, mux))
}

func readIconsIntoMemory() {
	// Read every single malibu icon into memory
	files, err := ioutil.ReadDir(customConfig.IconDirectory)
	allIcons.Icons = make(map[string]icon)
	check(err)

	for _, f := range files {
		// Only match .svg
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".svg") {
			var thisIcon icon
			dat, err := ioutil.ReadFile(customConfig.IconDirectory + f.Name())
			check(err)

			thisIcon.Source = string(dat)
			thisIcon.Name = strings.TrimSuffix(f.Name(), ".svg")

			allIcons.Icons[thisIcon.Name] = thisIcon
		}
	}

	print("Icons read into memory...\n")
}

func buildSite() {
	print("Building documentation...\n")

	err := ioutil.WriteFile("./docs/index.html", []byte(buildDocumentationWithTemplate()), 0644)
	check(err)
}

func buildDocumentationWithTemplate() string {
	t, _ := template.ParseFiles("./templates/index.html")
	var tpl bytes.Buffer
	err := t.Execute(&tpl, allIcons)
	check(err)

	result := tpl.String()
	return result
}

func loadConfig() {
	filename, _ := filepath.Abs("./royo_config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config config

	err = yaml.Unmarshal(yamlFile, &config)
	check(err)

	customConfig = config
	allIcons.Name = config.ServiceName
	allIcons.Tagline = config.ServiceTagline
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
