// Package web showcases Go's built-in net/http package for building web servers.
//
// For a Java developer:
// - Go's `net/http` package is similar to the Java Servlet API or Spring MVC.
// - `http.Handler` is like a `Servlet`. It has a `ServeHTTP` method.
// - Middleware in Go is a common pattern, similar to `Filter` in Java.
// - JSON marshaling is built-in with `encoding/json` (similar to Jackson/Gson).
// - There's no separate Tomcat/Jetty server; the binary IS the server.
package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ResponseData represents a generic JSON response structure.
//
// For a Java developer:
// - Similar to a Generic Response wrapper used in Spring controllers.
type ResponseData struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// Task represents a simple to-do item.
//
// For a Java developer:
// - Similar to a POJO or a Record in Java.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TaskStore is an in-memory storage for tasks, demonstrating thread safety.
//
// For a Java developer:
// - Similar to a Repository or Service in Spring.
// - It uses a `sync.RWMutex` to protect the internal slice, making it thread-safe.
// - `RWMutex` is like `ReentrantReadWriteLock` in Java.
type TaskStore struct {
	sync.RWMutex
	tasks  []Task
	nextID int
}

// List returns a copy of all tasks in the store.
//
// For a Java developer:
// - Java equivalent: `return new ArrayList<>(tasks)` (with a synchronized block).
func (s *TaskStore) List() []Task {
	s.RLock()
	defer s.RUnlock()
	return append([]Task{}, s.tasks...)
}

// Get retrieves a specific task by its ID.
//
// For a Java developer:
// - Java equivalent: `return tasks.stream().filter(t -> t.getId() == id).findFirst()`
func (s *TaskStore) Get(id int) (Task, bool) {
	s.RLock()
	defer s.RUnlock()
	for _, t := range s.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}

// Add creates a new task with the given title and adds it to the store.
//
// For a Java developer:
// - Java equivalent: `public Task save(String title)`
func (s *TaskStore) Add(title string) Task {
	s.Lock()
	defer s.Unlock()
	s.nextID++
	t := Task{ID: s.nextID, Title: title, Completed: false}
	s.tasks = append(s.tasks, t)
	return t
}

// Toggle flips the completion status of a task.
//
// For a Java developer:
// - Java equivalent: `public boolean toggle(int id)`
func (s *TaskStore) Toggle(id int) bool {
	s.Lock()
	defer s.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Completed = !s.tasks[i].Completed
			return true
		}
	}
	return false
}

