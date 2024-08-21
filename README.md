# Surge
A simple identity/authentication server written in Go

Following is core features provided by surge
- Authenticate user and manage sessions (revoking refresh token as well)
- Exposed JWKs endpoint (.well-known/jwks.json)
- Supports major OAuth2 providers out of box
- Automatic database migration with go-migrate
- Pre configured docker compose
