package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func main() {
	router := chi.NewRouter()

	router.Get("/tasks", getTasks)

	router.Get("/task/{id}", getTask)

	router.Post("/tasks", createTasks)

	router.Delete("/task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

func getTasks(response http.ResponseWriter, request *http.Request) {
	var jsonTasks, marshalError = json.Marshal(tasks)

	if marshalError != nil {
		http.Error(response, marshalError.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(jsonTasks)
}

func getTask(response http.ResponseWriter, request *http.Request) {
	var taskID = chi.URLParam(request, "id")

	var task, isExist = tasks[taskID]

	if !isExist {
		http.Error(response, "There is no task with id "+taskID, http.StatusBadRequest)
		return
	}

	var jsonTask, marshalError = json.Marshal(task)

	if marshalError != nil {
		http.Error(response, marshalError.Error(), http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonTask)
}

func createTasks(response http.ResponseWriter, request *http.Request) {
	var task Task
	var buffer bytes.Buffer

	var _, readError = buffer.ReadFrom(request.Body)

	if readError != nil {
		http.Error(response, readError.Error(), http.StatusBadRequest)
		return
	}

	var unmarshalError = json.Unmarshal(buffer.Bytes(), &task)

	if unmarshalError != nil {
		http.Error(response, unmarshalError.Error(), http.StatusBadRequest)
		return
	}

	if _, isExist := tasks[task.ID]; isExist {
		http.Error(response, "There is already such a task in the task list", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
}

func deleteTask(response http.ResponseWriter, request *http.Request) {
	var taskID = chi.URLParam(request, "id")

	var _, isExist = tasks[taskID]

	if !isExist {
		http.Error(response, "There is no task with id "+taskID, http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}
