package i18n

import (
	"testing"

	"golang.org/x/text/language"
)

func TestLoad(t *testing.T) {
	i := New()

	// Test case 1: Valid language file
	b1 := []byte(`{
		"_.code": "en",
		"_.name": "English",
		"hello": "Hello",
		"goodbye": "Goodbye"
	}`)
	err := i.Load(b1, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify language is loaded correctly
	if len(i.locales) != 1 {
		t.Errorf("expected 1 locale, got %d", len(i.locales))
	}
	if len(i.localesIndex) != 1 {
		t.Errorf("expected 1 locale index, got %d", len(i.localesIndex))
	}

	// Test case 2: Missing _.code field
	b2 := []byte(`{
		"_.name": "English",
		"hello": "Hello",
		"goodbye": "Goodbye"
	}`)
	err = i.Load(b2, false)
	if err == nil {
		t.Error("expected error for missing _.code field")
	} else if err.Error() != "missing _.code field in language file" {
		t.Errorf("unexpected error message: %v", err)
	}

	// Test case 3: Missing _.name field
	b3 := []byte(`{
		"_.code": "en",
		"hello": "Hello",
		"goodbye": "Goodbye"
	}`)
	err = i.Load(b3, false)
	if err == nil {
		t.Error("expected error for missing _.name field")
	} else if err.Error() != "missing _.name field in language file" {
		t.Errorf("unexpected error message: %v", err)
	}

	// Test case 4: Invalid language code
	b4 := []byte(`{
		"_.code": "invalid",
		"_.name": "Invalid",
		"hello": "Hello",
		"goodbye": "Goodbye"
	}`)
	err = i.Load(b4, false)
	if err == nil {
		t.Error("expected error for invalid language code")
	}

	// Test case 5: Overwriting existing language
	b5 := []byte(`{
		"_.code": "en",
		"_.name": "English (Updated)",
		"hello": "Hello",
		"goodbye": "Goodbye"
	}`)
	err = i.Load(b5, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify language is updated correctly
	if len(i.locales) != 1 {
		t.Errorf("expected 1 locale, got %d", len(i.locales))
	}
	if len(i.localesIndex) != 1 {
		t.Errorf("expected 1 locale index, got %d", len(i.localesIndex))
	}
	if i.locales["en"].name != "English (Updated)" {
		t.Errorf("expected updated name, got %s", i.locales["en"].name)
	}
}

func TestParseLanguage(t *testing.T) {
	i := &I18n{
		locales:      make(map[string]locale),
		localesIndex: []string{"en", "fr"},
		matcher: language.NewMatcher([]language.Tag{
			language.English,
			language.French,
		}),
	}

	// Test case 1: Accepting English
	acceptLanguage1 := "en-US,en;q=0.9,fr;q=0.8"
	code := i.ParseLanguage(acceptLanguage1)
	if code != "en" {
		t.Errorf("expected en, got %s", code)
	}

	// Test case 2: Accepting French
	acceptLanguage2 := "fr-FR,fr;q=0.9,en;q=0.8"
	code = i.ParseLanguage(acceptLanguage2)
	if code != "fr" {
		t.Errorf("expected fr, got %s", code)
	}

	// Test case 3: Accepting a language not in the list
	acceptLanguage3 := "es-ES,es;q=0.9,de;q=0.8"
	code = i.ParseLanguage(acceptLanguage3)
	if code != "en" {
		t.Errorf("expected en, got %s", code)
	}
}

func TestName(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				name: "English",
			},
			"fr": {
				name: "French",
			},
		},
	}

	// Test case 1: English
	name := i.Name("en")
	if name != "English" {
		t.Errorf("expected English, got %s", name)
	}

	// Test case 2: French
	name = i.Name("fr")
	if name != "French" {
		t.Errorf("expected French, got %s", name)
	}
}

func TestTag(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				tag: language.English,
			},
			"fr": {
				tag: language.French,
			},
		},
	}

	// Test case 1: English
	tag := i.Tag("en")
	if tag != language.English {
		t.Errorf("expected English, got %s", tag)
	}

	// Test case 2: French
	tag = i.Tag("fr")
	if tag != language.French {
		t.Errorf("expected French, got %s", tag)
	}
}

