package i18nrewrite_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WLun001/i18nrewrite"
)

const headerAcceptLanguage = "Accept-Language"

func TestDefaultConfigEn(t *testing.T) {
	cfg := i18nrewrite.CreateConfig()
	cfg.LangCodes = []string{"en", "zh-cn"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := i18nrewrite.New(ctx, next, cfg, "i18n-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAcceptLanguage, "en-US,en;q=0.5")

	handler.ServeHTTP(recorder, req)

	assertPath(t, req, "")
	assertHeader(t, req, cfg.LangCodeMatchedHeader, "en")
	assertHeader(t, req, cfg.LangCodeMatchedConfidence, "Exact")
}

func TestDefaultConfigZh(t *testing.T) {
	cfg := i18nrewrite.CreateConfig()
	cfg.LangCodes = []string{"en-us", "zh-cn"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := i18nrewrite.New(ctx, next, cfg, "i18n-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAcceptLanguage, "zh,en-US;q=0.7,en;q=0.3")

	handler.ServeHTTP(recorder, req)

	assertPath(t, req, "/zh-cn")
	assertHeader(t, req, cfg.LangCodeMatchedHeader, "zh-CN")
	assertHeader(t, req, cfg.LangCodeMatchedConfidence, "Exact")
}

func TestSkipDefaultLang(t *testing.T) {
	cfg := i18nrewrite.CreateConfig()
	cfg.LangCodes = []string{"zh", "en"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := i18nrewrite.New(ctx, next, cfg, "i18n-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAcceptLanguage, "zh,en-US;q=0.7,en;q=0.3")

	handler.ServeHTTP(recorder, req)

	assertPath(t, req, "")
	assertHeader(t, req, cfg.LangCodeMatchedHeader, "zh")
	assertHeader(t, req, cfg.LangCodeMatchedConfidence, "Exact")
}

func assertPath(t *testing.T, req *http.Request, expected string) {
	t.Helper()

	if req.URL.Path != expected {
		t.Errorf("invalid path value: %s, expected: %s", req.URL.Path, expected)
	}
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s, expected: %s", req.Header.Get(key), expected)
	}
}
