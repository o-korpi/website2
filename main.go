package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"website/src"
	"website/templates"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
	"github.com/rickb777/servefiles/v3"

	"website/src/models"
	p_ "website/src/parser"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	ctx := r.Context()
	_ = component.Render(ctx, w)
}

func handleArticles(w http.ResponseWriter, r *http.Request) {
	// Get all known articles for navigation purposes
	folder, err := models.FileTree("public", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Render the articles page with the folder structure
	component := templates.Articles(folder)
	ctx := r.Context()
	_ = component.Render(ctx, w)
}

type CustomRenderer struct {
	*html.Renderer
}

func (r *CustomRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.CodeBlock:
		if entering {
			// Custom code block rendering logic
			lang := string(node.Info)
			code := string(node.Literal)

			var lexer chroma.Lexer
			if lang == "" {
				lexer = lexers.Analyse(code)
			} else {
				lexer = lexers.Get(lang)
			}

			if lexer == nil {
				lexer = lexers.Fallback
			}

			style := styles.Get(styles.CatppuccinFrappe.Name) // Default style
			if style == nil {
				style = styles.Fallback
			}
			formatter := formatters.Get("html")
			if formatter == nil {
				formatter = formatters.Fallback
			}
			reader := strings.NewReader(code)
			contents, _ := io.ReadAll(reader)
			iterator, _ := lexer.Tokenise(nil, string(contents))

			var formattedCode string
			formattedCodeWriter := &strings.Builder{}
			err := formatter.Format(formattedCodeWriter, style, iterator)
			if err != nil {
				log.Println("Error formatting code:", err)
				return ast.GoToNext
			}
			formattedCode = formattedCodeWriter.String()

			// Postprocess code block to remove body tags, html tags
			codeBlock := strings.ReplaceAll(formattedCode, "<body>", "")
			codeBlock = strings.ReplaceAll(codeBlock, "</body>", "")
			codeBlock = strings.ReplaceAll(codeBlock, "<body class=\"bg\">", "")
			codeBlock = strings.ReplaceAll(codeBlock, "<html>", "")
			codeBlock = strings.ReplaceAll(codeBlock, "</html>", "")
			// Remove unnecessary styles
			bgLine := ""
			bodyLine := ""
			lines := strings.SplitSeq(codeBlock, "\n")
			for line := range lines {
				if strings.Contains(line, "/* Background */") {
					bgLine = line
				} else if strings.Contains(line, "body {") { // todo: use more robust method
					bodyLine = line
				}
			}

			codeBlock = strings.ReplaceAll(codeBlock, bgLine, "")
			codeBlock = strings.ReplaceAll(codeBlock, bodyLine, "")

			// fmt.Fprintf(w, `<div class="mockup-code bg-[#303446]">`)
			fmt.Fprintf(w, `<div>`)
			fmt.Fprintf(w, "%s", codeBlock)
			fmt.Fprintf(w, `</div>`)
			return ast.GoToNext
		}
		return ast.GoToNext
	default:
		return r.Renderer.RenderNode(w, node, entering)
	}
}

func handleDynamic(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")

	// Get all known dynamic files for navigation purposes
	folder, err := models.FileTree("public", resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// MD from static
	// TODO: check if HTML file exists
	mdPath := filepath.Join("public", resource+".md")
	md, err := os.ReadFile(mdPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Remove frontmatter before parsing
	mdNoFrontmatter := src.RemoveFrontmatter(string(md))

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(mdNoFrontmatter))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	customRenderer := &CustomRenderer{Renderer: renderer}

	renderedBytes := markdown.Render(doc, customRenderer)

	// Handle sidenotes
	renderedBytes = p_.ProcessSidenotes(renderedBytes)

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
	resource = strings.Replace(r.PathValue("resource"), ".md", "", 1)

	path := filepath.Join("public", resource+".html")
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	content := string(contentBytes)

	// Fancy breadcrumbs stuff
	splitResource := strings.Split(resource, "/")
	// for i, part := range splitResource {  // todo: consider removing this
	// 	if i != len(splitResource)-1 {
	// 		splitResource[i] = strings.ToUpper(part[:1]) + part[1:] // Capitalize first letter
	// 	}
	// }

	// Frontmatter
	frontmatter, err := src.ScanFrontmatter(mdPath)
	if err != nil {
		log.Fatal("Error scanning frontmatter:", err)
	}
	parsedFm, err := src.ParseFrontmatter(frontmatter)
	if err != nil {
		log.Println("Error parsing frontmatter:", err)
	}

	// TODO: table of contents

	component := templates.Page(folder, splitResource, *parsedFm, content)
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
	router.HandleFunc("GET /articles", handleArticles)
	router.HandleFunc("GET /page/{resource...}", handleDynamic)
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
