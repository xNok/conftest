package document

import (
	"embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed resources/*
var resources embed.FS

// TemplateKind helps us to select where to find the template. It can either be embedded or on the host filesystem
type TemplateKind int

const ( // iota is reset to 0
	FS   TemplateKind = iota
	FSYS              // fsys is used for embedded templates
)

type TemplateConfig struct {
	kind TemplateKind
	path string
}

func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		kind: FSYS,
		path: "resources/document.md",
	}
}

type RenderDocumentOption func(*TemplateConfig)

// WithTemplate is a functional option to override the documentation template
func WithTemplate(tpl string) RenderDocumentOption {
	return func(c *TemplateConfig) {
		c.kind = FS
		c.path = tpl
	}
}

// RenderDocument takes a slice of Section and generate the markdown documentation either using the default
// embedded template or the user provided template
func RenderDocument(out io.Writer, s []Section, opts ...RenderDocumentOption) error {
	var tpl = NewTemplateConfig()

	// Apply all the functional options to the template configurations
	for _, opt := range opts {
		opt(tpl)
	}

	err := renderTemplate(tpl, s, out)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(tpl *TemplateConfig, args interface{}, out io.Writer) error {
	var t *template.Template
	var err error

	switch tpl.kind {
	case FSYS:
		// read the embedded template
		t, err = template.ParseFS(resources, tpl.path)
		if err != nil {
			return err
		}
	case FS:
		t, err = template.ParseFiles(tpl.path)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown template kind: %v", tpl.kind)
	}

	// we render the template
	err = t.Execute(out, args)
	if err != nil {
		return err
	}

	return nil
}