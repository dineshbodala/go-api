package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

var tasks []Task
var currentID int

type App struct {
	Router *mux.Router
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/tasks", app.getTasks).Methods("GET")
	app.Router.HandleFunc("/task/{id}", app.readTask).Methods("GET")
	app.Router.HandleFunc("/task", app.createTask).Methods("POST")
	app.Router.HandleFunc("/task/{id}", app.updateTask).Methods("PUT")
	app.Router.HandleFunc("/task/{id}", app.deleteTask).Methods("DELETE")
}

func (app *App) Initialise(initialTasks []Task, id int) {
	tasks = initialTasks
	currentID = id
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
}
func main() {
	app := App{}
	
	tasks, id := CreateInitialTasks()
	fmt.Println(tasks)
	app.Initialise(tasks, id)
	app.Run("localhost:10000")
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)	

}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func (app *App) getTasks(writer http.ResponseWriter, request *http.Request) {
	tasks, err := getTasks()
	if err != nil{
		sendError(writer, http.StatusBadRequest, err.Error())
		return  
	}
	sendResponse(writer, http.StatusOK, tasks)
}

func (app *App) createTask(writer http.ResponseWriter, r *http.Request) {
    var p Task

    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        sendError(writer, http.StatusBadRequest, "Invalid request payload")
        return
    }
    err = p.createTask()
    if err != nil {
        sendError(writer, http.StatusInternalServerError, err.Error())
        return
    }
    sendResponse(writer, http.StatusCreated, p)
}

func (app *App) readTask(writer http.ResponseWriter, request *http.Request) {
	// fmt.Println(tasks)
	vars := mux.Vars(request)
    key, _ := strconv.Atoi(vars["id"])
	t := Task{ID: key}
	err := t.getTask()
	if err != nil {
		sendError(writer, http.StatusBadRequest, err.Error())
	}
	sendResponse(writer, http.StatusOK, t)



}

func (app *App) updateTask(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
    key, err:= strconv.Atoi(vars["id"])
	t := Task{ID: key}
    if err != nil {
        sendError(writer, http.StatusBadRequest, "Invalid product ID")
        return
    }
	err = json.NewDecoder(request.Body).Decode(&t)
	if err != nil {
        fmt.Println(err)
        sendError(writer, http.StatusBadRequest, "Invalid request body")
        return
    }
	err = t.updateTask()
	if err != nil {
        sendError(writer, http.StatusInternalServerError, err.Error())
        return
    }

    sendResponse(writer, http.StatusOK, t)
}
	// your code goes here




func (app *App) deleteTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
    key, _ := strconv.Atoi(vars["id"])
	t := Task{ID: key}
	err := t.deleteTask()
	if err != nil {
		sendError(writer, http.StatusBadRequest, err.Error())
	}
	sendResponse(writer, http.StatusOK, t)



}

	





