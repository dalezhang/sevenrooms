package env

type Redis struct {
	DNS          string
	Database     string
	ReadTimeout  int //ms
	WriteTimeout int //ms
}
