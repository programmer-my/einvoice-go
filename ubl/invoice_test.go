package ubl_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/go-playground/validator"
	"github.com/programmer-my/einvoice-go/ubl"
)

func Test_Marshal_UBL_Invoice(t *testing.T) {

}

func Test_Marshal_CAC_InvoicePeriod(t *testing.T) {

}

func Test_Marshal_CAC_OrderReference(t *testing.T) {

}

func Test_Marshal_CAC_BillingReference(t *testing.T) {

}

func Test_Marshal_CAC_InvoiceDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_DespatchDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_ReceiptDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_OriginatorDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_ContractDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_AdditionalDocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_Attachment(t *testing.T) {

}

func Test_Marshal_CAC_ExternalReference(t *testing.T) {

}

func Test_Marshal_CAC_ResultOfVerification(t *testing.T) {

}

func Test_Marshal_CAC_ProjectReference(t *testing.T) {

}

func Test_Marshal_CAC_AccountingSupplierParty(t *testing.T) {

}

func Test_Marshal_CAC_Party(t *testing.T) {

}

func Test_Marshal_CAC_PartyIdentification(t *testing.T) {

}

func Test_Marshal_CAC_PartyName(t *testing.T) {

}

func Test_Marshal_CAC_Address(t *testing.T) {

}

func Test_Marshal_CAC_PostalAddress(t *testing.T) {
	postalAddr := ubl.CAC_PostalAddress{
		CAC_Address: ubl.CAC_Address{
			Country: ubl.CAC_Country{
				IdentificationCode: "MY",
			},
		},
	}

	b, err := xml.Marshal(postalAddr)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}

	fmt.Println(string(b))
}

func Test_Marshal_CAC_AddressLine(t *testing.T) {

}

func Test_Marshal_CAC_Country(t *testing.T) {

}

func Test_Marshal_CAC_PartyTaxScheme(t *testing.T) {

}

func Test_Marshal_CAC_TaxScheme(t *testing.T) {

}

func Test_Marshal_CAC_PartyLegalEntity(t *testing.T) {

}

func Test_Marshal_CAC_Contact(t *testing.T) {

}

func Test_Marshal_CAC_AccountingCustomerParty(t *testing.T) {

}

func Test_Marshal_CAC_PayeeParty(t *testing.T) {

}

func Test_Marshal_CAC_TaxRepresentativeParty(t *testing.T) {

}

func Test_Marshal_CAC_Delivery(t *testing.T) {

}

func Test_Marshal_CAC_DeliveryLocation(t *testing.T) {

}

func Test_Marshal_CAC_DeliveryParty(t *testing.T) {

}

func Test_Marshal_CAC_Shipment(t *testing.T) {

}

func Test_Marshal_CAC_Consignment(t *testing.T) {

}

func Test_Marshal_CAC_DeliveryTerms(t *testing.T) {

}

func Test_Marshal_CAC_PaymentMeans(t *testing.T) {

}

func Test_Marshal_CAC_CardAccount(t *testing.T) {

}

func Test_Marshal_CAC_PayeeFinancialAccount(t *testing.T) {

}

func Test_Marshal_CAC_PaymentMandate(t *testing.T) {

}

func Test_Marshal_CAC_PayerFinancialAccount(t *testing.T) {

}

func Test_Marshal_CAC_PaymentTerms(t *testing.T) {

}

func Test_Marshal_CAC_PrepaidPayment(t *testing.T) {

}

func Test_Marshal_CAC_AllowanceCharge(t *testing.T) {

}

func Test_Marshal_CAC_TaxCategory(t *testing.T) {

}

func Test_Marshal_CAC_TaxExchangeRate(t *testing.T) {

}

func Test_Marshal_CAC_TaxTotal(t *testing.T) {
	taxTotal := ubl.CAC_TaxTotal{
		TaxAmount: "1.00",
		TaxSubtotal: []ubl.CAC_TaxSubtotal{
			{
				TaxableAmount: "1.00",
				TaxAmount:     "1.00",
				TaxCategory: ubl.CAC_TaxCategory{
					ID:      "T",
					Percent: "100",
					TaxScheme: ubl.CAC_TaxScheme{
						ID: "VAT",
					},
				},
			},
		},
	}

	v := validator.New()
	err := v.Struct(taxTotal)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			fmt.Printf("TAG: %s\n", e.Tag())
		}
		fmt.Printf("%+v\n", errs)
		t.Errorf("schema validation error: %s", err)
	}
}

func Test_Marshal_CAC_TaxSubtotal(t *testing.T) {

}

func Test_Marshal_CAC_LegalMonetaryTotal(t *testing.T) {

}

func Test_Marshal_CAC_InvoiceLine(t *testing.T) {

}

func Test_Marshal_CAC_OrderLineReference(t *testing.T) {

}

func Test_Marshal_CAC_DocumentReference(t *testing.T) {

}

func Test_Marshal_CAC_Item(t *testing.T) {

}

func Test_Marshal_CAC_BuyersItemIdentification(t *testing.T) {

}

func Test_Marshal_CAC_SellersItemIdentification(t *testing.T) {

}

func Test_Marshal_CAC_StandardItemIdentification(t *testing.T) {

}

func Test_Marshal_CAC_OriginCountry(t *testing.T) {

}

func Test_Marshal_CAC_CommodityClassification(t *testing.T) {

}

func Test_Marshal_CAC_ClassifiedTaxCategory(t *testing.T) {

}

func Test_Marshal_CAC_AdditionalItemProperty(t *testing.T) {

}

func Test_Marshal_CAC_ItemInstance(t *testing.T) {

}

func Test_Marshal_CAC_LotIdentification(t *testing.T) {
	li := ubl.CAC_LotIdentification{}

	b, err := xml.Marshal(li)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	fmt.Println(string(b))
}

func Test_Marshal_CAC_Price(t *testing.T) {
	price := ubl.CAC_Price{
		PriceAmount: "1.00",
		// BaseQuantity: 2,
		// AllowanceCharge: &ubl.CAC_AllowanceCharge{
		// 	ChargeIndicator: true,
		// },
	}

	v := validator.New()
	err := v.Struct(price)
	if err != nil {
		t.Errorf("schema error: %s", err)
	}

	b, err := xml.Marshal(price)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	fmt.Println(string(b))
}
