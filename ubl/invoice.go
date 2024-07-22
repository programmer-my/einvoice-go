package ubl

import (
	"encoding/xml"
	"fmt"

	"github.com/Rhymond/go-money"
)

// Invoice type code
// https://docs.peppol.eu/poacc/billing/3.0/codelist/UNCL1001-inv/

// EndpointID.SchemeID
// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/syntax/cac-AccountingSupplierParty/cac-Party/cbc-EndpointID/schemeID/
// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/codelist/eas/
// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/codelist/ICD/

// Signature generation
// https://sdk.myinvois.hasil.gov.my/signature/

func CurrencyMarshaler(m *money.Money, e *xml.Encoder, s xml.StartElement) error {
	toEncode := struct {
		Amount     string `xml:",innerxml"`
		CurrencyID string `xml:"currencyID,attr"`
	}{
		Amount:     fmt.Sprintf("%.2f", m.AsMajorUnits()),
		CurrencyID: m.Currency().Code,
	}

	return e.EncodeElement(toEncode, s)
}

func NewInvoice() *UBL_Invoice {
	inv := UBL_Invoice{
		CACEnv: "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		EXTEnv: "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2",
		CBCEnv: "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		QDTEnv: "urn:oasis:names:specification:ubl:schema:xsd:QualifiedDataTypes-2",
		UDTEnv: "urn:oasis:names:specification:ubl:schema:xsd:UnqualifiedDataTypes-2",
		CNEnv:  "urn:oasis:names:specification:ubl:schema:xsd:CreditNote-2",
		UBLEnv: "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
	}
	return &inv
}

type UBL_Invoice struct {
	/** XML Namespace Setup Start **/
	XMLName xml.Name `xml:"Invoice"`
	CACEnv  string   `xml:"xmlns:cac,attr"`
	EXTEnv  string   `xml:"xmlns:ext,attr"`
	CBCEnv  string   `xml:"xmlns:cbc,attr"`
	QDTEnv  string   `xml:"xmlns:qdt,attr"`
	UDTEnv  string   `xml:"xmlns:udt,attr"`
	CNEnv   string   `xml:"xmlns:cn,attr"`
	UBLEnv  string   `xml:"xmlns:ubl,attr"`
	/** XML Namespace Setup End **/

	Signature                   string                           `xml:"cbc:Signature"`
	CustomizationID             string                           `xml:"cbc:CustomizationID"`                       // [1..1] 	Specification identifier
	ProfileID                   string                           `xml:"cbc:ProfileID"`                             // [1..1] 	Business process type
	ID                          string                           `xml:"cbc:ID"`                                    // [1..1] 	Invoice number
	UUID                        string                           `xml:"cbc:UUID"`                                  // [1..1] 	IRBM Unique Identifier Number
	IssueDate                   string                           `xml:"cbc:IssueDate"`                             // [1..1] 	Invoice issue date
	IssueTime                   *string                          `xml:"cbc:IssueTime,omitempty"`                   // [0..1] 	Invoice issue time
	DueDate                     *string                          `xml:"cbc:DueDate,omitempty"`                     // [0..1] 	Payment due date
	InvoiceTypeCode             string                           `xml:"cbc:InvoiceTypeCode"`                       // [1..1] 	Invoice type code
	Note                        *string                          `xml:"cbc:Note,omitempty"`                        // [0..1] 	Invoice note
	TaxPointDate                string                           `xml:"cbc:TaxPointDate"`                          // [0..0] 	TAX point date ???? 0..0?
	Currency                    money.Currency                   `xml:"-"`                                         // for copying around in invoice line, legal monetary values, etc. not for serialization
	DocumentCurrencyCode        string                           `xml:"cbc:DocumentCurrencyCode"`                  // [1..1] 	Invoice currency code
	TaxCurrencyCode             *money.Currency                  `xml:"cbc:TaxCurrencyCode,omitempty"`             // [0..1] 	Tax accounting currency
	AccountingCost              *string                          `xml:"cbc:AccountingCost,omitempty"`              // [0..1] 	Buyer accounting reference
	BuyerReference              *string                          `xml:"cbc:BuyerReference,omitempty"`              // [0..1] 	Buyer reference
	InvoicePeriod               *CAC_InvoicePeriod               `xml:"cac:InvoicePeriod,omitempty"`               // [0..1] 	INVOICING PERIOD
	OrderReference              *CAC_OrderReference              `xml:"cac:OrderReference,omitempty"`              // [0..1] 	Order and sales order reference
	BillingReference            []CAC_BillingReference           `xml:"cac:BillingReference,omitempty"`            // [0..n] 	PRECEDING INVOICE REFERENCE
	DespatchDocumentReference   *CAC_DespatchDocumentReference   `xml:"cac:DespatchDocumentReference,omitempty"`   // [0..1] 	Despatch advice reference
	ReceiptDocumentReference    *CAC_ReceiptDocumentReference    `xml:"cac:ReceiptDocumentReference,omitempty"`    // [0..1] 	Receipt advice reference
	OriginatorDocumentReference *CAC_OriginatorDocumentReference `xml:"cac:OriginatorDocumentReference,omitempty"` // [0..1] 	Tender or Lot reference
	ContractDocumentReference   *CAC_ContractDocumentReference   `xml:"cac:ContractDocumentReference,omitempty"`   // [0..1] 	Contract reference
	AdditionalDocumentReference *CAC_AdditionalDocumentReference `xml:"cac:AdditionalDocumentReference,omitempty"` // [0..1] [cbc:DocumentTypeCode = "938"] [cbc:DocumentTypeCode = "914"] [cbc:DocumentTypeCode = 130], [0..n] if [cbc:DocumentTypeCode != 130]
	ProjectReference            *CAC_ProjectReference            `xml:"cac:ProjectReference,omitempty"`            // [0..1] 	Project reference
	AccountingSupplierParty     CAC_AccountingSupplierParty      `xml:"cac:AccountingSupplierParty"`               // [1..1] 	SELLER
	AccountingCustomerParty     CAC_AccountingCustomerParty      `xml:"cac:AccountingCustomerParty"`               // [1..1] 	BUYER
	PayeeParty                  *CAC_PayeeParty                  `xml:"cac:PayeeParty,omitempty"`                  // [0..1] 	Payee
	TaxRepresentativeParty      *CAC_TaxRepresentativeParty      `xml:"cac:TaxRepresentativeParty,omitempty"`      // [0..1] 	SELLER INVOICING REPRESENTATIVE PARTY
	Delivery                    *CAC_Delivery                    `xml:"cac:Delivery,omitempty"`                    // [0..1] 	DELIVERY INFORMATION
	PaymentMeans                []CAC_PaymentMeans               `xml:"cac:PaymentMeans,omitempty"`                // [0..n] 	PAYMENT INSTRUCTIONS
	PaymentTerms                []CAC_PaymentTerms               `xml:"cac:PaymentTerms,omitempty"`                // [0..n] 	INVOICE TERMS
	PrepaidPayment              []CAC_PrepaidPayment             `xml:"cac:PrepaidPayment,omitempty"`              // [0..n] 	PAID AMOUNTS
	AllowanceCharge             []CAC_AllowanceCharge            `xml:"cac:AllowanceCharge,omitempty"`             // [0..n] [cbc:ChargeIndicator = false] [0..n] [cbc:ChargeIndicator = true]
	TaxExchangeRate             *CAC_TaxExchangeRate             `xml:"cac:TaxExchangeRate,omitempty"`             // [0..1] 	TAX EXCHANGE RATE
	TaxTotal                    CAC_TaxTotal                     `xml:"cac:TaxTotal"`                              // [1..1] [cac:TaxTotal/TaxAmount/@currency = cbc:DocumentCurrencyCode] | [0..1] [cac:TaxTotal/TaxAmount/@currency = cbc:TaxCurrencyCode]
	LegalMonetaryTotal          CAC_LegalMonetaryTotal           `xml:"cac:LegalMonetaryTotal"`                    // [1..1] 	DOCUMENT TOTALS
	InvoiceLine                 []CAC_InvoiceLine                `xml:"cac:InvoiceLine"`                           // [1..n] 	INVOICE LINE
}

