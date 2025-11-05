package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jibitesh/state-service/internal/svc"
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type McpHandler struct {
	log *logger.Logger
	tr  trace.Tracer
	svc svc.McpService
}

func (mh *McpHandler) GetMcpHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	mcpId := c.Param("mcpId")
	mcp, err := mh.svc.GetByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

func (mh *McpHandler) CreateMcpHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	var req runtime.MCP
	if err := c.ShouldBindJSON(&req); err != nil {
		mh.log.Errorf("Error while binding request data to MCP: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mcp, err := mh.svc.CreateByInteractionIdAndWorkflowId(interactionId, workflowId, &req)
	if err != nil {
		mh.log.Errorf("Error while creating MCP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mcp)
}

func (mh *McpHandler) UpdateMcpHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	mcpId := c.Param("mcpId")
	var req runtime.MCP
	if err := c.ShouldBindJSON(&req); err != nil {
		mh.log.Errorf("Error while binding request data to MCP: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if mcpId != req.ID {
		mh.log.Errorf("Invalid mcp id: %s and MCP json ID: %s", mcpId, req.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mcp id"})
		return
	}
	mcp, err := mh.svc.UpdateByInteractionIdAndWorkflowId(interactionId, workflowId, &req)
	if err != nil {
		mh.log.Errorf("Error while updating MCP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

func (mh *McpHandler) AddToolHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	mcpId := c.Param("mcpId")
	var req runtime.Tool
	if err := c.ShouldBindJSON(&req); err != nil {
		mh.log.Errorf("Error while binding request data to Tool: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mcp, err := mh.svc.AddTool(interactionId, workflowId, mcpId, &req)
	if err != nil {
		mh.log.Errorf("Error while adding tool to MCP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

func (mh *McpHandler) DeleteMcpHandler(c *gin.Context) {
	interactionId := c.Param("interactionId")
	workflowId := c.Param("workflowId")
	mcpId := c.Param("mcpId")
	if err := mh.svc.DeleteByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId); err != nil {
		mh.log.Errorf("Error while deleting MCP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func NewMcpHandler(log *logger.Logger, tr trace.Tracer, svc svc.McpService) *McpHandler {
	return &McpHandler{
		log: log,
		tr:  tr,
		svc: svc,
	}
}
