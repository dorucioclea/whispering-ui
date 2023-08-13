package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"whispering-tiger-ui/Pages"
	"whispering-tiger-ui/Pages/Advanced"
	"whispering-tiger-ui/Resources"
	"whispering-tiger-ui/RuntimeBackend"
	"whispering-tiger-ui/Settings"
	"whispering-tiger-ui/UpdateUtility"
	"whispering-tiger-ui/Utilities"
	"whispering-tiger-ui/Websocket"
)

var WebsocketClient = Websocket.NewClient("127.0.0.1:5000")

func overwriteFyneFont() {
	pwd, _ := filepath.Abs("./")
	if _, err := os.Stat(pwd + "\\" + "GoNoto.ttf"); err == nil {
		if err := os.Setenv("FYNE_FONT", pwd+"\\"+"GoNoto.ttf"); err != nil {
			fmt.Printf("WARNING: failed to set FYNE_FONT=%s: %v\n", pwd+"\\"+"GoNoto.ttf", err)
		}
		return
	}

	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		winDir := os.Getenv("WINDIR")
		if len(winDir) == 0 {
			return
		}
		fontPath := determineWindowsFont(winDir + "\\Fonts")
		if err := os.Setenv("FYNE_FONT", fontPath); err != nil {
			fmt.Printf("WARNING: failed to set FYNE_FONT=%s: %v\n", fontPath, err)
		}
	}
}

func determineWindowsFont(fontsDir string) string {
	font := "YuGothM.ttc"
	if _, err := os.Stat(fontsDir + "\\" + font); err == nil {
		return fontsDir + "\\" + font
	}
	font = "meiryo.ttc"
	if _, err := os.Stat(fontsDir + "\\" + font); err == nil {
		return fontsDir + "\\" + font
	}
	font = "msgothic.ttc"
	if _, err := os.Stat(fontsDir + "\\" + font); err == nil {
		return fontsDir + "\\" + font
	}
	font = "segoeui.ttf"
	if _, err := os.Stat(fontsDir + "\\" + font); err == nil {
		return fontsDir + "\\" + font
	}
	return ""
}

