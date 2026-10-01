package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	// "google.golang.org/protobuf/types/descriptorpb"
	_ "modernc.org/sqlite"
	// "flag"
)


	var dbPath = "./tasks.db"

type Task struct {
	ID  int
	Title string
	Description string
	Status string
	Created_at time.Time
	Updated_at time.Time
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}

createTableSQL := `
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT,
	status TEXT NOT NULL DEFAULT  "pending",
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		fmt.Println("Error creating table:", err)
		os.Exit(1)
	}

	return db
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println(" go run . add [title] [description]")
		fmt.Println("go run . delete [task id]")
		fmt.Println("go run . update_status [task_id] [new_status]")
		fmt.Println("go run . edit_task [task_id] <args>")
		fmt.Println("go run . list_all [task_status] optional")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add" :
		if len (os.Args) < 3 {
			fmt.Println("Task must have a title")
			os.Exit(1)
		}

		title := os.Args[2]
		description := "" 
		if len(os.Args) >= 4 {
			description = os.Args[3]
		} 
		
		AddTask(title, description)

	case "list" :
		status := ""
		if len (os.Args) >= 3 {
			status = os.Args[2]
		}
		ListTasks(status)

	case "update_status" :
		if len(os.Args) < 4 {
			fmt.Println("update_status requires a task id and a status update")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid input for ID: '%s'\n Values must be integer", err)
		}

		newStatus := os.Args[3]
		UpdateTaskStatus(id, newStatus)

		case "delete" :
		if len(os.Args) <3 {
			fmt.Println("delete requires a task id")
			os.Exit(1)
		}	

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid input for ID: '%d'\n Values must be in integer", id)
		}

		DeleteTask(id)
	}




	// fmt.Println("Select your next choice: ")
}