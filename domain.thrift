enum EnumAddressUse {
	home
	work
	temp
	old
}
enum EnumPatientContactGender {
	male
	female
	other
	unknown
}
enum EnumPatientLinkType {
	replaced_by
	replaces
	refer
	seealso
}
enum EnumPatientGender {
	male
	female
	other
	unknown
}
enum EnumNarrativeStatus {
	generated
	extensions
	additional
	empty
}
enum EnumPatientResourceType {
	Patient
}
enum EnumContactPointUse {
	home
	work
	temp
	old
	mobile
}
enum EnumIdentifierUse {
	usual
	official
	temp
	secondary
}
enum EnumContactPointSystem {
	phone
	fax
	email
	pager
	url
	sms
	other
}
enum EnumAddressType {
	postal
	physical
	both
}
enum EnumHumanNameUse {
	usual
	official
	temp
	nickname
	anonymous
	old
	maiden
}
struct Extension {
	optional string valueId
	optional string url
}
struct Element {
	optional string id
	optional list<Extension> extension
}
struct Coding {
	optional string system
	optional string display
	optional string version
	optional list<Extension> extension
	optional string id
	optional string code
}
struct Meta {
	optional list<Coding> security
	optional list<Coding> tag
	optional string versionId
	optional string id
	optional list<Extension> extension
	optional string lastUpdated
	optional list<string> profile
}
struct Resource {
	optional string id
	optional list<Extension> extension
	optional string implicitRules
	optional string language
	optional Meta meta
}
struct Narrative {
	optional list<Extension> extension
	optional EnumNarrativeStatus status
	optional string div
	optional string id
}
struct DomainResource {
	optional string implicitRules
	optional string language
	optional Meta meta
	optional Narrative text
	optional list<Extension> modifierExtension
	optional string id
	optional list<Extension> extension
}
struct Period {
	optional string PeriodId
	optional list<Extension> PeriodExtension
	optional string PeriodStart
	optional string PeriodEnd
}
struct ContactPoint {
	optional string ContactPointId
	optional list<Extension> ContactPointExtension
	optional string ContactPointValue
	optional Period ContactPointPeriod
	optional EnumContactPointSystem ContactPointSystem
	optional EnumContactPointUse ContactPointUse
}
struct Attachment {
	optional string url
	optional string hash
	optional string title
	optional string contentType
	optional string language
	optional string id
	optional list<Extension> extension
	optional string creation
	optional string data
}
struct Address {
	optional list<string> AddressLine
	optional string AddressCity
	optional EnumAddressUse AddressUse
	optional string AddressText
	optional string AddressId
	optional string AddressPostalCode
	optional Period AddressPeriod
	optional EnumAddressType AddressType
	optional string AddressState
	optional string AddressCountry
	optional list<Extension> AddressExtension
	optional string AddressDistrict
}
struct BackboneElement {
	optional string id
	optional list<Extension> extension
	optional list<Extension> modifierExtension
}
struct CodeableConcept {
	optional string id
	optional list<Extension> extension
	optional list<Coding> coding
	optional string text
}
struct HumanName {
	optional list<string> HumanNamePrefix
	optional list<string> HumanNameSuffix
	optional Period HumanNamePeriod
	optional EnumHumanNameUse HumanNameUse
	optional string HumanNameText
	optional string HumanNameId
	optional list<Extension> HumanNameExtension
	optional string HumanNameFamily
	optional list<string> HumanNameGiven
}
struct Reference {
	optional string reference
	optional string display
	optional string id
	optional list<Extension> extension
}
struct Patient_Contact {
	optional list<ContactPoint> telecom
	optional Period period
	optional string id
	optional list<Extension> extension
	optional list<CodeableConcept> relationship
	optional HumanName name
	optional list<Extension> modifierExtension
	optional Address address
	optional EnumPatientContactGender gender
	optional Reference organization
}
struct Identifier {
	optional string IdentifierSystem
	optional list<Extension> IdentifierExtension
	optional string IdentifierId
	optional CodeableConcept IdentifierType
	optional string IdentifierValue
	optional Period IdentifierPeriod
	optional Reference IdentifierAssigner
	optional EnumIdentifierUse IdentifierUse
}
struct Patient_Communication {
	optional CodeableConcept language
	optional list<Extension> extension
	optional list<Extension> modifierExtension
	optional string id
}
struct Patient_Animal {
	optional CodeableConcept genderStatus
	optional string id
	optional list<Extension> extension
	optional list<Extension> modifierExtension
	optional CodeableConcept species
	optional CodeableConcept breed
}
struct Patient_Link {
	optional list<Extension> extension
	optional Reference other
	optional EnumPatientLinkType type
	optional list<Extension> modifierExtension
	optional string id
}
struct Patient {
	optional EnumPatientResourceType resourceType
	optional list<Patient_Contact> contact
	optional Reference managingOrganization
	optional EnumPatientGender gender
	optional string birthDate
	optional list<Extension> extension
	optional string id
	optional list<Attachment> photo
	optional list<Patient_Communication> communication
	optional list<Reference> generalPractitioner
	optional CodeableConcept maritalStatus
	optional list<HumanName> name
	optional string deceasedDateTime
	optional Meta meta
	optional list<Extension> modifierExtension
	optional list<Identifier> identifier
	optional Patient_Animal animal
	optional list<Patient_Link> link
	optional string language
	optional list<ContactPoint> telecom
	optional list<Address> address
	optional string implicitRules
	optional Narrative text
}
