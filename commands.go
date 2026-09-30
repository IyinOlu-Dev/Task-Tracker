package main

import (
	"fmt"
	"os"
	"time"
	"database/sql"
)

// AddTask adds a new task to the database
func AddTask(title string, description string) {

	TITLE := title
	DESCRIPTION := description
	STATUS := "Todo"
	CREATED_AT := time.Now()
	UPDATED_AT := time.Now()

	if TITLE == "" {
		fmt.Println("Task must have a title.")
		os.Exit(1)
	}

	db := initDB()
	defer db.Close()

	createTaskSQL := `
	INSERT INTO TASKS (TITLE, DESCRIPTION, STATUS, CREATED_AT, UPDATED_AT) VALUES (?,?,?,?,?)
	`

	statement, err := db.Prepare(createTaskSQL)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		os.Exit(1)
	}
	defer statement.Close()

	result, err := statement.Exec(TITLE, DESCRIPTION, STATUS, CREATED_AT, UPDATED_AT)
	if err != nil {
		fmt.Println("Error inserting task:", err)
		os.Exit(1)
	}
	fmt.Println("Results are: ", result)

	fmt.Println("Task added successfully!")
}

// ListTasks lists all tasks with the given status
func ListTasks(status string) {
	db := initDB()
	defer db.Close()

	QUERY := "SELECT * FROM tasks"
	var args [] any

	if status != "" {
		if status =="Todo" || status =="In progress" || status == "Done"{
			QUERY += " WHERE STATUS = ?"
			args = append(args, status)
		} else {
			fmt.Printf("Invalid status: %s\n Valid Status are 'Todo', 'In progress', 'Done'", status)
		}
	}
	rows, err := db.Query(QUERY, args...)
	if rows.Err()!= nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No tasks found with status '%s'\n", status)
			return
		}
		fmt.Println("error querying tasks database: ", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("List of tasks:")

	for rows.Next() {

		var t Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at)

		if err != nil {
			fmt.Println("error scanning row: ", err)
		}
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %s, Created At: %s, Updated At: %s\n", t.ID, t.Title, t.Description, t.Status, t.Created_at.Format(time.RFC3339), t.Updated_at.Format(time.RFC3339))
	}

	if err := rows.Err(); err != nil {
	fmt.Println("Error iterating over rows:", err)
	os.Exit(1)
	}

}
// UpdateTaskStatus updates the status of a task with the given ID
func UpdateTaskStatus(id int, newStatus string) {
	db := initDB()
	defer db.Close()

	var currentStatus string
	err := db.QueryRow("SELECT status FROM tasks WHERE id = ?", id).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No task found with ID %d\n", id)
			return
		}
		fmt.Println("Error querying task:", err)
		os.Exit(1)
	}

	updateTaskSQL := `
	UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?
	`

	statement, err := db.Prepare(updateTaskSQL)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		os.Exit(1)
	}
	defer statement.Close()

	if newStatus == currentStatus {
		fmt.Printf("Task with ID %d already has status '%s'\n", id, newStatus)
		return
	}

	if newStatus != "Todo" && newStatus != "In Progress" && newStatus != "Done" {
		fmt.Printf("Invalid status '%s'. Valid statuses are: Todo, In Progress, Done\n", newStatus)
		return
	}

	result, err := statement.Exec(newStatus, time.Now(), id)
	if err != nil {
		fmt.Println("Error updating task:", err)
		os.Exit(1)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		os.Exit(1)
	}

	if rowsAffected == 0 {
		fmt.Printf("No task found with ID %d\n", id)
	} else {
		fmt.Printf("Task with ID: %d updated successfully from '%s' to status '%s'\n", id,currentStatus, newStatus)
	}
}

func DeleteTask(id int) {
	db := initDB()
	defer db.Close()

	deleteTaskSQL :=  `DELETE FROM tasks WHERE id = ?`

	statement, err := db.Prepare(deleteTaskSQL)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		os.Exit(1)
	}
	defer statement.Close()

	result, err := statement.Exec(id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		os.Exit(1)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		os.Exit(1)
	}

	if rowsAffected == 0 {
		fmt.Printf("No task found with ID %d\n", id)
	} else {
		fmt.Printf("Task with ID %d deleted successfully\n", id)
	}
}

func EditTask(id int, newTitle string, newDescription string, newStatus string, ) {
	db := initDB()
	defer db.Close()

	currentInfo := `SELECT TITLE, DESCRIPTION, STATUS FROM TASKS WHERE ID = ?`

	var currentTitle, currentDescription, currentStatus string

	err := db.QueryRow(currentInfo, id).Scan(&currentTitle, &currentDescription, &currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No task with ID %d found", id)
		}
		fmt.Println("Error querying tasks database", err)
	}

	if newTitle == "" || newTitle == currentTitle{
		newTitle = currentTitle
		fmt.Println("No changes to Task Title were made")
	}
	if newDescription == "" || newDescription == currentDescription {
		newDescription = currentDescription
		fmt.Println("No changes to Task Description were made")
	}
	if newStatus == "" || newStatus == currentStatus {
		newStatus = currentStatus
		fmt.Println("No changes to Task Status were made")
	}

	updateTaskSQL := `
	UPDATE tasks SET title = ?, description = ?, status = ?, updated_at = ? WHERE id = ?
	`

	statement, err := db.Prepare(updateTaskSQL)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		os.Exit(1)
	}
	defer statement.Close()

	result, err := statement.Exec(newTitle, newDescription, newStatus, time.Now(), id)
	if err != nil {
		fmt.Println("Error updating task:", err)
		os.Exit(1)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		os.Exit(1)
	}

	if rowsAffected == 0 {
		fmt.Printf("No task found with ID %d\n", id)
	} else {
		fmt.Printf("Task with ID %d updated successfully\n", id)
	}
}