type CBC_DocumentCurrencyCode struct {
	Code money.Currency
}

func (d *CBC_DocumentCurrencyCode) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	return e.EncodeToken(d.Code.Code)
}

func (i *UBL_Invoice) AddItem() *CAC_InvoiceLine {
	line := CAC_InvoiceLine{}
	i.InvoiceLine = append(i.InvoiceLine, line)
	return &line
}

type CAC_InvoicePeriod struct {
	XMLName     xml.Name `xml:"cac:InvoicePeriod"`
	StartDate   *string  `xml:"cbc:StartDate"`   // [0..1] Invoicing period start date
	EndDate     *string  `xml:"cbc:EndDate"`     // [0..1] Invoicing period end date
	Description *string  `xml:"cbc:Description"` // [0..1] Frequency of billing
}

type CAC_OrderReference struct {
	XMLName      xml.Name `xml:"cac:OrderReference"`
	ID           string   `xml:"cbc:ID"`                     // [1..1] Purchase order reference
	SalesOrderID *string  `xml:"cbc:SalesOrderID,omitempty"` // [0..1] Sales order reference
}

type CAC_BillingReference struct {
	XMLName                  xml.Name                     `xml:"cac:BillingReference"`
	InvoiceDocumentReference CAC_InvoiceDocumentReference `xml:"cac:InvoiceDocumentReference"` // [1..1] Invoice document reference
}

type CAC_InvoiceDocumentReference struct {
	XMLName   xml.Name `xml:"cac:InvoiceDocumentReference"`
	ID        string   `xml:"cbc:ID"`                  // [1..1] Preceding Invoice reference - The identification of an Invoice that was previously sent by the Seller.
	IssueDate *string  `xml:"cbc:IssueDate,omitempty"` // [0..1] Preceding Invoice issue date - The date when the Preceding Invoice was issued.
}

type CAC_DespatchDocumentReference struct {
	XMLName xml.Name `xml:"cac:DespatchDocumentReference"`
	ID      string   `xml:"cbc:ID"` // [1..1] Despatch advice reference
}

type CAC_ReceiptDocumentReference struct {
	XMLName xml.Name `xml:"cac:ReceiptDocumentReference"`
	ID      string   `xml:"cbc:ID"` // [1..1] Receiving advice reference
}

type CAC_OriginatorDocumentReference struct {
	XMLName xml.Name `xml:"cac:OriginatorDocumentReference"`
	ID      string   `xml:"cbc:ID"` // [1..1] Tender or lot reference
}

type CAC_ContractDocumentReference struct {
	XMLName xml.Name `xml:"cac:ContractDocumentReference"`
	ID      string   `xml:"cbc:ID"` // [1..1] Contract reference
}

