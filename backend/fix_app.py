import re

with open('internal/app/app.go', 'r') as f:
    content = f.read()

# Add imports
imports = """
	"github.com/cloudstorex/backend/internal/security/abac"
	"github.com/cloudstorex/backend/internal/security/engine"
	"github.com/cloudstorex/backend/internal/security/rbac"
	"github.com/cloudstorex/backend/internal/security/zerotrust"
"""
content = re.sub(r'("github.com/cloudstorex/backend/internal/workspace")', r'\1' + imports, content)

# Create engine
engine_init = """
	authEngine := engine.NewAuthorizationEngine(rbac.NewEvaluator(), abac.NewEvaluator(), zerotrust.NewEvaluator())
"""
content = re.sub(r'(tokenService := identity\.NewTokenService\(cfg\.JWTSecret\))', r'\1\n' + engine_init, content)

# Replace AuthMiddleware
content = content.replace("middleware.AuthMiddleware(tokenService)", "middleware.AuthMiddleware(tokenService, authEngine)")

with open('internal/app/app.go', 'w') as f:
    f.write(content)
