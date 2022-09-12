package utils

import "code.gitea.io/sdk/gitea"

func CreateGithook() error {
	client, err := gitea.NewClient(sets.GetGitAddress(), InitClient)
	if err != nil {
		sets.Logger.Error("error creating gitea client: %v", err)
		return err
	}

	client.SetBasicAuth("gitea_admin", "gitea_admin")

	_, resp, err := client.CreateOrgHook("test", gitea.CreateHookOption{
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
		sets.Logger.Error("error creating gitea hook: %v", err)
		return err
	}
	if resp.StatusCode != 201 {
		sets.Logger.Error("error creating gitea hook: %v", resp.Status)
		return err
	}

	return nil

}

func InitClient(client *gitea.Client) error {

	return nil
}
