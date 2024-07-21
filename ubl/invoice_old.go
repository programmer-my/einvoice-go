package ubl

// // http://www.datypic.com/sc/ubl21/e-ns39_Invoice.html

// type CAC_Language struct {
// 	ID         string `xml:"cbc:ID"`
// 	Name       string `xml:"cbc:Name"`
// 	LocaleCode string `xml:"cbc:LocaleCode"`
// }

// type CAC_AddressLine struct {
// 	Line string `xml:"cbc:Line"`
// }

// type CAC_Country struct {
// 	IdentificationCode string `xml:"cbc:IdentificationCode"` // [0..1]    A code signifying this country.
// 	Name               string `xml:"cbc:Name"`               // [0..1]    The name of this country.
// }

// type CAC_LocationCoordinate struct {
// 	CoordinateSystemCode    string `xml:"cbc:CoordinateSystemCode"`    // [0..1]    A code signifying the location system used.
// 	LatitudeDegreesMeasure  string `xml:"cbc:LatitudeDegreesMeasure"`  // [0..1]    The degree component of a latitude measured in degrees and minutes.
// 	LatitudeMinutesMeasure  string `xml:"cbc:LatitudeMinutesMeasure"`  // [0..1]    The minutes component of a latitude measured in degrees and minutes (modulo 60).
// 	LatitudeDirectionCode   string `xml:"cbc:LatitudeDirectionCode"`   // [0..1]    A code signifying the direction of latitude measurement from the equator (north or south).
// 	LongitudeDegreesMeasure string `xml:"cbc:LongitudeDegreesMeasure"` // [0..1]    The degree component of a longitude measured in degrees and minutes.
// 	LongitudeMinutesMeasure string `xml:"cbc:LongitudeMinutesMeasure"` // [0..1]    The minutes component of a longitude measured in degrees and minutes (modulo 60).
// 	LongitudeDirectionCode  string `xml:"cbc:LongitudeDirectionCode"`  // [0..1]    A code signifying the direction of longitude measurement from the prime meridian (east or west).
// 	AltitudeMeasure         string `xml:"cbc:AltitudeMeasure"`         // [0..1]    The altitude of the location.
// }

// // TODO: resolve recursion
// type CAC_SubsidiaryLocation struct {
// 	ID                   string                  `xml:"cbc:ID"`                   // [0..1]    An identifier for this location, e.g., the EAN Location Number, GLN.
// 	Description          string                  `xml:"cbc:Description"`          // [0..*]    Text describing this location.
// 	Conditions           string                  `xml:"cbc:Conditions"`           // [0..*]    Free-form text describing the physical conditions of the location.
// 	CountrySubentity     string                  `xml:"cbc:CountrySubentity"`     // [0..1]    A territorial division of a country, such as a county or state, expressed as text.
// 	CountrySubentityCode string                  `xml:"cbc:CountrySubentityCode"` // [0..1]    A territorial division of a country, such as a county or state, expressed as a code.
// 	LocationTypeCode     string                  `xml:"cbc:LocationTypeCode"`     // [0..1]    A code signifying the type of location.
// 	InformationURI       string                  `xml:"cbc:InformationURI"`       // [0..1]    The Uniform Resource Identifier (URI) of a document providing information about this location.
// 	Name                 string                  `xml:"cbc:Name"`                 // [0..1]    The name of this location.
// 	ValidityPeriod       string                  `xml:"cac:ValidityPeriod"`       // [0..*]    A period during which this location can be used (e.g., for delivery).
// 	Address              CAC_Address             `xml:"cac:Address"`              // [0..1]    The address of this location.
// 	SubsidiaryLocation   *CAC_SubsidiaryLocation `xml:"cac:SubsidiaryLocation"`   // [0..*]    A location subsidiary to this location.
// 	LocationCoordinate   *CAC_LocationCoordinate `xml:"cac:LocationCoordinate"`   // [0..*]    The geographical coordinates of this location.
// }

// type CAC_ValidityPeriod struct {
// 	StartDate       string `xml:"cbc:StartDate"`       // [0..1]    The date on which this period begins.
// 	StartTime       string `xml:"cbc:StartTime"`       // [0..1]    The time at which this period begins.
// 	EndDate         string `xml:"cbc:EndDate"`         // [0..1]    The date on which this period ends.
// 	EndTime         string `xml:"cbc:EndTime"`         // [0..1]    The time at which this period ends.
// 	DurationMeasure string `xml:"cbc:DurationMeasure"` // [0..1]    The duration of this period, expressed as an ISO 8601 code.
// 	DescriptionCode string `xml:"cbc:DescriptionCode"` // [0..*]    A description of this period, expressed as a code.
// 	Description     string `xml:"cbc:Description"`     // [0..*]    A description of this period, expressed as text.
// }

