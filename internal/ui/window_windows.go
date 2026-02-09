//go:build windows

package ui

/*
#cgo LDFLAGS: -luser32

#include <windows.h>
#include <stdio.h>

void SetWindowOnTop(HWND hwnd, BOOL onTop) {
    if (onTop) {
        SetWindowPos(hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
    } else {
        SetWindowPos(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
    }
}

BOOL CALLBACK SetWindowOnTopCallback(HWND hwnd, LPARAM lParam) {
    char title[256];
    GetWindowTextA(hwnd, title, sizeof(title));
    const char *targetTitle = (const char *)lParam;
    if (strcmp(title, targetTitle) == 0) {
        BOOL onTop = *((BOOL *)lParam + strlen(targetTitle) + 1);
        if (onTop) {
            SetWindowPos(hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
        } else {
            SetWindowPos(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
        }
        return FALSE;
    }
    return TRUE;
}

void SetWindowOnTopByTitle(const char *title, BOOL onTop) {
    // We need to pass both title and onTop through lParam
    // Create a buffer with title + onTop
    char buffer[512];
    strcpy(buffer, title);
    *((BOOL *)(buffer + strlen(title) + 1)) = onTop;
    EnumWindows(SetWindowOnTopCallback, (LPARAM)buffer);
}
*/
import "C"
import (
	"unsafe"
)

// SetWindowOnTop sets the window to stay on top (Windows implementation)
func SetWindowOnTop(windowPtr unsafe.Pointer, onTop bool) {
	C.SetWindowOnTop(C.HWND(windowPtr), C.BOOL(onTop))
}

// SetWindowOnTopByTitle sets the window to stay on top by finding it by title
func SetWindowOnTopByTitle(title string, onTop bool) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.SetWindowOnTopByTitle(cTitle, C.BOOL(onTop))
}
