package deucyber

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//5568fbba42d2524cb5d7e613fbe04171005b1875
var client *github.Client

func GitInit() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Conf.GithubApiKey},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func GetPrs() ([]*github.PullRequest, error) {
	Prs, _, err := client.PullRequests.List(context.Background(), Conf.GithubUsername, Conf.GithubRepo, nil)

	if err != nil {
		return nil, err
	}
	return Prs, nil
}

func Merge(prNum int) error {
	_, _, err := client.PullRequests.Merge(context.Background(), Conf.GithubUsername, Conf.GithubRepo, prNum, "", nil)
	if err != nil {
		return err
	}
	return nil
}
