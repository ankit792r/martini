package release

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const releaseOutputsFileName = "release-outputs.md"

type ReleaseArtifact struct {
	Name string
	Type string
	SHA1 string
	Path string
}

func collectReleaseArtifacts(projectPath string) ([]ReleaseArtifact, error) {
	var artifacts []ReleaseArtifact

	apkDir := filepath.Join(projectPath, "build", "app", "outputs", "flutter-apk")
	apks, err := filepath.Glob(filepath.Join(apkDir, "*-release.apk"))
	if err != nil {
		return nil, err
	}
	for _, path := range apks {
		artifact, err := artifactFromPath(path, "release")
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}

	aabDir := filepath.Join(projectPath, "build", "app", "outputs", "bundle", "release")
	aabs, err := filepath.Glob(filepath.Join(aabDir, "*.aab"))
	if err != nil {
		return nil, err
	}
	for _, path := range aabs {
		artifact, err := artifactFromPath(path, "release")
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}

	sort.Slice(artifacts, func(i, j int) bool {
		if artifacts[i].Name == artifacts[j].Name {
			return artifacts[i].Path < artifacts[j].Path
		}
		return artifacts[i].Name < artifacts[j].Name
	})

	return artifacts, nil
}

func artifactFromPath(path, artifactType string) (ReleaseArtifact, error) {
	sha, err := fileSHA1(path)
	if err != nil {
		return ReleaseArtifact{}, fmt.Errorf("hash %s: %w", path, err)
	}
	return ReleaseArtifact{
		Name: filepath.Base(path),
		Type: artifactType,
		SHA1: sha,
		Path: path,
	}, nil
}

func fileSHA1(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func releaseOutputsMarkdown(artifacts []ReleaseArtifact) string {
	var b strings.Builder
	b.WriteString("### Release outputs\n\n")
	b.WriteString("| apk | type | sha |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, artifact := range artifacts {
		fmt.Fprintf(&b, "| %s | %s | `%s` |\n", artifact.Name, artifact.Type, artifact.SHA1)
	}
	b.WriteString("\n")
	return b.String()
}

func writeReleaseOutputsMarkdown(path string, artifacts []ReleaseArtifact) error {
	return os.WriteFile(path, []byte(releaseOutputsMarkdown(artifacts)), 0o644)
}

func createReleaseZip(projectPath string, artifacts []ReleaseArtifact) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}

	zipPath := filepath.Join(home, releaseZipName(projectPath))
	if err := os.MkdirAll(home, 0o755); err != nil {
		return "", err
	}

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("create zip: %w", err)
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	if err := addStringToZip(writer, releaseOutputsFileName, releaseOutputsMarkdown(artifacts)); err != nil {
		return "", err
	}

	for _, artifact := range artifacts {
		if err := addFileToZip(writer, artifact.Path, artifact.Name); err != nil {
			return "", err
		}
	}

	if err := writer.Close(); err != nil {
		return "", err
	}
	if err := zipFile.Close(); err != nil {
		return "", err
	}

	return zipPath, nil
}

func releaseZipName(projectPath string) string {
	base := filepath.Base(filepath.Clean(projectPath))
	if base == "" || base == "." || base == string(os.PathSeparator) {
		base = "flutter-app"
	}
	return fmt.Sprintf("%s-release-%s.zip", base, time.Now().Format("20060102-150405"))
}

func addStringToZip(writer *zip.Writer, name, content string) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
	}
	header.SetModTime(time.Now())

	entry, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.WriteString(entry, content)
	return err
}

func addFileToZip(writer *zip.Writer, sourcePath, entryName string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = entryName
	header.Method = zip.Deflate

	entry, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, file)
	return err
}
