package gradle

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const DefaultUploadAlias = "upload"

type Properties struct {
	StorePassword string
	KeyPassword   string
	KeyAlias      string
	StoreFile     string
}

const (
	keystoreImportsKotlin = "import java.util.Properties\nimport java.io.FileInputStream\n"
	keystoreLoaderKotlin  = `val keystoreProperties = Properties()
val keystorePropertiesFile = rootProject.file("key.properties")
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(FileInputStream(keystorePropertiesFile))
}

`
	signingConfigsKotlin = `    signingConfigs {
        create("release") {
            keyAlias = keystoreProperties.getProperty("keyAlias")
            keyPassword = keystoreProperties.getProperty("keyPassword")
            storeFile = keystoreProperties.getProperty("storeFile")?.let { file(it) }
            storePassword = keystoreProperties.getProperty("storePassword")
        }
    }

`
)

func Apply(projectPath string, props Properties) error {
	keyPropertiesPath := filepath.Join(projectPath, "android", "key.properties")
	if err := WriteKeyProperties(keyPropertiesPath, props); err != nil {
		return err
	}
	return UpdateProjectSigning(projectPath)
}

func WriteKeyProperties(path string, props Properties) error {
	content := fmt.Sprintf(
		"storePassword=%s\nkeyPassword=%s\nkeyAlias=%s\nstoreFile=%s\n",
		props.StorePassword,
		props.KeyPassword,
		props.KeyAlias,
		storeFilePath(props.StoreFile),
	)
	return os.WriteFile(path, []byte(content), 0o600)
}

func UpdateProjectSigning(projectPath string) error {
	androidDir := filepath.Join(projectPath, "android")
	ktsPath := filepath.Join(androidDir, "app", "build.gradle.kts")
	if _, err := os.Stat(ktsPath); err == nil {
		return updateBuildGradleKTS(ktsPath)
	}

	gradlePath := filepath.Join(androidDir, "app", "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		return updateBuildGradleGroovy(gradlePath)
	}

	return fmt.Errorf("build.gradle.kts or build.gradle not found under android/app")
}

func storeFilePath(path string) string {
	clean := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(clean, `\`, `\\`)
	}
	return clean
}

func updateBuildGradleKTS(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updated := string(content)
	if !strings.Contains(updated, "keystoreProperties") {
		updated = injectBeforeAndroidBlock(updated, keystoreImportsKotlin+keystoreLoaderKotlin)
	}

	if !strings.Contains(updated, `signingConfigs {`) {
		updated = injectInsideAndroidBlock(updated, signingConfigsKotlin)
	}

	updated = strings.ReplaceAll(updated, `signingConfig = signingConfigs.getByName("debug")`, `signingConfig = signingConfigs.getByName("release")`)
	if !strings.Contains(updated, `signingConfig = signingConfigs.getByName("release")`) {
		updated = injectReleaseSigningConfig(updated, `            signingConfig = signingConfigs.getByName("release")`)
	}

	return os.WriteFile(path, []byte(updated), 0o644)
}

func updateBuildGradleGroovy(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updated := string(content)
	if !strings.Contains(updated, "keystoreProperties") {
		imports := "import java.util.Properties\nimport java.io.FileInputStream\n"
		loader := "def keystoreProperties = new Properties()\n" +
			"def keystorePropertiesFile = rootProject.file('key.properties')\n" +
			"if (keystorePropertiesFile.exists()) {\n" +
			"    keystoreProperties.load(new FileInputStream(keystorePropertiesFile))\n" +
			"}\n\n"
		updated = injectBeforeAndroidBlock(updated, imports+loader)
	}

	if !strings.Contains(updated, "signingConfigs {") {
		block := "    signingConfigs {\n" +
			"        release {\n" +
			"            keyAlias = keystoreProperties['keyAlias']\n" +
			"            keyPassword = keystoreProperties['keyPassword']\n" +
			"            storeFile = keystoreProperties['storeFile'] ? file(keystoreProperties['storeFile']) : null\n" +
			"            storePassword = keystoreProperties['storePassword']\n" +
			"        }\n" +
			"    }\n\n"
		updated = injectInsideAndroidBlock(updated, block)
	}

	updated = strings.ReplaceAll(updated, "signingConfig = signingConfigs.debug", "signingConfig = signingConfigs.release")
	if !strings.Contains(updated, "signingConfig = signingConfigs.release") {
		updated = injectReleaseSigningConfig(updated, "            signingConfig = signingConfigs.release")
	}

	return os.WriteFile(path, []byte(updated), 0o644)
}

func injectBeforeAndroidBlock(content, snippet string) string {
	idx := strings.Index(content, "android {")
	if idx == -1 {
		return content
	}
	return content[:idx] + snippet + content[idx:]
}

func injectInsideAndroidBlock(content, snippet string) string {
	idx := strings.Index(content, "android {")
	if idx == -1 {
		return content
	}

	buildTypesIdx := strings.Index(content[idx:], "buildTypes {")
	if buildTypesIdx == -1 {
		return content
	}

	insertAt := idx + buildTypesIdx
	return content[:insertAt] + snippet + content[insertAt:]
}

func injectReleaseSigningConfig(content, line string) string {
	releaseIdx := strings.Index(content, "release {")
	if releaseIdx == -1 {
		return content
	}

	closeIdx := strings.Index(content[releaseIdx:], "\n        }")
	if closeIdx == -1 {
		closeIdx = strings.Index(content[releaseIdx:], "\n    }")
	}
	if closeIdx == -1 {
		return content
	}

	insertAt := releaseIdx + closeIdx
	if strings.Contains(content[releaseIdx:insertAt], "signingConfig") {
		return content
	}

	return content[:insertAt] + "\n" + line + content[insertAt:]
}
