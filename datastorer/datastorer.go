package datastorer

import (
	"os/exec"
	"strconv"
	"sync"

)

type DataStorer interface {
	Set(ID uint, val any) error
	Get(ID uint) (any, error)
}

type BashCommandData struct {
	FilePath string
	Cmd *exec.Cmd
}

type BashDataStorer struct {
	Stored map[uint]BashCommandData
	StoredMtx sync.Mutex
}

type BashDataStorerMismatchType struct {}

func (bdsm BashDataStorerMismatchType) Error() string{
	return "Type of value and stored missmatch"
}

type BashDataStorerNotFoundError struct {
	ID uint
}

func (bdsnfe BashDataStorerNotFoundError) Error() string {
	return "Not found bash command with ID = " + strconv.FormatUint(uint64(bdsnfe.ID),10)
}

func NewBashDataSorer() *BashDataStorer {
	return &BashDataStorer{ Stored: make(map[uint]BashCommandData)}
}

func (bds *BashDataStorer) Set(ID uint, val any) error {
	bds.StoredMtx.Lock()
	defer bds.StoredMtx.Unlock()
	cmd, ok := val.(BashCommandData)
	if !ok {
		return BashDataStorerMismatchType{}
	}
	bds.Stored[ID] = cmd
	return nil
}

func (bds *BashDataStorer) Get(ID uint) (any, error) {
	bds.StoredMtx.Lock()
	defer bds.StoredMtx.Unlock()
	val, ok := bds.Stored[ID]
	if !ok {
		return nil, BashDataStorerNotFoundError{ID: ID}
	}
	return val, nil
}