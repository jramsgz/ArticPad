package i18n

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/text/language"
)

// I18n offers translation functions over different languages.
type I18n struct {
	locales      map[string]locale
	localesIndex []string
	matcher      language.Matcher
}

// locale represents a single language locale.
type locale struct {
	name    string
	tag     language.Tag
	langMap map[string]string
}

var reParam = regexp.MustCompile(`(?i)\{([a-z0-9-.]+)\}`)

// New returns an empty I18n instance.
func New() *I18n {
	return &I18n{
		locales: map[string]locale{},
		matcher: language.NewMatcher([]language.Tag{}),
	}
}

// Load loads a JSON language map into the instance overwriting
// existing keys that conflict.
func (i *I18n) Load(b []byte, isDefault bool) error {
	var l map[string]string
	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	code, ok := l["_.code"]
	if !ok {
		return errors.New("missing _.code field in language file")
	}

	name, ok := l["_.name"]
	if !ok {
		return errors.New("missing _.name field in language file")
	}

	tag, err := language.Parse(code)
	if err != nil {
		return err
	}
	overwrite := false
	if _, ok := i.locales[code]; ok {
		overwrite = true
	}

	i.locales[code] = locale{
		name:    name,
		tag:     tag,
		langMap: l,
	}

	if !overwrite {
		if isDefault {
			i.localesIndex = append([]string{code}, i.localesIndex...)
		} else {
			i.localesIndex = append(i.localesIndex, code)
		}

		tags := []language.Tag{}
		for _, code := range i.localesIndex {
			tags = append(tags, i.locales[code].tag)
		}
		i.matcher = language.NewMatcher(tags)
	}

	return nil
}

// ParseLanguage parses the language string and returns the locale code.
func (i *I18n) ParseLanguage(acceptLanguage string) string {
	_, index := language.MatchStrings(i.matcher, acceptLanguage)
	return i.localesIndex[index]
}

// Name returns the canonical name of the language.
func (i *I18n) Name(code string) string {
	return i.locales[code].name
}

// Tag returns the language tag.
func (i *I18n) Tag(code string) language.Tag {
	return i.locales[code].tag
}

// JSON returns the languagemap as raw JSON.
func (i *I18n) JSON(code string) []byte {
	b, _ := json.Marshal(i.locales[code].langMap)
	return b
}

// T returns the translation for the given key similar to vue i18n's t().
func (i *I18n) T(code, key string) string {
	s, ok := i.locales[code].langMap[key]
	if !ok {
		return key
	}

	return i.getSingular(s)
}

// Ts returns the translation for the given key similar to vue i18n's t()
// and substitutes the params in the given map in the translated value.
// In the language values, the substitutions are represented as: {key}
// The params and values are received as a pairs of succeeding strings.
// That is, the number of these arguments should be an even number.
// eg:
//
//	 Ts("globals.message.notFound",
//		"name", "notes",
//		"error", err)
func (i *I18n) Ts(code, key string, params ...string) string {
	if len(params)%2 != 0 {
		return key + `: Invalid arguments`
	}

	s, ok := i.locales[code].langMap[key]
	if !ok {
		return key
	}

	s = i.getSingular(s)
	for n := 0; n < len(params); n += 2 {
		// If there are {params} in the param values, substitute them.
		val := i.subAllParams(code, params[n+1])
		s = strings.ReplaceAll(s, `{`+params[n]+`}`, val)
	}

	return s
}

// Tc returns the translation for the given key similar to vue i18n's tc().
// It expects the language string in the map to be of the form `Singular | Plural` and
// returns `Plural` if n > 1, or `Singular` otherwise.
func (i *I18n) Tc(code, key string, n int) string {
	s, ok := i.locales[code].langMap[key]
	if !ok {
		return key
	}

	// Plural.
	if n > 1 {
		return i.getPlural(s)
	}

	return i.getSingular(s)
}

// Tsc returns the translation for the given key similar to vue i18n's tc()
// and substitutes the params in the given map in the translated value.
// In the language values, the substitutions are represented as: {key}
// The params and values are received as a pairs of succeeding strings.
// That is, the number of these arguments should be an even number.
// eg:
//
//	 Tsc("globals.message.notFound",
//		"name", "notes",
//		"error", err)
func (i *I18n) Tsc(code, key string, n int, params ...string) string {
	if len(params)%2 != 0 {
		return key + `: Invalid arguments`
	}

	s, ok := i.locales[code].langMap[key]
	if !ok {
		return key
	}

	s = i.Tc(code, key, n)
	for n := 0; n < len(params); n += 2 {
		// If there are {params} in the param values, substitute them.
		val := i.subAllParams(code, params[n+1])
		s = strings.ReplaceAll(s, `{`+params[n]+`}`, val)
	}

	return s
}

// getSingular returns the singular term from the vuei18n pipe separated value.
// singular term | plural term
func (i *I18n) getSingular(s string) string {
	if !strings.Contains(s, "|") {
		return s
	}

	return strings.TrimSpace(strings.Split(s, "|")[0])
}

// getPlural returns the plural term from the vuei18n pipe separated value.
// singular term | plural term
func (i *I18n) getPlural(s string) string {
	if !strings.Contains(s, "|") {
		return s
	}

	chunks := strings.Split(s, "|")
	if len(chunks) == 2 {
		return strings.TrimSpace(chunks[1])
	}

	return strings.TrimSpace(chunks[0])
}

// subAllParams recursively resolves and replaces all {params} in a string.
func (i *I18n) subAllParams(code, s string) string {
	if !strings.Contains(s, `{`) {
		return s
	}

	parts := reParam.FindAllStringSubmatch(s, -1)
	if len(parts) < 1 {
		return s
	}

	for _, p := range parts {
		s = strings.ReplaceAll(s, p[0], i.T(code, p[1]))
	}

	return i.subAllParams(code, s)
}
