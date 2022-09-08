package data

import (
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"reflect"
	"testing"
)

func Test_classInfoRepo_ParseClassInfo(t *testing.T) {
	type fields struct {
		data   *Data
		logger *logger.Logger
	}
	type args struct {
		blockNumber uint64
		txIndex     uint32
		classMeta   map[string]any
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantClass biz.ClassInfo
		wantErr   bool
	}{
		{
			name: "should return empty string when characteristic is nil",
			fields: fields{
				data:   nil,
				logger: nil,
			},
			args: args{
				blockNumber: 1000,
				txIndex:     0,
				classMeta:   map[string]any{"cota_id": "0x718a6223d13598926c1e093e82e18b98d148f373", "version": "1", "name": "Kernel", "symbol": "udt"},
			},
			wantClass: biz.ClassInfo{
				BlockNumber:    1000,
				CotaId:         "0x718a6223d13598926c1e093e82e18b98d148f373",
				Version:        "1",
				Name:           "Kernel",
				Symbol:         "udt",
				Description:    "",
				Image:          "",
				Audio:          "",
				Video:          "",
				Model:          "",
				Characteristic: "",
				Properties:     "",
				Localization:   "",
				TxIndex:        0,
			},
			wantErr: false,
		}, {
			name: "should return corresponding characteristic",
			fields: fields{
				data:   nil,
				logger: nil,
			},
			args: args{
				blockNumber: 1000,
				txIndex:     0,
				classMeta: map[string]any{
					"cota_id": "0x718a6223d13598926c1e093e82e18b98d148f373", "version": "1", "name": "Kernel", "symbol": "udt",
					"description": "nice token", "characteristic": [][]string{{"hp", "1"}, {"act", "3"}}, "properties": map[string]any{"power": map[string]any{"value": "1", "value1": "2"}}},
			},
			wantClass: biz.ClassInfo{
				BlockNumber:    1000,
				CotaId:         "0x718a6223d13598926c1e093e82e18b98d148f373",
				Version:        "1",
				Name:           "Kernel",
				Symbol:         "udt",
				Description:    "nice token",
				Image:          "",
				Audio:          "",
				Video:          "",
				Model:          "",
				Characteristic: "[[\"hp\",\"1\"],[\"act\",\"3\"]]",
				Properties:     "{\"power\":{\"value\":\"1\",\"value1\":\"2\"}}",
				Localization:   "",
				TxIndex:        0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := classInfoRepo{
				data:   tt.fields.data,
				logger: tt.fields.logger,
			}
			gotClass, err := repo.ParseClassInfo(tt.args.blockNumber, tt.args.txIndex, tt.args.classMeta)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClassInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClass, tt.wantClass) {
				t.Errorf("ParseClassInfo() gotClass = %v, want %v", gotClass, tt.wantClass)
			}
		})
	}
}
