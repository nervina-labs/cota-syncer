package biz

import (
	"encoding/json"
)

type CTMeta struct {
	Id       string
	ver      string
	metadata Metadata
}

type Metadata struct {
	Target string
	Type   string
	Data   string
}

type ClassInfoJson struct {
	CotaId       string
	Version      string
	Name         string
	Symbol       string
	Description  string
	Image        string
	Audio        string
	Video        string
	Model        string
	Schema       string
	Properties   string
	Localization string
}

type IssuerInfoJson struct {
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
}

func ParseMetadata(meta []byte) (isIssuer bool, metadata []byte, err error) {
	var ctMeta CTMeta
	err = json.Unmarshal(meta, &ctMeta)
	if err != nil {
		return
	}
	metadata = []byte(ctMeta.metadata.Data)
	isIssuer = ctMeta.metadata.Type == "issuer"
	return
}
