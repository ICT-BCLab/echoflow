// Code generated by MockGen. DO NOT EDIT.
// Source: snapshot_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	accesscontrol "chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	common "chainmaker.org/chainmaker/pb-go/v2/common"
	vm "chainmaker.org/chainmaker/pb-go/v2/vm"
	protocol "chainmaker.org/chainmaker/protocol/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockSnapshotManager is a mock of SnapshotManager interface.
type MockSnapshotManager struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotManagerMockRecorder
}

// MockSnapshotManagerMockRecorder is the mock recorder for MockSnapshotManager.
type MockSnapshotManagerMockRecorder struct {
	mock *MockSnapshotManager
}

// NewMockSnapshotManager creates a new mock instance.
func NewMockSnapshotManager(ctrl *gomock.Controller) *MockSnapshotManager {
	mock := &MockSnapshotManager{ctrl: ctrl}
	mock.recorder = &MockSnapshotManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSnapshotManager) EXPECT() *MockSnapshotManagerMockRecorder {
	return m.recorder
}

// ClearSnapshot mocks base method.
func (m *MockSnapshotManager) ClearSnapshot(block *common.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearSnapshot", block)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearSnapshot indicates an expected call of ClearSnapshot.
func (mr *MockSnapshotManagerMockRecorder) ClearSnapshot(block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearSnapshot", reflect.TypeOf((*MockSnapshotManager)(nil).ClearSnapshot), block)
}

// GetSnapshot mocks base method.
func (m *MockSnapshotManager) GetSnapshot(prevBlock, block *common.Block) protocol.Snapshot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnapshot", prevBlock, block)
	ret0, _ := ret[0].(protocol.Snapshot)
	return ret0
}

// GetSnapshot indicates an expected call of GetSnapshot.
func (mr *MockSnapshotManagerMockRecorder) GetSnapshot(prevBlock, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnapshot", reflect.TypeOf((*MockSnapshotManager)(nil).GetSnapshot), prevBlock, block)
}

// NewSnapshot mocks base method.
func (m *MockSnapshotManager) NewSnapshot(prevBlock, block *common.Block) protocol.Snapshot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSnapshot", prevBlock, block)
	ret0, _ := ret[0].(protocol.Snapshot)
	return ret0
}

// NewSnapshot indicates an expected call of NewSnapshot.
func (mr *MockSnapshotManagerMockRecorder) NewSnapshot(prevBlock, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSnapshot", reflect.TypeOf((*MockSnapshotManager)(nil).NewSnapshot), prevBlock, block)
}

// NotifyBlockCommitted mocks base method.
func (m *MockSnapshotManager) NotifyBlockCommitted(block *common.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyBlockCommitted", block)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyBlockCommitted indicates an expected call of NotifyBlockCommitted.
func (mr *MockSnapshotManagerMockRecorder) NotifyBlockCommitted(block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyBlockCommitted", reflect.TypeOf((*MockSnapshotManager)(nil).NotifyBlockCommitted), block)
}

// MockSnapshot is a mock of Snapshot interface.
type MockSnapshot struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotMockRecorder
}

// MockSnapshotMockRecorder is the mock recorder for MockSnapshot.
type MockSnapshotMockRecorder struct {
	mock *MockSnapshot
}

// NewMockSnapshot creates a new mock instance.
func NewMockSnapshot(ctrl *gomock.Controller) *MockSnapshot {
	mock := &MockSnapshot{ctrl: ctrl}
	mock.recorder = &MockSnapshotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSnapshot) EXPECT() *MockSnapshotMockRecorder {
	return m.recorder
}

// ApplyBlock mocks base method.
func (m *MockSnapshot) ApplyBlock(block *common.Block, txRWSetMap map[string]*common.TxRWSet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ApplyBlock", block, txRWSetMap)
}

