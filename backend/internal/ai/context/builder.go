package context

import (
	"fmt"
	"strings"
)

type Builder interface {
	AddContext(key string, value string)
	Build() string
}

type contextBuilder struct {
	contexts map[string]string
}

func NewBuilder() Builder {
	return &contextBuilder{
		contexts: make(map[string]string),
	}
}

func (b *contextBuilder) AddContext(key, value string) {
	b.contexts[key] = value
}

func (b *contextBuilder) Build() string {
	var sb strings.Builder
	for k, v := range b.contexts {
		sb.WriteString(fmt.Sprintf("[%s]: %s\n", k, v))
	}
	return sb.String()
}