type CAC_AdditionalDocumentReference struct {
	XMLName              xml.Name                 `xml:"cac:AdditionalDocumentReference"`
	Attachment           CAC_Attachment           `xml:"cac:Attachment"`           // [1..1] Attachment
	ResultOfVerification CAC_ResultOfVerification `xml:"cac:ResultOfVerification"` // [1..1] Tax validation date and time
}

type CAC_Attachment struct {
	XMLName           xml.Name              `xml:"cac:Attachment"`
	ExternalReference CAC_ExternalReference `xml:"cac:ExternalReference"` // [1..1] External reference
}

type CAC_ExternalReference struct {
	XMLName xml.Name `xml:"cac:ExternalReference"`
	URI     string   `xml:"cbc:URI"` // [1..1] Validation Link
}

type CAC_ResultOfVerification struct {
	XMLName        xml.Name `xml:"cac:ResultOfVerification"`
	ValidationDate string   `xml:"cbc:ValidationDate"` // [1..1] Date of Validation
	ValidationTime string   `xml:"cbc:ValidationTime"` // [1..1] Time of Validation
}

type CAC_ProjectReference struct {
	XMLName xml.Name `xml:"cac:ProjectReference"`
	ID      string   `xml:"cbc:ID"` // [1..1] Project reference - The identification of the project the invoice refers to
}

type CAC_AccountingSupplierParty struct {
	XMLName xml.Name  `xml:"cac:AccountingSupplierParty"`
	Party   CAC_Party `xml:"cac:Party"` // [1..1]
}

type CAC_Party struct {
	XMLName                    xml.Name                  `xml:"cac:Party"`
	EndpointID                 CBC_EndpointID            `xml:"cbc:EndpointID"`                 // [1..1] 	Seller electronic address - Identifies the Seller’s electronic address to which the application level response to the invoice may be delivered.
	IndustryClassificationCode *string                   `xml:"cbc:IndustryClassificationCode"` // [0..1] 	Suppliers Malaysia Standard Industrial Classification (MSIC) Code - 5-digit numeric code that represent the business nature and activity.
	PartyIdentification        []CAC_PartyIdentification `xml:"cac:PartyIdentification"`        // [0..n]
	PartyName                  *CAC_PartyName            `xml:"cac:PartyName"`                  // [0..1] 	Party name
	PostalAddress              CAC_PostalAddress         `xml:"cac:PostalAddress"`              // [1..1] 	SELLER POSTAL ADDRESS - A group of business terms providing information about the address of the Seller.
	PartyTaxScheme             *CAC_PartyTaxScheme       `xml:"cac:PartyTaxScheme"`             // [0..1] [cac:TaxScheme = "AAL"] 	Suppliers Tourism Tax Registration [cac:TaxScheme = "VAT"] 	PARTY TAX [cac:TaxScheme != "VAT"] 	PARTY TAX
	PartyLegalEntity           CAC_PartyLegalEntity      `xml:"cac:PartyLegalEntity"`           // [1..1] Party legal entity
	Contact                    *CAC_Contact              `xml:"cac:Contact"`                    // [0..1] SELLER CONTACT - A group of business terms providing contact information about the Seller.
}

// TODO
// The code is based on MSIC code list
// https://sdk.myinvois.hasil.gov.my/codes/msic-codes/
// TODO 2: conflicting schema
// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/syntax/cac-AccountingSupplierParty/cac-Party/cbc-IndustryClassificationCode/
// - contains only listID attribute
// - IndustryClassificationCode is not mandatory
// https://sdk.myinvois.hasil.gov.my/documents/invoice-v1-1/#supplier
// - contains name attribute & mandatory
// type CAC_Party_IndustryClassificationCode struct {
// 	Code string `xml:",innerxml"` // [1..1]
// 	Name string `xml:"name,attr"` // [1..1]
// }

type CBC_EndpointID struct {
	Value    string `xml:",innerxml"`     // [1..1]
	SchemeID string `xml:"schemeID,attr"` // [1..1]
}

type CAC_PartyIdentification struct {
	XMLName xml.Name                   `xml:"cac:PartyIdentification"`
	ID      CAC_PartyIdentification_ID `xml:"cbc:ID"` // [1..1] Seller identifier
}

// https://sdk.myinvois.hasil.gov.my/documents/invoice-v1-1/#supplier
// [@schemeID=’NRIC’]
// [@schemeID=’BRN’]
// [@schemeID=’PASSPORT’]
// [@schemeID=’ARMY’]
// TODO: spec conflict
// link tells us the above schemeID is accepted
// but schematron says otherways (schemeID "0230" for malaysia e-invoice)
type CAC_PartyIdentification_ID struct {
	Value    string `xml:",innerxml"`
	SchemeID string `xml:"schemeID,attr"`
}

type CAC_PartyName struct {
	XMLName xml.Name `xml:"cac:PartyName"`
	Name    *string  `xml:"cbc:Name"` // [0..1] Seller trading name
}

type CAC_Address struct {
	XMLName              xml.Name         `xml:"cac:Address"`
	StreetName           *string          `xml:"cbc:StreetName"`           // [0..1] Seller address line 1 - The main address line in an address.
	AdditionalStreetName *string          `xml:"cbc:AdditionalStreetName"` // [0..1] Seller address line 2 - An additional address line in an address that can be used to give further details supplementing the main line.
	CityName             *string          `xml:"cbc:CityName"`             // [0..1] Seller city - The common name of the city, town or village, where the Seller address is located.
	PostalZone           *string          `xml:"cbc:PostalZone"`           // [0..1] Seller post code - The identifier for an addressable group of properties according to the relevant postal service.
	CountrySubentity     *string          `xml:"cbc:CountrySubentity"`     // [0..1] Seller country subdivision - The subdivision of a country. Such as a region, a county, a state, a province etc..
	AddressLine          *CAC_AddressLine `xml:"cac:AddressLine"`          // [0..1] Address line
	Country              CAC_Country      `xml:"cac:Country"`              // [1..1] Country
}