// type CAC_Address struct {
// 	ID                   string                 `xml:"cbc:ID"`                   // [0..1]    An identifier for this address within an agreed scheme of address identifiers.
// 	AddressTypeCode      string                 `xml:"cbc:AddressTypeCode"`      // [0..1]    A mutually agreed code signifying the type of this address.
// 	AddressFormatCode    string                 `xml:"cbc:AddressFormatCode"`    // [0..1]    A mutually agreed code signifying the format of this address.
// 	Postbox              string                 `xml:"cbc:Postbox"`              // [0..1]    A post office box number registered for postal delivery by a postal service provider.
// 	Floor                string                 `xml:"cbc:Floor"`                // [0..1]    An identifiable floor of a building.
// 	Room                 string                 `xml:"cbc:Room"`                 // [0..1]    An identifiable room, suite, or apartment of a building.
// 	StreetName           string                 `xml:"cbc:StreetName"`           // [0..1]    The name of the street, road, avenue, way, etc. to which the number of the building is attached.
// 	AdditionalStreetName string                 `xml:"cbc:AdditionalStreetName"` // [0..1]    An additional street name used to further clarify the address.
// 	BlockName            string                 `xml:"cbc:BlockName"`            // [0..1]    The name of the block (an area surrounded by streets and usually containing several buildings) in which this address is located.
// 	BuildingName         string                 `xml:"cbc:BuildingName"`         // [0..1]    The name of a building.
// 	BuildingNumber       string                 `xml:"cbc:BuildingNumber"`       // [0..1]    The number of a building within the street.
// 	InhouseMail          string                 `xml:"cbc:InhouseMail"`          // [0..1]    The specific identifable location within a building where mail is delivered.
// 	Department           string                 `xml:"cbc:Department"`           // [0..1]    The department of the addressee.
// 	MarkAttention        string                 `xml:"cbc:MarkAttention"`        // [0..1]    The name, expressed as text, of a person or department in an organization to whose attention incoming mail is directed; corresponds to the printed forms "for the attention of", "FAO", and ATTN:".
// 	MarkCare             string                 `xml:"cbc:MarkCare"`             // [0..1]    The name, expressed as text, of a person or organization at this address into whose care incoming mail is entrusted; corresponds to the printed forms "care of" and "c/o".
// 	PlotIdentification   string                 `xml:"cbc:PlotIdentification"`   // [0..1]    An identifier (e.g., a parcel number) for the piece of land associated with this address.
// 	CitySubdivisionName  string                 `xml:"cbc:CitySubdivisionName"`  // [0..1]    The name of the subdivision of a city, town, or village in which this address is located, such as the name of its district or borough.
// 	CityName             string                 `xml:"cbc:CityName"`             // [0..1]    The name of a city, town, or village.
// 	PostalZone           string                 `xml:"cbc:PostalZone"`           // [0..1]    The postal identifier for this address according to the relevant national postal service, such as a ZIP code or Post Code.
// 	CountrySubentity     string                 `xml:"cbc:CountrySubentity"`     // [0..1]    The political or administrative division of a country in which this address is located, such as the name of its county, province, or state, expressed as text.
// 	CountrySubentityCode string                 `xml:"cbc:CountrySubentityCode"` // [0..1]    The political or administrative division of a country in which this address is located, such as a county, province, or state, expressed as a code (typically nationally agreed).
// 	Region               string                 `xml:"cbc:Region"`               // [0..1]    The recognized geographic or economic region or group of countries in which this address is located.
// 	District             string                 `xml:"cbc:District"`             // [0..1]    The district or geographical division of a country or region in which this address is located.
// 	TimezoneOffset       string                 `xml:"cbc:TimezoneOffset"`       // [0..1]    The time zone in which this address is located (as an offset from Universal Coordinated Time (UTC)) at the time of exchange.
// 	AddressLine          []CAC_AddressLine      `xml:"cac:AddressLine"`          // [0..*]    An unstructured address line.
// 	Country              CAC_Country            `xml:"cac:Country"`              // [0..1]    The country in which this address is situated.
// 	LocationCoordinate   CAC_LocationCoordinate `xml:"cac:LocationCoordinate"`   // [0..*]    The geographical coordinates of this address.
// }

