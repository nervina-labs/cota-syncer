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

func ParseMetadata(meta []byte) (CTMeta, error) {
	var ctMeta CTMeta
	if err := json.Unmarshal(meta, &ctMeta); err != nil {
		return ctMeta, err
	}
	metaType := ctMeta.Metadata.Type
	if metaType != "issuer" && metaType != "cota" {
		return ctMeta, errors.New("invalid meta type")
	}
	return ctMeta, nil
}
