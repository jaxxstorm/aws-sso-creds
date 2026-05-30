package config

type SSOConfig struct {
	StartURL  string
	Region    string
	AccountID string
	RoleName  string
}

type SSOCacheConfig struct {
	StartURL     string `json:"startUrl"`
	Region       string `json:"region"`
	AccessToken  string `json:"accessToken"`
	ExpiresAt    string `json:"expiresAt"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken"`
}
