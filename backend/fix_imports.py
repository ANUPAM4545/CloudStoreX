import re

with open('internal/app/app.go', 'r') as f:
    content = f.read()

# Add imports right after "github.com/cloudstorex/backend/internal/identity"
imports = """
	"github.com/cloudstorex/backend/internal/security/abac"
	authEnginePkg "github.com/cloudstorex/backend/internal/security/engine"
	"github.com/cloudstorex/backend/internal/security/rbac"
	"github.com/cloudstorex/backend/internal/security/zerotrust"
"""
content = re.sub(r'("github\.com/cloudstorex/backend/internal/identity"\n)', r'\1' + imports, content)

# Fix engine.NewAuthorizationEngine to authEnginePkg.NewAuthorizationEngine
content = content.replace("engine.NewAuthorizationEngine", "authEnginePkg.NewAuthorizationEngine")

with open('internal/app/app.go', 'w') as f:
    f.write(content)
