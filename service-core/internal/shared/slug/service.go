package slug

import (
	"regexp"
	"strings"
)

type generatorImpl struct{}

func NewGenerator() Generator {
	return &generatorImpl{}
}

func (g *generatorImpl) Generate(input string) string {
	slug := strings.ToLower(input)

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
