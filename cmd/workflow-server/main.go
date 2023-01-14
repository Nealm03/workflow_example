package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"workflow.example.com/internal"

	"github.com/labstack/echo"
)

func main() {
	logger := log.Default()

	cfg, err := NewConfig()
	if err != nil {
		logger.Fatal("failed to read server config %w", err)
		return
	}
	workflowEngine, err := internal.NewWorkflowEngine(context.Background(), internal.WorkflowEngineConfig{
		Region:   cfg.AwsRegion,
		Endpoint: cfg.Endpoint,
	})

	if err != nil {
		logger.Fatal("failed to build workflow engine %w", err)
		return
	}

	server := echo.New()
	server.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, cfg)
	})

	server.POST("/workflow", func(ctx echo.Context) error {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, ctx.Request().Body)

		if err != nil {
			logger.Println(err)
			return ctx.JSON(http.StatusBadRequest, nil)
		}
		workflowId := uuid.NewString()
		jsonInput := buf.String()
		executionId, err := workflowEngine.StartWorkflow(ctx.Request().Context(), internal.StartWorkflowCommand{
			Input:        &jsonInput,
			WorkflowId:   &workflowId,
			WorkflowName: cfg.StateMachineARN,
		})
		if err != nil {
			logger.Println(err)
			return ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to start worklfow %w", err))
		}

		id := base64.URLEncoding.EncodeToString([]byte(executionId))

		ctx.Response().Header().Add("Location", "/workflow/"+id)
		return ctx.String(http.StatusAccepted, "")
	})

	server.GET("/workflow/:id", func(ctx echo.Context) error {
		idParam, err := base64.URLEncoding.DecodeString(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err)
		}

		info, err := workflowEngine.GetWorkflowStatus(ctx.Request().Context(), string(idParam))
		if err != nil {
			logger.Println(err)
			return ctx.JSON(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusOK, info)
	})

	if err := server.Start(":9000"); err != nil {
		logger.Fatal("failed to start server %w", err)
	}
}
