package dpfm_api_processing_formatter

type HeaderUpdates struct {
	Invitation				int		`json:"Site"`
}

type ValidationCheckUpdates struct {
	InvitationObjectType	string	`json:"InvitationObjectType"`
	InvitationObject		int		`json:"InvitationObject"`
	InvitationOwner			int		`json:"InvitationOwner"`
	InvitationGuest			int		`json:"InvitationGuest"`
}
