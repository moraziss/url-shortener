package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

// MapHandler returns an http.HandlerFunc that redirects (302 Found) any
// request whose path is a key of pathsToUrls to the matching URL. Requests
// for paths that are not in the map are passed on to fallback.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler parses yamlBytes, a list of path/url pairs such as
//
//	- path: /docs
//	  url: https://example.com/docs
//
// and returns a handler that redirects like MapHandler, falling back to
// fallback for unknown paths. It returns an error if the YAML cannot be
// parsed into that shape.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.Handler, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
