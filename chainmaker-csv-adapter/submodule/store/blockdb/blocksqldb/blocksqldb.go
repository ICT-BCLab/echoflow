/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0

*/

package blocksqldb

import (
	"errors"
	"math"

	commonPb "chainmaker.org/chainmaker/pb-go/v2/common"
	storePb "chainmaker.org/chainmaker/pb-go/v2/store"
	"chainmaker.org/chainmaker/protocol/v2"
	"chainmaker.org/chainmaker/store/v2/serialization"
	"chainmaker.org/chainmaker/utils/v2"
)

var errNotImplement = errors.New("implement me")
var errNullPoint = errors.New("null point")

// BlockSqlDB provider a implementation of `blockdb.BlockDB`
// This implementation provides a mysql based data model
// @Description:
type BlockSqlDB struct {
	db protocol.SqlDBHandle
	//workersSemaphore *semaphore.Weighted
	logger protocol.Logger
	dbName string
}

// NewBlockSqlDB constructs a new `BlockSqlDB` given an chainId and engine type
//  @Description:
//  @param dbName
//  @param db
//  @param logger
//  @return *BlockSqlDB
func NewBlockSqlDB(dbName string, db protocol.SqlDBHandle, logger protocol.Logger) *BlockSqlDB {
	//nWorkers := runtime.NumCPU()
	blockDB := &BlockSqlDB{
		db: db,
		//workersSemaphore: semaphore.NewWeighted(int64(nWorkers)),
		logger: logger,
		dbName: dbName,
	}
	return blockDB
}

// GetHeightByHash 通过hash得到块高
//  @Description:
//  @receiver db
//  @param blockHash
//  @return uint64
//  @return error
func (db *BlockSqlDB) GetHeightByHash(blockHash []byte) (uint64, error) {
	sql := "SELECT block_height FROM block_infos WHERE block_hash=?"
	var height uint64
	res, err := db.db.QuerySingle(sql, blockHash)
	if err != nil {
		return math.MaxUint64, err
	}
	if res.IsEmpty() {
		return math.MaxUint64, nil
	}
	err = res.ScanColumns(&height)
	if err != nil {
		return math.MaxUint64, err
	}
	return height, nil
}

// GetBlockHeaderByHeight 通过块高得到header
//  @Description:
//  @receiver db
//  @param height
//  @return *commonPb.BlockHeader
//  @return error
func (db *BlockSqlDB) GetBlockHeaderByHeight(height uint64) (*commonPb.BlockHeader, error) {
	sql := "SELECT * from block_infos WHERE block_height=?"
	blockInfo, err := db.getBlockInfoBySql(sql, height)
	if err != nil {
		return nil, err
	}
	if blockInfo == nil && err == nil {
		return nil, nil
	}
	return blockInfo.GetBlockHeader(), nil
}

// GetTxHeight 通过txid得到块高
//  @Description:
//  @receiver db
//  @param txId
//  @return uint64
//  @return error
func (db *BlockSqlDB) GetTxHeight(txId string) (uint64, error) {
	sql := "SELECT block_height FROM tx_infos WHERE tx_id=?"
	var height uint64
	res, err := db.db.QuerySingle(sql, txId)
	if err != nil {
		return math.MaxUint64, err
	}
	if res.IsEmpty() {
		return math.MaxUint64, nil
	}
	err = res.ScanColumns(&height)
	if err != nil {
		return math.MaxUint64, err
	}
	return height, nil

}

// TxArchived 交易归档
//  @Description:
//  @receiver db
//  @param txId
//  @return bool
//  @return error
func (db *BlockSqlDB) TxArchived(txId string) (bool, error) {
	return false, nil
}

// GetArchivedPivot return default 0 and nil
//  @Description:
//  @receiver db
//  @return uint64
//  @return error
func (db *BlockSqlDB) GetArchivedPivot() (uint64, error) {
	return 0, nil
}

// ShrinkBlocks add NotImplement
//  @Description:
//  @receiver db
//  @param startHeight
//  @param endHeight
//  @return map[uint64][]string
//  @return error
func (db *BlockSqlDB) ShrinkBlocks(startHeight uint64, endHeight uint64) (map[uint64][]string, error) {
	return nil, errNotImplement
}

// RestoreBlocks NotImplement
//  @Description:
//  @receiver db
//  @param blockInfos
//  @return error
func (db *BlockSqlDB) RestoreBlocks(blockInfos []*serialization.BlockWithSerializedInfo) error {
	return errNotImplement
}

