package sound

import (
	"bytes"
	"fmt"
	"io"

	"td/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 48000

// Effect identifies one playable sound effect.
type Effect int

const (
	// EffectRaiderDefeated is played when combat damage defeats a raider.
	EffectRaiderDefeated Effect = iota
)

// Manager owns Ebitengine audio playback for one-shot sound effects.
type Manager struct {
	context *audio.Context
	effects map[Effect][]byte
	players []*audio.Player
}

// NewManager creates a runtime sound manager from loaded audio assets.
func NewManager(audioCatalog assets.AudioCatalog) (*Manager, error) {
	raiderDefeated, err := decodeWAVF32(audioCatalog.RaiderDefeated)
	if err != nil {
		return nil, fmt.Errorf("prepare raider defeated sound: %w", err)
	}
	return &Manager{
		context: audio.NewContext(sampleRate),
		effects: map[Effect][]byte{
			EffectRaiderDefeated: raiderDefeated,
		},
	}, nil
}

// Play starts a one-shot sound effect when the audio context is ready.
func (m *Manager) Play(effect Effect) {
	if m == nil || m.context == nil || !m.context.IsReady() {
		return
	}
	data, ok := m.effects[effect]
	if !ok || len(data) == 0 {
		return
	}
	m.prunePlayers()
	player := m.context.NewPlayerF32FromBytes(data)
	player.SetVolume(0.55)
	player.Play()
	m.players = append(m.players, player)
}

// PlayRaiderDefeated plays the raider-defeated sound effect.
func (m *Manager) PlayRaiderDefeated() {
	m.Play(EffectRaiderDefeated)
}

// Update releases finished one-shot players.
func (m *Manager) Update() {
	if m == nil {
		return
	}
	m.prunePlayers()
}

// decodeWAVF32 decodes WAV data into Ebitengine's 32-bit float stereo byte format.
func decodeWAVF32(data []byte) ([]byte, error) {
	stream, err := wav.DecodeF32(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	decoded, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	if len(decoded) == 0 {
		return nil, fmt.Errorf("decoded WAV is empty")
	}
	return decoded, nil
}

// prunePlayers keeps only players that are still audible.
func (m *Manager) prunePlayers() {
	active := m.players[:0]
	for _, player := range m.players {
		if player.IsPlaying() {
			active = append(active, player)
		}
	}
	m.players = active
}
