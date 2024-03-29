package chrome

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Executable returns a string which points to the preferred Chrome executable file.
var Executable = Locate

// Locate returns a path to the Chrome binary, or an empty string if
// Chrome installation is not found.
func Locate() string {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	case "windows":
		paths = []string{
			filepath.Join(os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe"),
			filepath.Join(os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe"),
			filepath.Join(os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe"),
			filepath.Join(os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe"),
			filepath.Join(os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe"),
		}
	default:
		paths = []string{
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

// PromptDownload asks user if he wants to download and install Chrome, and
// opens a download web page if the user agrees.
func PromptDownload() {
	title := "Chrome not found"
	text := "No Chrome/Chromium installation was found. Would you like to download and install it now?"

	// Ask user for confirmation
	if !messageBox(title, text) {
		return
	}

	// Open download page
	url := "https://www.google.com/chrome/"
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Run()
	case "darwin":
		exec.Command("open", url).Run()
	case "windows":
		r := strings.NewReplacer("&", "^&")
		exec.Command("cmd", "/c", "start", r.Replace(url)).Run()
	}
}
