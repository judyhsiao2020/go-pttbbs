package cache

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func TestNewSHM(t *testing.T) {
	shmSetupTest()
	defer shmTeardownTest()

	type args struct {
		key          types.Key_t
		isUseHugeTlb bool
		isCreate     bool
	}
	tests := []struct {
		name        string
		args        args
		wantVersion int32
		wantSize    int32
		isClose     bool
		wantErr     bool
	}{
		// TODO: Add test caseShm.
		{
			args: args{
				key:          TestShmKey,
				isUseHugeTlb: false,
				isCreate:     true,
			},
			isClose:     false,
			wantVersion: SHM_VERSION,
			wantSize:    int32(SHM_RAW_SZ),
		},
		{
			args: args{
				key:          TestShmKey,
				isUseHugeTlb: false,
				isCreate:     false,
			},
			isClose:     true,
			wantVersion: SHM_VERSION,
			wantSize:    int32(SHM_RAW_SZ),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSHM(tt.args.key, tt.args.isUseHugeTlb, tt.args.isCreate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSHM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			Shm.Reset()

			if !reflect.DeepEqual(Shm.Raw.Version, tt.wantVersion) {
				t.Errorf("NewSHM() version: %v expected: %v", Shm.Raw.Version, tt.wantVersion)
			}

			if !reflect.DeepEqual(Shm.Raw.Size, tt.wantSize) {
				t.Errorf("NewSHM() size: %v expected :%v", Shm.Raw.Size, tt.wantSize)
			}

			if tt.isClose {
				CloseSHM()
			} else {
				Shm = nil
			}
		})
	}
}

func TestSHM_ReadAt(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, false, true)
	if err != nil {
		log.Errorf("cache.TestSHM_ReadAt: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	// test1
	in1 := byte(1)

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.UTMPNeedSort),
		unsafe.Sizeof(Shm.Raw.UTMPNeedSort),
		unsafe.Pointer(&in1),
	)

	out1 := byte(0)
	want1 := byte(1)

	// test2
	in2 := [ptttype.MAX_USERS]ptttype.UserID_t{}
	copy(in2[0][:], []byte("test"))
	copy(in2[1][:], []byte("test1"))
	copy(in2[2][:], []byte("test2"))
	copy(in2[3][:], []byte("SYSOP"))
	copy(in2[4][:], []byte("test4"))

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Userid),
		unsafe.Sizeof(Shm.Raw.Userid),
		unsafe.Pointer(&in2),
	)

	out2 := ptttype.UserID_t{}
	want2 := ptttype.UserID_t{}
	copy(want2[:], []byte("SYSOP"))

	// test3
	in3 := &ptttype.MsgQueueRaw{}
	copy(in3.UserID[:], []byte("test33"))

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.LoginMsg)+unsafe.Offsetof(Shm.Raw.LoginMsg.UserID),
		unsafe.Sizeof(Shm.Raw.LoginMsg.UserID),
		unsafe.Pointer(&in3.UserID),
	)

	out3 := &ptttype.MsgQueueRaw{}
	want3 := ptttype.UserID_t{}
	copy(want3[:], []byte("test33"))

	type fields struct {
		Shmid   int
		IsNew   bool
		Shmaddr unsafe.Pointer
		Size    types.Size_t
		offset  uintptr
		SHMRaw  SHMRaw
	}
	type args struct {
		offsetOfSHMRawComponent uintptr
		size                    uintptr
		outptr                  unsafe.Pointer
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		out      interface{}
		expected interface{}
	}{
		// TODO: Add test caseShm.
		{
			name: "UTMPNeedSort (byte)",
			args: args{
				unsafe.Offsetof(Shm.Raw.UTMPNeedSort),
				unsafe.Sizeof(Shm.Raw.UTMPNeedSort),
				unsafe.Pointer(&out1),
			},
			out:      &out1,
			expected: &want1,
		},
		{
			name: "Userid[3]",
			args: args{
				unsafe.Offsetof(Shm.Raw.Userid) + (ptttype.USER_ID_SZ)*3,
				unsafe.Sizeof(Shm.Raw.Userid[3]),
				unsafe.Pointer(&out2),
			},
			out:      &out2,
			expected: &want2,
		},
		{
			name: "LoginMsg.UserID",
			args: args{
				unsafe.Offsetof(Shm.Raw.LoginMsg) + unsafe.Offsetof(Shm.Raw.LoginMsg.UserID),
				unsafe.Sizeof(Shm.Raw.LoginMsg.UserID),
				unsafe.Pointer(&out3.UserID),
			},
			out:      &out3.UserID,
			expected: &want3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Shm.ReadAt(tt.args.offsetOfSHMRawComponent, tt.args.size, tt.args.outptr)

			if !reflect.DeepEqual(tt.out, tt.expected) {
				t.Errorf("SHM.ReadAt() out: %v expected: %v", tt.out, tt.expected)
			}
		})
	}
}

