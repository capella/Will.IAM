package models

// ServiceAccount type
type ServiceAccount struct {
	ID        string `json:"id" pg:"id"`
	KeyID     string `json:"keyId" pg:"key_id"`
	KeySecret string `json:"keySecret" pg:"key_secret"`
	Email     string `json:"email" pg:"email"`
	CreatedUpdatedAt
}