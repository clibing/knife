package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v61/github"
)

func TestGithubClient(t *testing.T) {

	client := github.NewClient(nil)

	orgs, _, err := client.Organizations.List(context.Background(), "clibing", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	opt := &github.RepositoryListByOrgOptions{Type: "public"}

	for _, o := range orgs {
		fmt.Println(o.Name)
		repos, _, e := client.Repositories.ListByOrg(context.Background(), *o.Name, opt)
		if e != nil {
			continue
		}
		for _, repo := range repos {
			fmt.Println(repo.Name)
		}

	}

	rs, _, _ := client.Repositories.ListByUser(context.Background(), "clibing", nil)
	for _, r := range rs {
		fmt.Println(r.Name, r.FullName, r.Name)
	}

}

func TestClientString(t *testing.T) {
	n := "admin"
	p := &n
	fmt.Println(*p)

}
