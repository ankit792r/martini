package pkgrename

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Result struct {
	FilesUpdated int
	DirsRenamed  int
	UpdatedFiles []string
	RenamedDirs  []string
}

func Update(projectPath, oldName, newName string) (*Result, error) {
	info, err := os.Stat(projectPath)
	if err != nil {
		return nil, fmt.Errorf("project path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project path is not a directory: %s", projectPath)
	}

	if oldName == newName {
		return nil, fmt.Errorf("old and new package names are the same")
	}

	oldSlash := packageToPath(oldName)
	newSlash := packageToPath(newName)

	result := &Result{}
	if err := replaceInFiles(projectPath, oldName, newName, oldSlash, newSlash, result); err != nil {
		return nil, err
	}
	if err := renamePackageDirs(projectPath, oldSlash, newSlash, result); err != nil {
		return nil, err
	}

	return result, nil
}

func packageToPath(name string) string {
	return strings.ReplaceAll(name, ".", "/")
}

func removeEmptyParents(dir string) {
	for {
		if dir == "" || dir == "." || dir == string(os.PathSeparator) {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		parent := filepath.Dir(dir)
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = parent
	}
}

func replacePathSuffix(path, oldSuffix, newSuffix string) string {
	slashPath := filepath.ToSlash(path)
	if !strings.HasSuffix(slashPath, oldSuffix) {
		return path
	}
	prefix := strings.TrimSuffix(slashPath, oldSuffix)
	return filepath.FromSlash(prefix + newSuffix)
}

func replaceInFiles(root, oldName, newName, oldSlash, newSlash string, result *Result) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipDir(path, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipFile(d.Name()) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !isTextFile(data) {
			return nil
		}

		content := string(data)
		updated := applyReplacements(content, oldName, newName, oldSlash, newSlash)
		if updated == content {
			return nil
		}

		if err := os.WriteFile(path, []byte(updated), fileMode(d, path)); err != nil {
			return err
		}

		result.FilesUpdated++
		result.UpdatedFiles = append(result.UpdatedFiles, path)
		return nil
	})
}

func applyReplacements(content, oldName, newName, oldSlash, newSlash string) string {
	content = strings.ReplaceAll(content, oldName, newName)
	content = strings.ReplaceAll(content, oldSlash, newSlash)
	return content
}

func renamePackageDirs(root, oldSlash, newSlash string, result *Result) error {
	if oldSlash == newSlash {
		return nil
	}

	renamed := make(map[string]bool)
	if err := renameAndroidSourcePackageDirs(root, oldSlash, newSlash, result, renamed); err != nil {
		return err
	}

	return renameMatchingPackageDirs(root, oldSlash, newSlash, result, renamed)
}

func renameAndroidSourcePackageDirs(root, oldSlash, newSlash string, result *Result, renamed map[string]bool) error {
	androidDir := filepath.Join(root, "android")
	if _, err := os.Stat(androidDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var sourceRoots []string
	err := filepath.WalkDir(androidDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if shouldSkipDir(path, d.Name()) {
			return filepath.SkipDir
		}
		if isAndroidSourceRoot(path, d.Name()) {
			sourceRoots = append(sourceRoots, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, sourceRoot := range sourceRoots {
		oldDir := filepath.Join(sourceRoot, filepath.FromSlash(oldSlash))
		newDir := filepath.Join(sourceRoot, filepath.FromSlash(newSlash))
		if err := movePackageDir(oldDir, newDir, result, renamed); err != nil {
			return err
		}
	}

	return nil
}

func isAndroidSourceRoot(path, name string) bool {
	if name != "kotlin" && name != "java" {
		return false
	}
	return strings.Contains(filepath.ToSlash(path), "/src/")
}

func renameMatchingPackageDirs(root, oldSlash, newSlash string, result *Result, renamed map[string]bool) error {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if shouldSkipDir(path, d.Name()) {
				return filepath.SkipDir
			}
			if strings.HasSuffix(filepath.ToSlash(path), oldSlash) {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) > len(dirs[j])
	})

	for _, dir := range dirs {
		if renamed[dir] {
			continue
		}
		newDir := replacePathSuffix(dir, oldSlash, newSlash)
		if err := movePackageDir(dir, newDir, result, renamed); err != nil {
			return err
		}
	}

	return nil
}

func movePackageDir(oldDir, newDir string, result *Result, renamed map[string]bool) error {
	if renamed[oldDir] {
		return nil
	}

	if _, err := os.Stat(oldDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if oldDir == newDir {
		return nil
	}
	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("cannot rename %s: destination already exists", oldDir)
	}

	if err := os.MkdirAll(filepath.Dir(newDir), 0o755); err != nil {
		return fmt.Errorf("create parent dirs for %s: %w", newDir, err)
	}

	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("rename %s: %w", oldDir, err)
	}

	removeEmptyParents(filepath.Dir(oldDir))

	renamed[oldDir] = true
	result.DirsRenamed++
	result.RenamedDirs = append(result.RenamedDirs, oldDir+" -> "+newDir)
	return nil
}

func isTextFile(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	if len(data) > 1024*1024 {
		return false
	}
	return !strings.Contains(string(data), "\x00")
}

func fileMode(d fs.DirEntry, path string) fs.FileMode {
	info, err := d.Info()
	if err == nil {
		return info.Mode()
	}
	info, err = os.Stat(path)
	if err == nil {
		return info.Mode()
	}
	return 0o644
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}

func shouldSkipDir(path, name string) bool {
	if isHidden(name) {
		return true
	}

	switch name {
	case "build", "node_modules", "Pods", ".symlinks", "DerivedData", "coverage", "dist":
		return true
	}

	slashPath := filepath.ToSlash(path)
	for _, segment := range []string{"/build/", "/node_modules/", "/Pods/"} {
		if strings.Contains(slashPath, segment) {
			return true
		}
	}

	return false
}

func shouldSkipFile(name string) bool {
	if isHidden(name) {
		return true
	}

	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf",
		".zip", ".jar", ".aar", ".so", ".dylib", ".dll", ".exe",
		".keystore", ".jks", ".bin", ".dat", ".db", ".sqlite",
		".ttf", ".otf", ".woff", ".woff2", ".mp3", ".mp4", ".mov":
		return true
	}
	return false
}