type CAC_PostalAddress struct {
	XMLName xml.Name `xml:"cac:PostalAddress"`
	// for some fucking reason struct embedding does not work!!!!!
	// CAC_Address
	StreetName           *string          `xml:"cbc:StreetName"`           // [0..1] Seller address line 1 - The main address line in an address.
	AdditionalStreetName *string          `xml:"cbc:AdditionalStreetName"` // [0..1] Seller address line 2 - An additional address line in an address that can be used to give further details supplementing the main line.
	CityName             *string          `xml:"cbc:CityName"`             // [0..1] Seller city - The common name of the city, town or village, where the Seller address is located.
	PostalZone           *string          `xml:"cbc:PostalZone"`           // [0..1] Seller post code - The identifier for an addressable group of properties according to the relevant postal service.
	CountrySubentity     *string          `xml:"cbc:CountrySubentity"`     // [0..1] Seller country subdivision - The subdivision of a country. Such as a region, a county, a state, a province etc..
	AddressLine          *CAC_AddressLine `xml:"cac:AddressLine"`          // [0..1] Address line
	Country              CAC_Country      `xml:"cac:Country"`              // [1..1] Country
}

type CAC_AddressLine struct {
	XMLName xml.Name `xml:"cac:AddressLine"`
	Line    string   `xml:"cbc:Line"` //	[1..1] Seller address line 3
}

type CAC_Country struct {
	XMLName            xml.Name `xml:"cac:Country"`
	IdentificationCode string   `xml:"cbc:IdentificationCode"` // [1..1] Seller country code ISO3166-1
}

type CAC_PartyTaxScheme struct {
	XMLName   xml.Name      `xml:"cac:PartyTaxScheme"`
	CompanyID *string       `xml:"cbc:CompanyID"` // [0..1]	Suppliers Tourism Tax Registration Number
	TaxScheme CAC_TaxScheme `xml:"cac:TaxScheme"` // [1..1]	Tax scheme
}

type CAC_TaxScheme struct {
	XMLName xml.Name `xml:"cac:TaxScheme"`
	ID      string   `xml:"cbc:ID"` // [1..1]
}

type CAC_PartyLegalEntity struct {
	XMLName          xml.Name `xml:"cac:PartyLegalEntity"`
	RegistrationName string   `xml:"cbc:RegistrationName"` // [1..1] The full formal name by which the Seller is registered in the national registry of legal entities or as a Taxable person or otherwise trades as a person or persons.
	CompanyID        *string  `xml:"cbc:CompanyID"`        // [0..1] An identifier issued by an official registrar that identifies the Seller as a legal entity or person.
	CompanyLegalForm *string  `xml:"cbc:CompanyLegalForm"` // [0..1] Additional legal information relevant for the Seller.
}

type CAC_Contact struct {
	XMLName        xml.Name `xml:"cac:Contact"`
	Name           *string  `xml:"cbc:Name"`           // [0..1] 	Seller contact point
	Telephone      *string  `xml:"cbc:Telephone"`      // [0..1] 	Seller contact telephone number
	ElectronicMail *string  `xml:"cbc:ElectronicMail"` // [0..1] 	Seller contact email address
}

type CAC_AccountingCustomerParty struct {
	XMLName xml.Name  `xml:"cac:AccountingCustomerParty"`
	Party   CAC_Party `xml:"cac:Party"` // [1..1]
}

type CAC_PayeeParty struct {
	XMLName             xml.Name                 `xml:"cac:PayeeParty"`
	PartyIdentification *CAC_PartyIdentification `xml:"cac:PartyIdentification"` // [0..1] 	PARTY IDENTIFICATION
	PartyName           CAC_PartyName            `xml:"cac:PartyName"`           // [1..1] 	Party name
	PartyLegalEntity    *CAC_PartyLegalEntity    `xml:"cac:PartyLegalEntity"`    // [0..1] 	Party legal entity
}

type CAC_TaxRepresentativeParty struct {
	XMLName        xml.Name           `xml:"cac:TaxRepresentativeParty"`
	PartyName      CAC_PartyName      `xml:"cac:PartyName"`      //	[1..1] Party name
	PartyTaxScheme CAC_PartyTaxScheme `xml:"cac:PartyTaxScheme"` //	[1..1] PARTY TAX
}

type CAC_Delivery struct {
	XMLName            xml.Name              `xml:"cac:Delivery"`
	ActualDeliveryDate *string               `xml:"cbc:ActualDeliveryDate"` // [0..1] 	the date on which the supply of goods or services was made or completed.
	DeliveryLocation   *CAC_DeliveryLocation `xml:"cac:DeliveryLocation"`   // [0..1]
	DeliveryParty      *CAC_DeliveryParty    `xml:"cac:DeliveryParty"`      // [0..1] 	DELIVER PARTY
	Shipment           *CAC_Shipment         `xml:"cac:Shipment"`           // [0..1] 	SHIPMENT INFORMATION
}

type CAC_DeliveryLocation struct {
	XMLName xml.Name     `xml:"cac:DeliveryLocation"`
	ID      *string      `xml:"cbc:ID"`      // [0..1] 	Deliver to location identifier
	Address *CAC_Address `xml:"cac:Address"` // [0..1] 	DELIVER TO ADDRESS
}

