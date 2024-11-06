// Package entities contains main entities for all internal packages of the project.
package entities

// URL is a URL struct with ShortURL and OriginalURL versions.
type URL struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}

//type URLGot struct {
//	CorrelationID string `json:"correlation_id"`
//	OriginalURL   string `json:"original_url,omitempty"`
//	ShortURL      string `json:"short_url"`
//}
