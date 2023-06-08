package models_test

import (
	"fmt"
	"testing"

	"github.com/goto/meteor/models"
	assetsv1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/stretchr/testify/assert"
)

func TestNewURN(t *testing.T) {
	testCases := []struct {
		service  string
		scope    string
		kind     string
		id       string
		expected string
	}{
		{
			"metabase", "main-dashboard", "collection", "123",
			"urn:metabase:main-dashboard:collection:123",
		},
		{
			"bigquery", "p-godata-id", "table", "p-godata-id:mydataset.mytable",
			"urn:bigquery:p-godata-id:table:p-godata-id:mydataset.mytable",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("should return expected urn (#%d)", i+1), func(t *testing.T) {
			actual := models.NewURN(tc.service, tc.scope, tc.kind, tc.id)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestToJSON(t *testing.T) {
	type args struct {
		asset *assetsv1beta2.Asset
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "should return the json representation of the asset",
			args: args{
				asset: &assetsv1beta2.Asset{
					Name: "test",
				},
			},
			want:    []byte(`{"urn":"", "name":"test", "service":"", "type":"", "url":"", "description":"", "data":null, "owners":[], "lineage":null, "labels":{}, "event":null, "create_time":null, "update_time":null}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := models.ToJSON(tt.args.asset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
