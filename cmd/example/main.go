package main

import (
	"encoding/xml"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/programmer-my/einvoice-go/document"
)

type User struct {
	Name string `xml:"name,attr"`
	Age  int    `xml:"age,attr"`
}

func main() {
	// tin1 := "TIN1"
	// tin2 := "TIN2"
	// inv := ubl.NewInvoice()
	// inv.IssueDate = "2024-01-01"
	// inv.DocumentCurrencyCode = "MYR"
	// inv.TaxPointDate = "2024-01-01"
	// inv.InvoiceTypeCode = "71"

	// inv.AccountingSupplierParty.Party.EndpointID.SchemeID = "0230"
	// inv.AccountingSupplierParty.Party.EndpointID.Value = "http://supplierwebsite.com/"
	// inv.AccountingSupplierParty.Party.PostalAddress.Country.IdentificationCode = "MY"
	// inv.AccountingSupplierParty.Party.PartyTaxScheme = &ubl.CAC_PartyTaxScheme{
	// 	TaxScheme: ubl.CAC_TaxScheme{ID: "VAT"},
	// 	CompanyID: &tin1,
	// }

	// inv.AccountingCustomerParty.Party.EndpointID.SchemeID = "0230"
	// inv.AccountingCustomerParty.Party.EndpointID.Value = "http://customerwebsite.com/"
	// inv.AccountingCustomerParty.Party.PostalAddress.Country.IdentificationCode = "MY"
	// inv.AccountingCustomerParty.Party.PartyTaxScheme = &ubl.CAC_PartyTaxScheme{
	// 	TaxScheme: ubl.CAC_TaxScheme{ID: "VAT"},
	// 	CompanyID: &tin2,
	// }

	// inv.LegalMonetaryTotal = ubl.CAC_LegalMonetaryTotal{
	// 	LineExtensionAmount: ubl.CBC_LineExtensionAmount{Value: "10.00", CurrencyID: inv.DocumentCurrencyCode},
	// 	TaxExclusiveAmount:  ubl.CBC_TaxExclusiveAmount{Value: "10.00", CurrencyID: inv.DocumentCurrencyCode},
	// 	TaxInclusiveAmount:  ubl.CBC_TaxInclusiveAmount{Value: "10.00", CurrencyID: inv.DocumentCurrencyCode},
	// 	PayableAmount:       ubl.CBC_PayableAmount{Value: "10.00", CurrencyID: "MYR"},
	// }
	// inv.TaxTotal = ubl.CAC_TaxTotal{
	// 	TaxAmount: ubl.CBC_TaxAmount{
	// 		Value:      "10.00",
	// 		CurrencyID: "MYR",
	// 	},
	// 	TaxSubtotal: []ubl.CAC_TaxSubtotal{},
	// }
	// desc := "Barang Baek"
	// inv.InvoiceLine = append(inv.InvoiceLine, ubl.CAC_InvoiceLine{
	// 	ID: "1",
	// 	InvoicedQuantity: ubl.CBC_InvoicedQuantity{
	// 		Value:    "100000",
	// 		UnitCode: "1I",
	// 	},
	// 	LineExtensionAmount: ubl.CBC_LineExtensionAmount{Value: "10.00", CurrencyID: inv.DocumentCurrencyCode},
	// 	Price: ubl.CAC_Price{
	// 		BaseQuantity: 1,
	// 		PriceAmount:  ubl.CBC_PriceAmount{Value: "10.00", CurrencyId: inv.DocumentCurrencyCode},
	// 	},
	// 	Item: ubl.CAC_Item{
	// 		Name:        "Barang",
	// 		Description: &desc,
	// 	},
	// })

	// b, err := xml.Marshal(inv)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(b))

	// //////////////////////

	// // tt := ubl.CAC_TaxTotal{
	// // 	TaxAmount:   "10.00",
	// // 	Currency:    "MYR",
	// // 	TaxSubtotal: nil,
	// // }

	// // b, err := xml.Marshal(tt)
	// // if err != nil {
	// // 	panic(err)
	// // }

	// // fmt.Println(string(b))
	supplier := document.InvoiceSupplier{
		Name:    "Legit Supplier",
		TIN:     "TIN0001",
		IdType:  "BRN",
		IdValue: "SSM0001",
		SSTNo:   "SST0001",
		Email:   "supplier@example.com",
		Address: document.Address{
			Line1:    "Supplier Line 1",
			Line2:    "Supplier Line 2",
			Line0:    "",
			Postcode: "47000",
			City:     "Sungai Buloh",
			State:    "Selangor",
			Country:  "MY",
		},
		ContactNo:           "+0123456789",
		TourismTaxNo:        "",
		MSICCode:            "01111",
		BusinessDescription: "Growing of maize",
	}
	buyer := document.InvoiceBuyer{
		Name:    "John Doe",
		TIN:     "TIN0002",
		IdType:  "BRN",
		IdValue: "SSM0002",
		SSTNo:   "SST0002",
		Email:   "totally_legit@gmail.com",
		Address: document.Address{
			Line1:    "Customer Line 1",
			Line2:    "Customer Line 2",
			Line0:    "",
			Postcode: "63000",
			City:     "Cyberjaya",
			State:    "Selangor",
			Country:  "MY",
		},
		ContactNo: "+60141231234",
	}

	invoice := document.InvoiceDocument{
		Code:     "0001",
		Supplier: supplier,
		Buyer:    buyer,
		Items: []document.InvoiceLineItem{
			{
				Description:       "Barang Baek",
				Classification:    "022", // others
				UnitPrice:         *money.New(10000, money.MYR),
				TaxType:           "01", // sales tax
				TaxRate:           "6",  // in %
				TaxAmount:         *money.New(600, money.MYR),
				Subtotal:          *money.New(10600, money.MYR),
				TotalExcludingTax: *money.New(600, money.MYR),
				Quantity:          "1",
				Measurement:       "1I",
			},
		},
		Date:         "2024-01-01",
		Time:         "00:00:00",
		CurrencyCode: *money.GetCurrency("MYR"),
		Version:      "1.1",
		TypeCode:     "80", // invoice
	}

	ublInvoice := document.UblInvoiceBuilder(invoice)
	ublInvoice.TaxPointDate = invoice.Date

	// fmt.Printf("%+v\n", ublInvoice)
	// encoder := xml.NewEncoder(os.Stdout)
	// encoder.Encode(ublInvoice)

	invoiceBytes, err := xml.MarshalIndent(ublInvoice, " ", "  ")
	if err != nil {
		panic(fmt.Errorf("UBL Invoice marshal error: %s", err))
	}

	fmt.Printf("%s\n", string(invoiceBytes))
}
