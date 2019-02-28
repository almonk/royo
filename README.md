# Royo

### What is Royo?
Royo is a general purpose Go app to serve SVG iconsets to your webapp or website. It has an extremely simple API that allows users to request a specific icon and, optionally, in a specific hex value.

Check out a sample Royo install using the excellent Unicorns icon set:
[https://serene-wildwood-49104.herokuapp.com/](https://serene-wildwood-49104.herokuapp.com/)

### What isn't Royo?
* Royo isn't an icon pack. 
* Royo only serves monochromatic SVGs for now and will probably have unintended results if used otherwise

---

## Getting Started with development

### Dependencies
* Go 1.11.5+
* [Glide Package Manage for Go](https://github.com/Masterminds/glide)

### Installation

```console
git clone https://github.com/almonk/royo
cd royo/
glide install
go run royo.go
```

Royo will start automatically on port `8080` unless a `PORT` environment variable is set otherwise

### Customising
Royo is configurable via `royo_config.yaml`:

```yaml
# Name your icon service
service_name: "Unicorn Icon Set"
service_tagline: "A beautiful collection of open source icons"

# What hex value do you want to use if none is specified?
default_color: "0070D2"

# Where do your SVGs live?
icon_directory: "./imports/unicorns/"
```

Remember to restart the Go server after changing any configuration.

### Customising documentation
Royo produces its own documentation which is served as the index page. Customising the documentation can be done by editing `./templates/index.html`.

In this template the following variables can be used;
* `{{ .Icons }}` is a map of all the Icons being served
  * Within the range;
  * `{{ .Name }}` is the addressable name of a single icon
  * `{{ .Source }}` is the raw SVG content
* `{{ .Name }}` is the name of the iconset


---

## Deployment

Royo ships as a Go binary, so you don't need to install Go or build anything to use it as-is. You can simply;
* Add the icons you want to serve to `./imports/`
* Edit `royo_config.yml` to reflect your customisations

Then, find somewhere to deploy it!

#### Heroku
Royo comes with Heroku support out of the box. You just need an app to push to:

```console
heroku create --buildpack https://github.com/ph3nx/heroku-binary-buildpack.git
heroku config:set PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/app/bin
git push heroku master
```

### Todo
* Add documentation for other hosting providers
* Add CDN documentation