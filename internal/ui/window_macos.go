//go:build darwin

package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

void SetWindowOnTop(void *windowPtr, bool onTop) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (onTop) {
        [window setLevel:NSFloatingWindowLevel];
    } else {
        [window setLevel:NSNormalWindowLevel];
    }
}

void SetWindowOnTopByTitle(const char *title, bool onTop) {
    NSString *windowTitle = [NSString stringWithUTF8String:title];
    NSApplication *app = [NSApplication sharedApplication];
    NSArray *windows = [app windows];

    for (NSWindow *window in windows) {
        if ([[window title] isEqualToString:windowTitle]) {
            if (onTop) {
                [window setLevel:NSFloatingWindowLevel];
            } else {
                [window setLevel:NSNormalWindowLevel];
            }
            break;
        }
    }
}
*/
import "C"
import (
	"unsafe"
)

// SetWindowOnTop sets the window to stay on top (macOS implementation)
func SetWindowOnTop(windowPtr unsafe.Pointer, onTop bool) {
	C.SetWindowOnTop(windowPtr, C.bool(onTop))
}

// SetWindowOnTopByTitle sets the window to stay on top by finding it by title
func SetWindowOnTopByTitle(title string, onTop bool) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.SetWindowOnTopByTitle(cTitle, C.bool(onTop))
}
