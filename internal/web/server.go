package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pnocera/novaco/internal/gitops"
	"github.com/pnocera/novaco/internal/settings"
)

var (
	sets = settings.GetSettings()
)

func Serve() {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/webhook", managehook)

		v1.GET("/healthz", func(c *gin.Context) {
			//sets.Logger.Info("healthz ok")
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})

	}

	router.Run(sets.BindIPs + ":" + sets.APIPort)

}

func managehook(c *gin.Context) {
	var webhookBody gitops.WebhookBody
	if err := c.ShouldBindJSON(&webhookBody); err != nil {
		sets.Logger.Error("error parsing webhook body ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := gitops.ManageWebhook(webhookBody)
	if err != nil {
		sets.Logger.Error("error managing webhook ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
