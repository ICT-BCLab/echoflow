/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package accesscontrol

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"chainmaker.org/chainmaker/common/v2/msgbus"

	"chainmaker.org/chainmaker/common/v2/concurrentlru"
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/localconf/v2"
	pbac "chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/protocol/v2"
)

var _ protocol.AccessControlProvider = (*pkACProvider)(nil)

var nilPkACProvider ACProvider = (*pkACProvider)(nil)

const (
	//AdminPublicKey admin trust orgId
	AdminPublicKey = "public"
	//DposOrgId chainconfig the DPoS of orgId
	DposOrgId = "dpos_org_id"

	//PermissionConsensusOrgId chainconfig orgId for permission consensus, such as tbft
	PermissionConsensusOrgId = "public"
)

var (
	pubPolicyConsensus = newPolicy(
		protocol.RuleAny,
		nil,
		[]protocol.Role{
			protocol.RoleConsensusNode,
		},
	)
	pubPolicyManage = newPolicy(
		protocol.RuleAny,
		nil,
		[]protocol.Role{
			protocol.RoleAdmin,
		},
	)
	pubPolicyMajorityAdmin = newPolicy(
		protocol.RuleMajority,
		nil,
		[]protocol.Role{
			protocol.RoleAdmin,
		},
	)
	pubPolicyTransaction = newPolicy(
		protocol.RuleAny,
		nil,
		nil,
	)
	pubPolicyForbidden = newPolicy(
		protocol.RuleForbidden,
		nil,
		nil,
	)
)

type pkACProvider struct {

	//chainconfig authType
	authType string

	hashType string

	adminNum int32

	log protocol.Logger

	adminMember *sync.Map

	consensusMember *sync.Map

	memberCache *concurrentlru.Cache

	dataStore protocol.BlockchainStore

	resourceNamePolicyMap *sync.Map

	exceptionalPolicyMap *sync.Map

	resourceNamePolicyMap220 *sync.Map
	exceptionalPolicyMap220  *sync.Map
}

type publicAdminMemberModel struct {
	publicKey crypto.PublicKey
	pkBytes   []byte
}

func (p *pkACProvider) NewACProvider(chainConf protocol.ChainConf, localOrgId string,
	store protocol.BlockchainStore, log protocol.Logger, msgBus msgbus.MessageBus) (
	protocol.AccessControlProvider, error) {
	pkAcProvider, err := newPkACProvider(chainConf.ChainConfig(), store, log)
	if err != nil {
		return nil, err
	}

	msgBus.Register(msgbus.ChainConfig, pkAcProvider)
	//v220_compat Deprecated
	chainConf.AddWatch(pkAcProvider) //nolint: staticcheck
	return pkAcProvider, nil
}

func newPkACProvider(chainConfig *config.ChainConfig,
	store protocol.BlockchainStore, log protocol.Logger) (*pkACProvider, error) {
	pkAcProvider := &pkACProvider{
		adminNum:                 0,
		hashType:                 chainConfig.Crypto.Hash,
		authType:                 chainConfig.AuthType,
		adminMember:              &sync.Map{},
		consensusMember:          &sync.Map{},
		memberCache:              concurrentlru.New(localconf.ChainMakerConfig.NodeConfig.CertCacheSize),
		log:                      log,
		dataStore:                store,
		resourceNamePolicyMap:    &sync.Map{},
		exceptionalPolicyMap:     &sync.Map{},
		resourceNamePolicyMap220: &sync.Map{},
		exceptionalPolicyMap220:  &sync.Map{},
	}

	if chainConfig.Consensus.Type == consensus.ConsensusType_DPOS {
		pkAcProvider.createDefaultResourcePolicyForDPoS()
		pkAcProvider.createDefaultResourcePolicyForDPoS_220()
	} else {
		pkAcProvider.createDefaultResourcePolicy()
		pkAcProvider.createDefaultResourcePolicy_220()
	}

	err := pkAcProvider.initAdminMembers(chainConfig.TrustRoots)
	if err != nil {
		return nil, fmt.Errorf("new public AC provider failed: %s", err.Error())
	}
	err = pkAcProvider.initConsensusMember(chainConfig)
	if err != nil {
		return nil, fmt.Errorf("new public AC provider failed: %s", err.Error())
	}
	return pkAcProvider, nil
}

