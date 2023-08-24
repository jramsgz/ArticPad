package i18n

import "golang.org/x/text/language"

// Available languages.
var serverLangs = []language.Tag{
	language.English, // Default language.
	language.Spanish,
}

// Matcher for the languages.
var matcher = language.NewMatcher(serverLangs)

// ParseLanguage parses the language from the Accept-Language header.
func ParseLanguageHeader(acceptLanguage string) language.Tag {
	t, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	tag, _, _ := matcher.Match(t...)
	return tag
}

// ParseLanguage parses the language string and returns the language code.
func ParseLanguage(lang string) language.Tag {
	tag, _, _ := matcher.Match(language.Make(lang))
	return tag
}
