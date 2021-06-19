package i18nrewrite

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

// Config the plugin configuration.
type Config struct {
	LangCodes                 []string `json:"langCodes"`
	LangCodeMatchedHeader     string   `json:"langCodeMatchedHeader,omitempty"`
	LangCodeMatchedConfidence string   `json:"langCodeMatchedConfidence,omitempty"`
	DefaultLangRewrite        bool     `json:"defaultLangRewrite,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		LangCodes:                 []string{},
		LangCodeMatchedHeader:     "X-Matched-Lang",
		LangCodeMatchedConfidence: "X-Matched-Lang-Confidence",
		DefaultLangRewrite:        false,
	}
}

type I18nRewrite struct {
	next               http.Handler
	matcher            language.Matcher
	langTags           []language.Tag
	defaultLangRewrite bool
	langCodeHeader     string
	langCodeConfidence string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.LangCodes) == 0 {
		return nil, fmt.Errorf("langCodes cannot be empty")
	}

	langTags := make([]language.Tag, len(config.LangCodes))

	for i, code := range config.LangCodes {
		langTags[i] = language.Make(code)
	}

	return &I18nRewrite{
		next:               next,
		matcher:            language.NewMatcher(langTags),
		langTags:           langTags,
		defaultLangRewrite: config.DefaultLangRewrite,
		langCodeHeader:     config.LangCodeMatchedHeader,
		langCodeConfidence: config.LangCodeMatchedConfidence,
	}, nil
}

func (i *I18nRewrite) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		_, index := language.MatchStrings(i.matcher, req.Header.Get("Accept-Language"))
		lang := i.langTags[index]
		_, confidence := lang.Base()

		// not equal to default lang
		if lang != i.langTags[0] {
			req.URL.Path = fmt.Sprintf("%s/%s", req.URL.Path, strings.ToLower(lang.String()))
		} else if lang == i.langTags[0] && i.defaultLangRewrite {
			req.URL.Path = fmt.Sprintf("%s/%s", req.URL.Path, strings.ToLower(lang.String()))
		}
		req.Header.Set(i.langCodeHeader, lang.String())
		req.Header.Set(i.langCodeConfidence, confidence.String())
	}

	i.next.ServeHTTP(rw, req)
}
