package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mangudaigb/state-service/internal/handler"
)

func SetupRouter(ge *gin.Engine, ih *handler.InteractionHandler, sh *handler.StepHandler, mh *handler.McpHandler) {
	v1 := ge.Group("/api/v1")
	{
		interactionRouter := v1.Group("/interactions")
		{
			interactionRouter.POST("", ih.CreateInteractionHandler)
			interactionRouter.GET("/:interactionId", ih.GetInteractionHandler)
			interactionRouter.PUT("/:interactionId", ih.UpdateInteractionHandler)
			interactionRouter.DELETE("/:interactionId", ih.DeleteInteractionHandler)

			planRouter := interactionRouter.Group("/:interactionId/plans")
			{
				planRouter.PUT("/:planId", ih.UpdatePlanHandler)
			}

			workflowRouter := interactionRouter.Group("/:interactionId/workflows")
			{
				workflowRouter.PUT("/:workflowId", ih.UpdateWorkflowHandler)
				mcpRouter := workflowRouter.Group("/:workflowId/mcps")
				{
					mcpRouter.POST("", mh.CreateMcpHandler)
					mcpRouter.PUT("/:mcpId", mh.UpdateMcpHandler)
					mcpRouter.DELETE("/:mcpId", mh.DeleteMcpHandler)
					mcpRouter.POST("/:mcpId/tools", mh.AddToolHandler)
				}
				executionRouter := workflowRouter.Group("/:workflowId/executions")
				{
					executionRouter.PUT("/:executionId", ih.UpdateExecutionGraphHandler)
					stepRouter := executionRouter.Group("/:executionId/steps")
					{
						stepRouter.POST("", sh.CreateStepHandler)
						stepRouter.GET("/:stepId", sh.GetStepHandler)
						stepRouter.PUT("/:stepId", sh.UpdateStepHandler)
						stepRouter.POST("/:stepId/status", sh.UpdateStatusHandler)
						stepRouter.DELETE("/:stepId", sh.DeleteStepHandler)
					}
				}
			}
			//interactionRouter.POST("", ih.CreateInteractionHandler)
			//interactionRouter.GET("/:interactionId", ih.GetInteractionHandler)
			//interactionRouter.PUT("/:interactionId", ih.UpdateInteractionHandler)
			//interactionRouter.DELETE("/:interactionId", ih.DeleteInteractionHandler)
			//interactionRouter.POST("/:interactionId/plans/:planId", ih.UpdatePlanHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId", ih.UpdateWorkflowHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId/executions/:executionId", ih.UpdateExecutionGraphHandler)
			//interactionRouter.GET("/:interactionId/workflows/:workflowId/mcps/mcpId", mh.GetMcpHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId/mcps", mh.CreateMcpHandler)
			//interactionRouter.PUT("/:interactionId/workflows/:workflowId/mcps/mcpId", mh.UpdateMcpHandler)
			//interactionRouter.DELETE("/:interactionId/workflows/:workflowId/mcps/mcpId", mh.DeleteMcpHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId/mcps/mcpId/tools", mh.AddToolHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId/executions/:executionId/steps", sh.CreateStepHandler)
			//interactionRouter.GET("/:interactionId/workflows/:workflowId/executions/:executionId/steps/:stepId", sh.GetStepHandler)
			//interactionRouter.PUT("/:interactionId/workflows/:workflowId/executions/:executionId/steps/:stepId", sh.UpdateStepHandler)
			//interactionRouter.POST("/:interactionId/workflows/:workflowId/executions/:executionId/steps/:stepId/status", sh.UpdateStatusHandler)
			//interactionRouter.DELETE("/:interactionId/workflows/:workflowId/executions/:executionId/steps/:stepId", sh.DeleteStepHandler)
		}
	}
}