type CAC_DeliveryParty struct {
	XMLName          xml.Name              `xml:"cac:DeliveryParty"`
	PartyName        CAC_PartyName         `xml:"cac:PartyName"`        // [1..1] 	PARTY NAME
	PartyLegalEntity *CAC_PartyLegalEntity `xml:"cac:PartyLegalEntity"` // [0..1]
}

type CAC_Shipment struct {
	XMLName     xml.Name        `xml:"cac:Shipment"`
	Consignment CAC_Consignment `xml:"cac:Consignment"` // [1..1] CONSIGNMENT INFORMATION
}

type CAC_Consignment struct {
	XMLName       xml.Name          `xml:"cac:Consignment"`
	DeliveryTerms CAC_DeliveryTerms `xml:"cac:DeliveryTerms"` // [1..1]
}

type CAC_DeliveryTerms struct {
	XMLName xml.Name `xml:"cac:DeliveryTerms"`
	ID      *string  `xml:"cbc:ID"` // [0..1]	Incoterms
}

type CAC_PaymentMeans struct {
	XMLName               xml.Name                   `xml:"cac:PaymentMeans"`
	ID                    *string                    `xml:"cbc:ID"`                    // [0..1] Payment Instructions ID - An identifier for the payment instructions.
	PaymentMeansCode      string                     `xml:"cbc:PaymentMeansCode"`      // [1..1] Payment means type code - The means, expressed as code, for how a payment is expected to be or has been settled.
	PaymentID             []string                   `xml:"cbc:PaymentID"`             // [0..n] Remittance information - A textual value used for payment routing or to establish a link between the payment and the Invoice.
	CardAccount           *CAC_CardAccount           `xml:"cac:CardAccount"`           // [0..1] PAYMENT CARD INFORMATION - A group of business terms providing information about card used for payment contemporaneous with invoice issuance.
	PayeeFinancialAccount *CAC_PayeeFinancialAccount `xml:"cac:PayeeFinancialAccount"` // [0..1] CREDIT TRANSFER - A group of business terms to specify credit transfer payments.
	PaymentMandate        *CAC_PaymentMandate        `xml:"cac:PaymentMandate"`        // [0..1] DIRECT DEBIT - A group of business terms to specify a direct debit.
}

type CAC_CardAccount struct {
	XMLName                xml.Name `xml:"cac:CardAccount"`
	PrimaryAccountNumberID string   `xml:"cbc:PrimaryAccountNumberID"` // [1..1] Payment card primary account number - The Primary Account Number (PAN) of the card used for payment.
	NetworkID              string   `xml:"cbc:NetworkID"`              // [1..1]  - Syntax required element not mapped to a business term. Use value NA
	HolderName             *string  `xml:"cbc:HolderName"`             // [0..1] Payment card holder name - The name of the payment card holder.
}

type CAC_PayeeFinancialAccount struct {
	XMLName                    xml.Name `xml:"cac:PayeeFinancialAccount"`
	ID                         string   `xml:"cbc:ID"`                         // [1..1] A unique identifier of the financial payment account, at a payment service provider, to which payment should be made.
	Name                       *string  `xml:"cbc:Name"`                       // [0..1] The name of the payment account, at a payment service provider, to which payment should be made.
	FinancialInstitutionBranch *string  `xml:"cac:FinancialInstitutionBranch"` // [0..1] FINANCIAL INSTITUTION BRANCH
}

type CAC_PaymentMandate struct {
	XMLName               xml.Name                   `xml:"cac:PaymentMandate"`
	ID                    *string                    `xml:"cbc:ID"`                    // [0..1] Unique identifier assigned by the Payee for referencing the direct debit mandate.
	PayerFinancialAccount *CAC_PayerFinancialAccount `xml:"cac:PayerFinancialAccount"` // [0..1] PAYER FINANCIAL ACCOUNT
}

type CAC_PayerFinancialAccount struct {
	XMLName xml.Name `xml:"cac:PayerFinancialAccount"`
	ID      string   `xml:"cbc:ID"` // [1..1]  Debited account identifier
}

type CAC_PaymentTerms struct {
	XMLName xml.Name `xml:"cac:PaymentTerms"`
	Note    *string  `xml:"cbc:note"` // [0..1] Payment terms - A textual description of the payment terms that apply to the amount due for payment (Including description of possible penalties).
}

type CAC_PrepaidPayment struct {
	XMLName      xml.Name `xml:"cac:PrepaidPayment"`
	ID           *string  `xml:"cbc:ID"`           // [0..1] Payment identifier - An identifier that references the payment, such as bank transfer identifier.
	ReceivedDate *string  `xml:"cbc:ReceivedDate"` // [0..1] The date when the paid amount is debited to the invoice. - The date when the prepaid amount was received by the seller.
}

type CAC_AllowanceCharge struct {
	XMLName                   xml.Name         `xml:"cac:AllowanceCharge"`
	ChargeIndicator           bool             `xml:"cbc:ChargeIndicator"`           // [1..1] Syntax binding qualifier - Use “false” when informing about Allowances. Use “true” when informing about Charges.
	AllowanceChargeReasonCode *string          `xml:"cbc:AllowanceChargeReasonCode"` // [0..1] Document level allowance reason code - The reason for the document level allowance, expressed as a code.
	AllowanceChargeReason     *string          `xml:"cbc:AllowanceChargeReason"`     // [0..1] Document level allowance reason - The reason for the document level allowance, expressed as text.
	MultiplierFactorNumeric   *string          `xml:"cbc:MultiplierFactorNumeric"`   // [0..1] Document level allowance percentage - The percentage that may be used, in conjunction with the document level allowance base amount, to calculate the document level allowance amount.
	Amount                    CBC_Amount       `xml:"cbc:Amount"`                    // [1..1] Document level allowance amount - The amount of an allowance, without TAX.
	BaseAmount                *CBC_Amount      `xml:"cbc:BaseAmount"`                // [0..1] Document level allowance base amount - The base amount that may be used, in conjunction with the document level allowance percentage, to calculate the document level allowance amount.
	TaxCategory               *CAC_TaxCategory `xml:"cac:TaxCategory"`               // [0..1] TAX CATEGORY
}

