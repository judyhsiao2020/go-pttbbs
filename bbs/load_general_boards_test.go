package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralBoards(t *testing.T) {
	//setupTest moves in for-loop
	//teardownTest moves in for-loop
	type args struct {
		userID      string
		startIdxStr string
		nBoards     int
		keyword     string
	}
	tests := []struct {
		name               string
		args               args
		expectedSummary    []*BoardSummary
		expectedNextIdxStr string
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args:               args{userID: "SYSOP", startIdxStr: "", nBoards: 4},
			expectedSummary:    []*BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
			expectedNextIdxStr: "8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			gotSummary, gotNextIdx, err := LoadGeneralBoards(tt.args.userID, tt.args.startIdxStr, tt.args.nBoards, tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for idx, each := range gotSummary {
				if idx >= len(tt.expectedSummary) {
					t.Errorf("LoadGeneralBoards: (%v/%v): %v", idx, len(gotSummary), each)
					continue
				}
				types.TDeepEqual(t, each, tt.expectedSummary[idx])
			}
			if gotNextIdx != tt.expectedNextIdxStr {
				t.Errorf("LoadGeneralBoards() gotNextIdx = %v, want %v", gotNextIdx, tt.expectedNextIdxStr)
			}
		})
	}
}