func (p *pkACProvider) initAdminMembers(trustRootList []*config.TrustRootConfig) error {
	var (
		tempSyncMap sync.Map
	)

	if len(trustRootList) == 0 {
		p.log.Debugf("no super administrator is configured")
		return nil
	}

	var adminNum int32

	for _, trustRoot := range trustRootList {
		if strings.ToLower(trustRoot.OrgId) == AdminPublicKey {
			for _, root := range trustRoot.Root {
				pk, err := asym.PublicKeyFromPEM([]byte(root))
				if err != nil {
					return fmt.Errorf("init admin member failed: parse the public key from PEM failed")
				}
				pkBytes, err := pk.Bytes()
				if err != nil {
					return fmt.Errorf("init admin member failed: %s", err.Error())
				}
				adminMember := &publicAdminMemberModel{
					publicKey: pk,
					pkBytes:   pkBytes,
				}
				adminKey := hex.EncodeToString(pkBytes)
				tempSyncMap.Store(adminKey, adminMember)
				adminNum++
			}
		}
	}
	p.adminMember = &tempSyncMap
	atomic.StoreInt32(&p.adminNum, adminNum)
	return nil
}

func (p *pkACProvider) initConsensusMember(chainConfig *config.ChainConfig) error {
	if chainConfig.Consensus.Type == consensus.ConsensusType_DPOS {
		return p.initDPoSMember(chainConfig.Consensus.Nodes)
	} else if chainConfig.Consensus.Type == consensus.ConsensusType_TBFT {
		return p.initPermissionMember(chainConfig.Consensus.Nodes)
	}
	return fmt.Errorf("public chain mode does not support other consensus")
}

func (p *pkACProvider) initDPoSMember(consensusConf []*config.OrgConfig) error {
	if len(consensusConf) == 0 {
		return fmt.Errorf("update dpos consensus member failed: DPoS config can't be empty in chain config")
	}

	var consensusMember sync.Map
	if consensusConf[0].OrgId != DposOrgId {
		return fmt.Errorf("update dpos consensus member failed: DPoS node config orgId do not match")
	}
	for _, nodeId := range consensusConf[0].NodeId {
		consensusMember.Store(nodeId, struct{}{})
	}
	p.consensusMember = &consensusMember
	p.log.Infof("update consensus list: [%v]", p.consensusMember)
	return nil
}

func (p *pkACProvider) initPermissionMember(consensusConf []*config.OrgConfig) error {
	if len(consensusConf) == 0 {
		return fmt.Errorf("update permission consensus member failed: consensus node config can't be empty in chain config")
	}

	var consensusMember sync.Map
	if consensusConf[0].OrgId != PermissionConsensusOrgId {
		return fmt.Errorf("update permission consensus member failed: node config orgId do not match")
	}
	for _, nodeId := range consensusConf[0].NodeId {
		consensusMember.Store(nodeId, struct{}{})
	}
	p.consensusMember = &consensusMember
	p.log.Infof("update permission consensus list: [%v]", p.consensusMember)
	return nil
}

func (p *pkACProvider) lookUpMemberInCache(memberInfo string) (*memberCached, bool) {
	ret, ok := p.memberCache.Get(memberInfo)
	if ok {
		return ret.(*memberCached), true
	}
	return nil, false
}

func (p *pkACProvider) getMemberFromCache(member *pbac.Member) protocol.Member {
	cached, ok := p.lookUpMemberInCache(string(member.MemberInfo))
	if ok {
		p.log.Debugf("member found in local cache")
		return cached.member
	}
	// handle false positive when member cache is cleared
	if p.authType == protocol.Public {
		tmpMember, err := p.NewMemberFromAcs(member)
		if err != nil {
			p.log.Debugf("new member failed, authType = %s, err = %s", p.authType, err.Error())
			return nil
		}
		p.memberCache.Add(string(member.MemberInfo), &memberCached{
			member:    tmpMember,
			certChain: nil,
		})
		return tmpMember
	}
	return nil
}

