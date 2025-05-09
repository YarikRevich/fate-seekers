package imgui

import (
	"fmt"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the imgui implementation, performing initilization if needed.
	GetInstance = sync.OnceValue[*ImGUI](newImGUI)
)

// ImGUI represents a wrapper for imgui implementation.
type ImGUI struct {
}

// Update updates imgui state.
func (i *ImGUI) Update() {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.SetClipMask(!ebimgui.ClipMask())

	ebimgui.BeginFrame()

	imgui.Text(fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))

	if imgui.BeginMenu("States") {
		imgui.Text(fmt.Sprintf("active_screen: %s", store.GetActiveScreen()))
		imgui.Text(fmt.Sprintf("answer_input_question_updated: %s", store.GetAnswerInputQuestionUpdated()))
		imgui.Text(fmt.Sprintf("entry_handshake_started_networking: %s", store.GetEntryHandshakeStartedNetworking()))
		imgui.Text(fmt.Sprintf("event_name: %s", store.GetEventName()))
		imgui.Text(fmt.Sprintf("letter_image: %s", store.GetLetterImage()))
		imgui.Text(fmt.Sprintf("letter_name: %s", store.GetLetterName()))
		imgui.Text(fmt.Sprintf("letter_updated: %s", store.GetLetterUpdated()))
		imgui.Text(fmt.Sprintf("application_exit: %s", store.GetApplicationExit()))
		imgui.Text(fmt.Sprintf("application_loading: %s", store.GetApplicationLoading()))
		imgui.Text(fmt.Sprintf("prompt_text: %s", store.GetPromptText()))
		imgui.Text(fmt.Sprintf("prompt_updated: %s", store.GetPromptUpdated()))
		imgui.Text(fmt.Sprintf("event_name: %s", store.GetEventName()))
		imgui.Text(fmt.Sprintf("event_started: %s", store.GetEventStarted()))
		imgui.Text(fmt.Sprintf("event_ending: %s", store.GetEventEnding()))
		imgui.Text(fmt.Sprintf("sound_fx_updated: %s", store.GetSoundFXUpdated()))
		imgui.Text(fmt.Sprintf("sound_music_updated: %s", store.GetSoundMusicUpdated()))
		imgui.Text(fmt.Sprintf("prompt_submit_callback: %v", store.GetPromptSubmitCallback()))
		imgui.Text(fmt.Sprintf("prompt_cancel_callback: %v", store.GetPromptCancelCallback()))

		imgui.EndMenu()
	}

	if imgui.BeginMenu("Settings") {
		imgui.Text(fmt.Sprintf("sound_music: %d", config.GetSettingsSoundMusic()))
		imgui.Text(fmt.Sprintf("sound_fx: %d", config.GetSettingsSoundFX()))
		imgui.Text(fmt.Sprintf("language: %s", config.GetSettingsLanguage()))
		imgui.Text(fmt.Sprintf("receiver_port: %s", config.GetSettingsNetworkingReceiverPort()))
		imgui.Text(fmt.Sprintf("server_host: %s", config.GetSettingsNetworkingServerHost()))
		imgui.Text(fmt.Sprintf("encryption_key: %s", config.GetSettingsNetworkingEncryptionKey()))

		imgui.EndMenu()
	}

	ebimgui.EndFrame()
}

// Draw performs imgui interface rendering.
func (i *ImGUI) Draw(screen *ebiten.Image) {
	ebimgui.Draw(screen)
}

// Layout performs imgui interface resizing.
func (i *ImGUI) Layout(outsideWidth, outsideHeight int) {
	ebimgui.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
}

// newImGUI initializes ImGUI.
func newImGUI() *ImGUI {
	return new(ImGUI)
}
