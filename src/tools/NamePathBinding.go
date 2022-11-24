package tools

type NamePathBinding struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func (b NamePathBinding) String() string {
	return "name: " + b.Name + "\npath: " + b.Path
}

func AreNamePathBindingEqual(b1 NamePathBinding, b2 NamePathBinding) (bool, bool) {
	areNamesEquals := false
	arePathEquals := false

	if b1.Name == b2.Name {
		areNamesEquals = true
	}
	if b1.Path == b2.Path {
		arePathEquals = true
	}
	return areNamesEquals, arePathEquals
}
