package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/zachgharst/blog/internal/markdown"
)

// todo:
//   have link to all blogs (in some semblance) on this page
//   make sure route is / in URL when page loaded
func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET getAllBlogs\n")

    homepageTemplate, err := template.ParseFiles("web/templates/index.html")

    if err != nil {
        panic("Problem reading homepage template")
    }

    if homepageTemplate.Execute(w, "") != nil {
        fmt.Fprintf(w, "Problem using template file")
    }
}

// todo:
//   proper 404 if blog is not found
//   for SEO: should be /blog/[some string identifying the title of the blog]
//   don't allow FS access via route, names should match to a file via GUID
func getBlog(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/blog/")
    fileName := fmt.Sprintf("blogs/%s.md", id)
    fmt.Printf("GET getBlog: [%s] with filename [%s]\n", id, fileName)

    data, err := os.ReadFile(fileName)
    if err != nil {
        fmt.Fprint(w, "404")
        fmt.Printf("Problem reading file %s\n", fileName)
        return
    }

    html := string(markdown.ToHtml(data))
    blogTemplate, err := template.ParseFiles("web/templates/blog.html")

    if err != nil {
        panic("Problem reading blog template")
    }

    if blogTemplate.Execute(w, template.HTML(html)) != nil {
        fmt.Fprintf(w, "Problem using template file")
    }
}

func main() {
	http.HandleFunc("/", getAllBlogs)
	http.HandleFunc("/blog/", getBlog)

    static := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", static))

    fmt.Println("Starting web server on https://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

