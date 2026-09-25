package auth

// SuperTokens mounts and serves the "/auth/*" API itself via
// pkg/supertokens.Middleware (wired in internal/controller/restapi/router.go).
// No handler in this codebase executes the functions below: they exist
// solely so `swag` can generate accurate OpenAPI documentation for routes
// this service exposes but does not implement. The request/response types
// they reference live in doctypes.go, next to them, for the same reason.

// docSignUp documents SuperTokens' sign-up endpoint.
//
// @Summary     Sign up
// @Description Create an account with email and password (SuperTokens EmailPassword recipe). On success, an access token is returned in the "st-access-token" response header (or "sAccessToken" cookie) - paste it into the Authorize dialog to call protected endpoints. Served by the SuperTokens SDK middleware, not by this codebase.
// @ID          authSignUp
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     Credentials true "Email and password form fields"
// @Success     200  {object} SignUpResponse
// @Router      /auth/signup [post]
//
//nolint:unused // documentation-only stub for swag, see comment above.
func docSignUp() {}

// docSignIn documents SuperTokens' sign-in endpoint.
//
// @Summary     Sign in
// @Description Authenticate with email and password (SuperTokens EmailPassword recipe). On success, an access token is returned in the "st-access-token" response header (or "sAccessToken" cookie) - paste it into the Authorize dialog to call protected endpoints. Served by the SuperTokens SDK middleware, not by this codebase.
// @ID          authSignIn
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     Credentials true "Email and password form fields"
// @Success     200  {object} SignInResponse
// @Router      /auth/signin [post]
//
//nolint:unused // documentation-only stub for swag, see comment above.
func docSignIn() {}

// docSignOut documents SuperTokens' sign-out endpoint.
//
// @Summary     Sign out
// @Description Revoke the current session (SuperTokens Session recipe). Served by the SuperTokens SDK middleware, not by this codebase.
// @ID          authSignOut
// @Tags        auth
// @Produce     json
// @Success     200 {object} SignOutResponse
// @Security    BearerAuth
// @Router      /auth/signout [post]
//
//nolint:unused // documentation-only stub for swag, see comment above.
func docSignOut() {}
