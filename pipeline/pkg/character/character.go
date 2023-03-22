package character

type Config struct {
	PackageName string   `yaml:"package_name,omitempty"`
	GenshinID   string   `yaml:"genshin_id,omitempty"`
	Key         string   `yaml:"key,omitempty"`
	Shortcuts   []string `yaml:"shortcuts,omitempty"`

	//extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Pipeline struct {
}

func NewPipeline(c []Config) (*Pipeline, error) {
	p := &Pipeline{}

	return p, nil
}
