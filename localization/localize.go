// Package localization is responsible for localizing strings for different languages.
//
// Exported function(s): Init, Localize, Localizers
package localization

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizers [7]*i18n.Localizer
var bundle *i18n.Bundle

// Init initializes the localization bundle and localizers
//
// Parameters:
//
//	path (string) - The path to the localization files
//
// Returns: _
func Init(path string) { //3
	bundle = i18n.NewBundle(language.BritishEnglish)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(path + "localization/localizations_src/de/active.de.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/en-GB/active.en-GB.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/en-US/active.en-US.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/es/active.es.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/fr/active.fr.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/nl/active.nl.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/pt/active.pt.json")

	// Localizes for each language
	localizers[0] = i18n.NewLocalizer(bundle, language.German.String())
	localizers[1] = i18n.NewLocalizer(bundle, language.BritishEnglish.String())
	localizers[2] = i18n.NewLocalizer(bundle, language.AmericanEnglish.String())
	localizers[3] = i18n.NewLocalizer(bundle, language.Spanish.String())
	localizers[4] = i18n.NewLocalizer(bundle, language.French.String())
	localizers[5] = i18n.NewLocalizer(bundle, language.Dutch.String())
	localizers[6] = i18n.NewLocalizer(bundle, language.Portuguese.String())
}

// Localize returns the localized string for the given language and ID
//
// Parameters:
//
//	language (int) - The language to localize the message to
//	ID (string) - The ID of the message to localize
//
// Returns: localized string (string)
func Localize(language int, MessageID string) string {
	return localizers[language].MustLocalize(&i18n.LocalizeConfig{MessageID: MessageID}) // Is it confusing to have field MessageID and parameter MessageID?
}

// Localizers returns the localizers variable
//
// Parameters: _
//
// Returns: localizers ([7]*i18n.Localizer)
func Localizers() [7]*i18n.Localizer {
	return localizers
}
