package common

const (
	SANDBOX_API_BASE_URL      = "https://preprod-api.myinvois.hasil.gov.my"
	SANDBOX_IDENTITY_BASE_URL = "https://preprod-api.myinvois.hasil.gov.my"
)

// Reference: https://sdk.myinvois.hasil.gov.my/standard-error-response/

type ErrResponse struct {
	PropertyName string        `json:"propertyName"`
	PropertyPath string        `json:"propertyPath"`
	ErrorCode    string        `json:"errorCode"`
	ErrorMessage string        `json:"error"`
	ErrorMs      string        `json:"errorMS"`
	Target       string        `json:"target"`
	InnerErrors  []ErrResponse `json:"innerError"`
}

type StandardErrResponse struct {
	Status string      `json:"status"`
	Error  ErrResponse `json:"error"`
}
