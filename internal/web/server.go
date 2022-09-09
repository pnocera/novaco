package web

import (
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/gin-gonic/gin"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/utils"
)

var (
	logger = utils.NewKLogger("web")
	sets   = settings.GetSettings()
)

func Serve() {

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/webhook", managehook)
	}

	router.Run(sets.BindIPs + ":" + sets.APIPort)

}

func managehook(c *gin.Context) {
	var webhookBody WebhookBody
	if err := c.ShouldBindJSON(&webhookBody); err != nil {
		logger.Error("error parsing webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("webhook body: %v", webhookBody)

	client, err := gitea.NewClient(
		fmt.Sprintf("%s:%s", utils.IP(), sets.GitPort), InitClient)

	if err != nil {
		logger.Error("error creating gitea client: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client.SetBasicAuth("admin", "admin")

	repos, resp, err := client.SearchRepos(gitea.SearchRepoOptions{
		Keyword: webhookBody.Repository.Name,
	})
	if err != nil {
		logger.Error("error searching repos: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("error searching repos: %v", resp.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Status})
		return
	}

	for _, repo := range repos {
		logger.Info("repo: %v", repo)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func InitClient(client *gitea.Client) error {

	return nil
}
