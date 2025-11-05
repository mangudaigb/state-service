package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jibitesh/state-service/internal/svc"
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type InteractionHandler struct {
	log *logger.Logger
	tr  trace.Tracer
	svc svc.InteractionService
}

func (ih *InteractionHandler) GetInteractionHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	interaction, err := ih.svc.GetById(interactionId)
	if err != nil {
		ih.log.Errorf("Error while getting interaction by id: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interaction)
}

func (ih *InteractionHandler) CreateInteractionHandler(c *gin.Context) {
	var req runtime.Interaction
	if err := c.ShouldBindJSON(&req); err != nil {
		ih.log.Errorf("Error while binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	interaction, err := ih.svc.Create(&req)
	if err != nil {
		ih.log.Errorf("Error while creating interaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, interaction)
}

func (ih *InteractionHandler) UpdateInteractionHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	var req runtime.Interaction
	if err := c.ShouldBindJSON(&req); err != nil {
		ih.log.Errorf("Error while binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if interactionId != req.ID {
		ih.log.Errorf("Invalid interactionId: %s and Interaction json ID: %s", interactionId, req.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interaction id"})
		return
	}
	interaction, err := ih.svc.Update(&req)
	if err != nil {
		ih.log.Errorf("Error while updating interaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interaction)
}

func (ih *InteractionHandler) DeleteInteractionHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	if err := ih.svc.DeleteById(interactionId); err != nil {
		ih.log.Errorf("Error while deleting interaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ih *InteractionHandler) UpdatePlanHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	planId := c.Param("planId")
	var req runtime.Plan
	if err := c.ShouldBindJSON(&req); err != nil {
		ih.log.Errorf("Error while binding request data to Plan: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if planId != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan id"})
		return
	}
	interaction, err := ih.svc.UpdatePlan(interactionId, req.ID, &req)
	if err != nil {
		ih.log.Errorf("Error while updating plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interaction)
}

func (ih *InteractionHandler) UpdateWorkflowHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	var req runtime.Workflow
	if err := c.ShouldBindJSON(&req); err != nil {
		ih.log.Errorf("Error while binding request data to Workflow: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workflowId != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow id"})
		return
	}
	interaction, err := ih.svc.UpdateWorkflow(interactionId, req.ID, &req)
	if err != nil {
		ih.log.Errorf("Error while updating workflow: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interaction)
}

func (ih *InteractionHandler) UpdateExecutionGraphHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	executionId := c.Param("executionId")
	var req runtime.ExecutionGraph
	if err := c.ShouldBindJSON(&req); err != nil {
		ih.log.Errorf("Error while binding request data to ExecutionGraph: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if executionId != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid execution id"})
		return
	}
	interaction, err := ih.svc.UpdateExecutionGraph(interactionId, workflowId, executionId, &req)
	if err != nil {
		ih.log.Errorf("Error while updating execution graph: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interaction)
}

func NewInteractionHandler(log *logger.Logger, tr trace.Tracer, svc svc.InteractionService) *InteractionHandler {
	return &InteractionHandler{
		log: log,
		tr:  tr,
		svc: svc,
	}
}
