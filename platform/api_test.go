package platform

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalTaxPayerSuccessLoginResp(t *testing.T) {
	successJson := `{
		"access_token":"somejwttoken",
		"expires_in":3600,
		"token_type":"Bearer",
		"scope":"InvoicingAPI"
	}`

	var resp TaxPayerLoginResponse

	err := json.Unmarshal([]byte(successJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	if resp.AccessToken != "somejwttoken" {
		t.Errorf("expected %s as access token, got %s", "123123", resp.AccessToken)
	}

	if resp.TokenType != "Bearer" {
		t.Errorf("expected %s as access token, got %s", "Bearer", resp.TokenType)
	}

	if resp.ExpiresInSec != 3600 {
		t.Errorf("expected %s as access token, got %d", "3600", resp.ExpiresInSec)
	}

	if resp.Scope != "InvoicingAPI" {
		t.Errorf("expected %s as access token, got %s", "InvoicingAPI", resp.Scope)
	}

}

func TestUnmarshalTaxPayerErrorResp(t *testing.T) {
	errorJson := `{
		"error": "invalid_request",
		"error_description": "User blocked",
		"error_uri": ""
	}`

	var resp TaxPayerLoginResponse

	err := json.Unmarshal([]byte(errorJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedErrorMsg := "invalid_request"
	expectedErrorDesc := "User blocked"

	if resp.Error != expectedErrorMsg {
		t.Errorf("expected ErrorMsg to be %s, got %s", expectedErrorMsg, resp.Error)
	}

	if resp.ErrorDescription != expectedErrorDesc {
		t.Errorf("expected ErrorDescription to be %s, got %s", expectedErrorDesc, resp.ErrorDescription)
	}
}

func TestUnmarshalIntermLoginSuccessResp(t *testing.T) {
	successJson := `{
		"access_token":"somejwttoken",
		"expires_in":3600,
		"token_type":"Bearer",
		"scope":"InvoicingAPI"
	}`

	var resp IntermLoginResponse

	err := json.Unmarshal([]byte(successJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedAccessToken := "somejwttoken"
	if resp.AccessToken != expectedAccessToken {
		t.Errorf("expected access_token to be %s, got %s", expectedAccessToken, resp.AccessToken)
	}

	expectedTokenType := "Bearer"
	if resp.TokenType != expectedTokenType {
		t.Errorf("expected token_type to be %s, got %s", expectedTokenType, resp.TokenType)
	}

	expectedTokenScope := "InvoicingAPI"
	if resp.Scope != expectedTokenScope {
		t.Errorf("expected scope to be %s, got %s", expectedAccessToken, resp.Scope)
	}

	expectedExpiresIn := 3600
	if resp.ExpiresInSec != expectedExpiresIn {
		t.Errorf("expected expires_in to be %d, got %d", expectedExpiresIn, resp.ExpiresInSec)
	}
}

func TestUnmarshalIntermLoginErrorResp(t *testing.T) {
	errorJson := `{
		"error": "invalid_request",
		"error_description": "User blocked",
		"error_uri": ""
	}`

	var resp IntermLoginResponse

	err := json.Unmarshal([]byte(errorJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedErrorMsg := "invalid_request"
	expectedErrorDesc := "User blocked"

	if resp.Error != expectedErrorMsg {
		t.Errorf("expected ErrorMsg to be %s, got %s", expectedErrorMsg, resp.Error)
	}

	if resp.ErrorDescription != expectedErrorDesc {
		t.Errorf("expected ErrorDescription to be %s, got %s", expectedErrorDesc, resp.ErrorDescription)
	}
}

func TestUnmarshalGetDocumentTypesResponseSuccess(t *testing.T) {
	successJson := `{
		"result": [
			{
				"id": 45,
				"invoiceTypeCode": 1,
				"description": "Invoice",
				"activeFrom": "2015-02-13T13:15:00Z",
				"activeTo": "2027-03-01T00:00:00Z",
				"documentTypeVersions": [
					{
						"id": 3,
						"name": "version 1",
						"description": "document version 1",
						"activeFrom": "2015-02-13T13:15:00Z",
						"activeTo": "2027-03-01T00:00:00Z",
						"versionNumber": 1.1,
						"status": "draft"
					}
				]
			}
		]
	}`

	var resp GetDocumentTypesResponse

	err := json.Unmarshal([]byte(successJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedResultLen := 1
	if len(resp.Result) != expectedResultLen {
		t.Errorf("expected result length to be %d, got %d", expectedResultLen, len(resp.Result))
	}

	docType := resp.Result[0]
	expectedId := 45
	expectedInvoiceTypeCode := 1
	expectedDescription := "Invoice"
	expectedActiveFrom := "2015-02-13T13:15:00Z"
	expectedActiveTo := "2027-03-01T00:00:00Z"
	docTypeVersion1 := docType.Versions[0]
	expectedDocTypeVersionLen := 1

	if docType.Id != expectedId {
		t.Errorf("expected Id to be %d, got %d", expectedId, docType.Id)
	}

	if docType.InvoiceTypeCode != expectedInvoiceTypeCode {
		t.Errorf("expected InvoiceTypeCode to be %d, got %d", expectedInvoiceTypeCode, docType.InvoiceTypeCode)
	}

	if docType.Description != expectedDescription {
		t.Errorf("expected Description to be %s, got %s", expectedDescription, docType.Description)
	}

	if docType.ActiveFrom != expectedActiveFrom {
		t.Errorf("expected ActiveFrom to be %s, got %s", expectedActiveFrom, docType.ActiveFrom)
	}

	if docType.ActiveTo != expectedActiveTo {
		t.Errorf("expected ActiveTo to be %s, got %s", expectedActiveTo, docType.ActiveTo)
	}

	if len(docType.Versions) != expectedDocTypeVersionLen {
		t.Errorf("expected Versions length to be %d, got %d", expectedDocTypeVersionLen, len(docType.Versions))
	}

	expectedVersionId := 3
	expectedVersionName := "version 1"
	expectedVersionDescription := "document version 1"
	expectedVersionActiveFrom := "2015-02-13T13:15:00Z"
	expectedVersionActiveTo := "2027-03-01T00:00:00Z"
	expectedVersionVersionNumber := float32(1.1)
	expectedVersionStatus := "draft"

	if docTypeVersion1.Id != expectedVersionId {
		t.Errorf("expected version id to be %d, got %d", expectedVersionId, docTypeVersion1.Id)
	}

	if docTypeVersion1.Name != expectedVersionName {
		t.Errorf("expected version name to be %s, got %s", expectedVersionName, docTypeVersion1.Name)
	}

	if docTypeVersion1.Description != expectedVersionDescription {
		t.Errorf("expected version description to be %s, got %s", expectedVersionDescription, docTypeVersion1.Description)
	}

	if docTypeVersion1.ActiveFrom != expectedActiveFrom {
		t.Errorf("expected version active from to be %s, got %s", expectedVersionActiveFrom, docTypeVersion1.ActiveFrom)
	}

	if docTypeVersion1.ActiveTo != expectedVersionActiveTo {
		t.Errorf("expected version active to to be %s, got %s", expectedVersionActiveTo, docTypeVersion1.ActiveTo)
	}

	if docTypeVersion1.VersionNumber != expectedVersionVersionNumber {
		t.Errorf("expected version number to be %.2f, got %.2f", expectedVersionVersionNumber, docTypeVersion1.VersionNumber)
	}

	if docTypeVersion1.Status != expectedVersionStatus {
		t.Errorf("expected version status to be %s, got %s", expectedVersionStatus, docTypeVersion1.Status)
	}
}

func TestUnmarshalGetDocumentTypesResponseError(t *testing.T) {
	errorJson := `{
		"status": "Invalid",
		"name": "Step03-Duplicated Submission Validator",
		"error": {
			"propertyName": null,
			"propertyPath": null,
			"errorCode": "Error03",
			"error": "Duplicated Submission Validator",
			"errorMS": "Penduaan Sub Validator",
			"innerError": [
				{
					"propertyName": "document",
					"propertyPath": "document",
					"errorCode": "DS302",
					"error": "test error message in english",
					"errorMs": "ralat dalam bahasa melayu",
					"innerError": null
				}
			]
		}
	}`

	var resp GetDocumentTypesResponse

	err := json.Unmarshal([]byte(errorJson), &resp)
	if err != nil {
		t.Errorf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedStatus := "Invalid"

	if resp.Status != expectedStatus {
		t.Errorf("expected status to be %s, got %s", expectedStatus, resp.Status)
	}

	// expectedName := "Step03-Duplicated Submission Validator"
	// TODO: test resp.Name: takde dalam schema tapi ada dalam example JSON wtf?

	theError := resp.Error
	var expectedPropertyName string
	var expectedPropertyPath string
	expectedErrorCode := "Error03"
	expectedErrorEn := "Duplicated Submission Validator"
	expectedErrorMs := "Penduaan Sub Validator"

	if theError.PropertyName != expectedPropertyName {
		t.Errorf("expected PropertyName to be %s, got %s", expectedPropertyName, theError.PropertyName)
	}

	if theError.PropertyPath != expectedPropertyPath {
		t.Errorf("expected PropertyPath to be %s, got %s", expectedPropertyPath, theError.PropertyPath)
	}

	if theError.ErrorCode != expectedErrorCode {
		t.Errorf("expected ErrorCode to be %s, got %s", expectedErrorCode, theError.ErrorCode)
	}

	if theError.ErrorMessage != expectedErrorEn {
		t.Errorf("expected ErrorMessage to be %s, got %s", expectedErrorEn, theError.ErrorMessage)
	}

	if theError.ErrorMs != expectedErrorMs {
		t.Errorf("expected ErrorMs to be %s, got %s", expectedErrorMs, theError.ErrorMs)
	}

	innerErrors := resp.Error.InnerErrors
	expectedInnerErrorLen := 1

	if len(innerErrors) != expectedInnerErrorLen {
		t.Errorf("expected length of InnerErrors to be %d, got %d", expectedInnerErrorLen, len(innerErrors))
	}

	innerError := innerErrors[0]
	expectedInnerErrorPropertyName := "document"
	expectedInnerErrorPropertyPath := "document"
	expectedInnerErrorErrorCode := "DS302"
	expectedInnerErrorErrorEn := "test error message in english"
	expectedInnerErrorErrorMs := "ralat dalam bahasa melayu"

	if innerError.PropertyName != expectedInnerErrorPropertyName {
		t.Errorf("expected InnerError.PropertyName to be %s, got %s", expectedInnerErrorPropertyName, innerError.PropertyName)
	}

	if innerError.PropertyPath != expectedInnerErrorPropertyPath {
		t.Errorf("expected InnerError.PropertyPath to be %s, got %s", expectedInnerErrorPropertyPath, innerError.PropertyPath)
	}

	if innerError.ErrorCode != expectedInnerErrorErrorCode {
		t.Errorf("expected InnerError.ErrorCode to be %s, got %s", expectedInnerErrorErrorCode, innerError.ErrorCode)
	}

	if innerError.ErrorMessage != expectedInnerErrorErrorEn {
		t.Errorf("expected InnerError.ErrorMessage to be %s, got %s", expectedInnerErrorErrorEn, innerError.ErrorMessage)
	}

	if innerError.ErrorMs != expectedInnerErrorErrorMs {
		t.Errorf("expected InnerError.ErrorMs to be %s, got %s", expectedInnerErrorErrorMs, innerError.ErrorMs)
	}

	if innerError.InnerErrors != nil {
		t.Errorf("expected InnerError.InnerError to be %+v, got %+v", nil, innerError.InnerErrors)
	}
}

func TestUnmarshalGetDocumentTypeByIdResponseSuccess(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"invoiceTypeCode": 1,
		"description": "Invoice",
		"activeFrom": "2024-05-02T10:34:11.15Z",
		"activeTo": null,
		"documentTypeVersions": [
			{
				"id": 1,
				"name": "Version 1",
				"description": "Version 1",
				"activeFrom": "2024-04-26T17:47:49.2Z",
				"activeTo": null,
				"versionNumber": 1.0,
				"status": "Published"
			},
			{
				"id": 2,
				"name": "Version 2",
				"description": "Version 2",
				"activeFrom": "2024-06-21T15:13:17.29Z",
				"activeTo": null,
				"versionNumber": 1.1,
				"status": "Published"
			}
		],
		"workflowParameters": [
			{
				"id": 1,
				"parameter": "Submission",
				"value": 80,
				"activeFrom": "2024-03-14T04:59:24.076Z",
				"activeTo": null
			}
		]
	}`

	var resp GetDocumentTypeByIdResponse

	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedId := 1
	expectedInvoiceTypeCode := 1
	expectedDescription := "Invoice"
	expectedActiveFrom := "2024-05-02T10:34:11.15Z"
	var expectedActiveTo string

	if resp.Id != expectedId {
		t.Errorf("expected id to be %d, got %d", expectedId, resp.Id)
	}

	if resp.InvoiceTypeCode != expectedInvoiceTypeCode {
		t.Errorf("expected invoiceTypeCode to be %d, got %d", expectedInvoiceTypeCode, resp.InvoiceTypeCode)
	}

	if resp.Description != expectedDescription {
		t.Errorf("expected description to be %s, got %s", expectedDescription, resp.Description)
	}

	if resp.ActiveFrom != expectedActiveFrom {
		t.Errorf("expected activeFrom to be %s, got %s", expectedActiveFrom, resp.ActiveFrom)
	}

	if resp.ActiveTo != expectedActiveTo {
		t.Errorf("expected activeTo to be %s, got %s", expectedActiveTo, resp.ActiveTo)
	}

	versions := resp.Versions
	expectedVersionCount := 2

	if len(versions) != expectedVersionCount {
		t.Errorf("expected len(documentTypeVersions) to be %d, got %d", expectedVersionCount, len(versions))
	}

	// TODO: documentTypeVersions assertions

	workflowParams := resp.WorkflowParameters
	expectedWorkflowParamsCount := 1

	if len(workflowParams) != expectedWorkflowParamsCount {
		t.Errorf("expected len(workflowParameters) to be %d, got %d", expectedWorkflowParamsCount, len(workflowParams))
	}

	wfp := workflowParams[0]
	expectedWfpId := 1
	expectedWfpParameter := "Submission"
	expectedWfpValue := 80
	expectedWfpActiveFrom := "2024-03-14T04:59:24.076Z"
	var expectedWfpActiveTo string

	if wfp.Id != expectedWfpId {
		t.Errorf("expected workflowParam.id to be %d, got %d", expectedWfpId, wfp.Id)
	}

	if wfp.Parameter != expectedWfpParameter {
		t.Errorf("expected workflowParam.parameter to be %s, got %s", expectedWfpParameter, wfp.Parameter)
	}

	if wfp.Value != expectedWfpValue {
		t.Errorf("expected workflowParam.value to be %d, got %d", expectedWfpValue, wfp.Value)
	}

	if wfp.ActiveFrom != expectedWfpActiveFrom {
		t.Errorf("expected workflowParam.activeFrom to be %s, got %s", expectedWfpActiveFrom, wfp.ActiveFrom)
	}

	if wfp.ActiveTo != expectedWfpActiveTo {
		t.Errorf("expected workflowParam.activeTo to be %s, got %s", expectedWfpActiveTo, wfp.ActiveTo)
	}
}

func TestUnmarshalGetDocumentTypeVersionResponse(t *testing.T) {
	jsonStr := `{
		"invoiceTypeCode":4,
		"name":"1.0",
		"description":"Credit note version 1.0",
		"versionNumber":1.0,
		"status":"published",
		"activeFrom":"2015-02-13T13:15:00Z",
		"activeTo":"2027-03-01T00:00:00Z",
		"jsonschema":"dummyValue1",
		"xmlschema":"dummyValue2"
	}`

	var resp GetDocumentTypeVersionResponse

	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("unexpected error when unmarshalling JSON: %s", err)
	}

	expectedInvoiceTypeCode := 4
	expectedName := "1.0"
	expectedDescription := "Credit note version 1.0"
	expectedVersionNumber := float32(1.0)
	expectedStatus := "published"
	expectedActiveFrom := "2015-02-13T13:15:00Z"
	expectedActiveTo := "2027-03-01T00:00:00Z"
	expectedJsonSchema := "dummyValue1"
	expectedXmlSchema := "dummyValue2"

	if resp.InvoiceTypeCode != expectedInvoiceTypeCode {
		t.Errorf("expected invoiceTypeCode to be %d, got %d", expectedInvoiceTypeCode, resp.InvoiceTypeCode)
	}

	if resp.Name != expectedName {
		t.Errorf("expected name to be %s, got %s", expectedName, resp.Name)
	}

	if resp.Description != expectedDescription {
		t.Errorf("expected description to be %s, got %s", expectedDescription, resp.Description)
	}

	if resp.VersionNumber != expectedVersionNumber {
		t.Errorf("expected versionNumber to be %.2f, got %.2f", expectedVersionNumber, resp.VersionNumber)
	}

	if resp.Status != expectedStatus {
		t.Errorf("expected status to be %s, got %s", expectedStatus, resp.Status)
	}

	if resp.ActiveFrom != expectedActiveFrom {
		t.Errorf("expected activeFrom to be %s, got %s", expectedActiveFrom, resp.ActiveFrom)
	}

	if resp.ActiveTo != expectedActiveTo {
		t.Errorf("expected activeTo to be %s, got %s", expectedActiveTo, resp.ActiveTo)
	}

	if resp.JsonSchema != expectedJsonSchema {
		t.Errorf("expected jsonSchema to be %s, got %s", expectedJsonSchema, resp.JsonSchema)
	}

	if resp.XmlSchema != expectedXmlSchema {
		t.Errorf("expected xmlSchema to be %s, got %s", expectedXmlSchema, resp.XmlSchema)
	}
}
