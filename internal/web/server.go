package web

import (
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/gin-gonic/gin"
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
	var webhookBody WebhookBody
	if err := c.ShouldBindJSON(&webhookBody); err != nil {
		sets.Logger.Error("error parsing webhook body ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sets.Logger.Info("webhook body", webhookBody)

	client, err := gitea.NewClient(
		fmt.Sprintf("%s:%s", sets.IP(), sets.GitPort), InitClient)

	if err != nil {
		sets.Logger.Error("error creating gitea client", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client.SetBasicAuth("gitea_admin", "gitea_admin")

	repos, resp, err := client.SearchRepos(gitea.SearchRepoOptions{
		Keyword: webhookBody.Repository.Name,
	})
	if err != nil {
		sets.Logger.Error("error searching repos: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK {
		sets.Logger.Error("error searching repos: %v", resp.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Status})
		return
	}

	for _, repo := range repos {
		sets.Logger.Info("repo: %v", repo)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func InitClient(client *gitea.Client) error {

	return nil
}
