package handlers

type accessControlList struct {
	Grant []Grant `xml:"Grant,omitempty"`
}
type canonicalUser struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName,omitempty"`
}

// AccessControlPolicy <AccessControlPolicy>
//  <Owner>
//    <ID>75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a</ID>
//    <DisplayName>CustomersName@amazon.com</DisplayName>
//  </Owner>
//  <AccessControlList>
//    <Grant>
//      <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
//			xsi:type="CanonicalUser">
//        <ID>75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a</ID>
//        <DisplayName>CustomersName@amazon.com</DisplayName>
//      </Grantee>
//      <Permission>FULL_CONTROL</Permission>
//    </Grant>
//  </AccessControlList>
//</AccessControlPolicy>
type AccessControlPolicy struct {
	Owner             canonicalUser     `xml:"Owner"`
	AccessControlList accessControlList `xml:"AccessControlList"`
}

//Grant grant
type Grant struct {
	Grantee    Grantee    `xml:"Grantee"`
	Permission Permission `xml:"Permission"`
}

//Grantee grant
type Grantee struct {
	XMLNS       string `xml:"xmlns:xsi,attr"`
	XMLXSI      string `xml:"xsi:type,attr"`
	Type        string `xml:"Type"`
	ID          string `xml:"ID,omitempty"`
	DisplayName string `xml:"DisplayName,omitempty"`
	URI         string `xml:"URI,omitempty"`
}

// Permission May be one of READ, WRITE, READ_ACP, WRITE_ACP, FULL_CONTROL
type Permission string
