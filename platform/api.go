package platform

//
// LHDN Platform API
// Reference: https://sdk.myinvois.hasil.gov.my/api/
//

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/programmer-my/einvoice-go/common"
)

const TAXPAYER_LOGIN_ENDPOINT = "/connect/token"
const INTERM_LOGIN_ENDPOINT = "/connect/token"
const GET_DOCUMENT_TYPES_ENDPOINT = "/api/v1.0/documenttypes"
const GET_DOCUMENT_TYPE_BY_ID_ENDPOINT = "/api/v1.0/documenttypes/{id}"
const GET_DOCUMENT_TYPE_VERSION_ENDPOINT = "/api/v1.0/documenttypes/{id}/versions/{vid}"
const GET_NOTIFICATIONS_ENDPOINT = "/api/v1.0/notifications/taxpayer"

const VALIDATE_TIN_ENDPOINT = "/api/v1.0/taxpayer/validate/{tin}"
const SUBMIT_DOCUMENT_ENDPOINT = "/api/v1.0/documentsubmissions/"

type TinIdType string

const (
	ID_NRIC     TinIdType = "NRIC"
	ID_PASSPORT TinIdType = "PASSPORT"
	ID_BRN      TinIdType = "BRN"
	ID_ARMY     TinIdType = "ARMY"
)

type TaxPayerLoginRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

type TaxPayerLoginResponse struct {
	Success bool
	Data    interface{}
}

func (s *TaxPayerLoginResponse) UnmarshalJSON(b []byte) error {
	m := make(map[string]any)

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	successResp := &TaxPayerLoginSuccessResp{}
	errorResp := &TaxPayerLoginErrorResp{}

	accessToken, accessTokenFieldExists := m["access_token"]
	if accessTokenFieldExists {
		successResp.AccessToken = accessToken.(string)

		if tokenType, ok := m["token_type"]; ok {
			successResp.TokenType = tokenType.(string)
		}

		if expiresInSec, ok := m["expires_in"]; ok {
			expiresInSecStr := expiresInSec.(string)
			expiresInSecInt, err := strconv.ParseInt(expiresInSecStr, 10, 32)
			if err != nil {
				return fmt.Errorf("expected 'expires_in' to be numeric, got %s instead: %s", expiresInSecStr, err)
			}
			successResp.ExpiresInSec = int(expiresInSecInt)
		}

		if scope, ok := m["scope"]; ok {
			successResp.Scope = scope.(string)
		}

		s.Success = true
		s.Data = successResp

		return nil
	}

	errorType, errorFieldExists := m["error"]
	if errorFieldExists {
		errorResp.Error = errorType.(string)

		if errorDesc, ok := m["error_description"]; ok {
			errorResp.ErrorDescription = errorDesc.(string)
		}

		s.Success = false
		s.Data = errorResp

		return nil
	}

	return fmt.Errorf("JSON unmarshalling error: field 'accessToken' or 'error' expected")
}

type TaxPayerLoginSuccessResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresInSec int    `json:"expires_in,string"`
	Scope        string `json:"scope"`
}

type TaxPayerLoginErrorResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewTaxPayerLoginRequest(clientId string, clientSecret string, scope []string) *TaxPayerLoginRequest {
	return &TaxPayerLoginRequest{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
		Scope:        strings.Join(scope, ","),
	}
}

type IntermLoginRequest struct {
	OnBehalfOf string
	TaxPayerLoginRequest
}

type IntermLoginResponse struct {
	Success bool
	Data    interface{}
}

type IntermLoginSuccessResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresInSec int    `json:"expires_in,string"`
	Scope        string `json:"scope"`
}

type IntermLoginErrorResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *IntermLoginResponse) UnmarshalJSON(b []byte) error {
	m := make(map[string]any)

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	successResp := &IntermLoginSuccessResp{}
	errorResp := &IntermLoginErrorResp{}

	accessToken, accessTokenFieldExists := m["access_token"]
	if accessTokenFieldExists {
		successResp.AccessToken = accessToken.(string)

		if tokenType, ok := m["token_type"]; ok {
			successResp.TokenType = tokenType.(string)
		}

		if expiresInSec, ok := m["expires_in"]; ok {
			expiresInSecStr := expiresInSec.(string)
			expiresInSecInt, err := strconv.ParseInt(expiresInSecStr, 10, 32)
			if err != nil {
				return fmt.Errorf("expected 'expires_in' to be numeric, got %s instead: %s", expiresInSecStr, err)
			}
			successResp.ExpiresInSec = int(expiresInSecInt)
		}

		if scope, ok := m["scope"]; ok {
			successResp.Scope = scope.(string)
		}

		s.Success = true
		s.Data = successResp

		return nil
	}

	errorType, errorFieldExists := m["error"]
	if errorFieldExists {
		errorResp.Error = errorType.(string)

		if errorDesc, ok := m["error_description"]; ok {
			errorResp.ErrorDescription = errorDesc.(string)
		}

		s.Success = false
		s.Data = errorResp

		return nil
	}

	return fmt.Errorf("JSON unmarshalling error: field 'accessToken' or 'error' expected")
}

type Api struct {
}

func (a *Api) DoTaxPayerLogin(req *TaxPayerLoginRequest) (interface{}, error) {
	loginUrl, err := url.JoinPath(common.SANDBOX_IDENTITY_BASE_URL, TAXPAYER_LOGIN_ENDPOINT)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal failed: %s", err)
	}

	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("request to API failed: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return nil, errors.New("bad request. please check your client id and client secret")
	}

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ = io.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	return string(body), nil
}

