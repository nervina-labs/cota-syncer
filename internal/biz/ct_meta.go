package biz

import (
	"encoding/json"
)

type CTMeta struct {
	Id       string   `json:"id"`
	Ver      string   `json:"ver"`
	Metadata MetaData `json:"metadata"`
}

type MetaData struct {
	Target string `json:"target"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}

type ClassInfoJson struct {
	CotaId       string `json:"cota_id"`
	Version      string `json:"version"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Audio        string `json:"audio"`
	Video        string `json:"video"`
	Model        string `json:"model"`
	Schema       string `json:"schema"`
	Properties   string `json:"properties"`
	Localization string `json:"localization"`
}

type IssuerInfoJson struct {
	Version      string `json:"version"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	Description  string `json:"description"`
	Localization string `json:"localization"`
}

type MetaType int

const (
	NotMeta MetaType = iota
	Issuer
	Class
)

func ParseMetadata(meta []byte) (result MetaType, metadata []byte) {
	var ctMeta CTMeta
	result = NotMeta
	if err := json.Unmarshal(meta, &ctMeta); err != nil {
		return
	}
	metaType := ctMeta.Metadata.Type
	if metaType != "issuer" && metaType != "cota" {
		return
	}
	if metaType == "issuer" {
		result = Issuer
	} else if metaType == "cota" {
		result = Class
	}
	metadata = []byte(ctMeta.Metadata.Data)
	return
}
