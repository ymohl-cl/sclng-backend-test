// Package gitcrawler provides a driver to crawl data from the GitHub API
// It is used to list public repositories and enrich them with languages and license
package gitcrawler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/go-github/v60/github"
)

const (
	maxPageSize = 100
)

// Repository representation of a GitHub repository
type Repository struct {
	ID        int64
	Owner     string
	Name      string
	Languages map[string]int
	License   string
}

// CreateEventPayload representation of a Github event type
type CreateEventPayload struct {
	RefType      string `json:"ref_type"`
	Ref          string `json:"ref"`
	MasterBranch string `json:"master_branch"`
	Description  string `json:"description"`
	PusherType   string `json:"pusher_type"`
}

// GitCrawler driver to crawl data from the GitHub API
//
//go:generate mockery --name=GitCrawler --output=mocks --filename=gitcrawler.go --outpkg=mocks
type GitCrawler interface {
	// PublicRepository list public repositories from oldest or newest
	PublicRepository(ctx context.Context, nb int32, oldest bool) ([]Repository, error)
	// EnrichRepositories enrich repositories with languages and license information (if available)
	EnrichRepositories(ctx context.Context, repositories []Repository) ([]Repository, error)
}

type gitCrawler struct {
	client *github.Client
}

func New(token string) GitCrawler {
	var g gitCrawler

	g.client = github.NewClient(nil)
	if token != "" {
		g.client = g.client.WithAuthToken(token)
	}

	return &g
}

func (g *gitCrawler) PublicRepository(ctx context.Context, nb int32, oldest bool) ([]Repository, error) {
	if oldest {
		return g.oldestPublicRepository(ctx, nb)
	}

	return g.newestPublicRepository(ctx, nb)
}

func (g *gitCrawler) newestPublicRepository(ctx context.Context, nb int32) ([]Repository, error) {
	var repositories []Repository
	var err error
	var identifier int64

	// indexRepositories is used to avoid duplicates when id repository between Since-100 and Since was removed
	indexRepositories := make(map[int64]bool)
	if identifier, err = g.newestPublicRepositoryIdentifier(ctx); err != nil {
		return nil, err
	}

	opt := &github.RepositoryListAllOptions{
		Since: identifier - (maxPageSize - 1),
	}

	for len(repositories) < int(nb) {
		repos, resp, err := g.client.Repositories.ListAll(ctx, opt)
		if err != nil {
			return nil, err
		}
		resp.Body.Close() // resp is not used

		for _, repo := range repos {
			repoID := repo.GetID()
			if _, ok := indexRepositories[repoID]; ok {
				continue
			}

			indexRepositories[repoID] = true
			repositories = append(repositories, Repository{
				Owner: repo.GetOwner().GetLogin(),
				ID:    repo.GetID(),
				Name:  repo.GetName(),
			})
		}

		opt.Since -= maxPageSize
	}

	return repositories[:nb], nil
}

func (g *gitCrawler) newestPublicRepositoryIdentifier(ctx context.Context) (int64, error) {
	for {
		events, resp, err := g.client.Activity.ListEvents(ctx, &github.ListOptions{
			PerPage: maxPageSize,
		})
		if err != nil {
			return 0, err
		}
		resp.Body.Close() // resp is not used

		for _, event := range events {
			if *event.Type != "CreateEvent" {
				continue
			}

			var payload CreateEventPayload
			if err := json.Unmarshal(event.GetRawPayload(), &payload); err != nil {
				return 0, err
			}
			if payload.RefType == "repository" {
				return event.Repo.GetID(), nil
			}
		}
	}
}

func (g *gitCrawler) oldestPublicRepository(ctx context.Context, nb int32) ([]Repository, error) {
	var repositories []Repository
	opt := &github.RepositoryListAllOptions{
		Since: 0,
	}

	for len(repositories) < int(nb) {
		repos, resp, err := g.client.Repositories.ListAll(ctx, opt)
		if err != nil {
			return nil, err
		}
		resp.Body.Close() // resp is not used

		for _, repo := range repos {
			repositories = append(repositories, Repository{
				Owner: repo.GetOwner().GetLogin(),
				ID:    repo.GetID(),
				Name:  repo.GetName(),
			})
		}

		opt.Since = repos[len(repos)-1].GetID()
	}

	return repositories[:nb], nil
}

func (g *gitCrawler) EnrichRepositories(ctx context.Context, repositories []Repository) ([]Repository, error) {
	var wg sync.WaitGroup
	errs := make(chan error, len(repositories))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range repositories {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			languages, resp, err := g.client.Repositories.ListLanguages(
				ctx,
				repositories[idx].Owner,
				repositories[idx].Name)
			if err != nil && resp.StatusCode != http.StatusNotFound {
				cancel()
				errs <- err

				return
			}
			resp.Body.Close() // resp is not used

			license, resp, err := g.client.Repositories.License(
				ctx,
				repositories[idx].Owner,
				repositories[idx].Name)
			if err != nil && resp.StatusCode != http.StatusNotFound {
				cancel()
				errs <- err

				return
			}
			resp.Body.Close() // resp is not used

			repositories[idx].License = license.GetLicense().GetName()
			repositories[idx].Languages = languages
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return repositories, nil
}
