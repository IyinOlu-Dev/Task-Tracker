# Task CLI

A lightweight command-line task management application written in Go, powered by SQLite.

---

## Features

* **Create Tasks:** Add tasks with a mandatory title and optional description.


* **List Tasks:** View all tasks or filter them by status (`Todo`, `In progress`, `Done`).


* **Update Status:** Transition task statuses between `Todo`, `In Progress`, and `Done`.


* **Delete Tasks:** Remove tasks by their ID.


* **Pure Go SQLite:** Uses [`modernc.org/sqlite`](https://gitlab.com/cznic/sqlite), requiring no CGo or external C compilers.



---

## Prerequisites

* [Go](https://go.dev/dl/) (1.18 or higher recommended)

---

## Installation & Setup

1. **Clone the repository:**
```bash
git clone https://github.com/<your-username>/<repo-name>.git
cd <repo-name>

```


2. **Download dependencies:**
```bash
go mod tidy

```


3. **Build the binary (optional):**
```bash
go build -o task-cli .

```



---

## Usage

You can run commands using `go run . <command>` or via the compiled binary `./task-cli <command>`.

### 1. Add a Task

Add a new task by specifying a title and an optional description:

```bash
go run . add "Buy groceries" "Milk, eggs, and bread"

```

```bash
go run . add "Read documentation"

```

*Newly created tasks default to the `Todo` status.*

### 2. List Tasks

List all tasks stored in the database:

```bash
go run . list

```

Filter tasks by status (`Todo`, `In progress`, `Done`):

```bash
go run . list "Todo"

```

```bash
go run . list "In progress"

```

```bash
go run . list "Done"
```

### 3. Update Task Status

Update the status of an existing task using its ID:

```bash
go run . update_status <task_id> <new_status>

```

*Example:*

```bash
go run . update_status 1 "In Progress"
go run . update_status 1 "Done"

```

### 4. Delete a Task

Delete a task by ID:

```bash
go run . delete <task_id>

```

*Example:*

```bash
go run . delete 1

```

---

## Command Reference

| Command | Arguments | Description |
| --- | --- | --- |
| `add` | `<title> [description]` | Add a new task

 |
| `list` | `[status]` | List all tasks or filter by status

 |
| `update_status` | `<task_id> <status>` | Update status (`Todo`, `In Progress`, `Done`)

 |
| `delete` | `<task_id>` | Remove a task by ID

 |

---

## Database Schema

The application automatically creates a local SQLite database file named `tasks.db` on first run with the following schema:

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT "Todo",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

```

*(Note: Add `.tasks.db` to your `.gitignore` to avoid checking the local database into version control).*

---

## Running Tests

Run the unit test suite:

```bash
go test -v ./...

```
