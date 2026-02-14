# Step 3 – HTTP server and GET all notes

## Goal

Start an HTTP server and respond to **GET /notes** with all notes as JSON.

---

## Theory

### 1. `http.Handler` (interface)

In Go’s `net/http` package, a **handler** is anything that can serve an HTTP request.

- **Type:** `interface` with one method:
  ```go
  type Handler interface {
      ServeHTTP(ResponseWriter, *Request)
  }
  ```
- **Meaning:** If your type has a `ServeHTTP(w http.ResponseWriter, r *http.Request)` method, it implements `http.Handler` and can be used to handle HTTP requests.
- **Your job in the handler:** read `r` (method, path, headers, body) and write the response using `w` (status, headers, body).

**Examples of handlers:**

- A struct that implements `ServeHTTP`:
  ```go
  type MyHandler struct{}
  func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "Hello")
  }
  ```
- A function wrapped with `http.HandlerFunc`:
  ```go
  func myHandler(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "Hello")
  }
  // http.HandlerFunc(myHandler) is an http.Handler
  ```
- The **default router** Go uses when you call `http.HandleFunc` or `http.Handle` (see below): that router is also an `http.Handler`.

---

### 2. `http.DefaultServeMux`

- **What it is:** The **default multiplexer** (router) provided by `net/http`. It is a global variable of type `*http.ServeMux`.
- **Role:** It is an `http.Handler` that **dispatches** requests to different handlers based on the **URL path** (and optionally method if you check it yourself).
- **How you use it:** You don’t create it; you **register** handlers on it:
  - `http.Handle(pattern string, handler http.Handler)` — register an `http.Handler` for a path.
  - `http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))` — register a function for a path; Go wraps it so it becomes an `http.Handler`.
- **Important:** When you call `http.ListenAndServe(addr, nil)`, passing `nil` as the handler means “use `http.DefaultServeMux`”. So any route you registered with `http.Handle` or `http.HandleFunc` is used.

**Example:**

```go
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello")
})
http.HandleFunc("/notes", notesHandler)
http.ListenAndServe(":8080", nil)  // nil = use DefaultServeMux
```

Here, requests to `/notes` are sent to `notesHandler`.

---

### 3. `http.ListenAndServe`

- **Signature:** `func ListenAndServe(addr string, handler Handler) error`
- **Meaning:**
  - **addr:** e.g. `":8080"` — start a TCP server on port 8080 (all interfaces).
  - **handler:** The top-level `http.Handler` that receives every request. If it is `nil`, Go uses `http.DefaultServeMux`, so you can register routes with `http.Handle` / `http.HandleFunc` and they will be used.
- **Behavior:** The call is **blocking**: it runs until the process exits or the server fails. It accepts connections and, for each request, calls `handler.ServeHTTP(w, r)` (or the mux’s `ServeHTTP`, which then calls the right registered handler).

**Examples:**

```go
// Use default mux (register with http.HandleFunc, then):
http.ListenAndServe(":8080", nil)

// Use your own handler (e.g. your own mux or single handler):
mux := http.NewServeMux()
mux.HandleFunc("/notes", notesHandler)
http.ListenAndServe(":8080", mux)
```

---

## Quick reference

| Concept                | Role                                                                 |
|------------------------|----------------------------------------------------------------------|
| `http.Handler`         | Interface: something that can serve a request via `ServeHTTP`.       |
| `http.DefaultServeMux` | Default router; use `http.Handle` / `http.HandleFunc` to register.   |
| `http.ListenAndServe`  | Starts the server; pass `nil` to use `DefaultServeMux`.              |

---

## Step-by-step: what you need to do (no code, just steps)

Do these in order; you write the code yourself.

### 1. Start the HTTP server in `main`

- Call `http.ListenAndServe` with an address (e.g. `":8080"`) and a handler.
- If you use `http.DefaultServeMux`, pass `nil` as the second argument.
- Register your route **before** calling `ListenAndServe`, otherwise the server starts with no route for `/notes`.

### 2. Register a handler for `/notes`

- Use either:
  - `http.HandleFunc("/notes", yourHandlerFunction)`, or  
  - `http.Handle("/notes", yourHandler)` if `yourHandler` implements `http.Handler`.
- Decide if you also want to handle `/notes/` (with trailing slash); the default mux treats `/notes` and `/notes/` as different unless you handle both or redirect.

### 3. In the handler: allow only GET

- In your handler, read the method: `r.Method`.
- If it is not `"GET"`, set the status to `405` (e.g. `w.WriteHeader(http.StatusMethodNotAllowed)`) and return. Optionally write a short body or leave it empty.

### 4. Build the list of notes and encode JSON

- You have `store = map[string]Note{}`. For “all notes” you need a **slice** (e.g. `[]Note`) to get a JSON array.
- Loop over `store` and collect all `Note` values into a slice.
- Use `json.Marshal(yourSlice)` to get `[]byte` and handle the error.
- Set the response header: `w.Header().Set("Content-Type", "application/json")`.
- Write the status (e.g. `200`) if you didn’t already, then write the body: `w.Write(data)`.

### 5. Run and test

- Run the program. It should listen on port 8080.
- In another terminal: `curl http://localhost:8080/notes` (or open that URL in a browser).
- You should see:
  - `Content-Type: application/json`
  - Body: a JSON array of notes (e.g. `[{"ID":1,"Title":"First Note","Content":"Just a note"}]`), or `[]` if the store is empty.

---

## Checklist (for you to verify)

- [ ] Server listens on a port (e.g. 8080).
- [ ] GET /notes returns JSON with `Content-Type: application/json`.
- [ ] Response body is the list of notes (or empty array `[]`).

---

## Hints (without giving code away)

- Import `encoding/json` for `json.Marshal`.
- Set `Content-Type` **before** writing the body.
- For 405, call `w.WriteHeader` (or use `w.WriteHeader` before any write); default is 200 if you don’t set another status.
- Your `Note` struct will marshal to JSON using the field names (`ID`, `Title`, `Content`) unless you use struct tags. That’s enough for this step.

Once you’re done, paste your handler code and how you registered it so it can be reviewed.
