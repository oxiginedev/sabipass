package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/config"
	"github.com/oxiginedev/sabipass/internal/database"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sabipass/internal/pkg/jwt"
	"github.com/oxiginedev/sabipass/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauthHandler struct {
	cfg          *config.Config
	googleConfig *oauth2.Config
	tokenManager jwt.TokenManager
	userRepo     models.UserRepository
}

func NewOauthHandler(cfg *config.Config, tokenManager jwt.TokenManager, userRepo models.UserRepository) *oauthHandler {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.Oauth.Google.ClientID,
		ClientSecret: cfg.Oauth.Google.ClientSecret,
		RedirectURL:  cfg.Oauth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &oauthHandler{
		cfg:          cfg,
		googleConfig: googleConfig,
		tokenManager: tokenManager,
		userRepo:     userRepo,
	}
}

func (o *oauthHandler) HandleGoogleLoginRedirect(c *gin.Context) {
	state, err := generateOauthStateCookie(c)
	if err != nil {
		slog.Error("could not generate oauth state cookie", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	url := o.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (o *oauthHandler) HandleGoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	oauthState, err := c.Cookie("oauth_state")
	if err != nil {
		slog.Error("could not get oauth state cookie", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	if state != oauthState {
		slog.Error("oauth state does not match", slog.String("state", state), slog.String("oauthState", oauthState))
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorResponse("invalid oauth state", nil))
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	code := c.Query("code")
	token, err := o.googleConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		slog.Error("could not exchange oauth code", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorResponse("unable to verify sign in with google", nil))
		return
	}

	client := o.googleConfig.Client(c.Request.Context(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		slog.Error("could not get userinfo", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("could not get google userinfo", slog.Int("status", res.StatusCode))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	googleUser := models.GoogleUser{}
	err = json.NewDecoder(res.Body).Decode(&googleUser)
	if err != nil {
		slog.Error("could not decode userinfo", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	user, err := o.userRepo.FindOne(c.Request.Context(), &models.FindUserOptions{
		Email: googleUser.Email,
	})
	if err != nil && !errors.Is(err, database.ErrUserNotFound) {
		slog.Error("could not find user", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	if user == nil {
		user = &models.User{
			ID:              utils.Uuid(),
			Name:            utils.Ptr(googleUser.Name),
			Username:        strings.Split(googleUser.Email, "@")[0],
			Email:           googleUser.Email,
			EmailVerifiedAt: utils.Ptr(time.Now()),
			GoogleID:        utils.Ptr(googleUser.ID),
			Avatar:          utils.Ptr(googleUser.Picture),
		}

		err = o.userRepo.Create(c.Request.Context(), user)
		if err != nil {
			slog.Error("could not create user", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
			return
		}
	}

	accessToken, err := o.tokenManager.GenerateToken(user)
	if err != nil {
		slog.Error("could not generate access token", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("something went wrong", nil))
		return
	}

	resp := gin.H{
		"user":  user,
		"token": accessToken,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("user auth successful", resp))
}

func generateOauthStateCookie(c *gin.Context) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauth_state", state, 60*60*24, "/", "", false, true)
	return state, nil
}
