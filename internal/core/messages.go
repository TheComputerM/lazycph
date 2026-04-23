package core

// NavigateMsg is sent to navigate to a new path.
// If Path is empty, the app reverts to the active model
// derived from the current state, simulating popup dismissal.
// If Path is non-empty, the app state updates and the
// active model is derived from the new path (dir → filepicker, file → workspace).
type NavigateMsg struct {
	Path string
}
