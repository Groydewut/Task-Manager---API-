package TaskController

import (
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request)

func TaskList(w http.ResponseWriter, r *http.Request)

func TaskListByID(w http.ResponseWriter, r *http.Request)

func UpdateTask(w http.ResponseWriter, r *http.Request)

func DeleateTask(w http.ResponseWriter, r *http.Request)
