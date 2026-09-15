package flutter

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	keystore "github.com/pavlo-v-chernykh/keystore-go/v4"
)

const (
	defaultKeyAlias    = "upload"
	defaultValidityDay = 10000
	defaultKeySize     = 2048
)

type KeystoreConfig struct {
	Name          string
	StorePassword string
	KeyPassword   string
	Alias         string
	CommonName    string
}

type KeystoreResult struct {
	HomePath      string
	ProjectPath   string
	Credentials   string
	CredentialsPath string
	Alias         string
	StorePassword string
	KeyPassword   string
}

func generateUploadKeystore(projectPath string, cfg KeystoreConfig) (*KeystoreResult, error) {
	if strings.TrimSpace(cfg.KeyPassword) == "" {
		cfg.KeyPassword = cfg.StorePassword
	}
	if strings.TrimSpace(cfg.Alias) == "" {
		cfg.Alias = defaultKeyAlias
	}
	if strings.TrimSpace(cfg.CommonName) == "" {
		cfg.CommonName = "Android Upload"
	}
	if err := validateKeystoreConfig(cfg); err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home directory: %w", err)
	}

	androidDir := filepath.Join(projectPath, "android")
	if _, err := os.Stat(androidDir); err != nil {
		return nil, fmt.Errorf("android directory not found at %s", androidDir)
	}

	fileName := cfg.Name + ".jks"
	homeKeystorePath := filepath.Join(homeDir, fileName)
	projectKeystorePath := filepath.Join(androidDir, fileName)

	if _, err := os.Stat(homeKeystorePath); err == nil {
		return nil, fmt.Errorf("keystore already exists: %s", homeKeystorePath)
	}
	if _, err := os.Stat(projectKeystorePath); err == nil {
		return nil, fmt.Errorf("keystore already exists: %s", projectKeystorePath)
	}

	keystoreBytes, cert, err := createJKS(cfg)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(homeKeystorePath, keystoreBytes, 0o600); err != nil {
		return nil, fmt.Errorf("write home keystore: %w", err)
	}
	if err := copyFile(homeKeystorePath, projectKeystorePath); err != nil {
		return nil, fmt.Errorf("write project keystore: %w", err)
	}

	credentialsPath := filepath.Join(homeDir, cfg.Name+"-upload-key.txt")
	credentials := formatCredentials(cfg, homeKeystorePath, projectKeystorePath, cert)
	if err := os.WriteFile(credentialsPath, []byte(credentials), 0o600); err != nil {
		return nil, fmt.Errorf("write credentials file: %w", err)
	}

	return &KeystoreResult{
		HomePath:        homeKeystorePath,
		ProjectPath:     projectKeystorePath,
		Credentials:     credentials,
		CredentialsPath: credentialsPath,
		Alias:           cfg.Alias,
		StorePassword:   cfg.StorePassword,
		KeyPassword:     cfg.KeyPassword,
	}, nil
}

func validateKeystoreConfig(cfg KeystoreConfig) error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("keystore name is required")
	}
	if strings.TrimSpace(cfg.StorePassword) == "" {
		return fmt.Errorf("store password is required")
	}
	if strings.TrimSpace(cfg.KeyPassword) == "" {
		return fmt.Errorf("key password is required")
	}
	if strings.TrimSpace(cfg.Alias) == "" {
		return fmt.Errorf("key alias is required")
	}
	if strings.TrimSpace(cfg.CommonName) == "" {
		return fmt.Errorf("certificate common name is required")
	}
	return nil
}

func createJKS(cfg KeystoreConfig) ([]byte, *x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultKeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("generate rsa key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(defaultValidityDay * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: cfg.CommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("create certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, fmt.Errorf("parse certificate: %w", err)
	}

	pkcs8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal private key: %w", err)
	}

	ks := keystore.New()
	entry := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:     pkcs8,
		CertificateChain: []keystore.Certificate{
			{Type: "X509", Content: certDER},
		},
	}
	if err := ks.SetPrivateKeyEntry(cfg.Alias, entry, []byte(cfg.KeyPassword)); err != nil {
		return nil, nil, fmt.Errorf("set private key entry: %w", err)
	}

	var buf bytes.Buffer
	if err := ks.Store(&buf, []byte(cfg.StorePassword)); err != nil {
		return nil, nil, fmt.Errorf("encode jks: %w", err)
	}

	return buf.Bytes(), cert, nil
}

func formatCredentials(cfg KeystoreConfig, homePath, projectPath string, cert *x509.Certificate) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Upload keystore credentials\n")
	fmt.Fprintf(&b, "===========================\n\n")
	fmt.Fprintf(&b, "Keystore (home): %s\n", homePath)
	fmt.Fprintf(&b, "Keystore (project): %s\n", projectPath)
	fmt.Fprintf(&b, "Store password: %s\n", cfg.StorePassword)
	fmt.Fprintf(&b, "Key password: %s\n", cfg.KeyPassword)
	fmt.Fprintf(&b, "Key alias: %s\n", cfg.Alias)
	fmt.Fprintf(&b, "Key algorithm: RSA\n")
	fmt.Fprintf(&b, "Key size: %d\n", defaultKeySize)
	fmt.Fprintf(&b, "Store type: JKS\n")
	fmt.Fprintf(&b, "Validity days: %d\n", defaultValidityDay)
	fmt.Fprintf(&b, "Common name: %s\n", cfg.CommonName)
	if cert != nil {
		fmt.Fprintf(&b, "Valid from: %s\n", cert.NotBefore.Format(time.RFC3339))
		fmt.Fprintf(&b, "Valid until: %s\n", cert.NotAfter.Format(time.RFC3339))
	}
	fmt.Fprintf(&b, "\nKeep this file private. Do not commit keystore files or passwords to source control.\n")
	return b.String()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func keystoreConfigFromAnswers(answers []string) KeystoreConfig {
	cfg := KeystoreConfig{
		Name:       "upload-keystore",
		Alias:      defaultKeyAlias,
		CommonName: "Android Upload",
	}

	if len(answers) > 0 && strings.TrimSpace(answers[0]) != "" {
		cfg.Name = strings.TrimSpace(strings.TrimSuffix(answers[0], ".jks"))
	}
	if len(answers) > 1 {
		cfg.StorePassword = strings.TrimSpace(answers[1])
	}
	if len(answers) > 2 && strings.TrimSpace(answers[2]) != "" {
		cfg.KeyPassword = strings.TrimSpace(answers[2])
	} else {
		cfg.KeyPassword = cfg.StorePassword
	}
	if len(answers) > 3 && strings.TrimSpace(answers[3]) != "" {
		cfg.Alias = strings.TrimSpace(answers[3])
	}
	if len(answers) > 4 && strings.TrimSpace(answers[4]) != "" {
		cfg.CommonName = strings.TrimSpace(answers[4])
	}

	return cfg
}