//func (p *pkACProvider) Module() string {
//	return ModuleNameAccessControl
//}
//
//
//func (p *pkACProvider) Watch(chainConfig *config.ChainConfig) error {
//
//	p.hashType = chainConfig.GetCrypto().GetHash()
//	err := p.initAdminMembers(chainConfig.TrustRoots)
//	if err != nil {
//		return fmt.Errorf("new public AC provider failed: %s", err.Error())
//	}
//
//	err = p.initConsensusMember(chainConfig)
//	if err != nil {
//		return fmt.Errorf("new public AC provider failed: %s", err.Error())
//	}
//	p.memberCache.Clear()
//	return nil
//}

// NewMember creates a member from pb Member
func (p *pkACProvider) NewMember(pbMember *pbac.Member) (protocol.Member, error) {
	cache := p.getMemberFromCache(pbMember)
	if cache != nil {
		return cache, nil
	}
	member, err := publicNewPkMemberFromAcs(pbMember, p.adminMember, p.consensusMember, p.hashType)
	if err != nil {
		return nil, fmt.Errorf("new member failed: %s", err.Error())
	}
	p.memberCache.Add(string(pbMember.MemberInfo), &memberCached{
		member:    member,
		certChain: nil,
	})
	return member, nil
}

// NewMember creates a member from pb Member
func (p *pkACProvider) NewMemberFromAcs(pbMember *pbac.Member) (protocol.Member, error) {
	member, err := publicNewPkMemberFromAcs(pbMember, p.adminMember, p.consensusMember, p.hashType)
	if err != nil {
		return nil, fmt.Errorf("new member failed: %s", err.Error())
	}
	return member, nil
}

func (p *pkACProvider) createDefaultResourcePolicy() {
	p.resourceNamePolicyMap.Store(protocol.ResourceNameReadData, policyRead)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameWriteData, policyWrite)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameUpdateSelfConfig, policySelfConfig)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameUpdateConfig, policyConfig)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameConsensusNode, policyConsensus)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameP2p, policyP2P)

	// for txtype
	p.resourceNamePolicyMap.Store(common.TxType_QUERY_CONTRACT.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_INVOKE_CONTRACT.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_SUBSCRIBE.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_ARCHIVE.String(), pubPolicyManage)

	// exceptional resourceName
	p.exceptionalPolicyMap.Store(protocol.ResourceNamePrivateCompute, pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PRIVATE_COMPUTE.String()+"-"+
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PRIVATE_COMPUTE.String()+"-"+
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_DELETE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_UPDATE.String(), pubPolicyForbidden)

	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_ADD.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_DELETE.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_UPDATE.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_UPDATE.String(), pubPolicyMajorityAdmin)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_FREEZE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_UNFREEZE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_DELETE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_REVOKE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ALIAS_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ALIAS_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_ALIAS_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PUBKEY_MANAGE.String()+"-"+
		syscontract.PubkeyManageFunction_PUBKEY_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PUBKEY_MANAGE.String()+"-"+
		syscontract.PubkeyManageFunction_PUBKEY_DELETE.String(), pubPolicyForbidden)

	// disable trust root add & delete for public mode
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_DELETE.String(), pubPolicyForbidden)

	// disable contract access for public mode
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_GRANT_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_REVOKE_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_VERIFY_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractQueryFunction_GET_DISABLED_CONTRACT_LIST.String(), pubPolicyForbidden)

	// forbidden charge gas by go sdk
	//p.exceptionalPolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
	//	syscontract.GasAccountFunction_CHARGE_GAS.String(), pubPolicyForbidden)

	// forbidden refund gas vm by go sdk
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_REFUND_GAS_VM.String(), pubPolicyForbidden)

	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_INIT_CONTRACT.String(), pubPolicyManage)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_UPGRADE_CONTRACT.String(), pubPolicyManage)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_FREEZE_CONTRACT.String(), pubPolicyManage)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_UNFREEZE_CONTRACT.String(), pubPolicyManage)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_REVOKE_CONTRACT.String(), pubPolicyManage)

	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CORE_UPDATE.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_BLOCK_UPDATE.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_UPDATE_VERSION.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_ENABLE_OR_DISABLE_GAS.String(), pubPolicyMajorityAdmin)
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_ALTER_ADDR_TYPE.String(), pubPolicyMajorityAdmin)

	// for admin management
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_UPDATE.String(), pubPolicyMajorityAdmin)

	// for gas admin
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_SET_ADMIN.String(), pubPolicyMajorityAdmin)
	// for charge gas in optimize mode
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_CHARGE_GAS_FOR_MULTI_ACCOUNT.String(), policyConsensus)
	// for set invoke base gas
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_SET_INVOKE_BASE_GAS.String(), pubPolicyMajorityAdmin)
	// move set admin method to chain config module
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_SET_ACCOUNT_MANAGER_ADMIN.String(), pubPolicyMajorityAdmin)
}