//  initDb 如果数据库不存在，则创建数据库，然后切换到这个数据库，创建表
//  @Description:
// 如果数据库存在，则切换数据库，检查表是否存在，不存在则创建表。
//  @receiver db
//  @param dbName
func (db *BlockSqlDB) initDb(dbName string) {
	_, err := db.db.CreateDatabaseIfNotExist(dbName)
	if err != nil {
		panic("init state sql db fail")
	}

	err = db.db.CreateTableIfNotExist(&BlockInfo{})
	if err != nil {
		panic("init state sql db table `block_infos` fail" + err.Error())
	}
	err = db.db.CreateTableIfNotExist(&TxInfo{})
	if err != nil {
		panic("init state sql db table `tx_infos` fail")
	}
}

//func getDbName(dbConfig *localconf.SqlDbConfig, chainId string) string {
//	return dbConfig.DbPrefix + "blockdb_" + chainId
//}

// InitGenesis 创世块初始化
//  @Description:
//  @receiver b
//  @param genesisBlock
//  @return error
func (b *BlockSqlDB) InitGenesis(genesisBlock *serialization.BlockWithSerializedInfo) error {
	b.initDb(b.dbName)
	return b.CommitBlock(genesisBlock, false)
}

// CommitBlock  commits the block and the corresponding rwsets in an atomic operation
//  @Description:
//  @receiver b
//  @param blocksInfo
//  @param isCache
//  @return error
func (b *BlockSqlDB) CommitBlock(blocksInfo *serialization.BlockWithSerializedInfo, isCache bool) error {
	block := blocksInfo.Block
	dbTxKey := block.GetTxKey()
	startCommitTxs := utils.CurrentTimeMillisSeconds()
	dbtx, err := b.db.BeginDbTransaction(dbTxKey)
	if err != nil {
		return err
	}
	//save txs
	for index, tx := range block.Txs {
		var txInfo *TxInfo
		txInfo, err = NewTxInfo(tx, uint64(block.Header.BlockHeight), block.Header.BlockHash, uint32(index))
		if err != nil {
			b.logger.Errorf("failed to init txinfo, err:%s", err)
			if err2 := b.db.RollbackDbTransaction(dbTxKey); err2 != nil {
				b.logger.Errorf("failed to rollback db transaction[%s],error:%s", dbTxKey, err2.Error())
				return err2
			}
			return err
		}
		_, err = dbtx.Save(txInfo)
		if err != nil {
			b.logger.Errorf("failed to commit txinfo info, height:%d, tx:%s,err:%s",
				block.Header.BlockHeight, txInfo.TxId, err)
			if err2 := b.db.RollbackDbTransaction(dbTxKey); err2 != nil {
				b.logger.Errorf("failed to rollback db transaction[%s],error:%s", dbTxKey, err2.Error())
				return err2
			}
			return err
		}
	}

	elapsedCommitTxs := utils.CurrentTimeMillisSeconds() - startCommitTxs
	//save block info
	startCommitBlockInfo := utils.CurrentTimeMillisSeconds()
	blockInfo, err := NewBlockInfo(block)
	if err != nil {
		b.logger.Errorf("failed to init blockinfo, err:%s", err)
		if err2 := b.db.RollbackDbTransaction(dbTxKey); err2 != nil {
			b.logger.Errorf("failed to rollback db transaction[%s],error:%s", dbTxKey, err2.Error())
			return err2
		}
		return err
	}
	_, err = dbtx.Save(blockInfo)
	if err != nil {
		b.logger.Errorf("failed to commit block info, height:%d, err:%s",
			block.Header.BlockHeight, err)
		_ = b.db.RollbackDbTransaction(dbTxKey) //rollback tx
		return err
	}
	err = b.db.CommitDbTransaction(dbTxKey)
	if err != nil {
		b.logger.Errorf("failed to commit tx, err:%s", err)
		return err
	}
	elapsedCommitBlockInfos := utils.CurrentTimeMillisSeconds() - startCommitBlockInfo
	b.logger.Debugf("chain[%s]: commit block[%d] sql blockdb, time used (commit_txs:%d, commit_block:%d, total:%d)",
		block.Header.ChainId, block.Header.BlockHeight, elapsedCommitTxs, elapsedCommitBlockInfos,
		utils.CurrentTimeMillisSeconds()-startCommitTxs)
	return nil
}

