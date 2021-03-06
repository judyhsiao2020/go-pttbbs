package cmbbs

import (
	"os"
	"sync"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

var (
	origBBSHOME = ""
	IsTest      = false
	TestMutex   sync.Mutex

	TestPASSWDSEM_KEY = 32763

	origPASSWDSEM_KEY = 0
)

func setupTest() {
	SetIsTest()
	cache.SetIsTest()

	origBBSHOME = ptttype.SetBBSHOME("./testcase")

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	err := cache.NewSHM(cache.TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("setupTest: unable to NewSHM: e: %v", err)
		return
	}

	cache.Shm.Reset()

	_ = cache.LoadUHash()
}

func teardownTest() {
	_ = cache.CloseSHM()

	os.Remove("./testcase/.PASSWDS")

	ptttype.SetBBSHOME(origBBSHOME)

	cache.UnsetIsTest()
	UnsetIsTest()
}

func SetIsTest() {
	TestMutex.Lock()
	IsTest = true
	origPASSWDSEM_KEY = ptttype.PASSWDSEM_KEY
	ptttype.PASSWDSEM_KEY = TestPASSWDSEM_KEY

	log.Infof("After set sem: TestPASSWDSEM_KEY: %v ptttype.PASSWDSEM_KEY: %v", TestPASSWDSEM_KEY, ptttype.PASSWDSEM_KEY)
}

func UnsetIsTest() {
	IsTest = false
	TestMutex.Unlock()
}
