package sub_func_complementer

type SDC struct {
	ConnectionKey       string   `json:"connection_key"`
	RedisKey            string   `json:"redis_key"`
	Filepath            string   `json:"filepath"`
	APIStatusCode       int      `json:"api_status_code"`
	RuntimeSessionID    string   `json:"runtime_session_id"`
	BusinessPartnerID   *int     `json:"business_partner"`
	ServiceLabel        string   `json:"service_label"`
	APIType             string   `json:"api_type"`
	Message             Message  `json:"message"`
	APISchema           string   `json:"api_schema"`
	Accepter            []string `json:"accepter"`
	Deleted             bool     `json:"deleted"`
	SQLUpdateResult     *bool    `json:"sql_update_result"`
	SQLUpdateError      string   `json:"sql_update_error"`
	SubfuncResult       *bool    `json:"subfunc_result"`
	SubfuncError        string   `json:"subfunc_error"`
	ExconfResult        *bool    `json:"exconf_result"`
	ExconfError         string   `json:"exconf_error"`
	APIProcessingResult *bool    `json:"api_processing_result"`
	APIProcessingError  string   `json:"api_processing_error"`
}

type Message struct {
	Header             *Header				`json:"Header"`
	ValidationCheck    *[]ValidationCheck	`json:"ValidationCheck"`
}

type Header struct {
	Invitation				int		`json:"Invitation"`
	InvitationType			string	`json:"InvitationType"`
	InvitationOwner			int		`json:"InvitationOwner"`
	InvitationGuest			int		`json:"InvitationGuest"`
	InvitationObjectType	string	`json:"InvitationObjectType"`
	InvitationObject		int		`json:"InvitationObject"`
	OwnerParticipation		int		`json:"OwnerParticipation"`
	OwnerAttendance			*int	`json:"OwnerAttendance"`
	GuestParticipation		*int	`json:"GuestParticipation"`
	GuestAttendance			*int	`json:"GuestAttendance"`
	ValidityStartDate		string	`json:"ValidityStartDate"`
	ValidityEndDate			string	`json:"ValidityEndDate"`
	CreationDate			string	`json:"CreationDate"`
	CreationTime			string	`json:"CreationTime"`
	IsCancelled				*bool	`json:"IsCancelled"`
}

type ValidationCheck struct {
	InvitationObjectType	string	`json:"InvitationObjectType"`
	InvitationObject		int		`json:"InvitationObject"`
	InvitationOwner			int		`json:"InvitationOwner"`
	InvitationGuest			int		`json:"InvitationGuest"`
	Invitation				int		`json:"Invitation"`
	CreationDate			string	`json:"CreationDate"`
	CreationTime			string	`json:"CreationTime"`
	IsCancelled				*bool	`json:"IsCancelled"`
}

type NumberRange struct {
	NumberRangeID            string `json:"NumberRangeID"`
	ServiceLabel             string `json:"ServiceLabel"`
	FieldNameWithNumberRange string `json:"FieldNameWithNumberRange"`
	LatestNumber             int    `json:"LatestNumber"`
}

type CalculateInvitationQueryGets struct {
	NumberRangeID            string `json:"NumberRangeID"`
	ServiceLabel             string `json:"service_label"`
	FieldNameWithNumberRange string `json:"FieldNameWithNumberRange"`
	InvitationLatestNumber	 int    `json:"InvitationLatestNumber"`
}
