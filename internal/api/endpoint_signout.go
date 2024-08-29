package api

import (
	"database/sql"
	"github.com/google/uuid"
	"net/http"
	"surge/internal/schema"
)

type LogoutScope string

const (
	LogoutScopeGlobal LogoutScope = "global"
	LogoutScopeLocal  LogoutScope = "local"
	LogoutScopeOthers LogoutScope = "others"
)

// EndpointSignOut processes request for logging out a current user
func (a *SurgeAPI) EndpointSignOut(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	config := a.config

	//scope := LogoutScopeGlobal
	//
	//if r.URL.Query() != nil {
	//	switch r.URL.Query().Get("scope") {
	//	//case "", "global":
	//	//	scope = LogoutScopeGlobal
	//	//case "local":
	//	//	scope = LogoutScopeLocal
	//	//case "others":
	//	//	scope = LogoutScopeOthers
	//	default:
	//		return BadRequestError(ErrorCodeInvalidLogoutScope, fmt.Sprintf("Unsupported logout scope %q", r.URL.Query().Get("scope")))
	//	}
	//}

	c := getClaims(ctx)
	userId, err := c.GetSubjectUUID()
	if err != nil {
		return err
	}

	err = a.Transaction(ctx, func(tx *sql.Tx, queries *schema.Queries) error {
		return queries.RevokeRefreshTokensOfUser(ctx, uuid.NullUUID{UUID: userId, Valid: true})
	})
	if err != nil {
		return InternalServerError("Error logging out user: %+v", err)
	}

	a.clearCookieTokens(config, w)
	w.WriteHeader(http.StatusNoContent)

	return nil
}
