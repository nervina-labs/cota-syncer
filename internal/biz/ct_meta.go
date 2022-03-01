package biz

import (
	"encoding/json"
	"github.com/pkg/errors"
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
	if err = json.Unmarshal(meta, &ctMeta); err != nil {
		return
	}
	metadata = []byte(ctMeta.metadata.Data)
	metaType := ctMeta.metadata.Type
	if metaType != "issuer" && metaType != "class" {
		err = errors.New("Cota metadata type error")
	} else {
		isIssuer = metaType == "issuer"
	}
	return
}
