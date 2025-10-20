package config

var (
	Username string
	Password string
)

// GetCredentials returns the username and password from the config
func GetCredentials() (string, string) {
	return Username, Password
}