// ApplyBlock indicates an expected call of ApplyBlock.
func (mr *MockSnapshotMockRecorder) ApplyBlock(block, txRWSetMap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyBlock", reflect.TypeOf((*MockSnapshot)(nil).ApplyBlock), block, txRWSetMap)
}

// ApplyTxSimContext mocks base method.
func (m *MockSnapshot) ApplyTxSimContext(arg0 protocol.TxSimContext, arg1 protocol.ExecOrderTxType, arg2, arg3 bool) (bool, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyTxSimContext", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// ApplyTxSimContext indicates an expected call of ApplyTxSimContext.
func (mr *MockSnapshotMockRecorder) ApplyTxSimContext(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyTxSimContext", reflect.TypeOf((*MockSnapshot)(nil).ApplyTxSimContext), arg0, arg1, arg2, arg3)
}

// BuildDAG mocks base method.
func (m *MockSnapshot) BuildDAG(isSql bool, txRWSetTable []*common.TxRWSet) *common.DAG {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildDAG", isSql, txRWSetTable)
	ret0, _ := ret[0].(*common.DAG)
	return ret0
}

// BuildDAG indicates an expected call of BuildDAG.
func (mr *MockSnapshotMockRecorder) BuildDAG(isSql, txRWSetTable interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildDAG", reflect.TypeOf((*MockSnapshot)(nil).BuildDAG), isSql, txRWSetTable)
}

// GetBlockFingerprint mocks base method.
func (m *MockSnapshot) GetBlockFingerprint() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockFingerprint")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBlockFingerprint indicates an expected call of GetBlockFingerprint.
func (mr *MockSnapshotMockRecorder) GetBlockFingerprint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockFingerprint", reflect.TypeOf((*MockSnapshot)(nil).GetBlockFingerprint))
}

// GetBlockHeight mocks base method.
func (m *MockSnapshot) GetBlockHeight() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeight")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetBlockHeight indicates an expected call of GetBlockHeight.
func (mr *MockSnapshotMockRecorder) GetBlockHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeight", reflect.TypeOf((*MockSnapshot)(nil).GetBlockHeight))
}

// GetBlockProposer mocks base method.
func (m *MockSnapshot) GetBlockProposer() *accesscontrol.Member {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockProposer")
	ret0, _ := ret[0].(*accesscontrol.Member)
	return ret0
}

// GetBlockProposer indicates an expected call of GetBlockProposer.
func (mr *MockSnapshotMockRecorder) GetBlockProposer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockProposer", reflect.TypeOf((*MockSnapshot)(nil).GetBlockProposer))
}

// GetBlockTimestamp mocks base method.
func (m *MockSnapshot) GetBlockTimestamp() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockTimestamp")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetBlockTimestamp indicates an expected call of GetBlockTimestamp.
func (mr *MockSnapshotMockRecorder) GetBlockTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockTimestamp", reflect.TypeOf((*MockSnapshot)(nil).GetBlockTimestamp))
}

// GetBlockchainStore mocks base method.
func (m *MockSnapshot) GetBlockchainStore() protocol.BlockchainStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockchainStore")
	ret0, _ := ret[0].(protocol.BlockchainStore)
	return ret0
}

// GetBlockchainStore indicates an expected call of GetBlockchainStore.
func (mr *MockSnapshotMockRecorder) GetBlockchainStore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockchainStore", reflect.TypeOf((*MockSnapshot)(nil).GetBlockchainStore))
}

// GetKey mocks base method.
func (m *MockSnapshot) GetKey(txExecSeq int, contractName string, key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", txExecSeq, contractName, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey.
func (mr *MockSnapshotMockRecorder) GetKey(txExecSeq, contractName, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockSnapshot)(nil).GetKey), txExecSeq, contractName, key)
}

