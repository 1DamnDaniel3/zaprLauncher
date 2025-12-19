package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ReleaseResp struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

const (
	release = "https://api.github.com/repos/1DamnDaniel3/ZaprUI/releases/latest"
)

// Downloading exe of zaprUi
func DownloadReleaseExe(client *http.Client, release *ReleaseResp, dir string) error {
	url := findExeReleaseURL(release)
	if url == "" {
		return fmt.Errorf("empty URL after request")
	}
	// request to download
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("download request failed with error: %v", err)
	}

	// Exe file from HTTP
	respExe, err := client.Do(req)
	if err != nil {
		return err
	}
	if respExe.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", respExe.Status)
	}
	defer respExe.Body.Close()

	path := filepath.Join(dir, "ZaprUi.exe")

	// Creating new empty file
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, respExe.Body)
	return nil
}

// Finding if version file exist and create it if it don't
func EnsureVersionFileExist(dir string, release *ReleaseResp) error {
	path := filepath.Join(dir, "zaprUI_version.txt")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(path, []byte(release.TagName), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// If we have outdated version
func CorrectVersionFile(dir string, release *ReleaseResp) error {
	path := filepath.Join(dir, "zaprUI_version.txt")
	err := os.WriteFile(path, []byte(release.TagName), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Trying if we have actual version
func IsLatestVersion(pathToFile string, release *ReleaseResp) (bool, error) {
	version, err := os.ReadFile(pathToFile)
	if err != nil {
		return false, fmt.Errorf("cannot read version file: %w", err)
	}
	versionStr := strings.TrimSpace(string(version))
	return versionStr == release.TagName, nil
}
func IsReleaseReady(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("cannot read release directory: %w", err)
	}

	return len(entries) > 0, nil
}

// Make request to releases and asking for URLs of needed files
func ParceLatestRelease(client *http.Client) (*ReleaseResp, error) {
	req, err := http.NewRequest("GET", release, nil)
	if err != nil {
		return nil, fmt.Errorf("request to `zapret-discord-youtube/releases` fail with error: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "go-client")

	// catching response
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// trying bad response
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad response %d: %v", resp.StatusCode, string(body))
	}

	// decode response in struct
	var release ReleaseResp
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("cant parse resource. error: %v", err)
	}

	return &release, nil
}

// ==================== local functions ======================

// Taking .exe file from release and saving version
func findExeReleaseURL(release *ReleaseResp) string {
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".exe") {
			return asset.URL
		}
	}
	return ""
}
