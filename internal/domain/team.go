package domain

type Team struct {
	TeamName string
}

func NewTeam(teamName string) *Team {
	return &Team{teamName}
}
