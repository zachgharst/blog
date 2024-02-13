package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/zachgharst/blog/internal/markdown"
)

type Blog struct {
	Url     string
	Title   string
	Content string
}

var blogs []Blog

//	have link to all blogs (in some semblance) on this page
//	make sure route is / in URL when page loaded
func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET getAllBlogs\n")

	homepageTemplate, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		panic("Problem reading homepage template")
	}

    if err := homepageTemplate.Execute(w, blogs); err != nil {
        panic(fmt.Sprintf("Problem using template file: %s", err))
    }
}

// todo:
//
//	proper 404 if blog is not found
//	for SEO: should be /blog/[some string identifying the title of the blog]
//	don't allow FS access via route, names should match to a file via GUID
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
    err := preloadBlogs()
    if err != nil {
        fmt.Printf("Problem preloading blogs: %s\n", err)
        return
    }

	http.HandleFunc("/", getAllBlogs)
	http.HandleFunc("/blog/", getBlog)

	static := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	fmt.Println("Starting web server on http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func preloadBlogs() error {
	blogFiles, err := os.ReadDir("blogs")
	if err != nil {
		return err
	}

	for _, blog := range blogFiles {
        fileContent, err := os.ReadFile("blogs/" + blog.Name())
        if err != nil {
            return err
        }
        content := string(markdown.ToHtml(fileContent))

        title := strings.Split(string(fileContent), "\n")[0][2:]
        newBlog := Blog{
            Url: blog.Name()[:len(blog.Name())-3],
            Title: title,
            Content: content,
        }

        blogs = append(blogs, newBlog)
	}

	return nil
}
