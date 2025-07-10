package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"website/templates"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
	"github.com/rickb777/servefiles/v3"

	"website/src/models"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	ctx := r.Context()
	_ = component.Render(ctx, w)
}

func handleDynamic(w http.ResponseWriter, r *http.Request) {
	// Get all known dynamic files for navigation purposes
	folder, err := models.FileTree("public")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// MD from static
	// TODO: check if HTML file exists
	mdPath := filepath.Join("public", "example.md")
	md, err := os.ReadFile(mdPath)
	if err != nil {
		log.Fatal(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	renderedBytes := markdown.Render(doc, renderer)

	// Handle $$ inline latex
	// Replaces $$...$$ with $<div class="inline-latex-block">...</div>$
	parsedBytes := make([]byte, 0)
	inInlineLatex := false
	prevWasDollarSign := false
	for _, b := range renderedBytes {
		if b == '$' && !prevWasDollarSign {
			prevWasDollarSign = true
		}

		if b != '$' && prevWasDollarSign {
			prevWasDollarSign = false
		} else if prevWasDollarSign && !inInlineLatex {
			inInlineLatex = true
			startInlineBlock := "<div class=\"inline-latex-block\">"
			parsedBytes = append(parsedBytes, []byte(startInlineBlock)...)
			continue
		} else if prevWasDollarSign && inInlineLatex {
			inInlineLatex = false
			endInlineBlock := "</div>"
			parsedBytes = append(parsedBytes, []byte(endInlineBlock)...)
			continue
		}

		parsedBytes = append(parsedBytes, b)
	}

	// Save the HTML
	err = os.WriteFile(strings.Replace(mdPath, ".md", "", 1)+".html", parsedBytes, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	// HTML from static
	resource := strings.Replace(r.PathValue("resource"), ".md", "", 1)

	path := filepath.Join("public", resource+".html")
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	content := string(contentBytes)

	component := templates.Page(folder, content)
	ctx := r.Context()
	_ = component.Render(ctx, w)
}

func login(w http.ResponseWriter, r *http.Request) {

	sessionId := uuid.New().String()

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	})
}

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)
		log.Println(wrappedWriter.statusCode, r.Method, r.URL, time.Since(start))
	})
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session")
		if err != nil || sessionToken.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", handleIndex)
	router.HandleFunc("GET /page/{resource}", handleDynamic)
	static := servefiles.NewAssetHandler("./static/").WithMaxAge(time.Second) // todo: different time on deploy, ex hour
	router.Handle("GET /static/", http.StripPrefix("/static/", static))

	stack := CreateStack(Logging)

	server := &http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Server running on port :8080")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
