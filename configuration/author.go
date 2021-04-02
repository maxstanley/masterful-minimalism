package configuration

type Author struct {
	FirstName string `yaml:"firstName"`
	LastName  string `yaml:"lastName"`
	Email     string `yaml:"email"`
	Twitter   string `yaml:"twitter"`
}

var Authors map[string]Author
