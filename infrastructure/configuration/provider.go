package configuration

func New() (Provider, error) {
	return NewFile(FileLookingDirs)
}

type Provider interface {
	Unmarshal(dest interface{}) error
	Sources() []Source
}

type Source struct {
	Looking string
	Found   string
	Used    bool
}