// need to consistent with 2.1.0 for dpos
func (p *pkACProvider) createDefaultResourcePolicyForDPoS() {
	p.resourceNamePolicyMap.Store(protocol.ResourceNameReadData, policyRead)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameWriteData, policyWrite)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameUpdateSelfConfig, policySelfConfig)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameUpdateConfig, policyConfig)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameConsensusNode, policyConsensus)
	p.resourceNamePolicyMap.Store(protocol.ResourceNameP2p, policyP2P)

	// for txtype
	p.resourceNamePolicyMap.Store(common.TxType_QUERY_CONTRACT.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_INVOKE_CONTRACT.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_SUBSCRIBE.String(), pubPolicyTransaction)
	p.resourceNamePolicyMap.Store(common.TxType_ARCHIVE.String(), pubPolicyManage)

	// exceptional resourceName
	p.exceptionalPolicyMap.Store(protocol.ResourceNamePrivateCompute, pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PRIVATE_COMPUTE.String()+"-"+
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PRIVATE_COMPUTE.String()+"-"+
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_DELETE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_MEMBER_UPDATE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_DELETE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ID_UPDATE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_NODE_ORG_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CONSENSUS_EXT_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_PERMISSION_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_FREEZE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_UNFREEZE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_DELETE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_REVOKE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ALIAS_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERT_ALIAS_UPDATE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CERT_MANAGE.String()+"-"+
		syscontract.CertManageFunction_CERTS_ALIAS_DELETE.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PUBKEY_MANAGE.String()+"-"+
		syscontract.PubkeyManageFunction_PUBKEY_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_PUBKEY_MANAGE.String()+"-"+
		syscontract.PubkeyManageFunction_PUBKEY_DELETE.String(), pubPolicyForbidden)

	// disable trust root add & delete for public mode
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_ADD.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_DELETE.String(), pubPolicyForbidden)

	// disable multisign for public mode
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_MULTI_SIGN.String()+"-"+
		syscontract.MultiSignFunction_REQ.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_MULTI_SIGN.String()+"-"+
		syscontract.MultiSignFunction_VOTE.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_MULTI_SIGN.String()+"-"+
		syscontract.MultiSignFunction_QUERY.String(), pubPolicyForbidden)

	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_CORE_UPDATE.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_BLOCK_UPDATE.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_UPDATE_VERSION.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_UPGRADE_CONTRACT.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_FREEZE_CONTRACT.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_UNFREEZE_CONTRACT.String(), pubPolicyManage)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_REVOKE_CONTRACT.String(), pubPolicyManage)
	// disable contract access for public mode
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_GRANT_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_REVOKE_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractManageFunction_VERIFY_CONTRACT_ACCESS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CONTRACT_MANAGE.String()+"-"+
		syscontract.ContractQueryFunction_GET_DISABLED_CONTRACT_LIST.String(), pubPolicyForbidden)

	// disable gas related native contract
	//p.exceptionalPolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
	//	syscontract.GasAccountFunction_CHARGE_GAS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_REFUND_GAS_VM.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_SET_ADMIN.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_ENABLE_OR_DISABLE_GAS.String(), pubPolicyForbidden)
	p.exceptionalPolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_ALTER_ADDR_TYPE.String(), pubPolicyForbidden)
	// for charge gas in optimize mode
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_ACCOUNT_MANAGER.String()+"-"+
		syscontract.GasAccountFunction_CHARGE_GAS_FOR_MULTI_ACCOUNT.String(), policyConsensus)

	// for admin management
	p.resourceNamePolicyMap.Store(syscontract.SystemContract_CHAIN_CONFIG.String()+"-"+
		syscontract.ChainConfigFunction_TRUST_ROOT_UPDATE.String(), pubPolicyMajorityAdmin)
}

