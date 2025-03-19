# Knovel Assessment Project

**Author:** Ha Phong (siransbach)

This project is built using Go and Docker. It includes a PostgreSQL database with seeded data and a set of scripts to
manage the Docker environment.

## Seeded Users

The following users are seeded in the database:

| Username  | Password  | Role     |
|-----------|-----------|----------|
| Radahn    | password1 | EMPLOYEE |
| Malenia   | password2 | EMPLOYEE |
| Tarnished | password3 | EMPLOYER |

## API Definition

### Employee API

#### Get All Tasks

- **Endpoint:** `/api/v1/employee/tasks`
- **Method:** `GET`
- **Description:** Retrieves a list of all tasks assigned to the employee.

#### Update Task Status

- **Endpoint:** `/api/v1/employee/tasks/{id}/status/{status}`
- **Method:** `PUT`
- **Description:** Update status for the employee's task.
- **Path Parameters:**
    - `id`: Task ID.
    - `status`: New status for the task. Possible values: `PENDING`, `IN_PROGRESS`, `COMPLETED`.

### Employer API

#### Create Task

- **Endpoint:** `/api/v1/employer/tasks`
- **Method:** `POST`
- **Description:** Creates a new task and assigns it to an employee.
- **Request Body:**
  ```json
  {
    "title": "string",
    "description": "string",
    "assigned_user_id": "integer",
    "due_date": "string" // Format: RFC3339
  }
  ```

#### Get Tasks

- **Endpoint:** `/api/v1/employer/tasks`
- **Method:** `GET`
- **Description:** Retrieves a list of all tasks.
- **Query Parameters:**
    - `status`: Filter tasks by status. Possible values: `PENDING`, `IN_PROGRESS`, `COMPLETED`.
    - `assignedUserId`: Filter tasks by assigned user ID.
    - `sortBy`: Sort tasks by any field. Possible values: `id`, `title`, `description`, `due_date`, `status`,
      `created_at`, `assigned_user_id`.
    - `sortOrder`: Sort order. Possible values: `asc`, `desc`.

#### Get Task Summary

- **Endpoint:** `/api/v1/employer/tasks/summary`
- **Method:** `GET`
- **Description:** Retrieves a summary of tasks grouped by employees, showing total number of tasks assigned and
  completed.

## Running the Project

### Prerequisites

- Docker
- Docker Compose
- Go 1.24 (optional, if you want to run the project without Docker)
- PostgreSQL client (optional, if you want to connect to the database)

### Running with Docker

1. Run the following command to start the Docker environment:

   ```bash
   $ cd scripts && ./boot.sh
   ```

2. The API will be available at `http://localhost:8000`.
3. To stop the Docker environment, run:

   ```bash
   $ cd scripts && ./shutdown.sh
   ```
4. When encountering permission denied error, run the following commands

   ```bash
   $ chmod +x scripts/boot.sh
   $ chmod +x scripts/shutdown.sh
   ```
   
   

