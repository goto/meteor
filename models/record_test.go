package models_test

import (
	"reflect"
	"testing"

	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
)

func TestNewRecord(t *testing.T) {
	type args struct {
		data *v1beta2.Asset
	}
	tests := []struct {
		name string
		args args
		want models.Record
	}{
		{
			name: "should return a new record",
			args: args{
				data: &v1beta2.Asset{},
			},
			want: models.NewRecord(&v1beta2.Asset{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.NewRecord(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_Data(t *testing.T) {
	type fields struct {
		data *v1beta2.Asset
	}
	tests := []struct {
		name   string
		fields fields
		want   *v1beta2.Asset
	}{
		{
			name: "should return the record data",
			fields: fields{
				data: &v1beta2.Asset{
					Name: "test",
				},
			},
			want: &v1beta2.Asset{
				Name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := models.NewRecord(tt.fields.data)
			if got := r.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}