func (a *Api) DoIntemediarySystemLogin(req *IntermLoginRequest) (interface{}, error) {
	endpointUrl, err := url.JoinPath(common.SANDBOX_IDENTITY_BASE_URL, INTERM_LOGIN_ENDPOINT)
	if err != nil {
		return nil, err
	}

	httpBody := map[string]string{
		"client_id":     req.ClientID,
		"client_secret": req.ClientSecret,
		"grant_type":    req.GrantType,
		"scope":         req.Scope,
	}

	bodyBytes, err := json.Marshal(httpBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, endpointUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("onbehalfof", req.OnBehalfOf)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

type DocumentTypeVersion struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ActiveFrom    string `json:"activeFrom"` // TODO: timestamp
	ActiveTo      string `json:"activeTo"`   // TODO: timestamp
	VersionNumber string `json:"versionNumber"`
	Status        string `json:"status"` // possible values: draft, published, deactivated
}

type DocumentType struct {
	Id              string                `json:"id"`
	InvoiceTypeCode string                `json:"invoiceTypeCode"` // possible values: 1,2,3,4,11,12,13,14
	Description     string                `json:"description"`
	ActiveFrom      string                `json:"activeFrom"` // TODO: timestamp
	ActiveTo        string                `json:"activeTo"`   // TODO: timestamp
	Versions        []DocumentTypeVersion `json:"documentTypeVersions"`
}

type GetDocumentTypesResponse struct {
	Result []DocumentType `json:"result"`
	common.StandardErrResponse
}

func (a *Api) GetDocumentTypes() (any, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + "/api/v1.0/documenttypes"

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

type WorkflowParameter struct {
	Id         string `json:"id"`
	Parameter  string `json:"parameter"`
	Value      int    `json:"value"`
	ActiveFrom string `json:"activeFrom"`         // TODO: timestamp
	ActiveTo   string `json:"activeTo,omitempty"` // TODO: timestamp
}

type GetDocumentTypeByIdResponse struct {
	DocumentType
	WorkflowParameters []WorkflowParameter
}

func (a *Api) GetDocumentTypeById(id string) (any, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + strings.Replace(GET_DOCUMENT_TYPE_BY_ID_ENDPOINT, "{id}", id, 1)

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned with status code %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var retval GetDocumentTypeByIdResponse

	err = json.Unmarshal(respBytes, &retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

type GetDocumentTypeVersionResponse struct {
	InvoiceTypeCode string
	Name            string
	Description     string
	VersionNumber   string
	Status          string // published, deactivated
	ActiveFrom      string // TODO: timestamp
	ActiveTo        string // TODO: timestamp
	JsonSchema      string
	XmlSchema       string
}

func (a *Api) GetDocumentTypeVersion(id string, version string) (any, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL +
		strings.Replace(
			strings.Replace(GET_DOCUMENT_TYPE_VERSION_ENDPOINT, "{id}", id, 1),
			"{vid}", version, 1,
		)

	resp, err := http.Post(endpointUrl, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

type Notification struct {
	Notificationid    string
	ReceivedDateTime  string // TODO: timestamp
	DeliveredDateTime string // TODO: timestamp, optional
	TypeId            string
	TypeName          string
	FinalMessage      string // optional
	Channel           string // possible values: email, push
	Address           string
	Language          string // possible values: ms,en
	Status            string // possible values: pending, batched, delivered, error
	DeliveryAttempts  []NotificationDeliveryAttempt
}

type NotificationDeliveryAttempt struct {
	AttemptDateTime string
	Status          string // possible values - delivered, error
	StatusDetails   string
}

type NotificationMetadata struct {
	TotalPages int `json:"totalPages,string"`
	TotalCount int `json:"totalCount,string"`
}

type GetNotificationsResponse struct {
	Result   []Notification       `json:"result"`
	Metadata NotificationMetadata `json:"metadata"`
}

func (a *Api) GetNotifications(dateFrom string, dateTo string, notifType string,
	language string, status string, channel string,
	pageNo string, pageSize string) (any, error) {

	endpointUrl := common.SANDBOX_API_BASE_URL + GET_NOTIFICATIONS_ENDPOINT

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	queryParam := httpReq.URL.Query()
	queryParam.Add("dateFrom", dateFrom)
	queryParam.Add("dateTo", dateTo)
	queryParam.Add("type", notifType)
	queryParam.Add("language", language)
	queryParam.Add("status", status)
	queryParam.Add("channel", channel)
	queryParam.Add("pageNo", pageNo)
	queryParam.Add("pageSize", pageSize)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval GetNotificationsResponse

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, &retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func (a *Api) ValidateTIN(tin string, idType TinIdType, idValue string) (bool, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + strings.ReplaceAll(VALIDATE_TIN_ENDPOINT, "{tin}", tin)

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return false, err
	}

	query := httpReq.URL.Query()
	query.Add("idType", string(idType))
	query.Add("idValue", idValue)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, fmt.Errorf("TIN %s not found", tin)
	} else if resp.StatusCode == http.StatusBadRequest {
		return false, fmt.Errorf("bad request")
	}

	return false, fmt.Errorf("unexpected error code %d: ", resp.StatusCode)
}

type Document struct {
	Format         string `json:"format"`   // TODO: valid: XML, JSON
	Document       string `json:"document"` // in base64
	DocumentSHA256 string `json:"documentHash"`
	CodeNumber     string `json:"codeNumber"`
}

type SubmitDocumentRequest struct { // max: 5MB
	Documents []Document `json:"documents"` // max: 100 items, 300KB per item
}

type SubmitDocumentResponse struct {
	SubmissionUID     string
	AcceptedDocuments []AcceptedDocuments
	RejectedDocuments []RejectedDocuments
}

type AcceptedDocuments struct {
	UUID              string `json:"uuid"`
	InvoiceCodeNumber string `json:"invoiceCodeNumber"`
}

type RejectedDocuments struct {
	InvoiceCodeNumber string `json:"invoiceCodeNumber"`
	Error             common.ErrResponse
}

func (a *Api) SubmitDocument() (any, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + SUBMIT_DOCUMENT_ENDPOINT

	var httpBody SubmitDocumentRequest // TODO: fill struct

	bodyBytes, err := json.Marshal(httpBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpointUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval SubmitDocumentResponse

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, &retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
