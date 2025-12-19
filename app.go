package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	openzapr "zaprLauncher/backend/openZapr"
	"zaprLauncher/backend/update"
	"zaprLauncher/backend/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	projectDir  string
	exeDir      string
	ExeFilePath string

	versionFilePath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resCh := make(chan UpdateResult, 1)

	go func() {

		projectDir := utils.GetAppDataPath("ZaprUI")
		exeDir := projectDir + "/bin"

		// creating ProjectDir in User/AppData/Roaming/
		if err := ensureAppDir(projectDir); err != nil {
			resCh <- UpdateResult{err: fmt.Errorf("❗error creating app directory in AppData: %w", err)}
			return
		}
		if err := ensureAppDir(exeDir); err != nil { // temp dir for sessions data files
			resCh <- UpdateResult{err: fmt.Errorf("❗error creating gitrepo directory in project dir: %w", err)}
			return
		}

		// =============================================== fetching zaprUI

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		release, err := update.ParceLatestRelease(client) // Asking GitHub Releases about latest
		if err != nil {
			resCh <- UpdateResult{err: fmt.Errorf("❗error parce latest release: %v", err)}
			return
		}

		if err := update.EnsureVersionFileExist(projectDir, release); err != nil { //  Check VersionFile
			resCh <- UpdateResult{err: fmt.Errorf("❗version file ensure error: %v", err)}
			return
		}
		versionFilePath := filepath.Join(projectDir, "zaprUI_version.txt")

		// CHECKING Latest and Ready
		latest, err := update.IsLatestVersion(versionFilePath, release) // Trying version
		if err != nil {
			resCh <- UpdateResult{err: fmt.Errorf("❗failed to check version: %v", err)}
			return
		}
		ready, err := update.IsReleaseReady(exeDir)
		if err != nil {
			resCh <- UpdateResult{err: fmt.Errorf("❗failed to check release files: %v", err)}
			return
		}

		if latest && ready {
			fmt.Println("You use actual version!")
			ExeFilePath := filepath.Join(exeDir, "ZaprUi.exe")
			resCh <- UpdateResult{exePath: ExeFilePath}

		} else {
			if err := update.DownloadReleaseExe(client, release, exeDir); err != nil {
				resCh <- UpdateResult{err: fmt.Errorf("❗Downloading failed because of: %v", err)}
				return
			}
			ExeFilePath := filepath.Join(exeDir, "ZaprUi.exe")
			resCh <- UpdateResult{exePath: ExeFilePath}

		}

	}()

	select {
	case res := <-resCh:
		if res.err != nil {
			panic(res.err)
		}

		fmt.Println("Update finished")
		a.ExeFilePath = res.exePath

		if !openzapr.IsAdmin() {
			openzapr.RunZaprAsAdmin(a.ExeFilePath)
		}

	case <-timeoutCtx.Done():
		fmt.Println("Timeout reached")
		runtime.Quit(a.ctx)
	}

}

// Getting sure that ProjectDir created
func ensureAppDir(path string) error {
	return os.MkdirAll(path, 0755)
}

type UpdateResult struct {
	exePath string
	err     error
}

// ===================================== WAILS API ==========================

// OpenURL opens the specified URL in the default browser
func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
