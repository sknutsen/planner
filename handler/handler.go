package handler

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sknutsen/planner/auth"
	"golang.org/x/oauth2"
)

type Handler struct {
	Host          string
	Port          string
	TursoConfig   TursoConfig
	AuthConfig    AuthConfig
	Authenticator auth.Authenticator
	// DB is a process-wide sql.DB pool; opened in main and closed on shutdown.
	DB *sql.DB
}

type TursoConfig struct {
	DBName     string
	PrimaryUrl string
	AuthToken  string
}

type AuthConfig struct {
	Domain       string
	Audience     string // AUTH0_AUDIENCE: Auth0 API identifier sent on /authorize for the web app (access token audience).
	APIAudience  string // AUTH0_API_AUDIENCE: verifies Bearer JWTs on /api/v1 only.
	ClientId     string
	ClientSecret string
	CallbackUrl  string
}

func (h *Handler) Setup() {
	provider, err := oidc.NewProvider(
		context.Background(),
		h.AuthConfig.Domain+"/",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed creating authenticator. Error: %s\n", err)
		os.Exit(1)
	}

	conf := oauth2.Config{
		ClientID:     h.AuthConfig.ClientId,
		ClientSecret: h.AuthConfig.ClientSecret,
		RedirectURL:  h.AuthConfig.CallbackUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	h.Authenticator = auth.Authenticator{
		Provider: provider,
		Config:   conf,
	}
}
