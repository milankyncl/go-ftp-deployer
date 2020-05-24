package deployer

type Config struct {
	Host      string
	User      string
	Password  string
	LocalPath string
	ExtPath   string
	Exclude   []string
}
