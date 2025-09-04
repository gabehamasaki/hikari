package hikari

import (
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var (
	// Regex para remover barras duplas ou múltiplas
	duplicateSlashRegex = regexp.MustCompile(`/+`)
	// Regex para validar o formato do padrão (letras, números, :param, *, -)
	validPatternRegex = regexp.MustCompile(`^[a-zA-Z0-9/:*_-]*$`)
)

func normalizedPattern(pattern string) string {
	if pattern == "" {
		return "/"
	}

	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	pattern = duplicateSlashRegex.ReplaceAllString(pattern, "/")

	if len(pattern) > 1 && strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimSuffix(pattern, "/")
	}

	return pattern
}

func isValidPattern(pattern string) bool {
	if !validPatternRegex.MatchString(pattern) {
		return false
	}

	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, ":") {
			if len(part) == 1 {
				return false // ":"
			}

			paramName := part[1:]
			if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(paramName) {
				return false
			}
		}

		// Wildcard só pode ser "*" sozinho
		if strings.Contains(part, "*") && part != "*" {
			return false
		}
	}

	return true
}

func buildPattern(prefix, pattern string, logger *zap.Logger) string {
	fullPattern := normalizedPattern(prefix + pattern)

	if !isValidPattern(fullPattern) {
		logger.Error("Invalid route pattern", zap.String("pattern", fullPattern))
		return "" // ou lançar um erro, dependendo do design da aplicação
	}

	return fullPattern
}
