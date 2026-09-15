package flutter

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

func writeReleaseOutputsMarkdown(path string, artifacts []ReleaseArtifact) error {
	var b strings.Builder
	b.WriteString("### Release outputs\n\n")
	b.WriteString("| apk | type | sha |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, artifact := range artifacts {
		fmt.Fprintf(&b, "| %s | %s | `%s` |\n", artifact.Name, artifact.Type, artifact.SHA1)
	}
	b.WriteString("\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
