package build

var (
	commitSHA string
	date string
)

// Info holds build information
type Info struct {
	GitCommit  string `json:"gitCommit"`
	BuildDate  string `json:"BuildDate"`
}

// Get creates and initialized Info object
func Get() Info {
	return Info{
		GitCommit:  commitSHA,
		BuildDate:  date,
	}
}
