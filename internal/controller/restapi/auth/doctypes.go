package auth

// The types below model the JSON payloads of SuperTokens' "/auth/*" HTTP
// API. They exist only to document that API for Swagger (see docs.go);
// SuperTokens' SDK middleware (see pkg/supertokens) encodes/decodes the
// real requests and responses itself.

// FormField -.
type FormField struct {
	ID    string `json:"id"    example:"email"`
	Value string `json:"value" example:"john@example.com"`
} // @name auth.FormField

// Credentials -.
type Credentials struct {
	FormFields []FormField `json:"formFields"`
} // @name auth.Credentials

// User -.
type User struct {
	Email      string   `json:"email"      example:"john@example.com"`
	ID         string   `json:"id"         example:"550e8400-e29b-41d4-a716-446655440000"`
	TenantIDs  []string `json:"tenantIds"  example:"public"`
	TimeJoined uint64   `json:"timeJoined" example:"1700000000000"`
} // @name auth.User

// SignUpResponse -.
type SignUpResponse struct {
	Status string `json:"status" example:"OK"`
	User   User   `json:"user"`
} // @name auth.SignUpResponse

// SignInResponse -.
type SignInResponse struct {
	Status string `json:"status" example:"OK"`
	User   User   `json:"user"`
} // @name auth.SignInResponse

// SignOutResponse -.
type SignOutResponse struct {
	Status string `json:"status" example:"OK"`
} // @name auth.SignOutResponse
