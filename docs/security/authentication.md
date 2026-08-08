# Authentication

Authentication uses an extensible Provider interface mapped via `/api/v1/auth/login/:provider`.
Supported mock providers include:
- Local
- Google
- GitHub
- Microsoft Entra ID
- OIDC
- SAML 2.0

Resulting identities map seamlessly into the central session management engine.
