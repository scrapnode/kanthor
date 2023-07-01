package configuration

func New() (Provider, error) {
	return NewFile(FileLookingDirs)
}

type Provider interface {
	Unmarshal(dest interface{}) error
	Sources() []Source
}

type Source struct {
	Source string
	Found  bool
}