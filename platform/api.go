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
const CANCEL_DOCUMENT_ENDPOINT = "/api/v1.0/documents/state/{UUID}/state"
const REJECT_DOCUMENT_ENDPOINT = "/api/v1.0/documents/state/{UUID}/state"
const GET_RECENT_DOCUMENTS_ENDPOINT = "/api/v1.0/documents/recent"

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
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInSec     int    `json:"expires_in"`
	Scope            string `json:"scope"`
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
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInSec     int    `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Api struct {
	clientId     string
	clientSecret string
	AccessToken  string
}

func NewApi(clientId string, clientSecret string) *Api {
	return &Api{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (a *Api) DoTaxPayerLogin(req *TaxPayerLoginRequest) (interface{}, error) {
	loginUrl, err := url.JoinPath(common.SANDBOX_IDENTITY_BASE_URL, TAXPAYER_LOGIN_ENDPOINT)
	if err != nil {
		return "", err
	}

	body := url.Values{}
	body.Set("client_id", req.ClientID)
	body.Set("client_secret", req.ClientSecret)
	body.Set("grant_type", req.GrantType)
	body.Set("scope", req.Scope)

	httpReq, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request to API failed: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		// { "statusCode": 400, "message": "Bad Request" }
		// or {"error":"invalid_request"}
		return nil, errors.New("bad request. please check your client id and client secret")
	} else if resp.StatusCode == http.StatusOK {
		var retval TaxPayerLoginResponse

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

	return nil, fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
}

func (a *Api) DoIntemediarySystemLogin(req *IntermLoginRequest) (*IntermLoginResponse, error) {
	endpointUrl, err := url.JoinPath(common.SANDBOX_IDENTITY_BASE_URL, INTERM_LOGIN_ENDPOINT)
	if err != nil {
		return nil, err
	}

	httpBody := url.Values{}
	httpBody.Set("client_id", req.ClientID)
	httpBody.Set("client_secret", req.ClientSecret)
	httpBody.Set("grant_type", req.GrantType)
	httpBody.Set("scope", req.Scope)

	httpReq, err := http.NewRequest(http.MethodPost, endpointUrl, strings.NewReader(httpBody.Encode()))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("onbehalfof", req.OnBehalfOf)
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("bad request: please check your client ID and client secret")
	} else if resp.StatusCode == http.StatusOK {
		var retval IntermLoginResponse

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, &retval); err != nil {
			return nil, err
		}

		return &retval, nil
	}

	return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
}

type DocumentTypeVersion struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ActiveFrom    string  `json:"activeFrom"`    // TODO: timestamp
	ActiveTo      string  `json:"activeTo"`      // TODO: timestamp
	VersionNumber float32 `json:"versionNumber"` // the API returns a float for some reason
	Status        string  `json:"status"`        // possible values: draft, published, deactivated
}

type DocumentType struct {
	Id              int                   `json:"id"`
	InvoiceTypeCode int                   `json:"invoiceTypeCode"` // possible values: 1,2,3,4,11,12,13,14
	Description     string                `json:"description"`
	ActiveFrom      string                `json:"activeFrom"` // TODO: timestamp
	ActiveTo        string                `json:"activeTo"`   // TODO: timestamp
	Versions        []DocumentTypeVersion `json:"documentTypeVersions"`
}

type GetDocumentTypesResponse struct {
	Result []DocumentType `json:"result"`
	common.StandardErrResponse
}

func (a *Api) GetDocumentTypes() (*GetDocumentTypesResponse, error) {

	endpointUrl := common.SANDBOX_API_BASE_URL + "/api/v1.0/documenttypes"

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+a.AccessToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var retval GetDocumentTypesResponse

	if err := json.Unmarshal(respBody, &retval); err != nil {
		return nil, err
	}

	return &retval, nil
}

type WorkflowParameter struct {
	Id         int    `json:"id"`
	Parameter  string `json:"parameter"`
	Value      int    `json:"value"`
	ActiveFrom string `json:"activeFrom"`         // TODO: timestamp
	ActiveTo   string `json:"activeTo,omitempty"` // TODO: timestamp
}

type GetDocumentTypeByIdResponse struct {
	DocumentType
	WorkflowParameters []WorkflowParameter
}

func (a *Api) GetDocumentTypeById(id string) (*GetDocumentTypeByIdResponse, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + strings.Replace(GET_DOCUMENT_TYPE_BY_ID_ENDPOINT, "{id}", id, 1)

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+a.AccessToken)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var retval GetDocumentTypeByIdResponse

		if err = json.Unmarshal(respBytes, &retval); err != nil {
			return nil, err
		}

		return &retval, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("invalid id '%s'", id)
	}

	return nil, fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
}

type GetDocumentTypeVersionResponse struct {
	InvoiceTypeCode int
	Name            string
	Description     string
	VersionNumber   float32
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

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+a.AccessToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var retval GetDocumentTypeVersionResponse

		if err := json.Unmarshal(respBody, &retval); err != nil {
			return nil, err
		}

		return &retval, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("document type id '%s' with vid '%s' not found", id, version)
	}

	return nil, fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
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

	httpReq.Header.Add("Authorization", "Bearer "+a.AccessToken)

	query := httpReq.URL.Query()
	query.Add("idType", string(idType))
	query.Add("idValue", idValue)

	httpReq.URL.RawQuery = query.Encode()

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

	return false, fmt.Errorf("unexpected error code %d", resp.StatusCode)
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
	common.StandardErrResponse
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

type CancelDocumentRequest struct {
	DesiredStatus string `json:"status"` // required
	Reason        string `json:"reason"` // required. max 300 chars
}

type CancelDocumentResponse struct {
	UUID   string             `json:"uuid"`
	Status string             `json:"status"`
	Error  common.ErrResponse `json:"error"`
}

func (a *Api) CancelDocument(docUuid string, reason string) (*CancelDocumentResponse, error) {
	endpointUrl := strings.ReplaceAll(CANCEL_DOCUMENT_ENDPOINT, "{UUID}", docUuid)

	reqBody := CancelDocumentRequest{
		DesiredStatus: "cancelled",
		Reason:        reason,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPut, endpointUrl, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval CancelDocumentResponse

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBytes, &retval); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		// TODO: handle success
		return &retval, nil
	} else {
		// TODO: handle error
		return &retval, nil
	}
}

type RejectDocumentRequest struct {
	DesiredStatus string `json:"status"` // required
	Reason        string `json:"reason"` // required. max 300 chars
}

type RejectDocumentResponse struct {
	UUID   string             `json:"uuid"`
	Status string             `json:"status"`
	Error  common.ErrResponse `json:"error"`
}

func (a *Api) RejectDocument(docUuid string, reason string) (*RejectDocumentResponse, error) {
	endpointUrl := strings.ReplaceAll(CANCEL_DOCUMENT_ENDPOINT, "{UUID}", docUuid)

	reqBody := RejectDocumentRequest{
		DesiredStatus: "rejected",
		Reason:        reason,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPut, endpointUrl, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval RejectDocumentResponse

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBytes, &retval); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		// TODO: handle success
		return &retval, nil
	} else {
		// TODO: handle error
		return &retval, nil
	}
}

type GetRecentDocumentsQuery struct {
	PageNo             int64   // Optional: number of the page to retrieve. Typically this parameter value is derived from initial parameter less call when caller learns total amount of page of certain size 	3 	Optional
	PageSize           int64   // Optional: number of the documents to retrieve per page. Page size cannot exceed system configured maximum page size for this API 	20 	Optional
	SubmissionDateFrom *string // Optional: The start date and time when the document was submitted to the e-Invoice API, Time to be supplied in UTC timezone. Mandatory when ‘submissionDateTo’ is provided 	2022-11-25T01:59:10Z 	Optional
	SubmissionDateTo   *string // Optional: The end date and time when the document was submitted to the e-Invoice API, Time to be supplied in UTC timezone. Mandatory when ‘submissionDateFrom’ is provided 	2022-12-22T23:59:59Z 	Optional
	IssueDateFrom      *string // Optional: The start date and time when the document was issued. Mandatory when ‘issueDateTo’ is provided 	2021-02-25T23:55:10Z 	Optional
	IssueDateTo        *string // Optional: The end date and time when the document was issued. Mandatory when ‘issueDateFrom’ is provided 	2021-03-10T01:59:10Z 	Optional
	Direction          *string // Optional: direction of the document. Possible values: (Sent, Received) 	Sent 	Optional
	Status             *string // Optional: status of the document. Possible values: (Valid, Invalid, Cancelled, Submitted) 	Valid 	Optional
	DocumentType       *string // Optional: Document type code. 	01 	Optional
	ReceiverId         *string // Optional: Document recipient identifier. Only can be used when ‘Direction’ filter is set to Sent. Possible values: (Business registration number, National ID(IC), Passport Number, Army ID) 	BRN example: 201901234567 - NRIC example: 770625015324 - Passport number example: A12345678 - Army number example: 551587706543 	Optional
	ReceiverIdType     *string // Optional: Document recipient identifier type. Only can be used when ‘Direction’ filter is set to Sent. Possible values: (BRN, PASSPORT, NRIC, ARMY) This is mandatory in case the receiverId is provided 	PASSPORT 	Optional
	IssuerIdType       *string // Optional: Document issuer identifier type. Only can be used when ‘Direction’ filter is set to Received. Possible values: (BRN, PASSPORT, NRIC, ARMY) This is mandatory in case the issuerId is provided 	PASSPORT 	Optional
	ReceiverTin        *string // Optional: Document recipient TIN. Only can be used when ‘Direction’ filter is set to Sent. 	C2584563200 	Optional
	IssuerTin          *string // Optional: Document issuer identifier. Only can be used when ‘Direction’ filter is set to Received. 	C2584563200 	Optional
	IssuerId           *string // Optional: Document issuer identifier. Only can be used when ‘Direction’ filter is set to Received. Possible values: (Business registration number, National ID(IC), Passport Number, Army ID)
}

type GetRecentDocumentsResponse struct {
	UUID                  string             `json:"uuid"`                  // 	Unique document ID in e-Invoice 	42S512YACQBRSRHYKBXBTGQG22
	SubmissionUid         string             `json:"submissionUID"`         // 	Unique ID of the submission that document was part of 	XYE60M8ENDWA7V9TKBXBTGQG10
	LongId                string             `json:"longId"`                // 	Unique long temporary Id that can be used to query document data anonymously 	YQH73576FY9VR57B…
	InternalId            string             `json:"internalId"`            // 	Internal ID used in submission for the document 	PZ-234-A
	TypeName              string             `json:"typeName"`              // 	Unique name of the document type that can be used in submission of the documents. 	invoice
	TypeVersionName       string             `json:"typeVersionName"`       // 	Name of the document type version within the document type that can be used in document submission to identify document type version being submitted 	1.0
	IssuerTIN             string             `json:"issuerTin"`             // 	TIN of issuer 	C2584563200
	IssuerName            string             `json:"issuerName"`            // 	Issuer company name 	AMS Setia Jaya Sdn. Bhd.
	ReceiverId            string             `json:"receiverId"`            // 	Optional: receiver registration number (can be national ID or foreigner ID). 	BRN example: 201901234567 - NRIC example: 770625015324 - Passport number example: A12345678 - Army number example: 551587706543
	ReceiverName          string             `json:"receiverName"`          // 	Optional: receiver name (can be company name or person’s name) 	AMS Setia Jaya Sdn. Bhd.
	DateTimeIssued        string             `json:"dateTimeIssued"`        // DateTime 	The date and time when the document was issued. 	2015-02-13T13:15:00Z
	DateTimeReceived      string             `json:"dateTimeReceived"`      // DateTime 	The date and time when the document was submitted. 	2015-02-13T14:20:00Z
	DateTimeValidated     string             `json:"dateTimeValidated"`     // DateTime 	The date and time when the document passed all validations and moved to the valid state. 	2015-02-13T14:20:00Z
	TotalSales            string             `json:"totalSales"`            // Decimal 	Total sales amount of the document in MYR. 	10.10
	TotalDiscount         string             `json:"totalDiscount"`         // Decimal 	Total discount amount of the document in MYR. 	50.00
	NetAmount             string             `json:"netAmount"`             // Decimal 	Total net amount of the document in MYR. 	100.70
	Total                 string             `json:"total"`                 // Decimal 	Total amount of the document in MYR. 	124.09
	Status                string             `json:"status"`                // 	Status of the document - Submitted, Valid, Invalid, Cancelled 	Valid
	CancelDateTime        string             `json:"cancelDateTime"`        // Date 	Refer to the document cancellation that has been initiated by the taxpayer “issuer” of the document on the system, will be in UTC format 	2021-02-25T01:59:10Z
	RejectRequestDateTime string             `json:"rejectRequestDateTime"` // 	Date 	Refer to the document rejection request that has been initiated by the taxpayer “receiver” of the document on the system, will be in UTC format 	2021-02-25T01:59:10Z
	DocumentStatusReason  string             `json:"documentStatusReason"`  // 	Mandatory: Reason of the cancellation or rejection of the document. 	Examples of reasons: Wrong buyer details or Wrong invoice details or any other reasons as appropriate
	CreatedByUserId       string             `json:"createdByUserId"`       // 	User created the document. Can be ERP ID or User Email C1XXXXXXXX00:9e21b10c-41c4-9323-c590-95abcb6e4e4d, general.ams@supplier.com
	SupplierTIN           string             `json:"supplierTIN"`           // 	TIN of issuer 	C2584563200
	SupplierName          string             `json:"supplierName"`          // 	Supplier company name 	AMS Setia Jaya Sdn. Bhd.
	SubmissionChannel     string             `json:"submissionChannel"`     // 	Channel through which document was introduced into the system 	possible values: ERP, Invoicing Portal, InvoicingMobileApp
	IntermediaryName      string             `json:"intermediaryName"`      // 	Intermediary company name 	AMS Setia Jaya Sdn. Bhd.
	IntermediaryTIN       string             `json:"intermediaryTIN"`       // 	TIN of intermediary 	C2584563200
	BuyerName             string             `json:"buyerName"`             // 	Buyer company name 	AMS Setia Jaya Sdn. Bhd.
	BuyerTIN              string             `json:"buyerTIN"`              // 	Tin of buyer 	C2584563200
	Metadata              RecentDocumentMeta `json:"metadata"`              // 	Information about the results retrieved or results matching the query 	See structure
	common.StandardErrResponse
}

type RecentDocumentMeta struct {
	TotalPages string `json:"totalPages"` // Number Total count of pages based on the supplied (or default) page size
	TotalCount string `json:"totalCount"` // Number Total count of matching objects
}

func (a *Api) GetRecentDocuments(query *GetRecentDocumentsQuery) (*GetRecentDocumentsResponse, error) {
	httpReq, err := buildGetRecentDocumentsRequest(query)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval GetRecentDocumentsResponse
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBytes, &retval); err != nil {
		return nil, err
	}

	return &retval, nil
}

func buildGetRecentDocumentsRequest(query *GetRecentDocumentsQuery) (*http.Request, error) {
	endpointUrl := common.SANDBOX_API_BASE_URL + GET_RECENT_DOCUMENTS_ENDPOINT

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	queryParam := httpReq.URL.Query()

	if query.PageNo > 0 {
		queryParam.Add("pageNo", fmt.Sprintf("%d", query.PageNo))
	}

	if query.PageSize > 0 {
		queryParam.Add("pageSize", fmt.Sprintf("%d", query.PageSize))
	}

	if query.SubmissionDateFrom != nil {
		queryParam.Add("submissionDateFrom", *query.SubmissionDateFrom)
	}

	if query.SubmissionDateTo != nil {
		queryParam.Add("submissionDateTo", *query.SubmissionDateTo)
	}

	if query.IssueDateFrom != nil {
		queryParam.Add("issueDateFrom", *query.IssueDateFrom)
	}

	if query.IssueDateTo != nil {
		queryParam.Add("issueDateTo", *query.IssueDateTo)
	}

	if query.Direction != nil {
		queryParam.Add("direction", *query.Direction)
	}

	if query.Status != nil {
		queryParam.Add("status", *query.Status)
	}

	if query.DocumentType != nil {
		queryParam.Add("documentType", *query.DocumentType)
	}

	if query.ReceiverId != nil {
		queryParam.Add("receiverId", *query.ReceiverId)
	}

	if query.ReceiverIdType != nil {
		queryParam.Add("receiverIdType", *query.ReceiverIdType)
	}

	if query.IssuerIdType != nil {
		queryParam.Add("issuerIdType", *query.IssuerIdType)
	}

	if query.ReceiverTin != nil {
		queryParam.Add("receiverTin", *query.ReceiverTin)
	}

	if query.IssuerTin != nil {
		queryParam.Add("issuerTin", *query.IssuerTin)
	}

	if query.IssuerId != nil {
		queryParam.Add("issuerId", *query.IssuerId)
	}

	httpReq.URL.RawQuery = queryParam.Encode()

	return httpReq, nil
}

const GET_SUBMISSION_ENDPOINT = "/api/v1.0/documentsubmissions/{submissionUid}"

type GetSubmissionQuery struct {
	PageNo   int64
	PageSize int64
}

type GetSubmissionResponse struct {
	SubmissionUid    string                      `json:"submissionUid"`    // String 	Unique document submission ID in e-Invoice 	HJSD135P2S7D8IU
	DocumentCount    string                      `json:"documentCount"`    // Number 	Total count of documents in submission that were accepted for processing 	234
	DateTimeReceived string                      `json:"dateTimeReceived"` // DateTime 	The date and time when the submission was received by e-Invoice. 	2015-02-13T14:20:10Z
	OverallStatus    string                      `json:"overallStatus"`    // String 	Overall status of the batch processing. Values: in progress, valid, partially valid, invalid 	valid
	DocumentSummary  []SubmissionDocumentSummary `json:"documentSummary"`  // Document Summary[] 	List of the retrieved batch documents in current page. 	See structure.
	common.StandardErrResponse
}

type SubmissionDocumentSummary struct {
	UUID                  string `json:"uuid"`                  // String 	Unique document ID in e-Invoice 	F9D425P6DS7D8IU
	SubmissionUid         string `json:"submissionUid"`         // String 	Unique ID of the submission the document was part of 	HJSD135P2S7D8IU
	LongId                string `json:"longId"`                // String 	Unique long temporary Id that can be used to query document data anonymously. The long id will be returned only for valid documents 	LIJAF97HJJKH 8298KHADH0990 8570FDKK9S2LSIU HB377373
	InternalId            string `json:"internalId"`            // String 	Internal ID used in submission for the document 	PZ-234-A
	TypeName              string `json:"typeName"`              // String 	Unique name of the document type that can be used in submission of the documents. 	invoice
	TypeVersionName       string `json:"typeVersionName"`       // String 	Name of the document type version within the document type that can be used in document submission to identify document type version being submitted 	1.0
	IssuerTIN             string `json:"issuerTin"`             // String 	TIN of issuer 	C2584563200
	IssuerName            string `json:"issuerName"`            // String 	Issuer company name 	AMS Setia Jaya Sdn. Bhd.
	ReceiverId            string `json:"receiverId"`            // String 	Optional: receiver registration number (can be national ID or foreigner ID). 	201901234567
	ReceiverName          string `json:"receiverName"`          // String 	Optional: receiver name (can be company name or person’s name) 	AMS Setia Jaya Sdn. Bhd.
	DateTimeIssued        string `json:"dateTimeIssued"`        // DateTime 	The date and time when the document was issued in the UTC format. 	2015-02-13T13:15:10Z
	DateTimeReceived      string `json:"dateTimeReceived"`      // DateTime 	The date and time when the document was submitted in the UTC format. 	2015-02-13T13:15:10Z
	DateTimeValidated     string `json:"dateTimeValidated"`     // DateTime 	The date and time when the document passed all validations and moved to the valid state. 	2015-02-13T13:15:10Z
	TotalExcludingTax     string `json:"totalExcludingTax"`     // Decimal 	Total sales amount of the document in MYR. 	10.10
	TotalDiscount         string `json:"totalDiscount"`         // Decimal 	Total discount amount of the document in MYR. 	50.00
	TotalNetAmount        string `json:"totalNetAmount"`        // Decimal 	Total net amount of the document in MYR. 	100.70
	TotalPayableAmount    string `json:"totalPayableAmount"`    // Decimal 	Total amount of the document in MYR. 	124.09
	Status                string `json:"status"`                // String 	Status of the document - Submitted, Valid, Invalid, Cancelled 	Valid
	CancelDateTime        string `json:"cancelDateTime"`        // Date 	Refer to the document cancellation that has been initiated by the taxpayer ‘issuer’ of the document on the system, will be in UTC format 	2021-02-25T01:59:10Z
	RejectRequestDateTime string `json:"rejectRequestDateTime"` // Date 	Refer to the document rejection request that has been initiated by the taxpayer ‘receiver’ of the document on the system, will be in UTC format 	2021-02-25T01:59:10Z
	DocumentStatusReason  string `json:"documentStatusReason"`  // String 	Mandatory: Reason of the cancellation or rejection of the document. 	Examples of reasons: Wrong buyer details or Wrong invoice details or any other reasons as appropriate
	CreatedByUserId       string `json:"createdByUserId"`       // String 	User created the document. Can be ERP ID or User Email 	1XXXXXXXX00:9e21b10c-41c4-9323-c590-95abcb6e4e4d general.ams@supplier.com
}

// This API allows caller to get details of a single submission to check
// its processing status after initially submitting it
// and getting back unique submission identifier.
//
// This API is available to submitter only as it might contain documents issued to multiple receivers.
func (a *Api) GetSubmission(submissionUid string, query *GetSubmissionQuery) (*GetSubmissionResponse, error) {
	endpointUrl := strings.ReplaceAll(GET_SUBMISSION_ENDPOINT, "{submissionUid}", submissionUid)

	httpReq, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}

	if query != nil {
		q := httpReq.URL.Query()

		if query.PageNo > 0 {
			q.Add("pageNo", fmt.Sprintf("%d", query.PageNo))
		}

		if query.PageSize > 0 {
			q.Add("pageSize", fmt.Sprintf("%d", query.PageSize))
		}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var retval GetSubmissionResponse

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBytes, &retval); err != nil {
		return nil, err
	}

	return &retval, nil
}

const GET_DOCUMENT_ENDPOINT = "/api/v1.0/documents/{uuid}/raw"

func (a *Api) GetDocument(docUuid string) (any, error) {
	// TODO
	return nil, nil
}

const GET_DOCUMENT_DETAILS_ENDPOINT = "/api/v1.0/documents/{uuid}/details"

func (a *Api) GetDocumentDetails(docUuid string) (any, error) {
	// TODO
	return nil, nil
}

const SEARCH_DOCUMENTS_ENDPOINT = "/api/v1.0/documents/search"

func (a *Api) SearchDocuments() (any, error) {
	// TODO
	return nil, nil
}
