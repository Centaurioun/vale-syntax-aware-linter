package lint

import (
	"strings"

	"github.com/errata-ai/vale/v2/internal/core"
	"github.com/jdkato/regexp"
	"github.com/niklasfasching/go-org/org"
)

var orgConverter = org.New()
var orgWriter = org.NewHTMLWriter()

var reOrgAttribute = regexp.MustCompile(`(#(?:\+| )[^\s]+:.+)`)
var reOrgProps = regexp.MustCompile(`(:PROPERTIES:\n.+\n:END:)`)
var reOrgSrc = regexp.MustCompile(`(?i)#\+BEGIN_SRC .+`)

func (l Linter) lintOrg(f *core.File) error {
	s := reOrgAttribute.ReplaceAllString(f.Content, "\n=$1=\n")
	s = reOrgProps.ReplaceAllString(s, "\n#+BEGIN_EXAMPLE\n$1\n#+END_EXAMPLE\n")

	// We don't want to find matches in `begin_src` lines.
	body := reOrgSrc.ReplaceAllStringFunc(f.Content, func(m string) string {
		return strings.Repeat("*", len(m))
	})

	doc := orgConverter.Parse(strings.NewReader(s), f.Path)
	// We don't want to introduce any *new* content into our HTML,
	// so we clear the outline.
	doc.Outline.Children = nil

	html, err := doc.Write(orgWriter)
	if err != nil {
		return err
	}

	f.Content = body
	return l.lintHTMLTokens(f, []byte(html), 0)
}
