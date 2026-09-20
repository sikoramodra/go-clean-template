package supertokens

import (
	"context"
	"fmt"

	"github.com/sikoramodra/go-clean-template/internal/usecase"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const authPath = "/auth"

// PostSignUpFunc -.
type PostSignUpFunc func(userID, email string) error

// Config holds SuperTokens configuration.
type Config struct {
	ConnectionURI string
	APIKey        string
	AppName       string
	APIDomain     string
	WebsiteDomain string
	PostSignUp    PostSignUpFunc
}

// New -.
func New(cfg *Config, u usecase.User) error {
	conn := &supertokens.ConnectionInfo{ConnectionURI: cfg.ConnectionURI}
	if cfg.APIKey != "" {
		conn.APIKey = cfg.APIKey
	}

	apiBasePath, webBasePath := authPath, authPath

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: conn,
		AppInfo: supertokens.AppInfo{
			AppName:         cfg.AppName,
			APIDomain:       cfg.APIDomain,
			WebsiteDomain:   cfg.WebsiteDomain,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &webBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(withPostSignUp(u)),
			session.Init(nil),
		},
	})
	if err != nil {
		return fmt.Errorf("supertokens - New: %w", err)
	}

	return nil
}

func withPostSignUp(u usecase.User) *epmodels.TypeInput {
	return &epmodels.TypeInput{
		Override: &epmodels.OverrideStruct{
			APIs: func(orig epmodels.APIInterface) epmodels.APIInterface {
				originalSignUpPOST := *orig.SignUpPOST

				signUpPOST := func(
					formFields []epmodels.TypeFormField,
					tenantID string,
					options epmodels.APIOptions,
					userContext supertokens.UserContext,
				) (epmodels.SignUpPOSTResponse, error) {
					resp, err := originalSignUpPOST(formFields, tenantID, options, userContext)
					if err != nil || resp.OK == nil {
						return resp, err
					}

					ctx := requestContext(userContext)

					if regErr := u.Register(ctx, resp.OK.User.ID, resp.OK.User.Email); regErr != nil {
						if delErr := supertokens.DeleteUser(resp.OK.User.ID); delErr != nil {
							return resp, fmt.Errorf(
								"supertokens - overrideSignUp: register: %w (rollback also failed: %w)",
								regErr, delErr,
							)
						}

						return resp, fmt.Errorf("supertokens - overrideSignUp - userUC.Register: %w", regErr)
					}

					return resp, nil
				}

				orig.SignUpPOST = &signUpPOST

				return orig
			},
		},
	}
}

// requestContext pulls the *http.Request's context out of SuperTokens'
// userContext so downstream calls use the real request context instead of
// whatever context app.Run happened to be holding at startup.
func requestContext(userContext supertokens.UserContext) context.Context {
	if req := supertokens.GetRequestFromUserContext(userContext); req != nil {
		return req.Context()
	}

	return context.Background()
}
