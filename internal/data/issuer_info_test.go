package data

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"reflect"
	"testing"
)

func Test_issuerInfoRepo_ParseIssuerInfo(t *testing.T) {
	type fields struct {
		data   *Data
		logger *logger.Logger
	}
	type args struct {
		blockNumber uint64
		txIndex     uint32
		lockScript  *ckbTypes.Script
		issuerMeta  map[string]any
	}
	scriptArgs, _ := hexutil.Decode("0xf9910364e0ca81a0e074f3aa42fe78cfcc880da6")
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantIssuer biz.IssuerInfo
		wantErr    bool
	}{
		{
			name: "should return empty string when localization is nil",
			fields: fields{
				data:   nil,
				logger: nil,
			},
			args: args{
				blockNumber: 100,
				txIndex:     0,
				lockScript: &ckbTypes.Script{
					CodeHash: ckbTypes.HexToHash("0x89cd8003a0eaf8e65e0c31525b7d1d5c1becefd2ea75bb4cff87810ae37764d8"),
					HashType: ckbTypes.HashTypeType,
					Args:     scriptArgs,
				},
				issuerMeta: map[string]any{"version": "0", "name": "kevin", "avatar": "https://i.loli.net/2021/04/28/ZCQPoxztsVHdNA9.jpg", "description": "just a man"},
			},
			wantIssuer: biz.IssuerInfo{
				BlockNumber:  100,
				LockHash:     "167fdfb476eb1f76f1dc2b2fe78f3afcf1185906daee966950317852d0c976db",
				Version:      "0",
				Name:         "kevin",
				Avatar:       "https://i.loli.net/2021/04/28/ZCQPoxztsVHdNA9.jpg",
				Description:  "just a man",
				Localization: "",
				TxIndex:      0,
			},
			wantErr: false,
		}, {
			name: "should return corresponding localization",
			fields: fields{
				data:   nil,
				logger: nil,
			},
			args: args{
				blockNumber: 100,
				txIndex:     0,
				lockScript: &ckbTypes.Script{
					CodeHash: ckbTypes.HexToHash("0x89cd8003a0eaf8e65e0c31525b7d1d5c1becefd2ea75bb4cff87810ae37764d8"),
					HashType: ckbTypes.HashTypeType,
					Args:     scriptArgs,
				},
				issuerMeta: map[string]any{
					"version": "0", "name": "kevin", "avatar": "https://i.loli.net/2021/04/28/ZCQPoxztsVHdNA9.jpg",
					"description":  "just a man",
					"localization": map[string]any{"uri": "https://abc.com", "default": "zh", "locales": []string{"en", "zh"}},
				},
			},
			wantIssuer: biz.IssuerInfo{
				BlockNumber:  100,
				LockHash:     "167fdfb476eb1f76f1dc2b2fe78f3afcf1185906daee966950317852d0c976db",
				Version:      "0",
				Name:         "kevin",
				Avatar:       "https://i.loli.net/2021/04/28/ZCQPoxztsVHdNA9.jpg",
				Description:  "just a man",
				Localization: "{\"uri\":\"https://abc.com\",\"default\":\"zh\",\"locales\":[\"en\",\"zh\"]}",
				TxIndex:      0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := issuerInfoRepo{
				data:   tt.fields.data,
				logger: tt.fields.logger,
			}
			gotIssuer, err := repo.ParseIssuerInfo(tt.args.blockNumber, tt.args.txIndex, tt.args.lockScript, tt.args.issuerMeta)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIssuerInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIssuer, tt.wantIssuer) {
				t.Errorf("ParseIssuerInfo() gotIssuer = %v, want %v", gotIssuer, tt.wantIssuer)
			}
		})
	}
}
