// Code generated by MockGen. DO NOT EDIT.
// Source: ./blockdb/blockdb.go

// Package blockdb is a generated GoMock package.
package blockdb

import (
	reflect "reflect"

	common "chainmaker.org/chainmaker/pb-go/v2/common"
	store "chainmaker.org/chainmaker/pb-go/v2/store"
	serialization "chainmaker.org/chainmaker/store/v2/serialization"
	gomock "github.com/golang/mock/gomock"
)

// MockBlockDB is a mock of BlockDB interface.
type MockBlockDB struct {
	ctrl     *gomock.Controller
	recorder *MockBlockDBMockRecorder
}

// MockBlockDBMockRecorder is the mock recorder for MockBlockDB.
type MockBlockDBMockRecorder struct {
	mock *MockBlockDB
}

// NewMockBlockDB creates a new mock instance.
func NewMockBlockDB(ctrl *gomock.Controller) *MockBlockDB {
	mock := &MockBlockDB{ctrl: ctrl}
	mock.recorder = &MockBlockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockDB) EXPECT() *MockBlockDBMockRecorder {
	return m.recorder
}

// BlockExists mocks base method.
func (m *MockBlockDB) BlockExists(blockHash []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockExists", blockHash)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockExists indicates an expected call of BlockExists.
func (mr *MockBlockDBMockRecorder) BlockExists(blockHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockExists", reflect.TypeOf((*MockBlockDB)(nil).BlockExists), blockHash)
}

// Close mocks base method.
func (m *MockBlockDB) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockBlockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBlockDB)(nil).Close))
}

// CommitBlock mocks base method.
func (m *MockBlockDB) CommitBlock(blockInfo *serialization.BlockWithSerializedInfo, isCache bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitBlock", blockInfo, isCache)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitBlock indicates an expected call of CommitBlock.
func (mr *MockBlockDBMockRecorder) CommitBlock(blockInfo, isCache interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitBlock", reflect.TypeOf((*MockBlockDB)(nil).CommitBlock), blockInfo, isCache)
}

// GetArchivedPivot mocks base method.
func (m *MockBlockDB) GetArchivedPivot() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArchivedPivot")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArchivedPivot indicates an expected call of GetArchivedPivot.
func (mr *MockBlockDBMockRecorder) GetArchivedPivot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArchivedPivot", reflect.TypeOf((*MockBlockDB)(nil).GetArchivedPivot))
}

// GetBlock mocks base method.
func (m *MockBlockDB) GetBlock(height uint64) (*common.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", height)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockBlockDBMockRecorder) GetBlock(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockBlockDB)(nil).GetBlock), height)
}

