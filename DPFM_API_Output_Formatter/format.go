package dpfm_api_output_formatter

import (
	dpfm_api_input_reader "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_processing_formatter "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Processing_Formatter"
	"data-platform-api-invitation-creates-rmq-kube/sub_func_complementer"
	"encoding/json"

	"golang.org/x/xerrors"
)

func ConvertToHeaderCreates(subfuncSDC *sub_func_complementer.SDC) (*Header, error) {
	data := subfuncSDC.Message.Header

	header, err := TypeConverter[*Header](data)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func ConvertToValidationCheckCreates(subfuncSDC *sub_func_complementer.SDC) (*[]ValidationCheck, error) {
	validationChecks := make([]ValidationCheck, 0)

	for _, data := range *subfuncSDC.Message.ValidationCheck {
		validationCheck, err := TypeConverter[*ValidationCheck](data)
		if err != nil {
			return nil, err
		}

		validationChecks = append(validationChecks, *validationCheck)
	}

	return &validationChecks, nil
}

func ConvertToHeaderUpdates(headerData dpfm_api_input_reader.Header) (*Header, error) {
	data := headerData

	header, err := TypeConverter[*Header](data)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func ConvertToValidationCheckUpdates(validationCheckUpdates *[]dpfm_api_processing_formatter.ValidationCheckUpdates) (*[]ValidationCheck, error) {
	validationChecks := make([]ValidationCheck, 0)

	for _, data := range *validationCheckUpdates {
		validationCheck, err := TypeConverter[*ValidationCheck](data)
		if err != nil {
			return nil, err
		}

		validationChecks = append(validationChecks, *validationCheck)
	}

	return &validationChecks, nil
}

func ConvertToHeader(
	input *dpfm_api_input_reader.SDC,
	subfuncSDC *sub_func_complementer.SDC,
) *sub_func_complementer.SDC {
	subfuncSDC.Message.Header = &sub_func_complementer.Header{
		Invitation:				*input.Header.Invitation,
		InvitationType:			input.Header.InvitationType,
		InvitationOwner:		input.Header.InvitationOwner,
		InvitationGuest:		input.Header.InvitationGuest,
		InvitationObjectType:	input.Header.InvitationObjectType,
		InvitationObject:		input.Header.InvitationObject,
		OwnerParticipation:		input.Header.OwnerParticipation,
		OwnerAttendance:		input.Header.OwnerAttendance,
		GuestParticipation:		input.Header.GuestParticipation,
		GuestAttendance:		input.Header.GuestAttendance,
		ValidityStartDate:		input.Header.ValidityStartDate,
		ValidityEndDate:		input.Header.ValidityEndDate,
		CreationDate:			input.Header.CreationDate,
		CreationTime:			input.Header.CreationTime,
		IsCancelled:			input.Header.IsCancelled,
	}

	return subfuncSDC
}

func ConvertToValidationCheck(
	input *dpfm_api_input_reader.SDC,
	subfuncSDC *sub_func_complementer.SDC,
) *sub_func_complementer.SDC {
	var validationChecks []sub_func_complementer.ValidationCheck

	validationChecks = append(
		validationChecks,
		sub_func_complementer.ValidationCheck{
			InvitationObjectType:		input.Header.ValidationCheck[0].InvitationObjectType,
			InvitationObject:			input.Header.ValidationCheck[0].InvitationObject,
			InvitationOwner:			input.Header.ValidationCheck[0].InvitationOwner,
			InvitationGuest:			input.Header.ValidationCheck[0].InvitationGuest,
			Invitation:					input.Header.ValidationCheck[0].Invitation,
			CreationDate:				input.Header.ValidationCheck[0].CreationDate,
			CreationTime:				input.Header.ValidationCheck[0].CreationTime,
			IsCancelled:				input.Header.ValidationCheck[0].IsCancelled,
		},
	)

	subfuncSDC.Message.ValidationCheck = &validationChecks

	return subfuncSDC
}

func TypeConverter[T any](data interface{}) (T, error) {
	var dist T
	b, err := json.Marshal(data)
	if err != nil {
		return dist, xerrors.Errorf("Marshal error: %w", err)
	}
	err = json.Unmarshal(b, &dist)
	if err != nil {
		return dist, xerrors.Errorf("Unmarshal error: %w", err)
	}
	return dist, nil
}
