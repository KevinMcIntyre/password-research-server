package jobs

type TrialUserImage struct {
	TrialID string
	Stage   string
	Row     string
	Column  string
	Image   []byte
}

type SubjectImage struct {
	ID          string
	SubjectName string
	Image       []byte
}

type ConfigImage struct {
	ConfigName string
	Stage      string
	Row        string
	Column     string
	Image      []byte
}