// GetKeys mocks base method.
func (m *MockSnapshot) GetKeys(txExecSeq int, keys []*vm.BatchKey) ([]*vm.BatchKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeys", txExecSeq, keys)
	ret0, _ := ret[0].([]*vm.BatchKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeys indicates an expected call of GetKeys.
func (mr *MockSnapshotMockRecorder) GetKeys(txExecSeq, keys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeys", reflect.TypeOf((*MockSnapshot)(nil).GetKeys), txExecSeq, keys)
}

// GetPreSnapshot mocks base method.
func (m *MockSnapshot) GetPreSnapshot() protocol.Snapshot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreSnapshot")
	ret0, _ := ret[0].(protocol.Snapshot)
	return ret0
}

// GetPreSnapshot indicates an expected call of GetPreSnapshot.
func (mr *MockSnapshotMockRecorder) GetPreSnapshot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreSnapshot", reflect.TypeOf((*MockSnapshot)(nil).GetPreSnapshot))
}

// GetSnapshotSize mocks base method.
func (m *MockSnapshot) GetSnapshotSize() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnapshotSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetSnapshotSize indicates an expected call of GetSnapshotSize.
func (mr *MockSnapshotMockRecorder) GetSnapshotSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnapshotSize", reflect.TypeOf((*MockSnapshot)(nil).GetSnapshotSize))
}

// GetSpecialTxTable mocks base method.
func (m *MockSnapshot) GetSpecialTxTable() []*common.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpecialTxTable")
	ret0, _ := ret[0].([]*common.Transaction)
	return ret0
}

// GetSpecialTxTable indicates an expected call of GetSpecialTxTable.
func (mr *MockSnapshotMockRecorder) GetSpecialTxTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpecialTxTable", reflect.TypeOf((*MockSnapshot)(nil).GetSpecialTxTable))
}

// GetTxRWSetTable mocks base method.
func (m *MockSnapshot) GetTxRWSetTable() []*common.TxRWSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxRWSetTable")
	ret0, _ := ret[0].([]*common.TxRWSet)
	return ret0
}

// GetTxRWSetTable indicates an expected call of GetTxRWSetTable.
func (mr *MockSnapshotMockRecorder) GetTxRWSetTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxRWSetTable", reflect.TypeOf((*MockSnapshot)(nil).GetTxRWSetTable))
}

// GetTxResultMap mocks base method.
func (m *MockSnapshot) GetTxResultMap() map[string]*common.Result {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxResultMap")
	ret0, _ := ret[0].(map[string]*common.Result)
	return ret0
}

// GetTxResultMap indicates an expected call of GetTxResultMap.
func (mr *MockSnapshotMockRecorder) GetTxResultMap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxResultMap", reflect.TypeOf((*MockSnapshot)(nil).GetTxResultMap))
}

// GetTxTable mocks base method.
func (m *MockSnapshot) GetTxTable() []*common.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxTable")
	ret0, _ := ret[0].([]*common.Transaction)
	return ret0
}

// GetTxTable indicates an expected call of GetTxTable.
func (mr *MockSnapshotMockRecorder) GetTxTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxTable", reflect.TypeOf((*MockSnapshot)(nil).GetTxTable))
}

// IsSealed mocks base method.
func (m *MockSnapshot) IsSealed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSealed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSealed indicates an expected call of IsSealed.
func (mr *MockSnapshotMockRecorder) IsSealed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSealed", reflect.TypeOf((*MockSnapshot)(nil).IsSealed))
}

// Seal mocks base method.
func (m *MockSnapshot) Seal() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Seal")
}

// Seal indicates an expected call of Seal.
func (mr *MockSnapshotMockRecorder) Seal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockSnapshot)(nil).Seal))
}

// SetPreSnapshot mocks base method.
func (m *MockSnapshot) SetPreSnapshot(arg0 protocol.Snapshot) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPreSnapshot", arg0)
}

// SetPreSnapshot indicates an expected call of SetPreSnapshot.
func (mr *MockSnapshotMockRecorder) SetPreSnapshot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPreSnapshot", reflect.TypeOf((*MockSnapshot)(nil).SetPreSnapshot), arg0)
}