type CBC_Amount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_Amount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)

	return CurrencyMarshaler(m, e, s)
}

// Valid values for ID: T E O (aligned-ibrp-cl-01-my)

type CAC_TaxCategory struct {
	XMLName            xml.Name      `xml:"cac:TaxCategory"`
	ID                 string        `xml:"cbc:ID"`                 // [1..1] Document level charge TAX category code - A coded identification of what TAX category applies to the document level charge.
	Percent            string        `xml:"cbc:Percent,omitempty"`  // [0..1] Document level charge TAX rate - The TAX rate, represented as percentage that applies to the document level charge.
	TaxExemptionReason *string       `xml:"cbc:TaxExemptionReason"` // [0..1] Document level charge TAX exemption reason text - A textual statement of the reason why the document level charge amount is exempted from TAX or why no TAX is being charged
	TaxScheme          CAC_TaxScheme `xml:"cac:TaxScheme"`          // [1..1] TAX SCHEME
}

// TODO: conflict in spec
// Invoice>Core section:
// - ubl:Invoice / cac:TaxExchangeRate / cbc:SourceCurrencyCode = DocumentCurrencyCode
// - ubl:Invoice / cac:TaxExchangeRate / cbc:TargetCurrencyCode = MYR
// - cac:TaxExchangeRate has child element cbc:SourceCurrencyCode and cbc:TargetCurrencyCode
// - this link does not mention these element https://docs.peppol.eu/poac/my/pint-my/trn-invoice/syntax/cac-TaxExchangeRate/
type CAC_TaxExchangeRate struct {
	XMLName         xml.Name `xml:"cac:TaxExchangeRate"`
	CalculationRate string   `xml:"cbc:CalculationRate"` // [1..1] Currency Exchange Rate - Rate at which non-Malaysian currency will be multiplied to convert it into Malaysian Ringgit. Applicable where the billing amount is in foreign currency.
}

type CAC_TaxTotal struct {
	XMLName     xml.Name          `xml:"cac:TaxTotal"`
	TaxAmount   CBC_TaxAmount     `xml:"cbc:TaxAmount"`                    // [1..1] Invoice total TAX amount - The total TAX amount for the Invoice.
	TaxSubtotal []CAC_TaxSubtotal `xml:"cac:TaxSubtotal" validate:"min=1"` // [1..n] TAX BREAKDOWN - A group of business terms providing information about TAX breakdown by different categories, rates and exemption reasons
}

type CBC_TaxAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_TaxAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CAC_TaxSubtotal struct {
	XMLName       xml.Name          `xml:"cac:TaxSubtotal"`
	TaxableAmount CBC_TaxableAmount `xml:"cbc:TaxableAmount"` // [1..1] TAX category taxable amount. - Sum of all taxable amounts subject to a specific TAX category code and TAX category rate (if the TAX category rate is applicable).
	TaxAmount     CBC_TaxAmount     `xml:"cbc:TaxAmount"`     // [1..1] TAX category tax amount. - The total TAX amount for a given TAX category.
	TaxCategory   CAC_TaxCategory   `xml:"cac:TaxCategory"`   // [1..1]
}

type CBC_TaxableAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_TaxableAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CAC_LegalMonetaryTotal struct {
	XMLName               xml.Name                   `xml:"cac:LegalMonetaryTotal"`
	LineExtensionAmount   CBC_LineExtensionAmount    `xml:"cbc:LineExtensionAmount"`   // [1..1] Sum of Invoice line net amount - Sum of all Invoice line net amounts in the Invoice.
	TaxExclusiveAmount    CBC_TaxExclusiveAmount     `xml:"cbc:TaxExclusiveAmount"`    // [1..1] Invoice total amount without TAX - The total amount of the Invoice without TAX.
	TaxInclusiveAmount    CBC_TaxInclusiveAmount     `xml:"cbc:TaxInclusiveAmount"`    // [1..1] Invoice total amount with TAX - The total amount of the Invoice with tax.
	AllowanceTotalAmount  *CBC_AllowanceTotalAmount  `xml:"cbc:AllowanceTotalAmount"`  // [0..1] Sum of allowances on document level - Sum of all allowances on document level in the Invoice.
	ChargeTotalAmount     *CBC_ChargeTotalAmount     `xml:"cbc:ChargeTotalAmount"`     // [0..1] Sum of charges on document level - Sum of all charges on document level in the Invoice.
	PrepaidAmount         *CBC_PrepaidAmount         `xml:"cbc:PrepaidAmount"`         // [0..1] Paid amount - The sum of amounts which have been paid in advance.
	PayableRoundingAmount *CBC_PayableRoundingAmount `xml:"cbc:PayableRoundingAmount"` // [0..1] Rounding amount - Syntax required attribute, value must equal invoice document currency (ibt-005)
	PayableAmount         CBC_PayableAmount          `xml:"cbc:PayableAmount"`         // [1..1] Amount due for payment - The outstanding amount that is requested to be paid.
}

