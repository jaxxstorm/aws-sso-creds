package config

type SSOConfig struct {
	StartUrl  string
	Region    string
	AccountID string
	RoleName  string
}

type SSOCacheConfig struct {
	StartUrl    string `json:"startUrl"`
	Region      string `json:"region"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}
