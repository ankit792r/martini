package keystore

import shared "martini/internal/keystore"

type Options struct {
	KeystoreName  string
	StorePassword string
	KeyPassword   string
	KeyAlias      string
	CommonName    string
}

func (o Options) Config() shared.Config {
	answers := []string{
		o.KeystoreName,
		o.StorePassword,
		o.KeyPassword,
		o.KeyAlias,
		o.CommonName,
	}
	return shared.ConfigFromAnswers(answers).Normalized()
}
