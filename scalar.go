package scalar

import (
	"encoding/json"
	"fmt"
	"strings"
)

func safeJSONConfiguration(options *Options) string {
	// Serializes the options to JSON
	jsonData, _ := json.Marshal(options)
	// Escapes double quotes into HTML entities
	escapedJSON := strings.ReplaceAll(string(jsonData), `"`, `&quot;`)
	return escapedJSON
}

func specContentHandler(specContent interface{}) string {
	switch spec := specContent.(type) {
	case func() map[string]interface{}:
		// If specContent is a function, it calls the function and serializes the return
		result := spec()
		jsonData, _ := json.Marshal(result)
		return string(jsonData)
	case map[string]interface{}:
		// If specContent is a map, it serializes it directly
		jsonData, _ := json.Marshal(spec)
		return string(jsonData)
	case string:
		// If it is a string, it returns directly
		return spec
	default:
		// Otherwise, returns empty
		return ""
	}
}

func GetScalarHTMLContent(options *Options) (string, error) {
	if options.SpecURL == "" && options.SpecContent == nil {
		return "", fmt.Errorf("specURL or specContent must be provided")
	}

	if options.SpecContent == nil && options.SpecURL != "" {
		if strings.HasPrefix(options.SpecURL, "http") {
			content, err := fetchContentFromURL(options.SpecURL)
			if err != nil {
				return "", err
			}
			options.SpecContent = content
		} else {
			urlPath, err := ensureFileURL(options.SpecURL)
			if err != nil {
				return "", err
			}

			content, err := readFileFromURL(urlPath)
			if err != nil {
				return "", err
			}

			options.SpecContent = string(content)
		}
	}

	return specContentHandler(options.SpecContent), nil
}

// GetScalarCDN returns the html script given an option.CDN
func GetScalarCDN(options *Options) string {
	return fmt.Sprintf(`<script src="%s"></script>`, options.CDN)
}

func GetCustomCSS(options *Options) string {
	var css string
	if options.Theme != "" {
		return ""
	}

	if options.CustomCss != "" {
		css = options.CustomCss
	} else {
		css = CustomThemeCSS
	}

	return fmt.Sprintf("<style>%s</style>", css)
}

func GetScalarScriptWithHTMLContent(config, specContent string) string {
	return fmt.Sprintf(`<script id="api-reference" type="application/json" data-configuration="%s">%s</script>`, config, specContent)
}

func ApiReferenceHTML(optionsInput *Options) (string, error) {
	options := DefaultOptions(*optionsInput)

	specContentHTML, err := GetScalarHTMLContent(options)
	if err != nil {
		return "", err
	}

	var pageTitle string
	if options.PageTitle != "" {
		pageTitle = options.PageTitle
	} else {
		pageTitle = "Scalar API Reference"
	}

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <title>%s</title>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
	%s
      </head>
      <body>
	%s
	%s
      </body>
    </html>
  `, pageTitle, GetCustomCSS(options), GetScalarScriptWithHTMLContent(safeJSONConfiguration(options), specContentHTML), GetScalarCDN(options)), nil
}