type CBC_LineExtensionAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_LineExtensionAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_TaxExclusiveAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_TaxExclusiveAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_TaxInclusiveAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_TaxInclusiveAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_AllowanceTotalAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_AllowanceTotalAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_ChargeTotalAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_ChargeTotalAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_PrepaidAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_PrepaidAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_PayableRoundingAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_PayableRoundingAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

type CBC_PayableAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyID money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_PayableAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyID.Code)
	return CurrencyMarshaler(m, e, s)
}

// TODO: conflicting spec
// Invoice Line Item section in https://sdk.myinvois.hasil.gov.my/documents/invoice-v1-1/#invoice-line-item
// mentions "cac:ItemPriceExtension" (ubl:Invoice / cac:InvoiceLine / cac:ItemPriceExtension / cbc:Amount [@currencyID=’MYR’])
// but this element is not mentioned in the peppol spec

// TODO: conflicting spec
// the link above mentions ubl:Invoice / cac:InvoiceLine / cac:TaxTotal / cbc:TaxAmount
// but element cac:TaxTotal is not mentioned in cac:InvoiceLine peppol spec

type CAC_InvoiceLine struct {
	XMLName             xml.Name                `xml:"cac:InvoiceLine"`
	ID                  string                  `xml:"cbc:ID"`                  // [1..1] Invoice line identifier - A unique identifier for the individual line within the Invoice.
	Note                *string                 `xml:"cbc:Note"`                // [0..1] Invoice line note - A textual note that gives unstructured information that is relevant to the Invoice line.
	InvoicedQuantity    CBC_InvoicedQuantity    `xml:"cbc:InvoicedQuantity"`    // [1..1] Invoiced quantity - The quantity of items (goods or services) that is charged in the Invoice line.
	LineExtensionAmount CBC_LineExtensionAmount `xml:"cbc:LineExtensionAmount"` // [1..1] Invoice line net amount - The total amount of the Invoice line (before tax).
	AccountingCost      *string                 `xml:"cbc:AccountingCost"`      // [0..1] Invoice line Buyer accounting reference - A textual value that specifies where to book the relevant data into the Buyer’s financial accounts.
	InvoicePeriod       *CAC_InvoicePeriod      `xml:"cac:InvoicePeriod"`       // [0..1] INVOICE LINE PERIOD - A group of business terms providing information about the period relevant for the Invoice line.
	OrderLineReference  *CAC_OrderLineReference `xml:"cac:OrderLineReference"`  // [0..1] ORDER LINE REFERENCE
	DocumentReference   *CAC_DocumentReference  `xml:"cac:DocumentReference"`   // [0..1]  [cbc:DocumentTypeCode = 130] 	LINE OBJECT IDENTIFIER
	// [0..n]  [cbc:ChargeIndicator = false] 	INVOICE LINE ALLOWANCES - A group of business terms providing information about allowances applicable to the individual Invoice line.
	// [0..n]  [cbc:ChargeIndicator = true] 	INVOICE LINE CHARGES - A group of business terms providing information about charges and taxes other than TAX applicable to the individual Invoice line.
	AllowanceCharge []CAC_AllowanceCharge `xml:"cac:AllowanceCharge"` // [0..n] INVOICE LINE CHARGES - A group of business terms providing information about charges and taxes other than TAX applicable to the individual Invoice line.
	Item            CAC_Item              `xml:"cac:Item"`            // [1..1] ITEM INFORMATION - A group of business terms providing information about the goods and services invoiced.
	Price           CAC_Price             `xml:"cac:Price"`           // [1..1] PRICE DETAILS - A group of business terms providing information about the price applied for the goods and services invoiced on the Invoice line.
}

// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/syntax/cac-InvoiceLine/cbc-InvoicedQuantity/unitCode/
// https://docs.peppol.eu/poac/my/pint-my/trn-invoice/codelist/UNECERec20/
type CBC_InvoicedQuantity struct {
	Value    string `xml:",innerxml"`
	UnitCode string `xml:"unitCode,attr"`
}

type CAC_OrderLineReference struct {
	XMLName xml.Name `xml:"cac:OrderLineReference"`
	LineID  string   `xml:"cbc:LineID"` // [1..1] Referenced purchase order line reference -  An identifier for a referenced line within a purchase order, issued by the Buyer.
}

type CAC_DocumentReference struct {
	XMLName          xml.Name `xml:"cac:DocumentReference"`
	ID               string   `xml:"cbc:ID"`               // [1..1] Invoice line object identifier - An identifier for an object on which the invoice line is based, given by the Seller.
	DocumentTypeCode string   `xml:"cbc:DocumentTypeCode"` // [1..1] Qualifier for syntax binding
}