// GetBlockByHash mocks base method.
func (m *MockBlockDB) GetBlockByHash(blockHash []byte) (*common.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByHash", blockHash)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByHash indicates an expected call of GetBlockByHash.
func (mr *MockBlockDBMockRecorder) GetBlockByHash(blockHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByHash", reflect.TypeOf((*MockBlockDB)(nil).GetBlockByHash), blockHash)
}

// GetBlockByTx mocks base method.
func (m *MockBlockDB) GetBlockByTx(txId string) (*common.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByTx", txId)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByTx indicates an expected call of GetBlockByTx.
func (mr *MockBlockDBMockRecorder) GetBlockByTx(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByTx", reflect.TypeOf((*MockBlockDB)(nil).GetBlockByTx), txId)
}

// GetBlockHeaderByHeight mocks base method.
func (m *MockBlockDB) GetBlockHeaderByHeight(height uint64) (*common.BlockHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeaderByHeight", height)
	ret0, _ := ret[0].(*common.BlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockHeaderByHeight indicates an expected call of GetBlockHeaderByHeight.
func (mr *MockBlockDBMockRecorder) GetBlockHeaderByHeight(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeaderByHeight", reflect.TypeOf((*MockBlockDB)(nil).GetBlockHeaderByHeight), height)
}

// GetBlockIndex mocks base method.
func (m *MockBlockDB) GetBlockIndex(height uint64) (*store.StoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockIndex", height)
	ret0, _ := ret[0].(*store.StoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockIndex indicates an expected call of GetBlockIndex.
func (mr *MockBlockDBMockRecorder) GetBlockIndex(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockIndex", reflect.TypeOf((*MockBlockDB)(nil).GetBlockIndex), height)
}

// GetBlockMetaIndex mocks base method.
func (m *MockBlockDB) GetBlockMetaIndex(height uint64) (*store.StoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockMetaIndex", height)
	ret0, _ := ret[0].(*store.StoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockMetaIndex indicates an expected call of GetBlockMetaIndex.
func (mr *MockBlockDBMockRecorder) GetBlockMetaIndex(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockMetaIndex", reflect.TypeOf((*MockBlockDB)(nil).GetBlockMetaIndex), height)
}

// GetFilteredBlock mocks base method.
func (m *MockBlockDB) GetFilteredBlock(height uint64) (*store.SerializedBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilteredBlock", height)
	ret0, _ := ret[0].(*store.SerializedBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilteredBlock indicates an expected call of GetFilteredBlock.
func (mr *MockBlockDBMockRecorder) GetFilteredBlock(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilteredBlock", reflect.TypeOf((*MockBlockDB)(nil).GetFilteredBlock), height)
}

// GetHeightByHash mocks base method.
func (m *MockBlockDB) GetHeightByHash(blockHash []byte) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeightByHash", blockHash)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHeightByHash indicates an expected call of GetHeightByHash.
func (mr *MockBlockDBMockRecorder) GetHeightByHash(blockHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeightByHash", reflect.TypeOf((*MockBlockDB)(nil).GetHeightByHash), blockHash)
}

// GetLastBlock mocks base method.
func (m *MockBlockDB) GetLastBlock() (*common.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastBlock")
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastBlock indicates an expected call of GetLastBlock.
func (mr *MockBlockDBMockRecorder) GetLastBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastBlock", reflect.TypeOf((*MockBlockDB)(nil).GetLastBlock))
}

// GetLastConfigBlock mocks base method.
func (m *MockBlockDB) GetLastConfigBlock() (*common.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastConfigBlock")
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastConfigBlock indicates an expected call of GetLastConfigBlock.
func (mr *MockBlockDBMockRecorder) GetLastConfigBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastConfigBlock", reflect.TypeOf((*MockBlockDB)(nil).GetLastConfigBlock))
}

// GetLastConfigBlockHeight mocks base method.
func (m *MockBlockDB) GetLastConfigBlockHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastConfigBlockHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastConfigBlockHeight indicates an expected call of GetLastConfigBlockHeight.
func (mr *MockBlockDBMockRecorder) GetLastConfigBlockHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastConfigBlockHeight", reflect.TypeOf((*MockBlockDB)(nil).GetLastConfigBlockHeight))
}

// GetLastSavepoint mocks base method.
func (m *MockBlockDB) GetLastSavepoint() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastSavepoint")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastSavepoint indicates an expected call of GetLastSavepoint.
func (mr *MockBlockDBMockRecorder) GetLastSavepoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastSavepoint", reflect.TypeOf((*MockBlockDB)(nil).GetLastSavepoint))
}

// GetTx mocks base method.
func (m *MockBlockDB) GetTx(txId string) (*common.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", txId)
	ret0, _ := ret[0].(*common.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx.
func (mr *MockBlockDBMockRecorder) GetTx(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockBlockDB)(nil).GetTx), txId)
}

// GetTxConfirmedTime mocks base method.
func (m *MockBlockDB) GetTxConfirmedTime(txId string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxConfirmedTime", txId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxConfirmedTime indicates an expected call of GetTxConfirmedTime.
func (mr *MockBlockDBMockRecorder) GetTxConfirmedTime(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxConfirmedTime", reflect.TypeOf((*MockBlockDB)(nil).GetTxConfirmedTime), txId)
}

// GetTxHeight mocks base method.
func (m *MockBlockDB) GetTxHeight(txId string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxHeight", txId)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxHeight indicates an expected call of GetTxHeight.
func (mr *MockBlockDBMockRecorder) GetTxHeight(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxHeight", reflect.TypeOf((*MockBlockDB)(nil).GetTxHeight), txId)
}

// GetTxIndex mocks base method.
func (m *MockBlockDB) GetTxIndex(txId string) (*store.StoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxIndex", txId)
	ret0, _ := ret[0].(*store.StoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxIndex indicates an expected call of GetTxIndex.
func (mr *MockBlockDBMockRecorder) GetTxIndex(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxIndex", reflect.TypeOf((*MockBlockDB)(nil).GetTxIndex), txId)
}

// GetTxInfoOnly mocks base method.
func (m *MockBlockDB) GetTxInfoOnly(txId string) (*store.TransactionStoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxInfoOnly", txId)
	ret0, _ := ret[0].(*store.TransactionStoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxInfoOnly indicates an expected call of GetTxInfoOnly.
func (mr *MockBlockDBMockRecorder) GetTxInfoOnly(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxInfoOnly", reflect.TypeOf((*MockBlockDB)(nil).GetTxInfoOnly), txId)
}

// GetTxWithBlockInfo mocks base method.
func (m *MockBlockDB) GetTxWithBlockInfo(txId string) (*store.TransactionStoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxWithBlockInfo", txId)
	ret0, _ := ret[0].(*store.TransactionStoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxWithBlockInfo indicates an expected call of GetTxWithBlockInfo.
func (mr *MockBlockDBMockRecorder) GetTxWithBlockInfo(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxWithBlockInfo", reflect.TypeOf((*MockBlockDB)(nil).GetTxWithBlockInfo), txId)
}

// InitGenesis mocks base method.
func (m *MockBlockDB) InitGenesis(genesisBlock *serialization.BlockWithSerializedInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitGenesis", genesisBlock)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitGenesis indicates an expected call of InitGenesis.
func (mr *MockBlockDBMockRecorder) InitGenesis(genesisBlock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitGenesis", reflect.TypeOf((*MockBlockDB)(nil).InitGenesis), genesisBlock)
}

// RestoreBlocks mocks base method.
func (m *MockBlockDB) RestoreBlocks(blockInfos []*serialization.BlockWithSerializedInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreBlocks", blockInfos)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreBlocks indicates an expected call of RestoreBlocks.
func (mr *MockBlockDBMockRecorder) RestoreBlocks(blockInfos interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreBlocks", reflect.TypeOf((*MockBlockDB)(nil).RestoreBlocks), blockInfos)
}

// ShrinkBlocks mocks base method.
func (m *MockBlockDB) ShrinkBlocks(startHeight, endHeight uint64) (map[uint64][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShrinkBlocks", startHeight, endHeight)
	ret0, _ := ret[0].(map[uint64][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShrinkBlocks indicates an expected call of ShrinkBlocks.
func (mr *MockBlockDBMockRecorder) ShrinkBlocks(startHeight, endHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShrinkBlocks", reflect.TypeOf((*MockBlockDB)(nil).ShrinkBlocks), startHeight, endHeight)
}

// TxArchived mocks base method.
func (m *MockBlockDB) TxArchived(txId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxArchived", txId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxArchived indicates an expected call of TxArchived.
func (mr *MockBlockDBMockRecorder) TxArchived(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxArchived", reflect.TypeOf((*MockBlockDB)(nil).TxArchived), txId)
}

// TxExists mocks base method.
func (m *MockBlockDB) TxExists(txId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxExists", txId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxExists indicates an expected call of TxExists.
func (mr *MockBlockDBMockRecorder) TxExists(txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxExists", reflect.TypeOf((*MockBlockDB)(nil).TxExists), txId)
}