// type CAC_PhysicalLocation struct {
// 	ID                   string                 `xml:"cbc:ID"`                   // [0..1]    An identifier for this location, e.g., the EAN Location Number, GLN.
// 	Description          string                 `xml:"cbc:Description"`          // [0..*]    Text describing this location.
// 	Conditions           string                 `xml:"cbc:Conditions"`           // [0..*]    Free-form text describing the physical conditions of the location.
// 	CountrySubentity     string                 `xml:"cbc:CountrySubentity"`     // [0..1]    A territorial division of a country, such as a county or state, expressed as text.
// 	CountrySubentityCode string                 `xml:"cbc:CountrySubentityCode"` // [0..1]    A territorial division of a country, such as a county or state, expressed as a code.
// 	LocationTypeCode     string                 `xml:"cbc:LocationTypeCode"`     // [0..1]    A code signifying the type of location.
// 	InformationURI       string                 `xml:"cbc:InformationURI"`       // [0..1]    The Uniform Resource Identifier (URI) of a document providing information about this location.
// 	Name                 string                 `xml:"cbc:Name"`                 // [0..1]    The name of this location.
// 	ValidityPeriod       CAC_ValidityPeriod     `xml:"cac:ValidityPeriod"`       // [0..*]    A period during which this location can be used (e.g., for delivery).
// 	Address              CAC_Address            `xml:"cac:Address"`              // [0..1]    The address of this location.
// 	SubsidiaryLocation   CAC_SubsidiaryLocation `xml:"cac:SubsidiaryLocation"`   // [0..*]    A location subsidiary to this location.
// 	LocationCoordinate   CAC_LocationCoordinate `xml:"cac:LocationCoordinate"`   // [0..*]    The geographical coordinates of this location.
// }

// type CAC_TaxScheme struct {
// 	ID                        string      `xml:"cbc:ID"`                        // [0..1]    An identifier for this taxation scheme.
// 	Name                      string      `xml:"cbc:Name"`                      // [0..1]    The name of this taxation scheme.
// 	TaxTypeCode               string      `xml:"cbc:TaxTypeCode"`               // [0..1]    A code signifying the type of tax.
// 	CurrencyCode              string      `xml:"cbc:CurrencyCode"`              // [0..1]    A code signifying the currency in which the tax is collected and reported.
// 	JurisdictionRegionAddress CAC_Address `xml:"cac:JurisdictionRegionAddress"` // [0..*]    A geographic area in which this taxation scheme applies.
// }

// type CAC_PartyTaxScheme struct {
// 	RegistrationName    string        `xml:"cbc:RegistrationName"`    // [0..1]    The name of the party as registered with the relevant fiscal authority.
// 	CompanyID           string        `xml:"cbc:CompanyID"`           // [0..1]    An identifier for the party assigned for tax purposes by the taxation authority.
// 	TaxLevelCode        string        `xml:"cbc:TaxLevelCode"`        // [0..1]    A code signifying the tax level applicable to the party within this taxation scheme.
// 	ExemptionReasonCode string        `xml:"cbc:ExemptionReasonCode"` // [0..1]    A reason for the party's exemption from tax, expressed as a code.
// 	ExemptionReason     string        `xml:"cbc:ExemptionReason"`     // [0..*]    A reason for the party's exemption from tax, expressed as text.
// 	RegistrationAddress CAC_Address   `xml:"cac:RegistrationAddress"` // [0..1]    The address of the party as registered for tax purposes.
// 	TaxScheme           CAC_TaxScheme `xml:"cac:TaxScheme"`           // [1..1]    The taxation scheme applicable to the party.
// }

// type CAC_PartyLegalEntity struct {
// }

// type CAC_Contact struct {
// }

// type CAC_Person struct {
// }

// type CAC_AgentParty struct {
// }

// type CAC_ServiceProviderParty struct {
// }

// type CAC_PowerOfAttorney struct {
// }

// type CAC_FinancialAccount struct {
// }

// type CAC_PartyIdentification struct {
// 	ID string `xml:"cbc:ID"` // [1..1]    An identifier for the party.
// }

// type CAC_PartyName struct {
// 	Name string `xml:"cbc:name"`
// }

