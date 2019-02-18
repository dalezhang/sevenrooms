package env

type DefaultSecrets struct {
	DefaultClientID              string
	DefaultClientSecret          string
	DefaultDashboardClientID     string
	DefaultDashboardClientSecret string
}

type Secret struct {
	CcPassPhrase    string
	SsnPassPhrase   string
	CcKeyPublicKey  string
	CcKeyPrivateKey string
}
