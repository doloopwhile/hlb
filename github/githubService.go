package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"errors"
	"net/url"
	"github.com/mpppk/hlb/project"
)

type Service struct {
	Client *github.Client
	ListOptions *github.ListOptions
}

func NewService(ctx context.Context, token string, baseUrlStrs ...string) (*Service, error) {
	if len(baseUrlStrs) > 1 {
		return nil, errors.New("too many base urls")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if len(baseUrlStrs) == 1 {
		baseUrl, err := url.Parse(baseUrlStrs[0])
		if err != nil {
			return nil, err
		}

		client.BaseURL = baseUrl
	}

	listOpt := &github.ListOptions{PerPage: 100}
	return &Service{Client: client, ListOptions: listOpt}, nil
}

func (s *Service) GetPullRequests(ctx context.Context, owner, repo string) (pullRequests []project.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *s.ListOptions}
	gitHubPullRequests, _, err := s.Client.PullRequests.List(ctx, owner, repo, &opt)

	if err != nil {
		return nil, err
	}

	for _, gitHubPullRequest := range gitHubPullRequests {
		pullRequests = append(pullRequests, &PullRequest{PullRequest: gitHubPullRequest})
	}

	return pullRequests, err
}

func (s *Service) GetIssues(ctx context.Context, owner, repo string) (issues []project.Issue, err error) {
	opt := &github.IssueListByRepoOptions{ListOptions: *s.ListOptions}
	gitHubIssues, err := s.getGitHubIssues(ctx, s.Client, owner, repo, opt)

	if err != nil {
		return nil, err
	}

	for _, gitHubIssue := range gitHubIssues {
		issues = append(issues, &Issue{Issue: gitHubIssue})
	}

	return issues, err
}

func (s *Service) getGitHubIssues(ctx context.Context, client *github.Client, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
	issuesAndPRs, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)

	if err != nil {
		return nil, err
	}

	for _, issueOrPR := range issuesAndPRs {
		if issueOrPR.PullRequestLinks == nil {
			issues = append(issues, issueOrPR)
		}
	}
	return issues, nil
}

func (s *Service) GetRepository(ctx context.Context, owner, repo string) (project.Repository, error) {
	githubRepo, _, err := s.Client.Repositories.Get(ctx, owner, repo)

	if err != nil {
		return nil, err
	}

	return &Repository{Repository: githubRepo}, err
}
