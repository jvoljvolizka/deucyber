package deucyber

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//5568fbba42d2524cb5d7e613fbe04171005b1875
var client *github.Client
var CurPrs []*github.PullRequest

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
	CurPrs = Prs
	return Prs, nil
}

func UpdatePrs() error {
	Prs, _, err := client.PullRequests.List(context.Background(), Conf.GithubUsername, Conf.GithubRepo, nil)

	if err != nil {
		return err
	}
	CurPrs = Prs
	return nil
}

func Merge(prNum int) error {
	_, _, err := client.PullRequests.Merge(context.Background(), Conf.GithubUsername, Conf.GithubRepo, prNum, "", nil)
	if err != nil {
		return err
	}
	go UpdatePrs()
	return nil
}