// BlockExists returns true if the block hash exist, or returns false if none exists.
//  @Description:
//  @receiver b
//  @param blockHash
//  @return bool
//  @return error
func (b *BlockSqlDB) BlockExists(blockHash []byte) (bool, error) {
	var count int64
	sql := "select count(*) from block_infos where block_hash = ?"
	res, err := b.db.QuerySingle(sql, blockHash)
	if err != nil {
		return false, err
	}
	err = res.ScanColumns(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetBlockByHash returns a block given it's hash, or returns nil if none exists.
//  @Description:
//  @receiver b
//  @param blockHash
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) GetBlockByHash(blockHash []byte) (*commonPb.Block, error) {

	return b.getFullBlockBySql("select * from block_infos where block_hash = ?", blockHash)
}

//  getBlockInfoBySql
//  @Description:
//  @receiver b
//  @param sql
//  @param values
//  @return *BlockInfo
//  @return error
func (b *BlockSqlDB) getBlockInfoBySql(sql string, values ...interface{}) (*BlockInfo, error) {
	//get block info from mysql
	var blockInfo BlockInfo
	res, _ := b.db.QuerySingle(sql, values...)

	if res == nil || res.IsEmpty() {
		b.logger.Infof("sql[%s] %v return empty result", sql, values)
		return nil, nil
	}
	err := blockInfo.ScanObject(res.ScanColumns)
	if err != nil {
		return nil, err
	}
	return &blockInfo, nil
}

//  getFullBlockBySql
//  @Description:
//  @receiver b
//  @param sql
//  @param values
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) getFullBlockBySql(sql string, values ...interface{}) (*commonPb.Block, error) {
	blockInfo, err := b.getBlockInfoBySql(sql, values...)
	if err != nil {
		return nil, err
	}
	if blockInfo == nil && err == nil {
		return nil, nil
	}
	block, err := blockInfo.GetBlock()
	if err != nil {
		return nil, err
	}
	txs, err := b.getTxsByBlockHeight(blockInfo.BlockHeight)
	if err != nil {
		return nil, err
	}
	block.Txs = txs
	return block, nil
}

// GetBlock returns a block given it's block height, or returns nil if none exists.
//  @Description:
//  @receiver b
//  @param height
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) GetBlock(height uint64) (*commonPb.Block, error) {
	return b.getFullBlockBySql("select * from block_infos where block_height =?", height)
}

// GetLastBlock  returns the last block.
//  @Description:
//  @receiver b
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) GetLastBlock() (*commonPb.Block, error) {
	return b.getFullBlockBySql(`select * 
from block_infos 
where block_height = (select max(block_height) from block_infos)`)
}

// GetLastConfigBlock returns the last config block.
//  @Description:
//  @receiver b
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) GetLastConfigBlock() (*commonPb.Block, error) {
	lastBlock, err := b.GetLastBlock()
	if err != nil {
		return nil, err
	}
	if utils.IsConfBlock(lastBlock) {
		return lastBlock, nil
	}
	return b.GetBlock(lastBlock.Header.PreConfHeight)
}

// GetLastConfigBlockHeight returns the last config block height.
//  @Description:
//  @receiver b
//  @return uint64
//  @return error
func (b *BlockSqlDB) GetLastConfigBlockHeight() (uint64, error) {
	return 0, errNotImplement
}

// GetFilteredBlock returns a filtered block given it's block height, or return nil if none exists.
//  @Description:
//  @receiver b
//  @param height
//  @return *storePb.SerializedBlock
//  @return error
func (b *BlockSqlDB) GetFilteredBlock(height uint64) (*storePb.SerializedBlock, error) {
	blockInfo, err := b.getBlockInfoBySql("select * from block_infos where block_height = ?", height)
	if err != nil {
		return nil, err
	}
	if blockInfo == nil && err == nil {
		return nil, nil
	}
	return blockInfo.GetFilteredBlock()
}

// GetLastSavepoint returns the last block height
//  @Description:
//  @receiver b
//  @return uint64
//  @return error
func (b *BlockSqlDB) GetLastSavepoint() (uint64, error) {
	sql := "select max(block_height) from block_infos"
	row, err := b.db.QuerySingle(sql)
	if err != nil {
		//b.logger.Errorf("get block sqldb save point error:%s", err.Error())
		return 0, err
	}
	if row.IsEmpty() {
		return 0, nil
	}
	var height uint64
	err = row.ScanColumns(&height)
	if err != nil {
		return 0, err
	}

	return height, nil
}

// GetBlockByTx returns a block which contains a tx.
//  @Description:
//  @receiver b
//  @param txId
//  @return *commonPb.Block
//  @return error
func (b *BlockSqlDB) GetBlockByTx(txId string) (*commonPb.Block, error) {
	sql := "select * from block_infos where block_height=(select block_height from tx_infos where tx_id=?)"
	return b.getFullBlockBySql(sql, txId)
}

// GetTx retrieves a transaction by txid, or returns nil if none exists.
//  @Description:
//  @receiver b
//  @param txId
//  @return *commonPb.Transaction
//  @return error
func (b *BlockSqlDB) GetTx(txId string) (*commonPb.Transaction, error) {
	if len(txId) == 0 {
		return nil, errors.New("parameter is null")
	}
	var txInfo TxInfo
	res, err := b.db.QuerySingle("select * from tx_infos where tx_id = ?", txId)
	if err != nil {
		return nil, err
	}
	if res.IsEmpty() {
		b.logger.Infof("tx[%s] not found in db", txId)
		return nil, nil
	}

	err = txInfo.ScanObject(res.ScanColumns)
	if err != nil {
		return nil, err
	}
	if len(txInfo.TxId) > 0 {
		return txInfo.GetTx()
	}
	b.logger.Errorf("tx data not found by txid:%s", txId)
	return nil, errors.New("data not found")
}

