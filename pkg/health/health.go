package healthcheck

type Health struct {
	Status      Status
	Version     string
	ReleaseId   string
	Notes       []string
	Output      string
	Checks      Checks
	Links       []string
	ServiceId   string
	Description string
}
