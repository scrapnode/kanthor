package configuration

func New() (Provider, error) {
	return NewFile(FileLookingDirs)
}

type Provider interface {
	Unmarshal(dest interface{}) error
	Sources() []Source
	SetDefault(key string, value interface{})
}

type Source struct {
	Looking string
	Found   string
	Used    bool
}
