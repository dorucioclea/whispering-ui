package SettingsMappings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"whispering-tiger-ui/CustomWidget"
	"whispering-tiger-ui/Settings"
)

var TextToSpeechSettingsMapping = SettingsMapping{
	Mappings: []SettingMapping{
		//{
		//	SettingsName:         "Integrated Text-to-Speech",
		//	SettingsInternalName: "tts_type",
		//	SettingsDescription:  "",
		//	_widget: func() fyne.CanvasObject {
		//		return widget.NewCheck("", func(b bool) {})
		//	},
		//},
		{
			SettingsName:         "Integrated Text-to-Speech",
			SettingsInternalName: "tts_type",
			SettingsDescription:  "",
			_widget: func() fyne.CanvasObject {
				settingWidget := CustomWidget.NewTextValueSelect("tts_type", []CustomWidget.TextValueOption{
					{Text: "Silero", Value: "silero"},
					{Text: "F5/E2", Value: "f5_e2"},
					{Text: lang.L("Disabled"), Value: ""},
				}, func(s CustomWidget.TextValueOption) {
					if Settings.Config.Tts_type != s.Value {
						dialog.ShowInformation(lang.L("App restart required"), lang.L("Changing the TTS Type requires a restart of the application to take effect."), fyne.CurrentApp().Driver().AllWindows()[1])
					}
				}, 0)
				return settingWidget
			},
		},
		{
			SettingsName:         "Speed/Rate",
			SettingsInternalName: "tts_prosody_rate",
			SettingsDescription:  "",
			_widget: func() fyne.CanvasObject {
				settingWidget := CustomWidget.NewTextValueSelect("tts_prosody_rate", []CustomWidget.TextValueOption{
					{Text: lang.L("Default"), Value: ""},
					{Text: "x-slow", Value: "x-slow"},
					{Text: "slow", Value: "slow"},
					{Text: "medium", Value: "medium"},
					{Text: "fast", Value: "fast"},
					{Text: "x-fast", Value: "x-fast"},
				}, func(s CustomWidget.TextValueOption) {}, 0)
				return settingWidget
			},
		},
		{
			SettingsName:         "Route text-to-speech to secondary audio device",
			SettingsInternalName: "tts_use_secondary_playback",
			SettingsDescription:  "Play text-to-speech on a secondary audio device at the same time as the selected output device.\n(By default uses windows default audio device)",
			_widget: func() fyne.CanvasObject {
				return widget.NewCheck("", func(b bool) {})
			},
		},
	},
}
