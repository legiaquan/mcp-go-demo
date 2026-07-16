package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-go-demo/internal/models"
)

type TaskTools struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func (t *TaskTools) Register(server *mcp.Server) {
	server.AddTool(mcp.NewTool("add_task",
		mcp.WithDescription("Add a new task to the database"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Title of the task")),
		mcp.WithString("status", mcp.Description("Status of the task (default: pending)")),
	), t.addTask)

	server.AddTool(mcp.NewTool("update_task",
		mcp.WithDescription("Update an existing task in the database"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID of the task to update")),
		mcp.WithString("title", mcp.Description("New title of the task")),
		mcp.WithString("status", mcp.Description("New status of the task")),
	), t.updateTask)

	server.AddTool(mcp.NewTool("delete_task",
		mcp.WithDescription("Delete a task from the database"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID of the task to delete")),
	), t.deleteTask)

	server.AddTool(mcp.NewTool("get_task",
		mcp.WithDescription("Get details of a specific task"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID of the task to retrieve")),
	), t.getTask)

	server.AddTool(mcp.NewTool("list_tasks",
		mcp.WithDescription("List all tasks in the database"),
	), t.listTasks)
}

func (t *TaskTools) addTask(ctx context.Context, req *mcp.CallToolRequest, input models.AddTaskInput) (*mcp.CallToolResult, any, error) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := t.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "add_task"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	status := input.Status
	if status == "" {
		status = "pending"
	}

	var id int
	reqLogger.Debug("Executing DB insert query")
	err := t.DB.QueryRowContext(ctx, "INSERT INTO tasks (title, status) VALUES ($1, $2) RETURNING id", input.Title, status).Scan(&id)
	if err != nil {
		reqLogger.Error("Database execution failed", slog.String("error", err.Error()), slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error adding task: %v", err)}},
		}, nil, nil
	}

	reqLogger.Info("Tool execution completed", slog.Int("inserted_task_id", id), slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task added successfully. ID: %d", id)}},
	}, nil, nil
}

func (t *TaskTools) updateTask(ctx context.Context, req *mcp.CallToolRequest, input models.UpdateTaskInput) (*mcp.CallToolResult, any, error) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := t.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "update_task"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	setClauses := []string{}
	args := []any{}
	argId := 1

	if input.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argId))
		args = append(args, input.Title)
		argId++
	}
	if input.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argId))
		args = append(args, input.Status)
		argId++
	}

	if len(setClauses) == 0 {
		reqLogger.Warn("No fields to update provided")
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "No fields to update provided."}},
		}, nil, nil
	}

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argId)
	args = append(args, input.ID)

	reqLogger.Debug("Executing DB update query", slog.String("query", query))
	res, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		reqLogger.Error("Database execution failed", slog.String("error", err.Error()), slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error updating task: %v", err)}},
		}, nil, nil
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		reqLogger.Warn("Task not found", slog.Int("task_id", input.ID))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task with ID %d not found.", input.ID)}},
		}, nil, nil
	}

	reqLogger.Info("Tool execution completed", slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task %d updated successfully.", input.ID)}},
	}, nil, nil
}

func (t *TaskTools) deleteTask(ctx context.Context, req *mcp.CallToolRequest, input models.DeleteTaskInput) (*mcp.CallToolResult, any, error) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := t.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "delete_task"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	reqLogger.Debug("Executing DB delete query")
	res, err := t.DB.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", input.ID)
	if err != nil {
		reqLogger.Error("Database execution failed", slog.String("error", err.Error()), slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error deleting task: %v", err)}},
		}, nil, nil
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		reqLogger.Warn("Task not found", slog.Int("task_id", input.ID))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task with ID %d not found.", input.ID)}},
		}, nil, nil
	}

	reqLogger.Info("Tool execution completed", slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task %d deleted successfully.", input.ID)}},
	}, nil, nil
}

func (t *TaskTools) getTask(ctx context.Context, req *mcp.CallToolRequest, input models.GetTaskInput) (*mcp.CallToolResult, any, error) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := t.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "get_task"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	var task models.Task
	reqLogger.Debug("Executing DB select query")
	err := t.DB.QueryRowContext(ctx, "SELECT id, title, status FROM tasks WHERE id = $1", input.ID).Scan(&task.ID, &task.Title, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			reqLogger.Warn("Task not found", slog.Int("task_id", input.ID))
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Task with ID %d not found.", input.ID)}},
			}, nil, nil
		}
		reqLogger.Error("Database execution failed", slog.String("error", err.Error()), slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching task: %v", err)}},
		}, nil, nil
	}

	data, _ := json.MarshalIndent(task, "", "  ")

	reqLogger.Info("Tool execution completed", slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

func (t *TaskTools) listTasks(ctx context.Context, req *mcp.CallToolRequest, input models.ListTasksInput) (*mcp.CallToolResult, any, error) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := t.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "list_tasks"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	reqLogger.Debug("Executing DB select all query")
	rows, err := t.DB.QueryContext(ctx, "SELECT id, title, status FROM tasks ORDER BY id ASC")
	if err != nil {
		reqLogger.Error("Database execution failed", slog.String("error", err.Error()), slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching tasks: %v", err)}},
		}, nil, nil
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Status); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		reqLogger.Info("Tool execution completed (no tasks)", slog.Duration("latency", time.Since(start)))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}},
		}, nil, nil
	}

	data, _ := json.MarshalIndent(tasks, "", "  ")

	reqLogger.Info("Tool execution completed", slog.Int("count", len(tasks)), slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