func (p *pkACProvider) verifyPrincipalPolicy(principal,
	refinedPrincipal protocol.Principal, pol *policy) (bool, error) {
	endorsements := refinedPrincipal.GetEndorsement()
	rule := pol.GetRule()
	switch rule {
	case protocol.RuleForbidden:
		return false, fmt.Errorf("public authentication fail: [%s] is forbidden to access",
			refinedPrincipal.GetResourceName())
	case protocol.RuleAny:
		return p.verifyRuleAnyCase(pol, endorsements)
	case protocol.RuleMajority:
		return p.verifyRuleMajorityCase(pol, endorsements)
	default:
		return false, fmt.Errorf("public authentication fail: [%s] is not supported", rule)
	}
}

func (p *pkACProvider) verifyRuleAnyCase(pol *policy, endorsements []*common.EndorsementEntry) (bool, error) {
	roleList := p.buildRoleListForVerifyPrincipal(pol)
	for _, endorsement := range endorsements {
		if len(roleList) == 0 {
			return true, nil
		}
		member := p.getMemberFromCache(endorsement.Signer)
		if member == nil {
			p.log.Infof(
				"authentication warning: the member is not in member cache, memberInfo[%s]",
				string(endorsement.Signer.MemberInfo))
			continue
		}

		if _, ok := roleList[member.GetRole()]; ok {
			return true, nil
		}
		p.log.Infof("authentication warning, the member role is not in roleList, role: [%s]",
			member.GetRole())
	}
	err := fmt.Errorf("authentication fail for any rule, policy: rule: [%v],roleList: [%v]",
		pol.rule, pol.roleList)
	return false, err
}

func (p *pkACProvider) verifyRuleMajorityCase(pol *policy, endorsements []*common.EndorsementEntry) (bool, error) {
	role := protocol.RoleAdmin
	refinedEndorsements := p.getValidEndorsements(map[string]bool{}, map[protocol.Role]bool{role: true}, endorsements)
	numOfValid := len(refinedEndorsements)
	p.log.Debugf("verifyRuleMajorityAdminCase: numOfValid=[%d], p.adminNum=[%d]", numOfValid, p.adminNum)
	if float64(numOfValid) > float64(p.adminNum)/2.0 {
		return true, nil
	}
	return false, fmt.Errorf("%s: %d valid endorsements required, %d valid endorsements received",
		notEnoughParticipantsSupportError, int(float64(p.adminNum)/2.0+1), numOfValid)
}

func (p *pkACProvider) buildRoleListForVerifyPrincipal(pol *policy) map[protocol.Role]bool {
	roleListRaw := pol.GetRoleList()
	roleList := map[protocol.Role]bool{}
	for _, roleRaw := range roleListRaw {
		roleList[roleRaw] = true
	}
	return roleList
}

func (p *pkACProvider) lookUpPolicyByResourceName(resourceName string) (*policy, error) {
	blockVersion, policyResourceName := getBlockVersionAndResourceName(resourceName)

	if blockVersion > 0 && blockVersion <= 220 {
		return p.lookUpPolicyByResourceName220(policyResourceName)
	}

	pol, ok := p.resourceNamePolicyMap.Load(resourceName)
	if !ok {
		if pol, ok = p.exceptionalPolicyMap.Load(resourceName); !ok {
			return nil, fmt.Errorf("look up access policy failed, did not configure access policy "+
				"for resource %s", resourceName)
		}
	}
	return pol.(*policy), nil
}

