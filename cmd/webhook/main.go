package main

import (
	"flag"
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/gin-gonic/gin"
)

type WebhookBody struct {
	Secret     string `json:"secret"`
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	CompareURL string `json:"compare_url"`
	Commits    []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Timestamp string `json:"timestamp"`
	} `json:"commits"`
	Repository struct {
		ID    int `json:"id"`
		Owner struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Username  string `json:"username"`
		} `json:"owner"`
		Name            string `json:"name"`
		FullName        string `json:"full_name"`
		Description     string `json:"description"`
		Private         bool   `json:"private"`
		Fork            bool   `json:"fork"`
		HTMLURL         string `json:"html_url"`
		SSHURL          string `json:"ssh_url"`
		CloneURL        string `json:"clone_url"`
		Website         string `json:"website"`
		StarsCount      int    `json:"stars_count"`
		ForksCount      int    `json:"forks_count"`
		WatchersCount   int    `json:"watchers_count"`
		OpenIssuesCount int    `json:"open_issues_count"`
		DefaultBranch   string `json:"default_branch"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
	} `json:"repository"`
	Pusher struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"pusher"`
	Sender struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"sender"`
}

func main() {

	addressport := flag.String("addressport", ":8080", "address and port to listen on")

	gittoken := flag.String("gittoken", "", "git token")
	giturl := flag.String("giturl", "", "git url")

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/webhook", func(c *gin.Context) {
			var webhookBody WebhookBody
			if err := c.ShouldBindJSON(&webhookBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			client, err := gitea.NewClient(*giturl, InitClient)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Sprintf("client: %v", client)
			fmt.Sprintf("gittoken : %v", *gittoken)

			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})
	}

	router.Run(*addressport)
}

func InitClient(client *gitea.Client) error {

	return nil
}
