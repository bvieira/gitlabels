package gitlabels

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/google/go-github/github"
)

// gitlabels errors
var (
	ErrEmptyUser = errors.New("at least one should be present: org or owner")
)

// Service rules for gitlabel
type Service struct {
	logger *log.Logger
	git    GitRepository
}

// NewService service constructor
func NewService(git GitRepository, logger *log.Logger) *Service {
	return &Service{git: git, logger: logger}
}

// Exec execute gitlabel functions based on config
func (s Service) Exec(ctx context.Context, cfg Config) error {
	if cfg.getUser() == "" {
		return ErrEmptyUser
	}

	regex := regexp.MustCompile(cfg.ProjectRegex)

	repos, rerr := s.listRepos(ctx, cfg.ORG, cfg.Owner)
	if rerr != nil {
		return rerr
	}

	for _, repo := range repos {
		if !regex.MatchString(repo.GetName()) {
			continue
		}

		s.removeLabels(ctx, cfg.getUser(), repo.GetName(), cfg.RemoveLabels)

		renamed := s.renameLabels(ctx, cfg.getUser(), repo.GetName(), cfg.RenameLabels, cfg.Labels)

		err := s.createLabels(ctx, cfg.getUser(), repo.GetName(), cfg.Labels, renamed)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) listRepos(ctx context.Context, org, owner string) ([]*github.Repository, error) {
	if org != "" {
		return s.git.ListORGRepos(ctx, org)
	}
	return s.git.ListOwnerRepos(ctx, owner)
}

func (s Service) removeLabels(ctx context.Context, user, repo string, labels []string) {
	for _, label := range labels {
		err := s.git.DeleteLabel(ctx, user, repo, label)
		if err == nil {
			s.logger.Printf("[%s/%s] removed label:[%s]", user, repo, label)
		} else {
			s.logger.Printf("[%s/%s] failed to remove label:[%s], error: %v", user, repo, label, err)
		}
	}
}

func (s Service) renameLabels(ctx context.Context, user, repo string, labels map[string]string, labelsInfo map[string]LabelConfig) map[string]interface{} {
	renamed := make(map[string]interface{})

	for label, newLabel := range labels {
		err := s.git.EditLabel(ctx, user, repo, label, newLabel, labelsInfo[newLabel])
		if err == nil {
			s.logger.Printf("[%s/%s] renaming label from:[%s] to:[%s]", user, repo, label, newLabel)
		} else {
			renamed[newLabel] = nil
			s.logger.Printf("[%s/%s] failed to rename label from:[%s] to:[%s], error: %v", user, repo, label, newLabel, err)
		}
	}
	return renamed
}

func (s Service) createLabels(ctx context.Context, user, repo string, labels map[string]LabelConfig, ignore map[string]interface{}) error {
	for name, labelCfg := range labels {
		if _, exists := ignore[name]; exists {
			continue
		}

		cerr := s.git.CreateLabel(ctx, user, repo, name, labelCfg)
		if cerr == nil {
			s.logger.Printf("[%s/%s] creating label:[%s]", user, repo, name)
			continue
		}

		if errResponse, ok := cerr.(*github.ErrorResponse); !ok || !findError(errResponse.Errors, "already_exists") {
			return cerr
		}

		err := s.git.EditLabel(ctx, user, repo, name, name, labelCfg)
		if err != nil {
			return err
		}
		s.logger.Printf("[%s/%s] updating label:[%s]", user, repo, name)
	}
	return nil
}

func findError(errors []github.Error, code string) bool {
	for _, err := range errors {
		if err.Code == code {
			return true
		}
	}
	return false
}
