# Adding a View

This guide walks through adding a new view mode to peep's TUI — a new way to
render the same parsed `Node` tree, selectable alongside tree/graph/raw/search/
filter. See [Architecture](../architecture/architecture.md) for how the pieces
referenced here fit together.

> In a hurry? Jump to the [Checklist](#checklist).

## Overview

Adding a view touches:

1. A new `Mode` constant in `model/app.go`
2. A new sub-model file in `model/`, implementing `SubModel`
3. A field on `App` and its construction in `NewApp`
4. Wiring in `App.Update` (window resize, mode entry key, key routing) and `App.View`
5. A keybinding entry in the status bar hints and the README's keybinding tables

### Flow

```
key press
     ↓
App.Update (model/app.go) ──→ mode-switch case sets a.mode
     ↓
App.View dispatches on a.mode ──→ yourmodel.View()
     ↓
subsequent key presses route to yourmodel.Update() while a.mode == ModeYours
```

## Steps

### 1. Add the `Mode` constant

**File:** `model/app.go`

```go
const (
    ModeTree Mode = iota
    ModeGraph
    ModeRaw
    ModeSearch
    ModeFilter
    ModeYours // add at the end to keep existing values stable
)
```

Also add a case to `modeLabel` for the status bar tag.

### 2. Implement `SubModel`

**New file:** `model/yours.go`

Every view implements the three methods in `model/submodel.go`:

```go
type SubModel interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (SubModel, tea.Cmd)
    View() string
}
```

Follow an existing view as a template depending on what your view needs:

- **Needs the parsed tree** (like `TreeModel`, `GraphModel`) — take a `*Node` (or
  the tree's root) in your constructor.
- **Needs the raw bytes / a scrollable viewport** (like `RawModel`) — wrap a
  `bubbles/viewport.Model`.
- **Needs a text input** (like `SearchModel`, `FilterModel`) — wrap a
  `bubbles/textinput.Model` and call `Init() tea.Cmd { return textinput.Blink }`.

Constructor signature should match the others: `NewYourModel(root *Node, width,
height int) *YourModel`, so it can be built once in `NewApp` alongside the rest.

### 3. Wire it into `App`

**File:** `model/app.go`

Add a field and construct it in `NewApp`, from the same `root`/`data` every other
view is built from — views are built once and kept alive for the process
lifetime, not rebuilt on mode switch:

```go
type App struct {
    // ...
    yours *YourModel
}

func NewApp(data []byte, source string, width, height int) (*App, error) {
    root, err := ParseJSON(data)
    // ...
    return &App{
        // ...
        yours: NewYourModel(root, width, height),
    }, nil
}
```

### 4. Route window-resize, mode entry, and key events

**File:** `model/app.go`, in `App.Update`

Add `a.yours.Update(msg)` to the `tea.WindowSizeMsg` case so it stays sized while
hidden, add a key to enter the mode (mirroring the `case "r":` toggle or the
`case "/":` mode-entry pattern), and add a case to the final `switch a.mode`
block that routes remaining key events to `a.yours.Update` while active.

If your view uses a text input, follow the `ModeSearch`/`ModeFilter` pattern so
`esc` returns to tree mode and `App.InputActive()` correctly reports focus (needed
so global `q` doesn't quit while your input is focused).

### 5. Render it

**File:** `model/app.go`, in `App.View`

Add a case to the mode `switch` that sets `content = a.yours.View()`.

### 6. Document the keybinding

Update the status bar hint string in `App.statusBar()` and the keybinding table
in the project [README](https://github.com/guisolski/peep/blob/main/README.md)
so the new mode is discoverable.

## Checklist

- [ ] `Mode` constant added in `model/app.go`, `modeLabel` updated
- [ ] New `model/*.go` file implements `SubModel` (`Init`, `Update`, `View`)
- [ ] Constructor mirrors existing views: `NewYourModel(root *Node, width, height int)`
- [ ] Field added to `App`, constructed in `NewApp` from the shared `root`/`data`
- [ ] `tea.WindowSizeMsg` forwarded to the new sub-model in `App.Update`
- [ ] A key enters the mode; `esc` exits back to tree if it's a text-input view
- [ ] Remaining keys route to the sub-model while `a.mode == ModeYours`
- [ ] `App.View` renders it
- [ ] Status bar hint string and README keybinding table updated
- [ ] `*_test.go` added for the new sub-model (see `model/tree_test.go` for the pattern)
