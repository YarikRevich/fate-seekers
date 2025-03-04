package translation

import (
	"encoding/json"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	// GetInstance retrieves instance of the answer input manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*TranslationManager](newTranslationManager)
)

// TranslationManager represents translation manager, which acts as a holder
type TranslationManager struct {
	// Represents localizer used for currently selected language.
	localizer *i18n.Localizer
}

// GetTranslation returns translated text for the given key, using currently
// selected language. It also accepts optional template arguments.
func (tm *TranslationManager) GetTranslation(key string, args ...map[string]interface{}) string {
	return tm.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	})
}

// newSubtitlesManager initializes SubtitlesManager.
func newTranslationManager() *TranslationManager {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	bundle.MustParseMessageFileBytes(
		loader.GetInstance().GetTemplate(loader.EnglishTemplate), loader.EnglishTemplate)

	bundle.MustParseMessageFileBytes(
		loader.GetInstance().GetTemplate(loader.UkrainianTemplate), loader.UkrainianTemplate)

	var localizer *i18n.Localizer

	switch config.GetSettingsLanguage() {
	case config.SETTINGS_LANGUAGE_ENGLISH:
		localizer = i18n.NewLocalizer(bundle, config.SETTINGS_LANGUAGE_ENGLISH)

	case config.SETTINGS_LANGUAGE_UKRAINIAN:
		localizer = i18n.NewLocalizer(bundle, config.SETTINGS_LANGUAGE_UKRAINIAN)

	}

	return &TranslationManager{localizer: localizer}
}
