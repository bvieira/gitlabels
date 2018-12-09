package gitlabels

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitRepository repository for git api
type GitRepository interface {
	ListOwnerRepos(ctx context.Context, owner string) ([]*github.Repository, error)
	ListORGRepos(ctx context.Context, org string) ([]*github.Repository, error)
	ListLabels(ctx context.Context, user, repo string) ([]string, error)

	CreateLabel(ctx context.Context, user, repo, name string, label LabelConfig) error
	EditLabel(ctx context.Context, user, repo, currentName, name string, label LabelConfig) error
	DeleteLabel(ctx context.Context, user, repo, name string) error
}

// GithubRepositoryImpl repository for github api
type GithubRepositoryImpl struct {
	client *github.Client
}

// NewGithubRepository GithubRepository constructor
func NewGithubRepository(token string) *GithubRepositoryImpl {
	return &GithubRepositoryImpl{
		client: github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))),
	}
}

// ListOwnerRepos list all repositories from owner
func (r GithubRepositoryImpl) ListOwnerRepos(ctx context.Context, owner string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	var result []*github.Repository
	for {
		repos, resp, err := r.client.Repositories.List(ctx, owner, opt)
		if err != nil {
			return nil, err
		}
		result = append(result, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

// ListORGRepos list all repositories from organization
func (r GithubRepositoryImpl) ListORGRepos(ctx context.Context, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	var result []*github.Repository
	for {
		repos, resp, err := r.client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			return nil, err
		}
		result = append(result, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

// ListLabels list all labels from 'user/repository'
func (r GithubRepositoryImpl) ListLabels(ctx context.Context, user, repo string) ([]string, error) {
	labels, _, err := r.client.Issues.ListLabels(ctx, user, repo, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, label := range labels {
		result = append(result, label.GetName())
	}
	return result, nil
}

// CreateLabel create label on 'user/repository'
func (r GithubRepositoryImpl) CreateLabel(ctx context.Context, user, repo, name string, label LabelConfig) error {
	_, _, err := r.client.Issues.CreateLabel(ctx, user, repo, &github.Label{
		Name:        prtstr(name),
		Color:       prtstr(label.Color),
		Description: prtstr(label.Description),
	})
	return err
}

// EditLabel edit label from 'user/repository'
func (r GithubRepositoryImpl) EditLabel(ctx context.Context, user, repo, currentName, name string, label LabelConfig) error {
	_, _, err := r.client.Issues.EditLabel(ctx, user, repo, currentName, &github.Label{
		Name:        prtstr(name),
		Color:       prtstr(label.Color),
		Description: prtstr(label.Description),
	})
	return err
}

// DeleteLabel delete label from 'user/repository'
func (r GithubRepositoryImpl) DeleteLabel(ctx context.Context, user, repo, name string) error {
	_, err := r.client.Issues.DeleteLabel(ctx, user, repo, name)
	return err
}

func prtstr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
