package artifactory

// UserAPIKey represents the JSON returned for a user's API Key in Artifactory
type UserAPIKey struct {
	APIKey string `json:"apiKey"`
}
