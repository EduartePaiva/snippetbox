package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (h *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		h.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
	}
}

func (h *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	numId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || numId < 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("Display a specific snippet with ID %d...", numId)))
}

func (h *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return

	}
	w.Write([]byte("Create a new snippet..."))
}