type CAC_Item struct {
	XMLName                    xml.Name                        `xml:"cac:Item"`
	Description                *string                         `xml:"cbc:Description"`                // [0..1]	Item description
	Name                       string                          `xml:"cbc:Name"`                       // [1..1]	Item name
	BuyersItemIdentification   *CAC_BuyersItemIdentification   `xml:"cac:BuyersItemIdentification"`   // [0..1]	BUYERS ITEM IDENTIFICATION
	SellersItemIdentification  *CAC_SellersItemIdentification  `xml:"cac:SellersItemIdentification"`  // [0..1]	SELLERS ITEM IDENTIFICATION
	StandardItemIdentification *CAC_StandardItemIdentification `xml:"cac:StandardItemIdentification"` // [0..1]	STANDARD ITEM IDENTIFICATION
	OriginCountry              *CAC_OriginCountry              `xml:"cac:OriginCountry"`              // [0..1]	ORIGIN COUNTRY
	CommodityClassification    []CAC_CommodityClassification   `xml:"cac:CommodityClassification"`    // [0..n]
	CommodityCode              *string                         `xml:"cbc:CommodityCode"`              // [0..1]	Product Tariff Code - Harmonized System code of the goods under the relevant Sales Tax Orders
	ClassifiedTaxCategory      []CAC_ClassifiedTaxCategory     `xml:"cac:ClassifiedTaxCategory"`      // [0..n]	LINE TAX INFORMATION - A group of business terms providing information about the TAX applicable for the goods and services invoiced on the Invoice line.
	AdditionalItemProperty     []CAC_AdditionalItemProperty    `xml:"cac:AdditionalItemProperty"`     // [0..n]	ITEM ATTRIBUTES - A group of business terms providing information about properties of the goods and services invoiced.
	ItemInstance               *CAC_ItemInstance               `xml:"cac:ItemInstance"`               // [0..1]	Item instance information
}

type CAC_BuyersItemIdentification struct {
	XMLName xml.Name `xml:"cac:BuyersItemIdentification"`
	ID      string   `xml:"cbc:ID"` // [1..1] Item Buyer's identifier - An identifier, assigned by the Buyer, for the item.
}

type CAC_SellersItemIdentification struct {
	XMLName xml.Name `xml:"cac:SellersItemIdentification"`
	ID      string   `xml:"cbc:ID"` // [1..1] Item Seller's identifier - An identifier, assigned by the Seller, for the item.
}

type CAC_StandardItemIdentification struct {
	XMLName xml.Name `xml:"cac:StandardItemIdentification"`
	ID      string   `xml:"cbc:ID"` // [1..1] Item standard identifier - An item identifier based on a registered scheme.
}

type CAC_OriginCountry struct {
	XMLName            xml.Name `xml:"cac:OriginCountry"`
	IdentificationCode string   `xml:"cbc:IdentificationCode"` // [1..1] Item country of origin - The code identifying the country from which the item originates.
}

type CAC_CommodityClassification struct {
	XMLName                xml.Name `xml:"cac:CommodityClassification"`
	ItemClassificationCode string   `xml:"cbc:ItemClassificationCode"` // [1..1] Item classification identifier - A code for classifying the item by its type or nature.
}

// Valid values for ID: T E O (aligned-ibrp-cl-01-my)

type CAC_ClassifiedTaxCategory struct {
	XMLName            xml.Name      `xml:"cac:ClassifiedTaxCategory"`
	ID                 string        `xml:"cbc:ID"`                 // [1..1] Invoiced item TAX category code - The TAX category code for the invoiced item.
	Percent            *string       `xml:"cbc:Percent"`            // [0..1] Invoiced item TAX rate - The TAX rate, represented as percentage that applies to the invoiced item.
	TaxExemptionReason *string       `xml:"cbc:TaxExemptionReason"` // [0..1] TAX exemption reason code - A coded statement of the reason for why the line amount is exempted from TAX.
	TaxScheme          CAC_TaxScheme `xml:"cac:TaxScheme"`          // [1..1] TAX SCHEME
}

type CAC_AdditionalItemProperty struct {
	XMLName xml.Name `xml:"cac:AdditionalItemProperty"`
	Name    string   `xml:"cbc:Name"`  // [1..1] The name of the attribute or property of the item.
	Value   string   `xml:"cbc:Value"` // [1..1] The value of the attribute or property of the item.
}

type CAC_ItemInstance struct {
	XMLName           xml.Name               `xml:"cac:ItemInstance"`
	LotIdentification *CAC_LotIdentification `xml:"cac:LotIdentification"` // [0..1]
}

type CAC_LotIdentification struct {
	XMLName     xml.Name `xml:"cac:LotIdentification"`
	LotNumberID *string  `xml:"cbc:LotNumberID"` // [0..1] Batch number - An identifier for the production bach or lot that the items come from. The information applies to the full quantity of items on the relevant invoice line.
	ExpiryDate  *string  `xml:"cbc:ExpiryDate"`  // [0..1] Expiry date - The date when the items in the relevant invoice line expire.
}

type CAC_Price struct {
	XMLName         xml.Name             `xml:"cac:Price"`
	PriceAmount     CBC_PriceAmount      `xml:"cbc:PriceAmount" validate:"required"` // [1..1] Item net price - The price of an item, exclusive of TAX, after subtracting item price discount.
	BaseQuantity    int                  `xml:"cbc:BaseQuantity,omitempty"`          // [0..1] Item price base quantity - The number of item units to which the price applies.
	AllowanceCharge *CAC_AllowanceCharge `xml:"cac:AllowanceCharge"`                 // [0..1] ALLOWANCE
}

type CBC_PriceAmount struct {
	Value      money.Amount   `xml:",innerxml"`       // required
	CurrencyId money.Currency `xml:"currencyID,attr"` // required
}

func (a *CBC_PriceAmount) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	m := money.New(a.Value, a.CurrencyId.Code)

	return CurrencyMarshaler(m, e, s)
}

// type Option[T any] func(*T)

// func NewCAC_Price(priceAmount string, options ...Option[CAC_Price]) *CAC_Price {
// 	price := CAC_Price{}
// 	price.PriceAmount = priceAmount

// 	for _, option := range options {
// 		option(&price)
// 	}

// 	return &price
// }
