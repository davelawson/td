# Add Raider Defeated Sound

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file as `plans/23-raider-defeated-sound.md`.

## Purpose / Big Picture

After this change, the prototype has its first sound-effect path. A player can start a Raid, let the starting towers defeat a raider, and hear a short prototype defeat sound when the enemy is removed by combat damage. This proves that sounds can be triggered from gameplay events without putting Ebitengine audio details inside combat rules.

## Progress

- [x] (2026-05-25 00:00Z) Created this ExecPlan for the accepted sound implementation.
- [x] (2026-05-25 00:00Z) Added a generated prototype WAV under `assets/audio/`.
- [x] (2026-05-25 00:00Z) Extended the asset catalog with embedded audio bytes and decoding validation tests.
- [x] (2026-05-25 00:00Z) Added an app-owned sound manager that uses Ebitengine audio and a no-op-safe gameplay sound sink.
- [x] (2026-05-25 00:00Z) Triggered the sound exactly when tower damage defeats a raider.
- [x] (2026-05-25 00:00Z) Updated control documents that describe current product behavior, architecture, roadmap, and game design.
- [x] (2026-05-25 00:00Z) Ran `go test ./...`, `git diff --check`, and a hand-written code-file line-count review.

## Surprises & Discoveries

- Observation: `newApp()` is used by tests, so constructing an Ebitengine `audio.Context` there would risk global-context panics across tests.
  Evidence: `cmd/td/main_test.go` calls `newApp()` in multiple tests, while Ebitengine documents one audio context per process.

## Decision Log

- Decision: Keep `internal/game` decoupled from Ebitengine audio with a small sound sink interface and default no-op sink.
  Rationale: Combat should report a gameplay event, while app runtime owns audio playback and tests remain deterministic.
  Date/Author: 2026-05-25 / Codex

- Decision: Use a generated embedded WAV for the first effect.
  Rationale: The user explicitly chose an embedded WAV, and generating it locally avoids third-party licensing uncertainty.
  Date/Author: 2026-05-25 / Codex

## Outcomes & Retrospective

Implemented a generated embedded WAV asset, typed audio catalog loading, a runtime `internal/sound` manager, a game-side sound sink, and the combat defeat trigger. `go test ./...` and `git diff --check` passed.

The final hand-written code-file line-count review found `internal/game/game_test.go` at 605 lines. This file was already over the 600-line preference before the sound change, and this plan did not add to it. Recommended response: defer a focused test-file split by behavior area, likely separating pause/menu tests or screenshot capture helpers from general game-state tests, before future changes add more test cases there. No unplanned refactor was performed.

## Context and Orientation

The project is a Go/Ebitengine local tower-defense prototype. `cmd/td/main.go` owns process startup and routes Ebitengine callbacks to menu or game state. `internal/game` owns testable gameplay state including Raids, enemies, and combat. `assets/catalog.go` embeds and loads runtime sprites and should also own static audio bytes. There is no current audio package or asset path.

Raider defeat currently happens in `internal/game/combat.go` inside `damageEnemy`, which subtracts health and removes an enemy when health reaches zero. Enemies removed by reaching the Sanctum are handled separately in `internal/game/raid.go` and must not trigger this sound.

`PRODUCT.md` and `README.md` describe current user-visible behavior and must mention the new audible defeat feedback. `ARCHITECTURE.md` must mention the app-owned sound manager and audio asset boundary. `GAME.md` should record the design decision that defeated raiders produce prototype sound feedback. `ROADMAP.md` should no longer imply audio is entirely future-only.

## Plan of Work

First create `assets/audio/raider-defeated.wav` as a short locally generated stereo WAV. Then update `assets/catalog.go` so the catalog includes an `AudioCatalog` with `RaiderDefeated []byte`, loaded from embedded files without decoding or playing in the asset package. Add an asset test that confirms the bytes are non-empty and can be decoded through Ebitengine's WAV decoder.

Create `internal/sound` to own Ebitengine audio playback. It should define `Effect`, `EffectRaiderDefeated`, `Manager`, `NewManager`, `Play`, `PlayRaiderDefeated`, and `Update`. `NewManager` receives `assets.AudioCatalog`, creates one `audio.Context`, decodes the WAV to 32-bit float stereo bytes once, and stores those bytes by effect. `Play` creates a fresh player from the decoded bytes so overlapping one-shot effects can play. `Update` prunes finished players.

Update `internal/game` with a tiny sound sink interface. `State` should default to a no-op sink in `New`, expose `SetSoundSink`, and call `PlayRaiderDefeated` only in `damageEnemy` after health reaches zero and before or after removing the enemy. Tests should use a fake sink to verify lethal combat damage records one event, nonlethal damage records none, and Sanctum contact records none.

Update `cmd/td/main.go` so `main` creates the runtime sound manager once using `assets.NewAudioCatalog()` and passes it to app construction. Keep `newApp()` no-op for tests, add a helper such as `newAppWithSound(game.SoundSink)`, and call `soundManager.Update()` from `app.Update()` when present. When `startGame` creates `game.State`, attach the sound sink.

Finally update the relevant control documents and run validation commands.

## Concrete Steps

From `/home/dave/dev/ai/td`, generate the WAV:

    python3 - <<'PY'
    # create assets/audio/raider-defeated.wav
    PY

Edit the Go files with `gofmt` afterward:

    gofmt -w assets/catalog.go assets/catalog_test.go cmd/td/main.go cmd/td/main_test.go internal/game/game.go internal/game/combat.go internal/game/combat_test.go internal/game/raid_test.go internal/sound/*.go

Run validation:

    go test ./...
    git diff --check
    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

If any hand-written code file exceeds 600 lines, report the file and count in `Outcomes & Retrospective`, recommend a concrete response, and do not perform an unplanned split or refactor without user approval.

## Validation and Acceptance

Acceptance is met when `go test ./...` passes, whitespace checks pass, and a manual run of `go run ./cmd/td` lets a player start a Raid and hear a short defeat sound each time a raider is defeated by tower damage. The sound must not play when enemies spend Barricade charges, when the Sanctum is breached and the active Raid is cleared, or when a projectile disappears because its target is already gone.

The documentation acceptance check is that `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, `GAME.md`, and `ROADMAP.md` all match the implemented behavior and do not claim a full settings, music, volume, or audio-options system exists.

## Idempotence and Recovery

The WAV generation is safe to rerun because it overwrites only `assets/audio/raider-defeated.wav`. Go tests can be run repeatedly. If audio manager construction fails because the WAV cannot decode, regenerate the WAV and rerun tests before changing gameplay code.

## Artifacts and Notes

No screenshot artifact is required because rendered output does not change. The important evidence is test output, `git diff --check`, and manual audible verification.

## Interfaces and Dependencies

Use the existing `github.com/hajimehoshi/ebiten/v2` dependency, specifically `github.com/hajimehoshi/ebiten/v2/audio` and `github.com/hajimehoshi/ebiten/v2/audio/wav`. Do not add a new Go module dependency.

Define the game-side sink in `internal/game`:

    type SoundSink interface {
        PlayRaiderDefeated()
    }

Define the runtime manager in `internal/sound`:

    type Effect int
    const EffectRaiderDefeated Effect = iota
    func NewManager(audioCatalog assets.AudioCatalog) (*Manager, error)
    func (m *Manager) Play(effect Effect)
    func (m *Manager) PlayRaiderDefeated()
    func (m *Manager) Update()
