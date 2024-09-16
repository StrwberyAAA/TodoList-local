// Purpose: web app for managing ToDo items

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID     int
	Task   string
	Status bool
}

var (
	todos     []Todo
	todoID    int
	todoMutex sync.Mutex
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func handleTodoDisplay(w http.ResponseWriter, r *http.Request) {
	todoMutex.Lock()
	defer todoMutex.Unlock()

	if err := templates.ExecuteTemplate(w, "index.html", todos); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

func handleTodoCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		if task == "" {
			handleError(w, fmt.Errorf("task cannot be empty"), http.StatusBadRequest)
			return
		}

		todoMutex.Lock()
		todos = append(todos, Todo{ID: todoID, Task: task, Status: false})
		todoID++
		todoMutex.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleTodoStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleError(w, fmt.Errorf("invalid ID"), http.StatusBadRequest)
			return
		}

		todoMutex.Lock()
		defer todoMutex.Unlock()

		for i := range todos {
			if todos[i].ID == id {
				todos[i].Status = !todos[i].Status
				break
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleTodoClear(w http.ResponseWriter, r *http.Request) {
	todoMutex.Lock()
	defer todoMutex.Unlock()

	todos = []Todo{}
	todoID = 0

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// open browser automatically when program is ran
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	}
	if err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/", handleTodoDisplay).Methods("GET")
	router.HandleFunc("/add", handleTodoCreate).Methods("POST")
	router.HandleFunc("/toggle", handleTodoStatus).Methods("POST")
	router.HandleFunc("/clear", handleTodoClear).Methods("POST")

	log.Println("Starting server on :PUT UR PORT HERE") //EDIT THIS
	go openBrowser("http://localhost:YOUR PORT") //EDIT THIS
	log.Fatal(http.ListenAndServe(":8080", router))
}
