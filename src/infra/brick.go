package infra

type Brick struct {
	// The brick's name. Usually the name of the parent directory
	Name string
	// The absolute path of the brick's directory
	Path string
	// The absolute path of the `brick.yml` file
	ConfigurationFilePath string
	// Wheither or not the brick contains a `brick.yml` file.
	// Meaning it does not contain any other brick.
	IsElementary bool
func (brick *Brick) SetElementary(cfp string) *Brick {
	brick.IsElementary = true
	brick.ConfigurationFilePath = cfp

	return brick
}

}