// all-in-one validation for signing members: signature, policies
func (p *pkACProvider) refinePrincipal(principal protocol.Principal) (protocol.Principal, error) {
	endorsements := principal.GetEndorsement()
	msg := principal.GetMessage()
	refinedEndorsement := p.refineEndorsements(endorsements, msg)
	if len(refinedEndorsement) <= 0 {
		return nil, fmt.Errorf("refine endorsements failed, all endorsers have failed verification")
	}

	refinedPrincipal, err := p.CreatePrincipal(principal.GetResourceName(), refinedEndorsement, msg)
	if err != nil {
		return nil, fmt.Errorf("create principal failed: [%s]", err.Error())
	}

	return refinedPrincipal, nil
}

func (p *pkACProvider) refineEndorsements(endorsements []*common.EndorsementEntry,
	msg []byte) []*common.EndorsementEntry {

	refinedSigners := map[string]bool{}
	var refinedEndorsement []*common.EndorsementEntry

	for _, endorsementEntry := range endorsements {
		endorsement := &common.EndorsementEntry{
			Signer: &pbac.Member{
				OrgId:      endorsementEntry.Signer.OrgId,
				MemberInfo: endorsementEntry.Signer.MemberInfo,
				MemberType: endorsementEntry.Signer.MemberType,
			},
			Signature: endorsementEntry.Signature,
		}
		memInfo := string(endorsement.Signer.MemberInfo)

		remoteMember, err := p.NewMember(endorsement.Signer)
		if err != nil {
			p.log.Infof("new member failed: [%s]", err.Error())
			continue
		}

		if err := remoteMember.Verify(p.hashType, msg, endorsement.Signature); err != nil {
			p.log.Infof("signer member verify signature failed: [%s]", err.Error())
			p.log.Debugf("information for invalid signature:\norganization: %s\npubkey: %s\nmessage: %s\n"+
				"signature: %s", endorsement.Signer.OrgId, memInfo, hex.Dump(msg), hex.Dump(endorsement.Signature))
			continue
		}
		if _, ok := refinedSigners[memInfo]; !ok {
			refinedSigners[memInfo] = true
			refinedEndorsement = append(refinedEndorsement, endorsement)
		}
	}
	return refinedEndorsement
}

func (p *pkACProvider) getValidEndorsements(orgList map[string]bool, roleList map[protocol.Role]bool,
	endorsements []*common.EndorsementEntry) []*common.EndorsementEntry {
	var refinedEndorsements []*common.EndorsementEntry
	for _, endorsement := range endorsements {
		if len(roleList) == 0 {
			refinedEndorsements = append(refinedEndorsements, endorsement)
			continue
		}

		member := p.getMemberFromCache(endorsement.Signer)
		if member == nil {
			p.log.Debugf(
				"authentication warning: the member is not in member cache, memberInfo[%s]",
				string(endorsement.Signer.MemberInfo))
			continue
		}

		p.log.Debugf("getValidEndorsements: signer's role [%v]", member.GetRole())

		if _, ok := roleList[member.GetRole()]; ok {
			refinedEndorsements = append(refinedEndorsements, endorsement)
		} else {
			p.log.Debugf("authentication warning: signer's role [%v] is not permitted, requires [%v]",
				member.GetRole(), roleList)
		}
	}

	return refinedEndorsements
}

// GetHashAlg return hash algorithm the access control provider uses
func (p *pkACProvider) GetHashAlg() string {
	return p.hashType
}

// ValidateResourcePolicy checks whether the given resource principal is valid
func (p *pkACProvider) ValidateResourcePolicy(resourcePolicy *config.ResourcePolicy) bool {
	return true
}

// LookUpPolicy returns corresponding policy configured for the given resource name
func (p *pkACProvider) LookUpPolicy(resourceName string) (*pbac.Policy, error) {
	blockVersion, policyResourceName := getBlockVersionAndResourceName(resourceName)

	if blockVersion > 0 && blockVersion <= 220 {
		return p.lookUpPolicy220(policyResourceName)
	}

	pol, ok := p.resourceNamePolicyMap.Load(resourceName)
	if !ok {
		return nil, fmt.Errorf("policy not found for resource %s", resourceName)
	}
	pbPolicy := pol.(*policy).GetPbPolicy()
	return pbPolicy, nil
}

