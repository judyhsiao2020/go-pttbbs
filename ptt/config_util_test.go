package ptt

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
)

func TestInitConfig(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer cmbbs.PasswdDestroy()
			if err := InitConfig(); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