// GetTxWithBlockInfo 通过交易id 得到交易Info信息
//  @Description:
//  @receiver b
//  @param txId
//  @return *storePb.TransactionStoreInfo
//  @return error
func (b *BlockSqlDB) GetTxWithBlockInfo(txId string) (*storePb.TransactionStoreInfo, error) {
	var txInfo TxInfo
	res, err := b.db.QuerySingle("select * from tx_infos where tx_id = ?", txId)
	if err != nil {
		return nil, err
	}
	if res.IsEmpty() {
		b.logger.Infof("tx[%s] not found in db", txId)
		return nil, nil
	}
	err = txInfo.ScanObject(res.ScanColumns)
	if err != nil {
		return nil, err
	}
	if len(txInfo.TxId) > 0 {
		return txInfo.GetTxInfo()
	}
	b.logger.Errorf("tx data not found by txid:%s", txId)
	return nil, errors.New("data not found")
}

// GetTxInfoOnly 获得除Tx之外的其他TxInfo信息
//  @Description:
//  @receiver b
//  @param txId
//  @return *storePb.TransactionStoreInfo
//  @return error
func (b *BlockSqlDB) GetTxInfoOnly(txId string) (*storePb.TransactionStoreInfo, error) {
	if len(txId) == 0 {
		return nil, errors.New("parameter is null")
	}
	var txInfo TxInfo
	res, err := b.db.QuerySingle("select * from tx_infos where tx_id = ?", txId)
	if err != nil {
		return nil, err
	}
	if res.IsEmpty() {
		b.logger.Infof("tx[%s] not found in db", txId)
		return nil, nil
	}

	err = txInfo.ScanObject(res.ScanColumns)
	if err != nil {
		return nil, err
	}
	if len(txInfo.TxId) > 0 {
		return txInfo.GetTxInfo()
	}
	b.logger.Errorf("tx data not found by txid:%s", txId)
	return nil, errors.New("data not found")
}

// TxExists returns true if the tx exist, or returns false if none exists.
//  @Description:
//  @receiver b
//  @param txId
//  @return bool
//  @return error
func (b *BlockSqlDB) TxExists(txId string) (bool, error) {
	var count int64
	sql := "select count(*) from tx_infos where tx_id = ?"
	res, err := b.db.QuerySingle(sql, txId)
	if err != nil {
		return false, err
	}
	err = res.ScanColumns(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//  getTxsByBlockHeight 获得某个区块高度下的所有交易
//  @Description:
//  @receiver b
//  @param blockHeight
//  @return []*commonPb.Transaction
//  @return error
func (b *BlockSqlDB) getTxsByBlockHeight(blockHeight uint64) ([]*commonPb.Transaction, error) {
	res, err := b.db.QueryMulti("select * from tx_infos where block_height = ? order by `offset`", blockHeight)
	if err != nil {
		return nil, err
	}
	result := []*commonPb.Transaction{}
	for res.Next() {
		var txInfo TxInfo
		err = txInfo.ScanObject(res.ScanColumns)
		if err != nil {
			return nil, err
		}
		tx, err := txInfo.GetTx()
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}

// GetTxConfirmedTime NotImplement
//  @Description:
//  @receiver b
//  @param txId
//  @return int64
//  @return error
func (b *BlockSqlDB) GetTxConfirmedTime(txId string) (int64, error) {
	return 0, errNotImplement
}

// GetBlockIndex NotImplement
//  @Description:
//  @receiver b
//  @param height
//  @return *storePb.StoreInfo
//  @return error
func (b *BlockSqlDB) GetBlockIndex(height uint64) (*storePb.StoreInfo, error) {
	return nil, errNotImplement
}

// GetBlockMetaIndex NotImplement
//  @Description:
//  @receiver b
//  @param height
//  @return *storePb.StoreInfo
//  @return error
func (b *BlockSqlDB) GetBlockMetaIndex(height uint64) (*storePb.StoreInfo, error) {
	return nil, errNotImplement
}

// GetTxIndex NotImplement
//  @Description:
//  @receiver b
//  @param txId
//  @return *storePb.StoreInfo
//  @return error
func (b *BlockSqlDB) GetTxIndex(txId string) (*storePb.StoreInfo, error) {
	return nil, errNotImplement
}

// Close is used to close database
//  @Description:
//  @receiver b
func (b *BlockSqlDB) Close() {
	b.logger.Info("close block sql db")
	b.db.Close()
}
