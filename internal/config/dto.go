package config

const AppName = "universal-proxy"

// Build information -ldflags .
var (
	branch     = "dev" //nolint
	commitHash = "-"   //nolint
	timeBuild  = "-"   //nolint
)

type Version struct {
	Name       string `json:"name,omitempty"`
	Branch     string `json:"branch,omitempty"`
	CommitHash string `json:"commitHash,omitempty"`
	TimeBuild  string `json:"timeBuild,omitempty"`
}

func GetVersion() Version {
	return Version{
		Name:       AppName,
		Branch:     branch,
		CommitHash: commitHash,
		TimeBuild:  timeBuild,
	}
}
