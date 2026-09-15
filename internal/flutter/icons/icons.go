package icons

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var iosIconMap = map[string]string{
	"AppIcon@2x.png":            "Icon-App-60x60@2x.png",
	"AppIcon@3x.png":            "Icon-App-60x60@3x.png",
	"AppIcon~ipad.png":          "Icon-App-76x76@1x.png",
	"AppIcon@2x~ipad.png":       "Icon-App-76x76@2x.png",
	"AppIcon-83.5@2x~ipad.png":  "Icon-App-83.5x83.5@2x.png",
	"AppIcon-40@2x.png":         "Icon-App-40x40@2x.png",
	"AppIcon-40@3x.png":         "Icon-App-40x40@3x.png",
	"AppIcon-40~ipad.png":       "Icon-App-40x40@1x.png",
	"AppIcon-40@2x~ipad.png":    "Icon-App-40x40@2x.png",
	"AppIcon-20@2x.png":         "Icon-App-20x20@2x.png",
	"AppIcon-20@3x.png":         "Icon-App-20x20@3x.png",
	"AppIcon-20~ipad.png":       "Icon-App-20x20@1x.png",
	"AppIcon-20@2x~ipad.png":    "Icon-App-20x20@2x.png",
	"AppIcon-29.png":            "Icon-App-29x29@1x.png",
	"AppIcon-29@2x.png":         "Icon-App-29x29@2x.png",
	"AppIcon-29@3x.png":         "Icon-App-29x29@3x.png",
	"AppIcon-29~ipad.png":       "Icon-App-29x29@1x.png",
	"AppIcon-29@2x~ipad.png":    "Icon-App-29x29@2x.png",
	"AppIcon~ios-marketing.png": "Icon-App-1024x1024@1x.png",
}

var webIconMap = map[string]string{
	"icon-192.png":          "icons/Icon-192.png",
	"icon-512.png":          "icons/Icon-512.png",
	"icon-192-maskable.png": "icons/Icon-maskable-192.png",
	"icon-512-maskable.png": "icons/Icon-maskable-512.png",
}

type Result struct {
	AndroidFiles int
	IOSFiles     int
	WebFiles     int
}

func ResolveIconsPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("icons path is required")
	}
	path = expandHome(path)

	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("icons path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("icons path is not a directory: %s", path)
	}
	return path, nil
}

func Update(projectPath, iconsPath string) (*Result, error) {
	if err := validateFlutterProject(projectPath); err != nil {
		return nil, err
	}

	resolved, err := ResolveIconsPath(iconsPath)
	if err != nil {
		return nil, err
	}

	result := &Result{}

	androidCount, err := updateAndroid(projectPath, resolved)
	if err != nil {
		return nil, err
	}
	result.AndroidFiles = androidCount

	iosCount, err := updateIOS(projectPath, resolved)
	if err != nil {
		return nil, err
	}
	result.IOSFiles = iosCount

	webCount, err := updateWeb(projectPath, resolved)
	if err != nil {
		return nil, err
	}
	result.WebFiles = webCount

	if result.AndroidFiles == 0 && result.IOSFiles == 0 && result.WebFiles == 0 {
		return nil, fmt.Errorf("no icon files copied; expected IconKitchen output under %s", resolved)
	}

	return result, nil
}

func validateFlutterProject(projectPath string) error {
	info, err := os.Stat(projectPath)
	if err != nil {
		return fmt.Errorf("project path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", projectPath)
	}
	pubspec := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspec); err != nil {
		return fmt.Errorf("pubspec.yaml not found in project root")
	}
	return nil
}

func updateAndroid(projectPath, iconsPath string) (int, error) {
	srcRoot := filepath.Join(iconsPath, "android", "res")
	if _, err := os.Stat(srcRoot); err != nil {
		return 0, fmt.Errorf("android icons not found at %s", srcRoot)
	}

	destRoot := filepath.Join(projectPath, "android", "app", "src", "main", "res")
	if _, err := os.Stat(destRoot); err != nil {
		return 0, fmt.Errorf("android res directory not found at %s", destRoot)
	}

	return copyTree(srcRoot, destRoot)
}

func updateIOS(projectPath, iconsPath string) (int, error) {
	srcRoot := filepath.Join(iconsPath, "ios")
	if _, err := os.Stat(srcRoot); err != nil {
		return 0, fmt.Errorf("ios icons not found at %s", srcRoot)
	}

	destRoot := filepath.Join(projectPath, "ios", "Runner", "Assets.xcassets", "AppIcon.appiconset")
	if _, err := os.Stat(destRoot); err != nil {
		return 0, fmt.Errorf("ios AppIcon.appiconset not found at %s", destRoot)
	}

	count := 0
	for srcName, destName := range iosIconMap {
		src := filepath.Join(srcRoot, srcName)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dest := filepath.Join(destRoot, destName)
		if err := copyFile(src, dest); err != nil {
			return count, fmt.Errorf("copy ios icon %s: %w", srcName, err)
		}
		count++
	}
	return count, nil
}

func updateWeb(projectPath, iconsPath string) (int, error) {
	srcRoot := filepath.Join(iconsPath, "web")
	if _, err := os.Stat(srcRoot); err != nil {
		return 0, fmt.Errorf("web icons not found at %s", srcRoot)
	}

	destRoot := filepath.Join(projectPath, "web")
	if _, err := os.Stat(destRoot); err != nil {
		return 0, fmt.Errorf("web directory not found at %s", destRoot)
	}

	count := 0
	for srcName, destRel := range webIconMap {
		src := filepath.Join(srcRoot, srcName)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dest := filepath.Join(destRoot, destRel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return count, err
		}
		if err := copyFile(src, dest); err != nil {
			return count, fmt.Errorf("copy web icon %s: %w", srcName, err)
		}
		count++
	}

	faviconSrc := filepath.Join(srcRoot, "icon-192.png")
	if _, err := os.Stat(faviconSrc); err == nil {
		if err := copyFile(faviconSrc, filepath.Join(destRoot, "favicon.png")); err != nil {
			return count, err
		}
		count++
	}

	icoSrc := filepath.Join(srcRoot, "favicon.ico")
	if _, err := os.Stat(icoSrc); err == nil {
		if err := copyFile(icoSrc, filepath.Join(destRoot, "favicon.ico")); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func copyTree(srcRoot, destRoot string) (int, error) {
	count := 0
	err := filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(destRoot, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := copyFile(path, dest); err != nil {
			return err
		}
		count++
		return nil
	})
	return count, err
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func expandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		if path == "~" {
			return home
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}
