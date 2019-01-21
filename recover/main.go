package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// listenAndServeFunction - using function as variable
var listenAndServeFunction = http.ListenAndServe

// FunctionCopy - using function as variable
var FunctionCopy = io.Copy

// FunctionConversion - using function as variable
var FunctionConversion = strconv.Atoi

// FunctionStyle - using function as variable
var FunctionStyle = styles.Get

// main This function is used to call the handler
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/debug/", sourceCodeHandler)
	log.Fatal(listenAndServeFunction(":4000", devMw(mux)))
}

// sourceCodeHandler This function is used highlight the links & line numbers
func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	lineStr := r.FormValue("line")
	line, err := FunctionConversion(lineStr) // converts string to integer
	if err != nil {
		line = -1
	}
	file, err := os.Open(path) // open the file i.e path stores the url
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bytes.NewBuffer(nil)
	_, err = FunctionCopy(b, file)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, b.String())
	style := FunctionStyle("github")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(), html.LineNumbersInTable(), html.HighlightLines(lines))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<style>pre { font-size: 1.2em; }</style>")
	formatter.Format(w, style, iterator)

}

// devMw This is an anonymous function
func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//log.Println(err)
				stack := debug.Stack()
				//log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

// panicDemo This is just the demo of how the panic function calls
func panicDemo(w http.ResponseWriter, r *http.Request) {
	functionThatPanics()
}

// funcThatPanics This function called whenever any panic occurs
func functionThatPanics() {
	panic("Panic called!!")
}

// hello Sample function to display Hello text
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

// makeLinks This function adds the hyperlink to URL so if we click on the same it will redirect to that page
func makeLinks(stack string) string {

	lines := strings.Split(stack, "\n")

	for li, line := range lines {

		if len(line) == 0 || line[0] != '\t' {
			continue
		}

		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}
		var lineStr strings.Builder
		// + 2 --  1 for first tab & 1 for : which is present at last
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineStr.WriteByte(line[i])
		}

		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineStr.String())

		lines[li] = "\t<a href=\"/debug?" + v.Encode() + "\">" + file + ":" + lineStr.String() + "</a>" + line[len(file)+2+len(lineStr.String()):]
	}
	return strings.Join(lines, "\n")
}
