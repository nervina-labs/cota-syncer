package biz

import (
	"encoding/json"
	"errors"
)

type CTMeta struct {
	Id       string   `json:"id"`
	Ver      string   `json:"ver"`
	Metadata MetaData `json:"metadata"`
}

type MetaData struct {
	Target string         `json:"target"`
	Type   string         `json:"type"`
	Data   map[string]any `json:"data"`
}

type Localization struct {
	Uri     string   `json:"uri,omitempty" mapstructure:",omitempty"`
	Default string   `json:"default,omitempty" mapstructure:",omitempty"`
	Locales []string `json:"locales,omitempty" mapstructure:",omitempty"`
}

type ClassInfoJson struct {
	CotaId         string         `json:"cota_id" mapstructure:"cota_id,omitempty"`
	Version        string         `json:"version" mapstructure:",omitempty"`
	Name           string         `json:"name" mapstructure:",omitempty"`
	Symbol         string         `json:"symbol" mapstructure:",omitempty"`
	Description    string         `json:"description" mapstructure:",omitempty"`
	Image          string         `json:"image" mapstructure:",omitempty"`
	Audio          string         `json:"audio"`
	Video          string         `json:"video" mapstructure:",omitempty"`
	Model          string         `json:"model" mapstructure:",omitempty"`
	Characteristic [][]any        `json:"characteristic" mapstructure:",omitempty"`
	Properties     map[string]any `json:"properties" mapstructure:",omitempty"`
	Localization   Localization   `json:"localization" mapstructure:",omitempty"`
}

type IssuerInfoJson struct {
	Version      string       `json:"version" mapstructure:",omitempty"`
	Name         string       `json:"name" mapstructure:",omitempty"`
	Avatar       string       `json:"avatar" mapstructure:",omitempty"`
	Description  string       `json:"description" mapstructure:",omitempty"`
	Localization Localization `json:"localization" mapstructure:",omitempty"`
}

type SubKeyInfoJson struct {
	PubKey       string `json:"pub_key" mapstructure:"pub_key,omitempty"`
	CredentialId string `json:"credential_id" mapstructure:"credential_id,omitempty"`
	Alg          string `json:"alg" mapstructure:"alg,omitempty"`
	FrontEnd     string `json:"front_end" mapstructure:"front_end,omitempty"`
}

type JoyIDInfoJson struct {
	Version      string           `json:"version" mapstructure:",omitempty"`
	PubKey       string           `json:"pub_key" mapstructure:"pub_key,omitempty"`
	CredentialId string           `json:"credential_id" mapstructure:"credential_id,omitempty"`
	Alg          string           `json:"alg" mapstructure:"alg,omitempty"`
	CotaCellId   string           `json:"cota_cell_id" mapstructure:"cota_cell_id,omitempty"`
	Name         string           `json:"name" mapstructure:",omitempty"`
	Avatar       string           `json:"avatar" mapstructure:",omitempty"`
	Description  string           `json:"description" mapstructure:",omitempty"`
	Extension    string           `json:"extension" mapstructure:",omitempty"`
	FrontEnd     string           `json:"front_end" mapstructure:"front_end,omitempty"`
	SubKeys      []SubKeyInfoJson `json:"sub_keys" mapstructure:"sub_keys,omitempty"`
}

type MetaType int

func ParseMetadata(meta []byte) (CTMeta, error) {
	var ctMeta CTMeta
	if err := json.Unmarshal(meta, &ctMeta); err != nil {
		return ctMeta, err
	}
	metaType := ctMeta.Metadata.Type
	if metaType != "issuer" && metaType != "cota" && metaType != "joy_id" {
		return ctMeta, errors.New("invalid meta type")
	}
	return ctMeta, nil
}
