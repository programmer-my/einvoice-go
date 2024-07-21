package document

import (
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/programmer-my/einvoice-go/ubl"
)

type InvoiceSupplier struct {
	Name                string  // required
	TIN                 string  // required
	IdType              string  // required. NRIC, BRN, PASSPORT, ARMY
	IdValue             string  // required. value depending on type, will be serialized as "BRN: 202001234567"
	SSTNo               string  // mandatory for SST registrants. "NA" if not provided, max 2 can be provided separated by ;
	Email               string  // optional
	Address             Address // required
	ContactNo           string  // required. E.164 format
	TourismTaxNo        string
	MSICCode            string // https://sdk.myinvois.hasil.gov.my/codes/msic-codes/
	BusinessDescription string // refer CAC_Party_IndustryClassificationCode
}

type InvoiceBuyer struct {
	Name      string  // required
	TIN       string  // required
	IdType    string  // required. NRIC, BRN, PASSPORT, ARMY
	IdValue   string  // required. value depending on type, will be serialized as "BRN: 202001234567"
	SSTNo     string  // mandatory for SST registrants. "NA" if not provided, max 2 can be provided separated by ;
	Email     string  // optional
	Address   Address // required
	ContactNo string  // required. E.164 format
}

type Address struct {
	Line0    string // optional, "NA" if empty
	Line1    string
	Line2    string
	Postcode string // postcode
	City     string
	State    string // https://sdk.myinvois.hasil.gov.my/codes/state-codes/
	Country  string // ISO3166-1 https://sdk.myinvois.hasil.gov.my/codes/countries/
}

type InvoiceLineItem struct {
	Classification         string      // required. max: 3 https://sdk.myinvois.hasil.gov.my/codes/classification-codes/
	Description            string      // required. max: 300
	UnitPrice              money.Money // required.
	TaxType                string      // required. https://sdk.myinvois.hasil.gov.my/codes/tax-types/
	TaxRate                string      // required where applicable (fixed rate, whatever that means) TODO: monetary type
	TaxAmount              money.Money // required.
	TaxExemptionInfo       string      // required if applicable
	TotalTaxAmountExempted money.Money // required if applicable.
	Subtotal               money.Money // required.
	TotalExcludingTax      money.Money // required.
	Quantity               string      // optional
	Measurement            string      // optional https://docs.peppol.eu/poac/my/pint-my/trn-invoice/codelist/UNECERec20/
	DiscountRate           string      // optional, percentage
	DiscountAmount         money.Money // optional
	// ChargeRate string
	// ChargeAmount string
	// ProductTariffCode string
	// OriginCountry string
}

type InvoiceDocument struct {
	Supplier             InvoiceSupplier
	Buyer                InvoiceBuyer
	Version              string
	TypeCode             string
	Code                 string
	Date                 string
	Time                 string // hh:mm:ss
	Signature            string // TODO
	CurrencyCode         money.Currency
	CurrencyExchangeRate string // only required if non-MYR
	// BillingFrequency string /* Daily, Weekly, Biweekly, Monthly, Bimonthly, Quarterly, Half-yearly, Yearly, Others / Not Applicable */
	// BillingPeriodStartDate string
	// BillingPeriodEndDate string
	Items []InvoiceLineItem
	// PaymentMode string
	// SupplierBankAccount string
	// PaymentTerms string
	// PrePaymentAmount string // TODO: monetary type
	// PrePaymentDate string
	// PrePaymentTime string
	// PrePaymentRefNo string
	// BillRefNo string
	TotalExcludingTax  money.Money
	TotalIncludingTax  money.Money
	TotalPayableAmount money.Money
	// TotalNetAmount money.Money
	// TotalDiscountValue money.Money
	// TotalFeeAmount money.Money
	TotalTaxAmount money.Money
	// RoundingAmount string
	// TotalTaxableAmountPerTaxType string
	TotalTaxAmountPerTaxType string
	TaxExemptionInfo         string
	TotalTaxAmountExempted   string
	TaxType                  string // https://sdk.myinvois.hasil.gov.my/codes/tax-types/
}