// type Party struct {
// 	MarkCareIndicator          string                     `xml:"cbc:MarkCareIndicator"`          // [0..1]    An indicator that this party is "care of" (c/o) (true) or not (false).
// 	MarkAttentionIndicator     string                     `xml:"cbc:MarkAttentionIndicator"`     // [0..1]    An indicator that this party is "for the attention of" (FAO) (true) or not (false).
// 	WebsiteURI                 string                     `xml:"cbc:WebsiteURI"`                 // [0..1]    The Uniform Resource Identifier (URI) that identifies this party's web site; i.e., the web site's Uniform Resource Locator (URL).
// 	LogoReferenceID            string                     `xml:"cbc:LogoReferenceID"`            // [0..1]    An identifier for this party's logo.
// 	EndpointID                 string                     `xml:"cbc:EndpointID"`                 // [0..1]    An identifier for the end point of the routing service (e.g., EAN Location Number, GLN).
// 	IndustryClassificationCode string                     `xml:"cbc:IndustryClassificationCode"` // [0..1]    This party's Industry Classification Code.
// 	PartyIdentification        []CAC_PartyIdentification  `xml:"cac:PartyIdentification"`        // [0..*]    An identifier for this party.
// 	PartyName                  []CAC_PartyName            `xml:"cac:PartyName"`                  // [0..*]    A name for this party.
// 	Language                   *CAC_Language              `xml:"cac:Language"`                   // [0..1]    The language associated with this party.
// 	PostalAddress              *CAC_Address               `xml:"cac:PostalAddress"`              // [0..1]    The party's postal address.
// 	PhysicalLocation           *CAC_PhysicalLocation      `xml:"cac:PhysicalLocation"`           // [0..1]    The physical location of this party.
// 	PartyTaxScheme             []CAC_PartyTaxScheme       `xml:"cac:PartyTaxScheme"`             // [0..*]    A tax scheme applying to this party.
// 	PartyLegalEntity           []CAC_PartyLegalEntity     `xml:"cac:PartyLegalEntity"`           // [0..*]    A description of this party as a legal entity.
// 	Contact                    *CAC_Contact               `xml:"cac:Contact"`                    // [0..1]    The primary contact for this party.
// 	Person                     []CAC_Person               `xml:"cac:Person"`                     // [0..*]    A person associated with this party.
// 	AgentParty                 *CAC_AgentParty            `xml:"cac:AgentParty"`                 // [0..1]    A party who acts as an agent for this party.
// 	ServiceProviderParty       []CAC_ServiceProviderParty `xml:"cac:ServiceProviderParty"`       // [0..*]    A party providing a service to this party.
// 	PowerOfAttorney            []CAC_PowerOfAttorney      `xml:"cac:PowerOfAttorney"`            // [0..*]    A power of attorney associated with this party.
// 	FinancialAccount           *CAC_FinancialAccount      `xml:"cac:FinancialAccount"`           // [0..1]    The financial account associated with this party.
// }

// type CAC_OtherCommunication struct {
// 	ChannelCode string `xml:"cbc:ChannelCode"` // [0..1]    The method of communication, expressed as a code.
// 	Channel     string `xml:"cbc:Channel"`     // [0..1]    The method of communication, expressed as text.
// 	Value       string `xml:"cbc:Value"`       // [0..1]    An identifying value (phone number, email address, etc.) for this channel of communication
// }

// type Contact struct {
// 	ID                 string                  `xml:"cbc:ID"`                 // [0..1]    An identifier for this contact.
// 	Name               string                  `xml:"cbc:Name"`               // [0..1]    The name of this contact. It is recommended that this be used for a functional name and not a personal name.
// 	Telephone          string                  `xml:"cbc:Telephone"`          // [0..1]    The primary telephone number of this contact.
// 	Telefax            string                  `xml:"cbc:Telefax"`            // [0..1]    The primary fax number of this contact.
// 	ElectronicMail     string                  `xml:"cbc:ElectronicMail"`     // [0..1]    The primary email address of this contact.
// 	Note               string                  `xml:"cbc:Note"`               // [0..*]    Free-form text conveying information that is not contained explicitly in other structures; in particular, a textual description of the circumstances under which this contact can be used (e.g., "emergency" or "after hours").
// 	OtherCommunication *CAC_OtherCommunication `xml:"cac:OtherCommunication"` // [0..*]    Another means of communication with this contact.
// }

// type CAC_DespatchContact struct {
// 	Contact
// }

// type CAC_AccountingContact struct {
// 	Contact
// }

// type CAC_SellerContact struct {
// 	Contact
// }

// type CAC_AccountingSupplierParty struct {
// 	CustomerAssignedAccountID string                 `xml:"cbc:CustomerAssignedAccountID,omitempty"` // [0..1]    An identifier for this supplier party, assigned by the customer.
// 	AdditionalAccountID       string                 `xml:"cbc:AdditionalAccountID,omitempty"`       // [0..*]    An additional identifier for this supplier party.
// 	DataSendingCapability     string                 `xml:"cbc:DataSendingCapability,omitempty"`     // [0..1]    Text describing the supplier's ability to send invoice data via a purchase card provider (e.g., VISA, MasterCard, American Express).
// 	Party                     *Party                 `xml:"cac:Party,omitempty"`                     // [0..1]    The supplier party itself.
// 	DespatchContact           *CAC_DespatchContact   `xml:"cac:DespatchContact,omitempty"`           // [0..1]    A contact at this supplier party for despatches (pickups).
// 	AccountingContact         *CAC_AccountingContact `xml:"cac:AccountingContact,omitempty"`         // [0..1]    A contact at this supplier party for accounting.
// 	SellerContact             *CAC_SellerContact     `xml:"cac:SellerContact,omitempty"`             // [0..1]
// }

// type Invoice struct {
// 	AccountingSupplierParty CAC_AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
// }
