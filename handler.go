package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Scalingo/sclng-backend-test-v1/pkg/gitcrawler"

	"github.com/Scalingo/go-utils/logger"
)

const (
	languageParam = "language"
	licenseParam  = "license"
)

type Handler struct {
	gitToken string
}

type PublicRepositoryParams struct {
	Language string
	License  string
}

func (h Handler) publicRepositoryHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	params := readParams(r)
	log := logger.Get(r.Context())
	git := gitcrawler.New(h.gitToken)

	repo, err := git.PublicRepository(r.Context(), 100, false)
	if err != nil {
		log.WithError(err).Error("Failed to list repositories")
		return err
	}

	repo, err = git.EnrichRepositories(r.Context(), repo)
	if err != nil {
		log.WithError(err).Error("Failed to enrich repositories")
		return err
	}

	var response PublicRepositoryResponse
	response.Repositories = convertToRepository(repo, params)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithError(err).Error("Fail to encode JSON")

		return err
	}

	return nil
}

func readParams(r *http.Request) PublicRepositoryParams {
	return PublicRepositoryParams{
		Language: r.URL.Query().Get(languageParam),
		License:  r.URL.Query().Get(licenseParam),
	}
}

func convertToRepository(repo []gitcrawler.Repository, params PublicRepositoryParams) []Repository {
	result := []Repository{}

	for i := range repo {
		if excludeRepository(repo[i], params) {
			continue
		}

		r := Repository{
			FullName: repo[i].Owner + "/" + repo[i].Name,
			Owner:    repo[i].Owner,
			Name:     repo[i].Name,
			License:  repo[i].License,
		}
		r.Languages = make(map[string]LanguagePart)
		for lang, bytes := range repo[i].Languages {
			r.Languages[lang] = LanguagePart{Bytes: bytes}
		}

		result = append(result, r)
	}

	return result
}

func excludeRepository(repo gitcrawler.Repository, params PublicRepositoryParams) bool {
	if params.Language != "" && repo.Languages[params.Language] == 0 {
		return true
	}

	if params.License != "" && !strings.Contains(
		strings.ToLower(repo.License),
		strings.ToLower(params.License)) {
		return true
	}

	return false
}
