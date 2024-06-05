package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Output_Formatter"
	dpfm_api_processing_formatter "data-platform-api-invitation-creates-rmq-kube/DPFM_API_Processing_Formatter"
	"data-platform-api-invitation-creates-rmq-kube/sub_func_complementer"
	"fmt"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	"golang.org/x/xerrors"
)

func (c *DPFMAPICaller) createSqlProcess(
	ctx context.Context,
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	subfuncSDC *sub_func_complementer.SDC,
	accepter []string,
	errs *[]error,
	log *logger.Logger,
) interface{} {
	var header *dpfm_api_output_formatter.Header
	var validationCheck *[]dpfm_api_output_formatter.ValidationCheck
	for _, fn := range accepter {
		switch fn {
		case "Header":
			var calculateInvitationQueryGets *sub_func_complementer.CalculateInvitationQueryGets
			var invitationIssuedID int

			calculateInvitationQueryGets = c.CalculateInvitation(errs)

			if calculateInvitationQueryGets == nil {
				err := xerrors.Errorf("calculateInvitationQueryGets is nil")
				*errs = append(*errs, err)
				return nil
			}

			invitationIssuedID = calculateInvitationQueryGets.InvitationLatestNumber + 1

			input.Header.Invitation = &invitationIssuedID

			header = c.headerCreateSql(nil, mtx, input, output, subfuncSDC, errs, log)

			if calculateInvitationQueryGets != nil {
				err := c.UpdateLatestNumber(errs, invitationIssuedID)
				if err != nil {
					*errs = append(*errs, err)
					return nil
				}
			}
		case "ValidationCheck":
			validationCheck = c.validationCheckCreateSql(nil, mtx, input, output, subfuncSDC, errs, log)
		default:

		}
	}

	data := &dpfm_api_output_formatter.Message{
		Header:				header,
		ValidationCheck:	validationCheck,
	}

	return data
}

func (c *DPFMAPICaller) updateSqlProcess(
	ctx context.Context,
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	accepter []string,
	errs *[]error,
	log *logger.Logger,
) interface{} {
	var header *dpfm_api_output_formatter.Header
	var validationCheck *[]dpfm_api_output_formatter.ValidationCheck
	for _, fn := range accepter {
		switch fn {
		case "Header":
			header = c.headerUpdateSql(mtx, input, output, errs, log)
		case "ValidationCheck":
			validationCheck = c.validationCheckUpdateSql(mtx, input, output, errs, log)
		default:

		}
	}

	data := &dpfm_api_output_formatter.Message{
		Header:				header,
		ValidationCheck:	validationCheck,
	}

	return data
}

