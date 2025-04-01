package providers

type AbstractProvider struct {
	Command          string
	InstallCommand   string
	UpdateCommand    string
	CleanCommand     string
	VersionSeparator string
	RequiresRoot     bool
}