func TestTranslate(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"hello":   "Hello",
					"goodbye": "Goodbye",
				},
			},
		},
	}

	// Test case 1: English
	translation := i.T("en", "hello")
	if translation != "Hello" {
		t.Errorf("expected Hello, got %s", translation)
	}

	// Test case 2: Missing translation
	translation = i.T("en", "missing")
	if translation != "missing" {
		t.Errorf("expected missing, got %s", translation)
	}
}

func TestTranslateFallback(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"hello": "Hello",
				},
			},
			"fr": {
				langMap: map[string]string{
					"hello": "Bonjour",
				},
			},
		},
	}

	// Test case 1: English
	translation := i.T("en", "hello")
	if translation != "Hello" {
		t.Errorf("expected Hello, got %s", translation)
	}

	// Test case 2: French
	translation = i.T("fr", "hello")
	if translation != "Bonjour" {
		t.Errorf("expected Bonjour, got %s", translation)
	}

	// Test case 3: Missing translation
	translation = i.T("en", "missing")
	if translation != "missing" {
		t.Errorf("expected missing, got %s", translation)
	}
}

func TestTranslatePlural(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"apple": "apple | apples",
				},
			},
		},
	}

	// Test case 1: Singular
	translation := i.Tc("en", "apple", 1)
	if translation != "apple" {
		t.Errorf("expected apple, got %s", translation)
	}

	// Test case 2: Plural
	translation = i.Tc("en", "apple", 2)
	if translation != "apples" {
		t.Errorf("expected apples, got %s", translation)
	}

	// Test case 3: Missing translation
	translation = i.Tc("en", "missing", 1)
	if translation != "missing" {
		t.Errorf("expected missing, got %s", translation)
	}
}

func TestTranslatePluralFallback(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"apple": "apple | apples",
				},
			},
			"fr": {
				langMap: map[string]string{
					"apple": "pomme | pommes",
				},
			},
		},
	}

	// Test case 1: Singular
	translation := i.Tc("en", "apple", 1)
	if translation != "apple" {
		t.Errorf("expected apple, got %s", translation)
	}

	// Test case 2: Plural
	translation = i.Tc("en", "apple", 2)
	if translation != "apples" {
		t.Errorf("expected apples, got %s", translation)
	}

	// Test case 3: French Singular
	translation = i.Tc("fr", "apple", 1)
	if translation != "pomme" {
		t.Errorf("expected pomme, got %s", translation)
	}

	// Test case 4: French Plural
	translation = i.Tc("fr", "apple", 2)
	if translation != "pommes" {
		t.Errorf("expected pommes, got %s", translation)
	}

	// Test case 5: Missing translation
	translation = i.Tc("en", "missing", 1)
	if translation != "missing" {
		t.Errorf("expected missing, got %s", translation)
	}
}

func TestTranslateParams(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"hello": "Hello {name}",
				},
			},
		},
	}

	// Test case 1: English
	translation := i.Ts("en", "hello", "name", "John")
	if translation != "Hello John" {
		t.Errorf("expected Hello John, got %s", translation)
	}

	// Test case 2: Missing params
	translation = i.Ts("en", "hello", "name")
	if translation != "hello: Invalid arguments" {
		t.Errorf("expected hello: Invalid arguments, got %s", translation)
	}
}

func TestTranslatePluralParams(t *testing.T) {
	i := &I18n{
		locales: map[string]locale{
			"en": {
				langMap: map[string]string{
					"apple": "apple {count} | apples {count}",
				},
			},
		},
	}

	// Test case 1: Singular
	translation := i.Tsc("en", "apple", 1, "count", "1")
	if translation != "apple 1" {
		t.Errorf("expected apple 1, got %s", translation)
	}

	// Test case 2: Plural
	translation = i.Tsc("en", "apple", 2, "count", "2")
	if translation != "apples 2" {
		t.Errorf("expected apples 2, got %s", translation)
	}

	// Test case 3: Missing params
	translation = i.Tsc("en", "apple", 1, "count")
	if translation != "apple: Invalid arguments" {
		t.Errorf("expected apple: Invalid arguments, got %s", translation)
	}
}
