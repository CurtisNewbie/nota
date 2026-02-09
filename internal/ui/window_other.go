//go:build !darwin && !windows

package ui

import "unsafe"

// SetWindowOnTop is a no-op on unsupported platforms
func SetWindowOnTop(windowPtr unsafe.Pointer, onTop bool) {
	// Not supported on this platform
}

// SetWindowOnTopByTitle is a no-op on unsupported platforms
func SetWindowOnTopByTitle(title string, onTop bool) {
	// Not supported on this platform
}