// Perform mapping of core data structures into UBL Invoice
//
// Reference: https://sdk.myinvois.hasil.gov.my/documents/invoice-v1-1/
func UblInvoiceBuilder(doc InvoiceDocument) *ubl.UBL_Invoice {
	supplier := doc.Supplier
	buyer := doc.Buyer

	// TODO: TaxCurrencyCode
	inv := ubl.NewInvoice()
	inv.DocumentCurrencyCode = doc.CurrencyCode.Code
	inv.Currency = doc.CurrencyCode
	inv.InvoiceTypeCode = doc.TypeCode
	inv.ID = doc.Code
	inv.IssueDate = doc.Date
	inv.IssueTime = &doc.Time
	// TODO: add Signature field to InvoiceDocument
	// inv.Signature = ""
	// TODO: refer comments (conflict in spec)
	inv.TaxExchangeRate = &ubl.CAC_TaxExchangeRate{
		CalculationRate: doc.CurrencyExchangeRate,
	}
	// TODO: SourceCurrencyCode, TargetCurrencyCode: view comments in CAC_TaxExchangeRate
	// TODO: BillingFreq (optional)
	// TODO: BillingPeriodStartDate (optional)
	// TODO: BillingPeriodEndDate (optional)

	ublLineItems := inv.InvoiceLine

	// map document.InvoiceLineItem -> ubl.CAC_InvoiceLine
	for i, item := range doc.Items {
		ublItem := ubl.CAC_InvoiceLine{
			ID: fmt.Sprintf("%d", i),
			// TODO: set note other than desc
			// Note: &item.Description,
			InvoicedQuantity: ubl.CBC_InvoicedQuantity{
				Value:    item.Quantity,
				UnitCode: item.Measurement,
			},
			LineExtensionAmount: ubl.CBC_LineExtensionAmount{
				CurrencyID: inv.Currency,
				Value:      item.TotalExcludingTax.Amount(),
			},
			Item: ubl.CAC_Item{
				Description: &item.Description,
				Name:        item.Description, // TODO: other than desc
			},
			Price: ubl.CAC_Price{
				PriceAmount: ubl.CBC_PriceAmount{
					CurrencyId: inv.Currency,
					Value:      item.UnitPrice.Amount(),
				},
			},

			// TODO: mapping of item.Subtotal to cac:InvoiceLine / cac:ItemPriceExtension / cbc:Amount
			// TODO: mapping of item.TaxAmount to cac:InvoiceLine / cac:TaxTotal / cbc:TaxAmount
			// refer comment in CAC_InvoiceLine
			// TODO: details of tax exemption, discount rate, discount amount, fee/charge rate, fee/charge amount, product tariff code, country of origin
		}

		ublLineItems = append(ublLineItems, ublItem)
	}

	inv.InvoiceLine = ublLineItems

	// TODO: mapping & calculation of ubl:Invoice / cac:TaxTotal
	inv.TaxTotal = ubl.CAC_TaxTotal{
		TaxAmount: ubl.CBC_TaxAmount{
			Value:      0,
			CurrencyID: inv.Currency,
		},
		// TODO: this is dummy value
		TaxSubtotal: []ubl.CAC_TaxSubtotal{
			{
				TaxableAmount: ubl.CBC_TaxableAmount{
					Value:      0,
					CurrencyID: inv.Currency,
				},
				TaxAmount: ubl.CBC_TaxAmount{
					Value:      0,
					CurrencyID: inv.Currency,
				},
				TaxCategory: ubl.CAC_TaxCategory{
					ID:        "SST",
					Percent:   "0",
					TaxScheme: ubl.CAC_TaxScheme{ID: "VAT"},
				},
			},
			{
				TaxableAmount: ubl.CBC_TaxableAmount{
					Value:      0,
					CurrencyID: inv.Currency,
				},
				TaxAmount: ubl.CBC_TaxAmount{
					Value:      0,
					CurrencyID: inv.Currency,
				},
				TaxCategory: ubl.CAC_TaxCategory{
					ID:        "GST",
					Percent:   "0",
					TaxScheme: ubl.CAC_TaxScheme{ID: "GST"},
				},
			},
		},
	}

	// All calculations must follow https://docs.peppol.eu/poac/my/pint-my/bis/#_calculations
	lmt_LineExtensionAmount := money.New(0, inv.Currency.Code)
	lmt_TaxExclusiveAmount := money.New(0, inv.Currency.Code)
	lmt_TaxInclusiveAmount := money.New(0, inv.Currency.Code)
	lmt_AllowanceTotalAmount := money.New(0, inv.Currency.Code)
	lmt_ChargeTotalAmount := money.New(0, inv.Currency.Code)
	lmt_PrepaidAmount := money.New(0, inv.Currency.Code)
	lmt_PayableRoundingAmount := money.New(0, inv.Currency.Code)
	lmt_PayableAmount := money.New(0, inv.Currency.Code)

	for _, line := range inv.InvoiceLine {
		// // TODO: bukan parse int!! kena decide what type nak guna untuk represent money value
		// lea, err := strconv.ParseInt(line.LineExtensionAmount.Value, 10, 64)
		// if err != nil {
		// 	// TODO: handle error
		// 	fmt.Printf("failed to convert line.LineExtensionAmount.Value: %s\n", err)
		// 	fmt.Printf("lea=%q\n", lea)
		// 	continue
		// }
		if leaNew, err := lmt_LineExtensionAmount.Add(money.New(line.LineExtensionAmount.Value, line.LineExtensionAmount.CurrencyID.Code)); err != nil {
			fmt.Printf("error when calculating lineextensionamount: %s", err)
		} else {
			lmt_LineExtensionAmount = leaNew
		}
	}

	for _, charge := range inv.AllowanceCharge {
		// // TODO: bukan parse int!! kena decide what type nak guna untuk represent money value
		// chargeAmount, err := strconv.ParseInt(charge.Amount.Value, 10, 64)
		// if err != nil {
		// 	// TODO: handle error
		// 	fmt.Printf("failed to convert charge.Amount.Value: %s\n", err)
		// 	continue
		// }

		amount := money.New(charge.Amount.Value, charge.Amount.CurrencyID.Code)

		if charge.ChargeIndicator {
			if ata, err := lmt_AllowanceTotalAmount.Add(amount); err != nil {
				fmt.Printf("error when calculating allowancetotalamount: %s", err)
			} else {
				lmt_AllowanceTotalAmount = ata
			}
		} else {
			if cta, err := lmt_ChargeTotalAmount.Add(amount); err != nil {
				fmt.Printf("error when calculating chargetotalamount: %s", err)
			} else {
				lmt_ChargeTotalAmount = cta
			}
		}
	}

	// MAKE SURE ALL OF THESE VARIABLES HAVE BEEN CALCULATED
	if tea, err := lmt_LineExtensionAmount.Subtract(money.New(lmt_AllowanceTotalAmount.Amount(), lmt_AllowanceTotalAmount.Currency().Code)); err != nil {
		fmt.Printf("error when calculating lmt_TaxExclusiveAmount - lmt_AllowanceTotalAmount: %s", err)
	} else {
		lmt_TaxExclusiveAmount = tea
		if tea, err := lmt_TaxExclusiveAmount.Subtract(money.New(lmt_ChargeTotalAmount.Amount(), lmt_ChargeTotalAmount.Currency().Code)); err != nil {
			fmt.Printf("error when calculating lmt_TaxExclusiveAmount - lmt_ChargeTotalAmount: %s", err)
		} else {
			lmt_TaxExclusiveAmount = tea
		}
	}
	// lmt_TaxExclusiveAmount = lmt_LineExtensionAmount - lmt_AllowanceTotalAmount + lmt_ChargeTotalAmount

	// // TODO: bukan parse int!! kena decide what type nak guna untuk represent money value
	// taxTotalAmount, err := strconv.ParseInt(inv.TaxTotal.TaxAmount.Value, 10, 64)
	// if err != nil {
	// 	// TODO: handle error
	// 	fmt.Printf("failed to convert inv.TaxTotal.TaxAmount.Value: %s\n", err)
	// 	fmt.Printf("inv.TaxTotal.TaxAmount.Value after = %d\n", taxTotalAmount)
	// }
	// lmt_TaxInclusiveAmount = lmt_TaxExclusiveAmount + taxTotalAmount

	if tia, err := lmt_TaxExclusiveAmount.Add(money.New(inv.TaxTotal.TaxAmount.Value, inv.TaxTotal.TaxAmount.CurrencyID.Code)); err != nil {
		fmt.Printf("error when calculating lmt_TaxExclusiveAmount + inv.TaxTotal.TaxAmount.Value: %s", err)
	} else {
		lmt_TaxInclusiveAmount = tia
	}

	// lmt_PayableAmount = lmt_TaxInclusiveAmount - lmt_PrepaidAmount + lmt_PayableRoundingAmount
	if pa, err := lmt_TaxInclusiveAmount.Subtract(lmt_PrepaidAmount); err != nil {
		fmt.Printf("error when calculating lmt_TaxInclusiveAmount - lmt_PrepaidAmount: %s", err)
	} else {
		lmt_PayableAmount = pa
	}

	if pa, err := lmt_PayableAmount.Add(lmt_PayableRoundingAmount); err != nil {
		fmt.Printf("error when calculating lmt_PayableAmount + lmt_PayableRoundingAmount: %s", err)
	} else {
		lmt_PayableAmount = pa
	}

	// Finally, set the calculated values into the invoice
	inv.LegalMonetaryTotal.LineExtensionAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.LineExtensionAmount.Value = lmt_LineExtensionAmount.Amount()

	inv.LegalMonetaryTotal.TaxExclusiveAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.TaxExclusiveAmount.Value = lmt_TaxExclusiveAmount.Amount()

	inv.LegalMonetaryTotal.TaxInclusiveAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.TaxInclusiveAmount.Value = lmt_TaxInclusiveAmount.Amount()

	inv.LegalMonetaryTotal.AllowanceTotalAmount = &ubl.CBC_AllowanceTotalAmount{}
	inv.LegalMonetaryTotal.AllowanceTotalAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.AllowanceTotalAmount.Value = lmt_AllowanceTotalAmount.Amount()

	inv.LegalMonetaryTotal.ChargeTotalAmount = &ubl.CBC_ChargeTotalAmount{}
	inv.LegalMonetaryTotal.ChargeTotalAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.ChargeTotalAmount.Value = lmt_ChargeTotalAmount.Amount()

	inv.LegalMonetaryTotal.PrepaidAmount = &ubl.CBC_PrepaidAmount{}
	inv.LegalMonetaryTotal.PrepaidAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.PrepaidAmount.Value = lmt_PrepaidAmount.Amount()

	inv.LegalMonetaryTotal.PayableRoundingAmount = &ubl.CBC_PayableRoundingAmount{}
	inv.LegalMonetaryTotal.PayableRoundingAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.PayableRoundingAmount.Value = lmt_PayableRoundingAmount.Amount()

	inv.LegalMonetaryTotal.PayableAmount.CurrencyID = inv.Currency
	inv.LegalMonetaryTotal.PayableAmount.Value = lmt_PayableAmount.Amount()

	// taxSubtotalByType := make(map[string]string)

	/** Start fill in supplier info **/

	inv.AccountingSupplierParty.Party.EndpointID.SchemeID = "0230"
	inv.AccountingSupplierParty.Party.EndpointID.Value = supplier.MSICCode
	inv.AccountingSupplierParty.Party.PartyLegalEntity.RegistrationName = supplier.Name
	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.TIN, SchemeID: "0230"},
	})

	// inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 	ID: ubl.CAC_PartyIdentification_ID{Value: supplier.TIN, SchemeID: "TIN"},
	// })

	// if supplier.IdType == "NRIC" {
	// 	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.IdValue, SchemeID: "NRIC"},
	// 	})
	// } else if supplier.IdType == "BRN" {
	// 	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.IdValue, SchemeID: "BRN"},
	// 	})
	// } else if supplier.IdType == "PASSPORT" {
	// 	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.IdValue, SchemeID: "PASSPORT"},
	// 	})
	// } else if supplier.IdType == "ARMY" {
	// 	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.IdValue, SchemeID: "ARMY"},
	// 	})
	// }

	// // TODO: only mandatory for SST registrants
	// inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 	ID: ubl.CAC_PartyIdentification_ID{Value: supplier.SSTNo, SchemeID: "SST"},
	// })

	// // TODO: only mandatory for Tourism Tax (TTX)
	// if supplier.TourismTaxNo != "" {
	// 	inv.AccountingSupplierParty.Party.PartyIdentification = append(inv.AccountingSupplierParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: supplier.TourismTaxNo, SchemeID: "TTX"},
	// 	})
	// }

	inv.AccountingSupplierParty.Party.Contact = &ubl.CAC_Contact{
		ElectronicMail: &supplier.Email,
		Telephone:      &supplier.ContactNo,
	}
	inv.AccountingSupplierParty.Party.IndustryClassificationCode = &supplier.MSICCode
	inv.AccountingSupplierParty.Party.PostalAddress.StreetName = &supplier.Address.Line1
	inv.AccountingSupplierParty.Party.PostalAddress.AdditionalStreetName = &supplier.Address.Line2
	if supplier.Address.Line0 != "" {
		inv.AccountingSupplierParty.Party.PostalAddress.AddressLine = &ubl.CAC_AddressLine{
			Line: supplier.Address.Line0,
		}
	}
	inv.AccountingSupplierParty.Party.PostalAddress.CityName = &supplier.Address.City
	inv.AccountingSupplierParty.Party.PostalAddress.PostalZone = &supplier.Address.Postcode
	inv.AccountingSupplierParty.Party.PostalAddress.CountrySubentity = &supplier.Address.State
	inv.AccountingSupplierParty.Party.PostalAddress.Country = ubl.CAC_Country{
		IdentificationCode: supplier.Address.Country,
	}

	/** End fill in supplier info **/

	/** Start fill in buyer info **/

	inv.AccountingCustomerParty.Party.EndpointID.SchemeID = "0230"
	inv.AccountingCustomerParty.Party.EndpointID.Value = supplier.MSICCode
	inv.AccountingCustomerParty.Party.PartyLegalEntity.RegistrationName = buyer.Name
	inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
		ID: ubl.CAC_PartyIdentification_ID{Value: buyer.TIN, SchemeID: "0230"},
	})

	// inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 	ID: ubl.CAC_PartyIdentification_ID{Value: buyer.TIN, SchemeID: "TIN"},
	// })

	// if buyer.IdType == "NRIC" {
	// 	inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: buyer.IdValue, SchemeID: "NRIC"},
	// 	})
	// } else if buyer.IdType == "BRN" {
	// 	inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: buyer.IdValue, SchemeID: "BRN"},
	// 	})
	// } else if buyer.IdType == "PASSPORT" {
	// 	inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: buyer.IdValue, SchemeID: "PASSPORT"},
	// 	})
	// } else if buyer.IdType == "ARMY" {
	// 	inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 		ID: ubl.CAC_PartyIdentification_ID{Value: buyer.IdValue, SchemeID: "ARMY"},
	// 	})
	// }

	// // TODO: only mandatory for SST registrants
	// inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 	ID: ubl.CAC_PartyIdentification_ID{Value: buyer.SSTNo, SchemeID: "SST"},
	// })

	// TODO: buyer no TTX?
	// inv.AccountingCustomerParty.Party.PartyIdentification = append(inv.AccountingCustomerParty.Party.PartyIdentification, ubl.CAC_PartyIdentification{
	// 	ID: ubl.CAC_PartyIdentification_ID{Value: buyer.TourismTaxNo, SchemeID: "TTX"},
	// })
	inv.AccountingCustomerParty.Party.Contact = &ubl.CAC_Contact{
		ElectronicMail: &buyer.Email,
		Telephone:      &buyer.ContactNo,
	}
	inv.AccountingCustomerParty.Party.PostalAddress.StreetName = &buyer.Address.Line1
	inv.AccountingCustomerParty.Party.PostalAddress.AdditionalStreetName = &buyer.Address.Line2
	if buyer.Address.Line0 != "" {
		inv.AccountingCustomerParty.Party.PostalAddress.AddressLine = &ubl.CAC_AddressLine{
			Line: buyer.Address.Line0,
		}
	}
	inv.AccountingCustomerParty.Party.PostalAddress.CityName = &buyer.Address.City
	inv.AccountingCustomerParty.Party.PostalAddress.PostalZone = &buyer.Address.Postcode
	inv.AccountingCustomerParty.Party.PostalAddress.CountrySubentity = &buyer.Address.State
	inv.AccountingCustomerParty.Party.PostalAddress.Country.IdentificationCode = buyer.Address.Country

	/** End fill in buyer info **/

	return inv
}