func (c *DPFMAPICaller) headerCreateSql(
	ctx context.Context,
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	subfuncSDC *sub_func_complementer.SDC,
	errs *[]error,
	log *logger.Logger,
) *dpfm_api_output_formatter.Header {
	if ctx == nil {
		ctx = context.Background()
	}
	sessionID := input.RuntimeSessionID

	dpfm_api_output_formatter.ConvertToHeader(input, subfuncSDC)

	headerData := subfuncSDC.Message.Header
	res, err := c.rmq.SessionKeepRequest(ctx, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": headerData, "function": "InvitationHeader", "runtime_session_id": sessionID})
	if err != nil {
		err = xerrors.Errorf("rmq error: %w", err)
		*errs = append(*errs, err)
		return nil
	}
	res.Success()
	if !checkResult(res) {
		output.SQLUpdateResult = getBoolPtr(false)
		output.SQLUpdateError = "Header Data cannot insert"
		return nil
	}

	if output.SQLUpdateResult == nil {
		output.SQLUpdateResult = getBoolPtr(true)
	}

	data, err := dpfm_api_output_formatter.ConvertToHeaderCreates(subfuncSDC)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) validationCheckCreateSql(
	ctx context.Context,
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	subfuncSDC *sub_func_complementer.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.ValidationCheck {
	if ctx == nil {
		ctx = context.Background()
	}
	sessionID := input.RuntimeSessionID
	for _, validationCheckData := range *subfuncSDC.Message.ValidationCheck {
		res, err := c.rmq.SessionKeepRequest(ctx, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": validationCheckData, "function": "InvitationValidationCheck", "runtime_session_id": sessionID})
		if err != nil {
			err = xerrors.Errorf("rmq error: %w", err)
			*errs = append(*errs, err)
			return nil
		}
		res.Success()
		if !checkResult(res) {
			output.SQLUpdateResult = getBoolPtr(false)
			output.SQLUpdateError = "ValidationCheck Data cannot insert"
			return nil
		}
	}

	if output.SQLUpdateResult == nil {
		output.SQLUpdateResult = getBoolPtr(true)
	}

	data, err := dpfm_api_output_formatter.ConvertToValidationCheckCreates(subfuncSDC)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) headerUpdateSql(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *dpfm_api_output_formatter.Header {
	header := input.Header
	headerData := dpfm_api_processing_formatter.ConvertToHeaderUpdates(header)

	sessionID := input.RuntimeSessionID
	if headerIsUpdate(headerData) {
		res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": headerData, "function": "InvitationHeader", "runtime_session_id": sessionID})
		if err != nil {
			err = xerrors.Errorf("rmq error: %w", err)
			*errs = append(*errs, err)
			return nil
		}
		res.Success()
		if !checkResult(res) {
			output.SQLUpdateResult = getBoolPtr(false)
			output.SQLUpdateError = "Header Data cannot update"
			return nil
		}
	}

	if output.SQLUpdateResult == nil {
		output.SQLUpdateResult = getBoolPtr(true)
	}

	data, err := dpfm_api_output_formatter.ConvertToHeaderUpdates(header)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) validationCheckUpdateSql(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.ValidationCheck {
	req := make([]dpfm_api_processing_formatter.ValidationCheckUpdates, 0)
	sessionID := input.RuntimeSessionID

	header := input.Header
	for _, validationCheck := range header.ValidationCheck {
		validationCheckData := *dpfm_api_processing_formatter.ConvertToValidationCheckUpdates(header, validationCheck)

		if validationCheckIsUpdate(&validationCheckData) {
			res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": validationCheckData, "function": "InvitationValidationCheck", "runtime_session_id": sessionID})
			if err != nil {
				err = xerrors.Errorf("rmq error: %w", err)
				*errs = append(*errs, err)
				return nil
			}
			res.Success()
			if !checkResult(res) {
				output.SQLUpdateResult = getBoolPtr(false)
				output.SQLUpdateError = "ValidationCheck Data cannot update"
				return nil
			}
		}
		req = append(req, validationCheckData)
	}

	if output.SQLUpdateResult == nil {
		output.SQLUpdateResult = getBoolPtr(true)
	}

	data, err := dpfm_api_output_formatter.ConvertToValidationCheckUpdates(&req)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func headerIsUpdate(header *dpfm_api_processing_formatter.HeaderUpdates) bool {
	invitation := header.Invitation

	return !(invitation == 0)
}

func validationCheckIsUpdate(validationCheck *dpfm_api_processing_formatter.ValidationCheckUpdates) bool {
	invitationObjectType := validationCheck.InvitationObjectType
	invitationObject := validationCheck.InvitationObject
	invitationOwner := validationCheck.InvitationOwner
	InvitationGuest := validationCheck.InvitationGuest

	return !(invitationObjectType == "" || invitationObject == 0 || invitationOwner == 0 || InvitationGuest == 0)
}

func (c *DPFMAPICaller) CalculateInvitation(
	errs *[]error,
) *sub_func_complementer.CalculateInvitationQueryGets {
	pm := &sub_func_complementer.CalculateInvitationQueryGets{}

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_number_range_latest_number_data
		WHERE (ServiceLabel, FieldNameWithNumberRange) = (?, ?);`, "INVITATION", "Invitation",
	)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	for i := 0; true; i++ {
		if !rows.Next() {
			if i == 0 {
				*errs = append(*errs, fmt.Errorf("'data_platform_number_range_latest_number_data'テーブルに対象のレコードが存在しません。"))
				return nil
			} else {
				break
			}
		}
		err = rows.Scan(
			&pm.NumberRangeID,
			&pm.ServiceLabel,
			&pm.FieldNameWithNumberRange,
			&pm.InvitationLatestNumber,
		)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}

	return pm
}

func (c *DPFMAPICaller) UpdateLatestNumber(
	errs *[]error,
	invitationIssuedID int,
) error {
	_, err := c.db.Exec(`
			UPDATE data_platform_number_range_latest_number_data SET LatestNumber=(?)
			WHERE (ServiceLabel, FieldNameWithNumberRange) = (?, ?);`,
		invitationIssuedID,
		"INVITATION",
		"Invitation",
	)
	if err != nil {
		return xerrors.Errorf("'data_platform_number_range_latest_number_data'テーブルの更新に失敗しました。")
	}

	return nil
}
