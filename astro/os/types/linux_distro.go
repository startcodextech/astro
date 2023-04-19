package types

type LinuxDistro struct {
	Name          string
	Version       string
	MajorVersion  string
	Unsupported   bool
	OnlySupported string
}
