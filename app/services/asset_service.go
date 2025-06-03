package services

import (
	"encoding/json"
	"fmt"
	"os"
)

type AssetService struct {
	manifest map[string]string
	isDev    bool
}

func NewAssetService() *AssetService {
	as := &AssetService{
		manifest: make(map[string]string),
		isDev:    os.Getenv("ENV") != "production",
	}

	as.loadManifest()
	return as
}

func (as *AssetService) loadManifest() {
	manifestPath := "web/static/dist/manifest.json"

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		fmt.Printf("Warning: Asset manifest not found at %s\n", manifestPath)
		return
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Printf("Error reading asset manifest: %v\n", err)
		return
	}

	if err := json.Unmarshal(data, &as.manifest); err != nil {
		fmt.Printf("Error parsing asset manifest: %v\n", err)
		return
	}

	fmt.Printf("Loaded asset manifest with %d entries\n", len(as.manifest))
}

// GetAssetPath returns the fingerprinted path for an asset
func (as *AssetService) GetAssetPath(assetName string) string {
	// In development, always reload manifest
	if as.isDev {
		as.loadManifest()
	}

	if fingerprintedPath, exists := as.manifest[assetName]; exists {
		return "/static/dist/" + fingerprintedPath
	}

	// Fallback to original path if not in manifest
	return "/static/dist/" + assetName
}

// GetCSSPath returns the path to the main CSS file
func (as *AssetService) GetCSSPath() string {
	return as.GetAssetPath("styles.css")
}

// GetJSPath returns the path to the main JS file
func (as *AssetService) GetJSPath() string {
	return as.GetAssetPath("app.js")
}
