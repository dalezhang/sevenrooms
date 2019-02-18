package env

var (
	KeyBindo   = "bindo"
	KeyGateway = "gateway"
)

type Faye struct {
	URL  string
	Keys map[string]string
}
