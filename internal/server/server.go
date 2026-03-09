package server

import (
	"context"
	"log"
	"sync"

	"github.com/chirag-bruno/mock-grpc/pkg/todo"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServer struct {
	todo.UnimplementedTodoServiceServer
	mu    sync.RWMutex
	todos map[string]*todo.Todo
}

func NewTodoServer() *TodoServer {
	return &TodoServer{
		todos: make(map[string]*todo.Todo),
	}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *todo.CreateTodoRequest) (*todo.CreateTodoResponse, error) {
	log.Printf("[CreateTodo] Request: title=%q, description=%q", req.Title, req.Description)

	s.mu.Lock()
	defer s.mu.Unlock()

	t := &todo.Todo{
		Id:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}
	s.todos[t.Id] = t

	log.Printf("[CreateTodo] Created todo with id=%s", t.Id)
	return &todo.CreateTodoResponse{Todo: t}, nil
}

func (s *TodoServer) GetTodo(ctx context.Context, req *todo.GetTodoRequest) (*todo.GetTodoResponse, error) {
	log.Printf("[GetTodo] Request: id=%q", req.Id)

	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.todos[req.Id]
	if !ok {
		log.Printf("[GetTodo] Todo not found: id=%s", req.Id)
		return nil, status.Errorf(codes.NotFound, "todo not found: %s", req.Id)
	}

	log.Printf("[GetTodo] Found todo: id=%s, title=%q", t.Id, t.Title)
	return &todo.GetTodoResponse{Todo: t}, nil
}

func (s *TodoServer) ListTodos(ctx context.Context, req *todo.ListTodosRequest) (*todo.ListTodosResponse, error) {
	log.Printf("[ListTodos] Request received")

	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*todo.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}

	log.Printf("[ListTodos] Returning %d todos", len(todos))
	return &todo.ListTodosResponse{Todos: todos}, nil
}

func (s *TodoServer) UpdateTodo(ctx context.Context, req *todo.UpdateTodoRequest) (*todo.UpdateTodoResponse, error) {
	log.Printf("[UpdateTodo] Request: id=%q, title=%q, completed=%v", req.Id, req.Title, req.Completed)

	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.todos[req.Id]
	if !ok {
		log.Printf("[UpdateTodo] Todo not found: id=%s", req.Id)
		return nil, status.Errorf(codes.NotFound, "todo not found: %s", req.Id)
	}

	t.Title = req.Title
	t.Description = req.Description
	t.Completed = req.Completed

	log.Printf("[UpdateTodo] Updated todo: id=%s", t.Id)
	return &todo.UpdateTodoResponse{Todo: t}, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, req *todo.DeleteTodoRequest) (*todo.DeleteTodoResponse, error) {
	log.Printf("[DeleteTodo] Request: id=%q", req.Id)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[req.Id]; !ok {
		log.Printf("[DeleteTodo] Todo not found: id=%s", req.Id)
		return nil, status.Errorf(codes.NotFound, "todo not found: %s", req.Id)
	}

	delete(s.todos, req.Id)
	log.Printf("[DeleteTodo] Deleted todo: id=%s", req.Id)
	return &todo.DeleteTodoResponse{Success: true}, nil
}
