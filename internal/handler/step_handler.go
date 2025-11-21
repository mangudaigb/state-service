package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/types/runtime"
	"github.com/mangudaigb/state-service/internal/svc"
	"go.opentelemetry.io/otel/trace"
)

type StepHandler struct {
	log *logger.Logger
	tr  trace.Tracer
	svc svc.StepService
}

func (sh *StepHandler) GetStepHandler(c *gin.Context) {
	ctx := c.Request.Context()
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	stepId := c.Param("stepId")
	step, err := sh.svc.GetByInteractionIdAndExecutionIdAndId(ctx, interactionId, workflowId, executionId, stepId)
	if err != nil {
		sh.log.Errorf("Error while getting step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (sh *StepHandler) CreateStepHandler(c *gin.Context) {
	ctx := c.Request.Context()
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	var req runtime.Step
	if err := c.ShouldBindJSON(&req); err != nil {
		sh.log.Errorf("Error while binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	step, err := sh.svc.CreateByInteractionIdAndExecutionId(ctx, interactionId, workflowId, executionId, &req)
	if err != nil {
		sh.log.Errorf("Error while creating step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, step)
}

func (sh *StepHandler) UpdateStepHandler(c *gin.Context) {
	ctx := c.Request.Context()
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	stepId := c.Param("stepId")
	var req runtime.Step
	if err := c.ShouldBindJSON(&req); err != nil {
		sh.log.Errorf("Error while binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if stepId != req.ID {
		sh.log.Errorf("Invalid step id: %s and Step json ID: %s", stepId, req.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step id"})
		return
	}
	step, err := sh.svc.UpdateByInteractionIdAndExecutionId(ctx, interactionId, workflowId, executionId, &req)
	if err != nil {
		sh.log.Errorf("Error while updating step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (sh *StepHandler) UpdateStatusHandler(c *gin.Context) {
	ctx := c.Request.Context()
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	stepId := c.Param("stepId")
	status := c.Query("status")
	var req runtime.Status
	if err := c.ShouldBindJSON(&req); err != nil {
		sh.log.Errorf("Error while binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := runtime.Status(status)
	step, err := sh.svc.UpdateStatusByInteractionIdAndExecutionIdAndId(ctx, interactionId, workflowId, executionId, stepId, s)
	if err != nil {
		sh.log.Errorf("Error while updating step status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (sh *StepHandler) DeleteStepHandler(c *gin.Context) {
	ctx := c.Request.Context()
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	stepId := c.Param("stepId")
	if err := sh.svc.DeleteByInteractionIdAndExecutionIdAndId(ctx, interactionId, workflowId, executionId, stepId); err != nil {
		sh.log.Errorf("Error while deleting step: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func NewStepHandler(log *logger.Logger, tr trace.Tracer, svc svc.StepService) *StepHandler {
	return &StepHandler{
		log: log,
		tr:  tr,
		svc: svc,
	}
}
