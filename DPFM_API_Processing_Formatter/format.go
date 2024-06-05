package dpfm_api_processing_formatter

import (
	dpfm_api_input_reader "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Input_Reader"
)

func ConvertToHeaderUpdates(header dpfm_api_input_reader.Header) *HeaderUpdates {
	data := header

	return &HeaderUpdates{
		Invitation:                   *data.Invitation,
	}
}

func ConvertToValidationCheckUpdates(header dpfm_api_input_reader.Header, partner dpfm_api_input_reader.ValidationCheck) *ValidationCheckUpdates {
	data := validationCheck

	return &ValidationCheckUpdates{
		InvitationObjectType:    data.InvitationObjectType,
		InvitationObject:        data.InvitationObject,
		InvitationOwner:         data.InvitationOwner,
		InvitationGuest:		 data.InvitationGuest,
	}
}
