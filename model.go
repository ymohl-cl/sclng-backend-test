package main

// PublicRepositoryResponse JSON API Response
type PublicRepositoryResponse struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	FullName  string                  `json:"full_name"`
	Owner     string                  `json:"owner"`
	Name      string                  `json:"repository"`
	Languages map[string]LanguagePart `json:"languages"`
	License   string                  `json:"license"`
}

type LanguagePart struct {
	Bytes int `json:"bytes"`
}