var (
	store = &TaskStore{
		tasks: []Task{
			{ID: 1, Title: "Learn Go Generics", Completed: true},
			{ID: 2, Title: "Explore Go 1.23 Features", Completed: false},
		},
		nextID: 2,
	}

	// TaskListTemplate for rendering HTML.
	taskListTemplate = template.Must(template.New("tasks").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Go Tasks</title>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 40px auto; line-height: 1.6; }
        .completed { text-decoration: line-through; color: #888; }
        form { margin-top: 20px; }
    </style>
</head>
<body>
    <h1>Task List</h1>
    <ul>
    {{range .}}
        <li class="{{if .Completed}}completed{{end}}">
            <strong>#{{.ID}}</strong>: {{.Title}}
            <form action="/tasks/{{.ID}}/toggle" method="POST" style="display:inline;">
                <button type="submit">Toggle</button>
            </form>
        </li>
    {{else}}
        <li>No tasks yet!</li>
    {{end}}
    </ul>

    <hr>
    <h3>Add New Task</h3>
    <form action="/tasks" method="POST">
        <input type="text" name="title" placeholder="What needs to be done?" required>
        <button type="submit">Add Task</button>
    </form>
</body>
</html>
`))
)

// Middleware defines a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain combines multiple middlewares into a single one.
//
// For a Java developer:
// - This is like chaining `Filter` instances in a `FilterChain`.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// NewWebServerHandler creates a new http.Handler with all routes and middlewares configured.
//
// For a Java developer:
// - Go's `net/http` is a production-grade web server (like Tomcat/Jetty combined with Spring MVC).
// - `http.ServeMux` is the router (equivalent to Spring's `DispatcherServlet` or `@RequestMapping`).
// - Handlers are functions or structs that satisfy the `http.Handler` interface.
// - Middleware is a common pattern in Go: a function that takes a Handler and returns a Handler.
func NewWebServerHandler() http.Handler {
	// Custom ServeMux for route definition
	mux := http.NewServeMux()

	// Registering handlers
	// Java equivalent: @GetMapping("/hello")
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Gopher"
		}
		fmt.Fprintf(w, "Hello, %s! Current URL: %s", name, r.URL.Path)
	})

	// Java equivalent: Returning a POJO from a @RestController.
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := ResponseData{
			Message: "Welcome to the Go Demo Web Server!",
			Time:    time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Java equivalent: @PostMapping("/echo") with @RequestBody Map<String, Object>.
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	})

	// --- Advanced Handlers (Go 1.22+ routing) ---

	// GET /tasks - List all tasks (HTML or JSON)
	// Java equivalent: A Spring Controller with Content Negotiation.
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks := store.List()
		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
			return
		}
		taskListTemplate.Execute(w, tasks)
	})

	// POST /tasks - Create a new task
	// Java equivalent: @PostMapping("/tasks") handling form data.
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		if title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		store.Add(title)
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	})

	// POST /tasks/{id}/toggle - Toggle task completion
	// Java equivalent: @PostMapping("/tasks/{id}/toggle") with @PathVariable.
	mux.HandleFunc("POST /tasks/{id}/toggle", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if !store.Toggle(id) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	})

	// GET /tasks/{id} - Get a task as JSON
	// Java equivalent: @GetMapping("/tasks/{id}") returning a Task object.
	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, _ := strconv.Atoi(idStr)
		task, ok := store.Get(id)
		if !ok {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	})

	// GET /long-op - Demonstrates context timeout
	// Java equivalent: Async Servlet or a Controller returning a CompletableFuture.
	mux.HandleFunc("GET /long-op", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println("   [Server] Long operation started...")

		select {
		case <-time.After(2 * time.Second):
			// Simulate work
			fmt.Fprintf(w, "Operation finished successfully after 2 seconds!")
		case <-ctx.Done():
			// Honoring client cancellation or server timeout
			fmt.Printf("   [Server] Operation cancelled: %v\n", ctx.Err())
			// Note: We can't really send a response if the client closed the connection,
			// but this shows how to detect it.
		}
	})

	// Wrap everything in middleware
	return Chain(mux,
		recoveryMiddleware,
		loggingMiddleware,
	)
}

// StartWebServer starts a simple HTTP server on the given port.
//
// For a Java developer:
// - No Tomcat or Jetty needed. `http.ListenAndServe` is the production-ready server.
func StartWebServer(port string) {
	fmt.Printf("--- Web Server Demo (starting on :%s) ---\n", port)
	fmt.Printf("   Visit http://localhost:%s/hello or http://localhost:%s/tasks\n", port, port)

	handler := NewWebServerHandler()

	// Use a goroutine to start the server so it doesn't block the main thread forever in a demo.
	go func() {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			fmt.Printf("Web server stopped: %v\n", err)
		}
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Println("   Web server is running in the background.")
}

// loggingMiddleware logs details about each incoming HTTP request.
//
// For a Java developer:
// - Similar to a Logback AccessLog or a Spring RequestLoggingFilter.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Internal print for demo visibility
		fmt.Printf("   [Server Log] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		fmt.Printf("   [Server Log] Completed in %v\n", time.Since(start))
	})
}

// recoveryMiddleware catches any panics during request handling to prevent the server from crashing.
//
// For a Java developer:
// - Similar to a global ExceptionHandler in Spring.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("   [Recovery Log] Recovered from panic: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
