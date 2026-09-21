// Package supertokens implements the SuperTokens Go SDK.
package supertokens

import (
	"context"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	st "github.com/supertokens/supertokens-golang/supertokens"
)

const authPath = "/auth"

// Config - holds SuperTokens configuration.
type Config struct {
	ConnectionURI string
	APIKey        string
	AppName       string
	APIDomain     string
	WebsiteDomain string
}

// Hooks - callbacks that let callers react to SuperTokens events.
type Hooks struct {
	OnSignUp func(ctx context.Context, userID, email string) error
}

// New -.
func New(cfg *Config, h Hooks) error {
	conn := &st.ConnectionInfo{ConnectionURI: cfg.ConnectionURI}
	if cfg.APIKey != "" {
		conn.APIKey = cfg.APIKey
	}

	apiBasePath, webBasePath := authPath, authPath

	err := st.Init(st.TypeInput{
		Supertokens: conn,
		AppInfo: st.AppInfo{
			AppName:         cfg.AppName,
			APIDomain:       cfg.APIDomain,
			WebsiteDomain:   cfg.WebsiteDomain,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &webBasePath,
		},
		RecipeList: []st.Recipe{
			emailpassword.Init(overrideSignUp(h)),
			session.Init(nil),
		},
	})
	if err != nil {
		return fmt.Errorf("supertokens - New: %w", err)
	}

	return nil
}

func overrideSignUp(h Hooks) *epmodels.TypeInput {
	if h.OnSignUp == nil {
		return nil
	}

	return &epmodels.TypeInput{
		Override: &epmodels.OverrideStruct{
			APIs: func(orig epmodels.APIInterface) epmodels.APIInterface {
				originalSignUpPOST := *orig.SignUpPOST

				signUpPOST := func(
					formFields []epmodels.TypeFormField,
					tenantID string,
					options epmodels.APIOptions,
					userContext st.UserContext,
				) (epmodels.SignUpPOSTResponse, error) {
					resp, err := originalSignUpPOST(formFields, tenantID, options, userContext)
					if err != nil || resp.OK == nil {
						return resp, err
					}

					ctx := requestContext(userContext)

					if hookErr := h.OnSignUp(ctx, resp.OK.User.ID, resp.OK.User.Email); hookErr != nil {
						return resp, hookErr
					}

					return resp, nil
				}

				orig.SignUpPOST = &signUpPOST

				return orig
			},
		},
	}
}

// Middleware -.
func Middleware(next http.Handler) http.Handler {
	return st.Middleware(next)
}

// AllCORSHeaders -.
func AllCORSHeaders() []string {
	return st.GetAllCORSHeaders()
}

// DeleteUser -.
func DeleteUser(userID string) error {
	return st.DeleteUser(userID)
}

func requestContext(userContext st.UserContext) context.Context {
	if req := st.GetRequestFromUserContext(userContext); req != nil {
		return req.Context()
	}

	return context.Background()
}
