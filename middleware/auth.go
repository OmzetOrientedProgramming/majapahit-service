package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"strings"
)

const (
	authorizationTypeBearer = "bearer"
)

// AuthMiddleware struct for middleware auth
type AuthMiddleware struct {
	firebaseAuth firebaseauth.Repo
}

// NewAuthMiddleware for creating AuthMiddleware instance
func NewAuthMiddleware(firebaseAuth firebaseauth.Repo) AuthMiddleware {
	return AuthMiddleware{firebaseAuth: firebaseAuth}
}

// AuthMiddleware function for handling auth middleware
func (a AuthMiddleware) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			//var token string
			for key, values := range ctx.Request().Header {
				if key == "Authorization" {
					if len(values) < 1 || values[0] == "" {
						err := errors.New("invalid authorization header format")
						return ctx.JSON(http.StatusUnauthorized, util.APIResponse{
							Status:  http.StatusUnauthorized,
							Message: "unauthorized",
							Errors:  []string{err.Error()},
						})
					}

					authHeader := strings.Split(values[0], " ")

					authorizationType := strings.ToLower(authHeader[0])
					if authorizationType != authorizationTypeBearer {
						err := fmt.Errorf("unsupported authorization type %s", authorizationType)
						return ctx.JSON(http.StatusUnauthorized, util.APIResponse{
							Status:  http.StatusUnauthorized,
							Message: "unauthorized",
							Errors:  []string{err.Error()},
						})
					}

					userData, err := a.firebaseAuth.GetUserDataFromToken(authHeader[1])
					if err != nil {
						if errors.Cause(err) == firebaseauth.ErrInputValidation {
							err, _ := util.ErrorUnwrap(err)
							return ctx.JSON(http.StatusUnauthorized, util.APIResponse{
								Status:  http.StatusUnauthorized,
								Message: "unauthorized",
								Errors:  err,
							})
						}
						logrus.Error("[failed to get data from firebase repo] ", err.Error())
						return ctx.JSON(http.StatusInternalServerError, util.APIResponse{
							Status:  http.StatusInternalServerError,
							Message: "internal server error",
						})
					}

					ctx.Set("userData", userData)
					return next(ctx)
				}
			}

			err := errors.New("authorization header is not provided")
			return ctx.JSON(http.StatusUnauthorized, util.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
				Errors:  []string{err.Error()},
			})
		}
	}
}

// ParseUserData is used to get the user data from middleware context
func ParseUserData(ctx echo.Context, status int) (*firebaseauth.UserDataFromToken, error) {
	userData := (ctx.Get("userData")).(*firebaseauth.UserDataFromToken)
	switch status {
	case util.StatusCustomer:
		if userData.Users[0].ProviderUserInfo[0].ProviderID == "phone" {
			return userData, nil
		}

		return nil, errors.Wrap(ErrForbidden, "user is not customer")
	case util.StatusBusinessAdmin:
		if userData.Users[0].ProviderUserInfo[0].ProviderID == "password" {
			return userData, nil
		}

		return nil, errors.Wrap(ErrForbidden, "user is not business admin")
	}

	return nil, errors.Wrap(ErrInputValidationError, "status must be 0 or 1")
}
