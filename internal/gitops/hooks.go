package gitops

import (
	"code.gitea.io/sdk/gitea"
	"github.com/pnocera/novaco/internal/consulapi"
	"github.com/pnocera/novaco/internal/utils"
)

func CreateGithook(passw string) error {
	var err error = nil

	consulexe := utils.Join(utils.BinPath("consul"), "consul.exe")

	consulparams := []string{"kv", "put", "gitea/admin/password", passw}

	err = ExecAndWait(consulexe, consulparams)
	if err != nil {
		sets.Logger.Error("error setting gitea password in consul", err)
	}

	client, err := gitea.NewClient(sets.GetGitAddress(), nil)
	if err != nil {
		sets.Logger.Error("error creating gitea client", err)
		return err
	}

	client.SetBasicAuth("gitea_admin", passw)

	// create gitea organization
	_, _, err = client.AdminCreateOrg("gitea_admin", gitea.CreateOrgOption{
		FullName:    "gitea_org",
		Name:        "gitea_org",
		Description: "gitea organization to deploy apps",
	})
	if err != nil {
		sets.Logger.Error("error creating gitea organization", err)
	}

	_, resp, _ := client.GetOrgHook("gitea_org", 1)
	if resp.StatusCode == 404 {
		_, _, err = client.CreateOrgHook("gitea_org", gitea.CreateHookOption{
			Type: "gitea",
			Config: map[string]string{
				"url":          sets.GetAPIAddress() + "/api/v1/webhook",
				"content_type": "json",
				"secret":       "secret",
			},
			Events: []string{"push", "release"},
			Active: true,
		})
		if err != nil {
			sets.Logger.Error("error creating gitea hook", err)

		}
	}

	//create token for gitea admin user
	token, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{
		Name: "gitea_admin",
	})
	if err != nil {
		sets.Logger.Error("error creating gitea token", err)
	}
	sets.Logger.Info("gitea token: ", token.Token)

	//store token in consul
	consulparams = []string{"kv", "put", "gitea/admin/token", token.Token}

	err = ExecAndWait(consulexe, consulparams)
	if err != nil {
		sets.Logger.Error("error setting gitea token in consul", err)
	}

	return err

}

func ManageWebhook(webhookBody WebhookBody) error {
	var err error = nil

	// get token from consul
	token, err := consulapi.GetKV("gitea/admin/token")
	if err != nil {
		sets.Logger.Error("error getting gitea token from consul", err)
		return err
	}

	client, err := gitea.NewClient(sets.GetGitAddress(), gitea.SetToken(token))
	if err != nil {
		sets.Logger.Error("error creating gitea client", err)
		return err
	}

	_, resp, err := client.SearchRepos(gitea.SearchRepoOptions{
		Keyword: webhookBody.Repository.Name,
	})
	if err != nil {
		sets.Logger.Error("error searching repos", err)
		return err
	}

	if resp.StatusCode == 404 {
		_, _, err = client.CreateOrgRepo("gitea_org", gitea.CreateRepoOption{
			Name:        webhookBody.Repository.Name,
			Description: webhookBody.Repository.Description,
			Private:     true,
			AutoInit:    true,
		})
		if err != nil {
			sets.Logger.Error("error creating repo", err)
			return err
		}
	}

	return err
}
