package ftc_client

type Applications struct {
	Apps []Application
}

type Application struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	EntityID      string              `json:"entity_id"`
	SsoUrl        string              `json:"sso_url"`
	SloUrl        string              `json:"slo_url"`
	RealmID       string              `json:"realm_id"`
	Type          int                 `json:"type"`
	Prefix        string              `json:"prefix"`
	BrandingID    string              `json:"branding_id"`
	TTL           int                 `json:"ttl"`
	AttrMapping   interface{}         `json:"attr_mapping"`
	SigningCertID string              `json:"signing_cert_id"`
	SpEntityID    string              `json:"sp_entity_id"`
	SpAcsUrl      string              `json:"sp_acs_url"`
	SpSloUrl      string              `json:"sp_slo_url"`
	SpNameID      string              `json:"sp_name_id"`
	SpSigningCert string              `json:"sp_signing_cert"`
	UserSources   []UserSourceElement `json:"user_sources"`
}

type SamlParams struct {
	SigningCertID string `json:"signing_cert_id"`
	SpEntityID    string `json:"sp_entity_id"`
	SpAcsUrl      string `json:"sp_acs_url"`
	SpSloUrl      string `json:"sp_slo_url"`
	SpNameID      string `json:"sp_name_id"`
	SpSigningCert string `json:"sp_signing_cert"`
}

type UserSourceElement struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Prefix string `json:"prefix"`
}

type AppUserMapping struct {
	ApplicationID string `json:"application_id"`
	UserSourceID  string `json:"user_source_id"`
}

type UserSourceList struct {
	UserSourceIDs []string `json:"user_source_ids"`
}

type UserSources struct {
	UserSources []UserSource
}

type UserSource struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	RealmID           string          `json:"realm_id"`
	Prefix            string          `json:"prefix"`
	Type              int             `json:"type"`
	EntityID          string          `json:"entity_id"`
	LoginUrl          string          `json:"login_url"`
	LogoutUrl         string          `json:"logout_url"`
	AuthUri           string          `json:"auth_uri"`
	TokenUri          string          `json:"token_uri"`
	UserInfoUri       string          `json:"userinfo_uri"`
	LogoutUri         string          `json:"logout_uri"`
	Issuer            string          `json:"issuer"`
	ClientID          string          `json:"client_id"`
	ClientSecret      string          `json:"client_secret"`
	SigningCert       string          `json:"signing_cert"`
	PostBinding       bool            `json:"post_binding"`
	IncludeSubject    bool            `json:"include_subject"`
	AttrMapping       interface{}     `json:"attr_mapping"`
	UsernameAssertion string          `json:"username_assertion"`
	LoginHint         string          `json:"login_hint"`
	FQDN              string          `json:"fqdn"`
	ProxySP           ProxySP         `json:"proxy_sp"`
	Domains           []DomainElement `json:"domains"`
}

type ProxySP struct {
	Prefix                string `json:"prefix"`
	EntityID              string `json:"entity_id"`
	AcsUrl                string `json:"acs_url"`
	SloUrl                string `json:"slo_url"`
	SsoUrl                string `json:"sso_url"`
	CallbackUrl           string `json:"callback_url"`
	PostLogoutRedirectUrl string `json:"post_logout_redirect_uri"`
	OidcLoginUrl          string `json:"oidc_login_url"`
}

type DomainElement struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type UserSourceDomainMapping struct {
	UserSourceID string `json:"user_source_id"`
	DomainID     string `json:"domain_id"`
}

type Domain struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	RealmID      string `json:"realm_id"`
	UserSourceID string `json:"user_source_id"`
}

type Realm struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