// LookUpExceptionalPolicy returns corresponding exceptional policy configured for the given resource name
func (p *pkACProvider) LookUpExceptionalPolicy(resourceName string) (*pbac.Policy, error) {
	blockVersion, policyResourceName := getBlockVersionAndResourceName(resourceName)

	if blockVersion > 0 && blockVersion <= 220 {
		return p.lookUpExceptionalPolicy220(policyResourceName)
	}

	pol, ok := p.exceptionalPolicyMap.Load(resourceName)
	if !ok {
		return nil, fmt.Errorf("exceptional policy not found for resource %s", resourceName)
	}
	pbPolicy := pol.(*policy).GetPbPolicy()
	return pbPolicy, nil
}

// CreatePrincipal creates a principal for one time authentication
func (p *pkACProvider) CreatePrincipal(resourceName string, endorsements []*common.EndorsementEntry,
	message []byte) (protocol.Principal, error) {

	if len(endorsements) == 0 || message == nil {
		return nil, fmt.Errorf("setup access control principal failed, a principal should contain valid (non-empty)" +
			" signer information, signature, and message")
	}
	if endorsements[0] == nil {
		return nil, fmt.Errorf("setup access control principal failed, signer-signature pair should not be nil")
	}
	return &principal{
		resourceName: resourceName,
		endorsement:  endorsements,
		message:      message,
		targetOrg:    "",
	}, nil
}

func (p *pkACProvider) CreatePrincipalForTargetOrg(resourceName string,
	endorsements []*common.EndorsementEntry, message []byte, targetOrgId string) (protocol.Principal, error) {

	return nil, fmt.Errorf("setup access control principal failed, CreatePrincipalForTargetOrg is not supported")
}

// VerifyPrincipal verifies if the principal for the resource is met
func (p *pkACProvider) VerifyPrincipal(principal protocol.Principal) (bool, error) {

	refinedPrincipal, err := p.refinePrincipal(principal)
	if err != nil {
		return false, fmt.Errorf("authentication failed, [%s]", err.Error())
	}

	if localconf.ChainMakerConfig.DebugConfig.IsSkipAccessControl {
		return true, nil
	}

	pol, err := p.lookUpPolicyByResourceName(principal.GetResourceName())
	if err != nil {
		return false, fmt.Errorf("authentication failed, [%s]", err.Error())
	}

	return p.verifyPrincipalPolicy(principal, refinedPrincipal, pol)
}

//GetMemberStatus get the status information of the member
func (p *pkACProvider) GetMemberStatus(member *pbac.Member) (pbac.MemberStatus, error) {
	return pbac.MemberStatus_NORMAL, nil
}

//VerifyRelatedMaterial verify the member's relevant identity material
func (p *pkACProvider) VerifyRelatedMaterial(verifyType pbac.VerifyType, data []byte) (bool, error) {
	return true, nil
}

//GetValidEndorsements filters all endorsement entries and returns all valid ones
func (p *pkACProvider) GetValidEndorsements(principal protocol.Principal) ([]*common.EndorsementEntry, error) {
	refinedPolicy, err := p.refinePrincipal(principal)
	if err != nil {
		return nil, fmt.Errorf("refinePrincipal fail in GetValidEndorsements: [%v]", err)
	}
	endorsements := refinedPolicy.GetEndorsement()

	pol, err := p.lookUpPolicyByResourceName(principal.GetResourceName())
	if err != nil {
		return nil, fmt.Errorf("lookUpPolicyByResourceName fail in GetValidEndorsements: [%v]", err)
	}
	roleListRaw := pol.GetRoleList()
	orgList := map[string]bool{}
	roleList := map[protocol.Role]bool{}
	for _, roleRaw := range roleListRaw {
		roleList[roleRaw] = true
	}
	return p.getValidEndorsements(orgList, roleList, endorsements), nil
}

//GetAllPolicy returns all default policies
func (p *pkACProvider) GetAllPolicy() (map[string]*pbac.Policy, error) {
	var policyMap = make(map[string]*pbac.Policy)
	p.resourceNamePolicyMap.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v, _ := value.(*policy)
		policyMap[k] = newPbPolicyFromPolicy(v)
		return true
	})
	p.exceptionalPolicyMap.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v, _ := value.(*policy)
		policyMap[k] = newPbPolicyFromPolicy(v)
		return true
	})
	return policyMap, nil
}
