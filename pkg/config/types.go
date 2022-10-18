package config

type SSOConfig struct {
	StartURL  string
	Region    string
	AccountID string
	RoleName  string
}

type SSOCacheConfig struct {
	StartURL    string `json:"startUrl"`
	Region      string `json:"region"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}