func main() {
	// 1. Set up a deferred function to handle panics
	defer func() {
		if r := recover(); r != nil {
			// 2. Create a log file when a crash occurs
			logFile, err := os.OpenFile("crash.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Printf("Failed to create log file: %v\n", err)
				os.Exit(1)
			}
			defer logFile.Close()

			// Use logFile as the output for the log package
			log.SetOutput(logFile)

			// 3. Capture the stack trace and log the error details
			stackTrace := debug.Stack()
			log.Printf("Panic occurred: %v\nStack trace:\n%s", r, stackTrace)

			fmt.Println("A crash occurred. Check the crash.log file for more information.")
		}
	}()

	// main application
	val, ok := os.LookupEnv("WT_SCALE")
	if ok {
		_ = os.Setenv("FYNE_SCALE", val)
	}

	a := app.NewWithID("tiger.whispering")
	a.SetIcon(Resources.ResourceAppIconPng)

	a.Settings().SetTheme(&AppTheme{})

	w := a.NewWindow("Whispering Tiger")
	w.SetMaster()
	w.CenterOnScreen()

	w.SetOnClosed(func() {
		fyne.CurrentApp().Preferences().SetFloat("MainWindowWidth", float64(w.Canvas().Size().Width))
		fyne.CurrentApp().Preferences().SetFloat("MainWindowHeight", float64(w.Canvas().Size().Height))
	})

	// initialize whisper process
	//var whisperProcess = RuntimeBackend.NewWhisperProcess()
	//whisperProcess.DeviceIndex = Settings.Config.Device_index.(string)
	//whisperProcess.DeviceOutIndex = Settings.Config.Device_out_index.(string)
	//whisperProcess.SettingsFile = Settings.Config.SettingsFilename
	//RuntimeBackend.BackendsList = append(RuntimeBackend.BackendsList, whisperProcess)

	profileWindow := a.NewWindow("Whispering Tiger Profiles")

	onProfileClose := func() {

		RuntimeBackend.BackendsList = append(RuntimeBackend.BackendsList, RuntimeBackend.NewWhisperProcess())
		RuntimeBackend.BackendsList[0].DeviceIndex = strconv.Itoa(Settings.Config.Device_index.(int))
		RuntimeBackend.BackendsList[0].DeviceOutIndex = strconv.Itoa(Settings.Config.Device_out_index.(int))
		RuntimeBackend.BackendsList[0].SettingsFile = Settings.Config.SettingsFilename
		// Setting this to use UTF-8 encoding for Python does not work when build using PyInstaller
		RuntimeBackend.BackendsList[0].AttachEnvironment("PYTHONIOENCODING", "UTF-8")
		RuntimeBackend.BackendsList[0].AttachEnvironment("PYTHONLEGACYWINDOWSSTDIO", "UTF-8")
		RuntimeBackend.BackendsList[0].AttachEnvironment("PYTHONUTF8", "1")
		// RuntimeBackend.BackendsList[0].AttachEnvironment("CUBLAS_WORKSPACE_CONFIG", ":4096:8")
		if Utilities.FileExists("ffmpeg/bin/ffmpeg.exe") {
			appExec, _ := os.Executable()
			appPath := filepath.Dir(appExec)

			RuntimeBackend.BackendsList[0].AttachEnvironment("Path", filepath.Join(appPath, "ffmpeg/bin"))
		}
		if Settings.Config.Run_backend {
			if !fyne.CurrentApp().Preferences().BoolWithFallback("DisableUiDownloads", false) {
				RuntimeBackend.BackendsList[0].UiDownload = true
			}
			RuntimeBackend.BackendsList[0].Start()
		}

		// initialize main window
		appTabs := container.NewAppTabs(
			container.NewTabItemWithIcon("Speech-to-Text", theme.NewThemedResource(Resources.ResourceSpeechToTextIconSvg), Pages.CreateSpeechToTextWindow()),
			container.NewTabItemWithIcon("Text-Translate", theme.NewThemedResource(Resources.ResourceTranslateIconSvg), Pages.CreateTextTranslateWindow()),
			container.NewTabItemWithIcon("Text-to-Speech", theme.NewThemedResource(Resources.ResourceTextToSpeechIconSvg), Pages.CreateTextToSpeechWindow()),
			container.NewTabItemWithIcon("Image-to-Text", theme.NewThemedResource(Resources.ResourceImageRecognitionIconSvg), Pages.CreateOcrWindow()),
			container.NewTabItemWithIcon("Plugins", theme.NewThemedResource(Resources.ResourcePluginsIconSvg), Advanced.CreatePluginSettingsPage()),
			container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), Pages.CreateSettingsWindow()),
			container.NewTabItemWithIcon("Advanced", theme.MoreVerticalIcon(), Pages.CreateAdvancedWindow()),
		)
		appTabs.SetTabLocation(container.TabLocationTop)

		appTabs.OnSelected = func(tab *container.TabItem) {
			if tab.Text == "Settings" {
				tab.Content = Pages.CreateSettingsWindow()
				tab.Content.Refresh()
			}
			if tab.Text == "Plugins" {
				tab.Content.(*container.Scroll).Content = Advanced.CreatePluginSettingsPage()
				tab.Content.(*container.Scroll).Content.Refresh()
				tab.Content.(*container.Scroll).Refresh()
			}
			if tab.Text == "Advanced" {
				// check if tab content is of type container.AppTabs
				if tabContent, ok := tab.Content.(*container.AppTabs); ok {
					if tabContent.Selected().Text == "Advanced Settings" {
						// force trigger onselect for (Advanced -> Settings) Tab
						tabContent.OnSelected(tabContent.Items[tabContent.SelectedIndex()])
					}
				}
			}
		}

		w.SetContent(appTabs)

		// set main window size
		mainWindowWidth := fyne.CurrentApp().Preferences().FloatWithFallback("MainWindowWidth", 1200)
		mainWindowHeight := fyne.CurrentApp().Preferences().FloatWithFallback("MainWindowHeight", 600)

		w.Resize(fyne.NewSize(float32(mainWindowWidth), float32(mainWindowHeight)))

		// set websocket client to configured ip+port
		WebsocketClient.Addr = Settings.Config.Websocket_ip + ":" + strconv.Itoa(Settings.Config.Websocket_port)
		go WebsocketClient.Start()

		// show main window
		w.Show()

		fyne.CurrentApp().Preferences().SetFloat("ProfileWindowWidth", float64(profileWindow.Canvas().Size().Width))
		fyne.CurrentApp().Preferences().SetFloat("ProfileWindowHeight", float64(profileWindow.Canvas().Size().Height))

		// close profile window
		profileWindow.Close()
	}

	profilePage := Pages.CreateProfileWindow(onProfileClose)
	profileWindow.SetContent(profilePage)

	// set profile window size
	profileWindowWidth := fyne.CurrentApp().Preferences().FloatWithFallback("ProfileWindowWidth", 1400)
	profileWindowHeight := fyne.CurrentApp().Preferences().FloatWithFallback("ProfileWindowHeight", 600)
	profileWindow.Resize(fyne.NewSize(float32(profileWindowWidth), float32(profileWindowHeight)))

	profileWindow.CenterOnScreen()
	profileWindow.Show()

	// check for updates
	if fyne.CurrentApp().Preferences().BoolWithFallback("CheckForUpdateAtStartup", true) {
		go func() {
			if len(fyne.CurrentApp().Driver().AllWindows()) == 2 {
				UpdateUtility.VersionCheck(fyne.CurrentApp().Driver().AllWindows()[1], false)
			}
		}()
	}

	a.Run()

	// after run (app exit), send whisper process signal to stop
	if len(RuntimeBackend.BackendsList) > 0 {
		RuntimeBackend.BackendsList[0].Stop()
		RuntimeBackend.BackendsList[0].WriterBackend.Close()
		RuntimeBackend.BackendsList[0].ReaderBackend.Close()
	}
}