func TestSHM_WriteAt(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, false, true)
	if err != nil {
		log.Errorf("cache.TestSHM_WriteAt: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	// in1
	in1 := byte(1)

	out1 := byte(0)
	out1ptr := unsafe.Pointer(&out1)
	want1 := byte(1)

	// in2
	in2 := [ptttype.MAX_USERS]ptttype.UserID_t{}
	copy(in2[0][:], []byte("test"))
	copy(in2[1][:], []byte("test1"))
	copy(in2[2][:], []byte("test2"))
	copy(in2[3][:], []byte("SYSOP"))
	copy(in2[4][:], []byte("test4"))

	out2 := ptttype.UserID_t{}
	out2ptr := unsafe.Pointer(&out2)
	want2 := ptttype.UserID_t{}

	copy(want2[:], []byte("SYSOP"))

	//in3
	in3 := &ptttype.MsgQueueRaw{}
	copy(in3.UserID[:], []byte("test33"))

	out3 := ptttype.UserID_t{}
	out3ptr := unsafe.Pointer(&out3)
	want3 := ptttype.UserID_t{}
	copy(want3[:], []byte("test33"))

	//in4
	in4 := int32(100)

	out4 := [ptttype.MAX_USERS]int32{}
	out4ptr := unsafe.Pointer(&out4)
	want4 := [ptttype.MAX_USERS]int32{}
	want4[4] = 100

	type fields struct {
		Shmid   int
		IsNew   bool
		Shmaddr unsafe.Pointer
		Size    types.Size_t
		offset  uintptr
		SHMRaw  SHMRaw
	}
	type args struct {
		offsetOfSHMRawComponent uintptr
		size                    uintptr
		inptr                   unsafe.Pointer
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		readOffset uintptr
		readSize   uintptr
		read       unsafe.Pointer
		out        interface{}
		expected   interface{}
	}{
		// TODO: Add test caseShm.
		{
			args: args{
				unsafe.Offsetof(Shm.Raw.UTMPNeedSort),
				unsafe.Sizeof(Shm.Raw.UTMPNeedSort),
				unsafe.Pointer(&in1),
			},
			readOffset: unsafe.Offsetof(Shm.Raw.UTMPNeedSort),
			readSize:   unsafe.Sizeof(Shm.Raw.UTMPNeedSort),
			read:       out1ptr,
			out:        &out1,
			expected:   &want1,
		},
		{
			name: "with 2d-list [][]interface{}, the unit size is unsafe.Sizeof []interface{}",
			args: args{
				unsafe.Offsetof(Shm.Raw.Userid) + unsafe.Sizeof(Shm.Raw.Userid[0])*3,
				unsafe.Sizeof(Shm.Raw.Userid[0]),
				unsafe.Pointer(&in2[3]),
			},
			readOffset: unsafe.Offsetof(Shm.Raw.Userid) + unsafe.Sizeof(Shm.Raw.Userid[0])*3,
			readSize:   unsafe.Sizeof(Shm.Raw.Userid[0]),
			read:       out2ptr,
			out:        &out2,
			expected:   &want2,
		},
		{
			name: "with nested-data",
			args: args{
				unsafe.Offsetof(Shm.Raw.LoginMsg) + unsafe.Offsetof(Shm.Raw.LoginMsg.UserID),
				unsafe.Sizeof(Shm.Raw.LoginMsg.UserID),
				unsafe.Pointer(&in3.UserID),
			},
			readOffset: unsafe.Offsetof(Shm.Raw.LoginMsg) + unsafe.Offsetof(Shm.Raw.LoginMsg.UserID),
			readSize:   unsafe.Sizeof(Shm.Raw.LoginMsg.UserID),
			read:       out3ptr,
			out:        &out3,
			expected:   &want3,
		},
		{
			name: "with list, remember to have unit-size",
			args: args{
				unsafe.Offsetof(Shm.Raw.Money) + types.INT32_SZ*4,
				types.INT32_SZ,
				unsafe.Pointer(&in4),
			},
			readOffset: unsafe.Offsetof(Shm.Raw.Money),
			readSize:   unsafe.Sizeof(Shm.Raw.Money),
			read:       out4ptr,
			out:        &out4,
			expected:   &want4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Shm.WriteAt(
				tt.args.offsetOfSHMRawComponent,
				tt.args.size,
				tt.args.inptr,
			)

			Shm.ReadAt(
				tt.readOffset,
				tt.readSize,
				tt.read,
			)
			if !reflect.DeepEqual(tt.out, tt.expected) {
				t.Errorf("SHM.WriteAt() out: %v expected: %v", tt.out, tt.expected)
			}
		})
	}
}
