package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	ID         int    `json:"id" bson:"id"`
	Title      string `json:"title" bson:"title"`
	Note       string `json:"note" bson:"note"`
	DueDate    string `json:"due_date,omitempty" bson:"due_date,omitempty"`
	CreatedAt  string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	IsComplete bool   `json:"is_complete" bson:"is_complete"`
}

type Handlers struct {
	Coll *mongo.Collection
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (h Handlers) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cur, err := h.Coll.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer cur.Close(ctx)

	var todos []Todo
	if err := cur.All(ctx, &todos); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, 200, todos)
}

func (h Handlers) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var todo Todo
	err = h.Coll.FindOne(ctx, bson.M{"id": id}).Decode(&todo)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, 200, todo)
}

func (h Handlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if todo.CreatedAt == "" {
		todo.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.Coll.InsertOne(ctx, todo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, 201, todo)
}

func (h Handlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	todo.ID = id // enforce path id

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.Coll.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": todo},
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, 200, todo)
}
