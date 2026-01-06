package common

const (
	// BasicAuth is SonarQube's BasicAuth method of authentification that needs a username and a password
	BasicAuth AuthType = "BasicAuth"

	// PersonalAccessToken is SonarQube's PersonalAccessToken method of authentification.
	PersonalAccessToken AuthType = "PersonalAccessToken"
)

// AuthType represents an authentication type within SonarQube.
type AuthType string
