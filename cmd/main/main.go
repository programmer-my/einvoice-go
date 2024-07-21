package main

import (
	"encoding/xml"
	"fmt"

	"github.com/Rhymond/go-money"
)

type InvoiceLine struct {
	Name   string
	Amount Money
}

type Invoice struct {
	Amount Money
	Items  []InvoiceLine
}

func (i *Invoice) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	if err := e.Encode(i.Amount); err != nil {
		return err
	}

	if err := e.Encode(i.Items); err != nil {
		return err
	}

	return nil
}

type Money struct {
	money.Money
}

func (m *Money) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	toEncode := struct {
		Amount     string `xml:",innerxml"`
		CurrencyID string `xml:"currencyID,attr"`
	}{
		Amount:     fmt.Sprintf("%.2f", m.AsMajorUnits()),
		CurrencyID: m.Currency().Code,
	}

	return e.EncodeElement(toEncode, s)
}

type XMLMarshaler func(e *xml.Encoder, s xml.StartElement) error

func MoneyMarshaler(m Money, tagName string) XMLMarshaler {
	return func(e *xml.Encoder, s xml.StartElement) error {
		e.EncodeElement(
			m.Amount,
			xml.StartElement{
				Name: xml.Name{Local: tagName},
				Attr: []xml.Attr{
					xml.Attr{Name: xml.Name{Local: "currencyID"}, Value: m.Currency().Code},
				},
			},
		)
		return nil
	}
}

func main() {
	mm := Money{Money: *money.New(15050, money.MYR)}
	i := Invoice{
		Amount: mm,
	}

	b, err := xml.MarshalIndent(i, " ", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
