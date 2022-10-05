package models

type Dictionary struct {
	BankName    string `json:"bankName"`
	LogoURL     string `json:"logoURL"`
	Schema      string `json:"schema"`
	PackageName string `json:"package_name"`
}

type Data struct {
	Version    string
	Dictionary []Dictionary
}
