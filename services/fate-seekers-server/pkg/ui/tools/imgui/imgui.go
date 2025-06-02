package imgui

import (
	"fmt"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/store"
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
		imgui.Text(fmt.Sprintf("listener_started_networking: %s", store.GetListenerStartedNetworking()))
		imgui.Text(fmt.Sprintf("application_exit: %s", store.GetApplicationExit()))
		imgui.Text(fmt.Sprintf("application_loading: %s", store.GetApplicationLoading()))
		imgui.Text(fmt.Sprintf("prompt_text: %s", store.GetPromptText()))
		imgui.Text(fmt.Sprintf("prompt_updated: %s", store.GetPromptUpdated()))
		imgui.Text(fmt.Sprintf("prompt_submit_callback: %v", store.GetPromptSubmitCallback()))
		imgui.Text(fmt.Sprintf("prompt_cancel_callback: %v", store.GetPromptCancelCallback()))

		imgui.EndMenu()
	}

	if imgui.BeginMenu("Settings") {
		imgui.Text(fmt.Sprintf("language: %s", config.GetSettingsLanguage()))
		imgui.Text(fmt.Sprintf("server_port: %s", config.GetSettingsNetworkingServerPort()))
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
