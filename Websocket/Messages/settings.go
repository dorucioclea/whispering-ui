package Messages

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"strings"
	"whispering-tiger-ui/Fields"
	"whispering-tiger-ui/Settings"
	"whispering-tiger-ui/Utilities"
)

type WhisperLanguage struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type TranslateSetting struct {
	WhisperLanguages         []WhisperLanguage `json:"whisper_languages"`
	OscAutoProcessingEnabled bool              `json:"osc_auto_processing_enabled"`
	Settings.Conf
}

var TranslateSettings TranslateSetting

func (res TranslateSetting) Update() *TranslateSetting {

	Settings.Form = Settings.BuildSettingsForm(nil, Settings.Config.SettingsFilename).(*widget.Form)
	Settings.Form.Refresh()

	log.Println("InstalledLanguages.GetNameByCode()")
	log.Println(InstalledLanguages.GetNameByCode(res.Trg_lang))
	log.Println(res.Trg_lang)

	// fill combo-box with whisper languages
	if len(Fields.Field.TranscriptionSpeakerLanguageCombo.Options) < len(TranslateSettings.WhisperLanguages) {
		Fields.Field.TranscriptionSpeakerLanguageCombo.Options = nil
		for _, element := range TranslateSettings.WhisperLanguages {
			Fields.Field.TranscriptionSpeakerLanguageCombo.Options = append(Fields.Field.TranscriptionSpeakerLanguageCombo.Options, cases.Title(language.English, cases.Compact).String(element.Name))
		}
	}

	// Set options to current settings
	if strings.Contains(res.Whisper_task, "translate") && Fields.Field.TranscriptionTaskCombo.Selected != "translate (to English)" {
		Fields.Field.TranscriptionTaskCombo.SetSelected("translate (to English)")
	}
	if strings.Contains(res.Whisper_task, "transcribe") && !strings.Contains(Fields.Field.TranscriptionTaskCombo.Selected, "transcribe") {
		Fields.Field.TranscriptionTaskCombo.SetSelected("transcribe")
	}
	if Fields.Field.TranscriptionSpeakerLanguageCombo.Selected != TranslateSettings.GetWhisperLanguageNameByCode(res.Current_language) {
		Fields.Field.TranscriptionSpeakerLanguageCombo.SetSelected(
			cases.Title(language.English, cases.Compact).String(TranslateSettings.GetWhisperLanguageNameByCode(res.Current_language)),
		)
	}

	// Set SourceLanguageCombo
	if strings.ToLower(Fields.Field.SourceLanguageCombo.Selected) != strings.ToLower(InstalledLanguages.GetNameByCode(res.Src_lang)) {
		Fields.Field.SourceLanguageCombo.SetSelected(cases.Title(language.English, cases.Compact).String(InstalledLanguages.GetNameByCode(res.Src_lang)))
	} else if Fields.Field.SourceLanguageCombo.Selected == "" && res.Src_lang == "auto" {
		Fields.Field.SourceLanguageCombo.SetSelected(cases.Title(language.English, cases.Compact).String(res.Src_lang))
	}

	// Set TargetLanguageCombo
	if strings.ToLower(Fields.Field.TargetLanguageCombo.Selected) != strings.ToLower(InstalledLanguages.GetNameByCode(res.Trg_lang)) {
		Fields.Field.TargetLanguageCombo.SetSelected(cases.Title(language.English, cases.Compact).String(InstalledLanguages.GetNameByCode(res.Trg_lang)))
		// Set TargetLanguageTxtTranslateCombo if it is not set
		if Fields.Field.TargetLanguageTxtTranslateCombo.Selected == "" {
			Fields.Field.TargetLanguageTxtTranslateCombo.SetSelected(cases.Title(language.English, cases.Compact).String(InstalledLanguages.GetNameByCode(res.Trg_lang)))
		}
	}

	checkValue, _ := Fields.DataBindings.SpeechToTextEnabledDataBinding.Get()
	if checkValue != res.Stt_enabled {
		Fields.DataBindings.SpeechToTextEnabledDataBinding.Set(res.Stt_enabled)
	}

	checkValue, _ = Fields.DataBindings.TextTranslateEnabledDataBinding.Get()
	if checkValue != res.Txt_translate {
		Fields.DataBindings.TextTranslateEnabledDataBinding.Set(res.Txt_translate)
	}

	checkValue, _ = Fields.DataBindings.TextToSpeechEnabledDataBinding.Get()
	if checkValue != res.Tts_answer {
		Fields.DataBindings.TextToSpeechEnabledDataBinding.Set(res.Tts_answer)
	}
	checkValue, _ = Fields.DataBindings.OSCEnabledDataBinding.Get()
	if checkValue != res.OscAutoProcessingEnabled {
		Fields.DataBindings.OSCEnabledDataBinding.Set(res.OscAutoProcessingEnabled)
	}

	// Set TtsModelCombo
	if len(res.Tts_model) > 0 && len(Fields.Field.TtsModelCombo.Options) > 0 && Fields.Field.TtsModelCombo.Selected != res.Tts_model[1] {
		Fields.Field.TtsModelCombo.SetSelected(res.Tts_model[1])
	}

	// Set TtsVoiceCombo
	// only set new tts voice if select is not received tts_voice and
	// if select is not empty and does not contain only one empty element
	if Fields.Field.TtsVoiceCombo.Selected != res.Tts_voice && (len(Fields.Field.TtsVoiceCombo.Options) > 0 &&
		(len(Fields.Field.TtsVoiceCombo.Options) == 1 && Fields.Field.TtsVoiceCombo.Options[0] != "")) {
		Fields.Field.TtsVoiceCombo.SetSelected(res.Tts_voice)
	}
	// Set OcrWindowCombo
	if Fields.Field.OcrWindowCombo.Selected != res.Ocr_window_name {
		if !Utilities.Contains(Fields.Field.OcrWindowCombo.Options, res.Ocr_window_name) {
			Fields.Field.OcrWindowCombo.Options = append(Fields.Field.OcrWindowCombo.Options, res.Ocr_window_name)
		}
		Fields.Field.OcrWindowCombo.SetSelected(res.Ocr_window_name)
	}
	//}
	// Set OcrLanguageCombo
	if Fields.Field.OcrLanguageCombo.Selected != res.Ocr_lang {
		Fields.Field.OcrLanguageCombo.SetSelected(OcrLanguagesList.GetNameByCode(res.Ocr_lang))
	}

	// set oscEnabledLabel Update function
	if res.OscAutoProcessingEnabled {
		Fields.Field.OscLimitHint.Show()
	} else {
		Fields.Field.OscLimitHint.Hide()
	}
	Fields.OscLimitHintUpdateFunc = func() {
		transcriptionInputCount := Utilities.CountUTF16CodeUnits(Fields.Field.TranscriptionInput.Text)
		transcriptionTranslationInputCount := Utilities.CountUTF16CodeUnits(Fields.Field.TranscriptionTranslationInput.Text)
		oscSplitCount := Utilities.CountUTF16CodeUnits(Settings.Config.Osc_type_transfer_split)
		maxCount := res.Conf.Osc_chat_limit

		Fields.Field.OscLimitHint.Text = fmt.Sprintf(Fields.OscLimitLabelConst, 0, maxCount)
		switch res.Conf.Osc_type_transfer {
		case "source":
			Fields.Field.OscLimitHint.Text = fmt.Sprintf(Fields.OscLimitLabelConst, transcriptionInputCount, maxCount)
		case "translation_result":
			Fields.Field.OscLimitHint.Text = fmt.Sprintf(Fields.OscLimitLabelConst, transcriptionTranslationInputCount, maxCount)
		case "both":
			Fields.Field.OscLimitHint.Text = fmt.Sprintf(Fields.OscLimitLabelConst, transcriptionInputCount+oscSplitCount+transcriptionTranslationInputCount, maxCount)
		}
		Fields.Field.OscLimitHint.Refresh()
	}
	Fields.OscLimitHintUpdateFunc()

	Settings.Config = res.Conf

	return &res
}

func (res TranslateSetting) GetWhisperLanguageCodeByName(name string) string {
	for _, entry := range res.WhisperLanguages {
		if strings.ToLower(entry.Name) == strings.ToLower(name) {
			return entry.Code
		}
	}
	return ""
}

func (res TranslateSetting) GetWhisperLanguageNameByCode(code string) string {
	for _, entry := range res.WhisperLanguages {
		if entry.Code == code {
			return entry.Name
		}
	}
	return ""
}
