package main

import (
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB sets a temporary database path for testing and cleans up afterwards.
func setupTestDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath = filepath.Join(tempDir, "test_tasks.db")
}

func TestAddTask(t *testing.T) {
	setupTestDB(t)

	AddTask("Test Task 1", "This is a test description")

	db := initDB()
	defer db.Close()

	var title, description, status string
	err := db.QueryRow("SELECT title, description, status FROM tasks WHERE id = 1").Scan(&title, &description, &status)
	if err != nil {
		t.Fatalf("Failed to query inserted task: %v", err)
	}

	if title != "Test Task 1" {
		t.Errorf("Expected title 'Test Task 1', got '%s'", title)
	}
	if description != "This is a test description" {
		t.Errorf("Expected description 'This is a test description', got '%s'", description)
	}
	if status != "Todo" {
		t.Errorf("Expected default status 'Todo', got '%s'", status)
	}
}

func TestListTasks(t *testing.T) {
	setupTestDB(t)

	AddTask("Task Todo", "Desc 1")
	AddTask("Task In Progress", "Desc 2")
	UpdateTaskStatus(2, "In Progress")

	t.Run("List All Tasks", func(t *testing.T) {
		ListTasks("")
	})

	t.Run("List By Status Todo", func(t *testing.T) {
		ListTasks("Todo")
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	setupTestDB(t)

	AddTask("Task to Update", "Description")
	UpdateTaskStatus(1, "Done")

	db := initDB()
	defer db.Close()

	var status string
	err := db.QueryRow("SELECT status FROM tasks WHERE id = 1").Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query task status: %v", err)
	}

	if status != "Done" {
		t.Errorf("Expected status 'Done', got '%s'", status)
	}
}

func TestEditTask(t *testing.T) {
	setupTestDB(t)

	AddTask("Original Title", "Original Desc")
	EditTask(1, "Updated Title", "Updated Desc", "In Progress")

	db := initDB()
	defer db.Close()

	var title, description, status string
	err := db.QueryRow("SELECT title, description, status FROM tasks WHERE id = 1").Scan(&title, &description, &status)
	if err != nil {
		t.Fatalf("Failed to query edited task: %v", err)
	}

	if title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", title)
	}
	if description != "Updated Desc" {
		t.Errorf("Expected description 'Updated Desc', got '%s'", description)
	}
	if status != "In Progress" {
		t.Errorf("Expected status 'In Progress', got '%s'", status)
	}
}

func TestDeleteTask(t *testing.T) {
	setupTestDB(t)

	AddTask("Task to Delete", "Description")
	DeleteTask(1)

	db := initDB()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = 1").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected task to be deleted, but it still exists in the database")
	}
}