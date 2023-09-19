package keeper_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/gogoproto/proto"

	sdkmath "cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
)

func getTestProposalMessage() sdk.Msg {
	_, _, addr := testdata.KeyTestPubAddr()
	return banktypes.NewMsgSend(authtypes.NewModuleAddress("gov"), addr, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))))
}

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	testedEncodings := []string{icatypes.EncodingProtobuf, icatypes.EncodingProto3JSON}
	var (
		path       *ibctesting.Path
		packetData []byte
	)

	testCases := []struct {
		msg      string
		malleate func(encoding string)
		expPass  bool
	}{
		{
			"interchain account successfully executes an arbitrary message type using the * (allow all message types) param",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				proposal, err := govtypesv1.NewProposal([]sdk.Msg{getTestProposalMessage()}, govtypesv1.DefaultStartingProposalID, time.Now(), time.Now().Add(time.Hour), "test proposal", "title", "Description", sdk.AccAddress(interchainAccountAddr), false)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().GovKeeper.SetProposal(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)
				err = suite.chainB.GetSimApp().GovKeeper.ActivateVotingPeriod(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)

				msg := &govtypesv1.MsgVote{
					ProposalId: govtypesv1.DefaultStartingProposalID,
					Voter:      interchainAccountAddr,
					Option:     govtypesv1.OptionYes,
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{"*"})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes banktypes.MsgSend",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &banktypes.MsgSend{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes stakingtypes.MsgDelegate",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				validatorAddr := (sdk.ValAddress)(suite.chainB.Vals.Validators[0].Address)
				msg := &stakingtypes.MsgDelegate{
					DelegatorAddress: interchainAccountAddr,
					ValidatorAddress: validatorAddr.String(),
					Amount:           sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(5000)),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes stakingtypes.MsgDelegate and stakingtypes.MsgUndelegate sequentially",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				validatorAddr := (sdk.ValAddress)(suite.chainB.Vals.Validators[0].Address)
				msgDelegate := &stakingtypes.MsgDelegate{
					DelegatorAddress: interchainAccountAddr,
					ValidatorAddress: validatorAddr.String(),
					Amount:           sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(5000)),
				}

				msgUndelegate := &stakingtypes.MsgUndelegate{
					DelegatorAddress: interchainAccountAddr,
					ValidatorAddress: validatorAddr.String(),
					Amount:           sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(5000)),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msgDelegate, msgUndelegate}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msgDelegate), sdk.MsgTypeURL(msgUndelegate)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes govtypes.MsgSubmitProposal",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				protoAny, err := codectypes.NewAnyWithValue(getTestProposalMessage())
				suite.Require().NoError(err)

				msg := &govtypesv1.MsgSubmitProposal{
					Messages:       []*codectypes.Any{protoAny},
					InitialDeposit: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(5000))),
					Proposer:       interchainAccountAddr,
					Title:          "test proposal",
					Summary:        "test proposal",
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes govtypes.MsgVote",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				proposal, err := govtypesv1.NewProposal([]sdk.Msg{getTestProposalMessage()}, govtypesv1.DefaultStartingProposalID, time.Now(), time.Now().Add(time.Hour), "test proposal", "title", "Description", sdk.AccAddress(interchainAccountAddr), false)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().GovKeeper.SetProposal(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)
				err = suite.chainB.GetSimApp().GovKeeper.ActivateVotingPeriod(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)

				msg := &govtypesv1.MsgVote{
					ProposalId: govtypesv1.DefaultStartingProposalID,
					Voter:      interchainAccountAddr,
					Option:     govtypesv1.OptionYes,
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes disttypes.MsgFundCommunityPool",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &disttypes.MsgFundCommunityPool{
					Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(5000))),
					Depositor: interchainAccountAddr,
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes disttypes.MsgSetWithdrawAddress",
			func(encoding string) {
				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &disttypes.MsgSetWithdrawAddress{
					DelegatorAddress: interchainAccountAddr,
					WithdrawAddress:  suite.chainB.SenderAccount.GetAddress().String(),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes transfertypes.MsgTransfer",
			func(encoding string) {
				transferPath := ibctesting.NewTransferPath(suite.chainB, suite.chainC)

				suite.coordinator.Setup(transferPath)

				interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &transfertypes.MsgTransfer{
					SourcePort:       transferPath.EndpointA.ChannelConfig.PortID,
					SourceChannel:    transferPath.EndpointA.ChannelID,
					Token:            sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)),
					Sender:           interchainAccountAddr,
					Receiver:         suite.chainA.SenderAccount.GetAddress().String(),
					TimeoutHeight:    suite.chainB.GetTimeoutHeight(),
					TimeoutTimestamp: uint64(0),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"unregistered sdk.Msg",
			func(encoding string) {
				msg := &banktypes.MsgSendResponse{}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{"/" + proto.MessageName(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			false,
		},
		{
			"cannot unmarshal interchain account packet data",
			func(encoding string) {
				packetData = []byte{}
			},
			false,
		},
		{
			"cannot deserialize interchain account packet data messages",
			func(encoding string) {
				data := []byte("invalid packet data")

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()
			},
			false,
		},
		{
			"invalid packet type - UNSPECIFIED",
			func(encoding string) {
				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{&banktypes.MsgSend{}}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.UNSPECIFIED,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()
			},
			false,
		},
		{
			"unauthorised: interchain account not found for controller port ID",
			func(encoding string) {
				path.EndpointA.ChannelConfig.PortID = "invalid-port-id"

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{&banktypes.MsgSend{}}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()
			},
			false,
		},
		{
			"unauthorised: message type not allowed", // NOTE: do not update params to explicitly force the error
			func(encoding string) {
				msg := &banktypes.MsgSend{
					FromAddress: suite.chainB.SenderAccount.GetAddress().String(),
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()
			},
			false,
		},
		{
			"unauthorised: signer address is not the interchain account associated with the controller portID",
			func(encoding string) {
				msg := &banktypes.MsgSend{
					FromAddress: suite.chainB.SenderAccount.GetAddress().String(), // unexpected signer
					ToAddress:   suite.chainB.SenderAccount.GetAddress().String(),
					Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))),
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []proto.Message{msg}, encoding)
				suite.Require().NoError(err)

				icaPacketData := icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				packetData = icaPacketData.GetBytes()

				params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			false,
		},
	}

	for _, encoding := range testedEncodings {
		for _, tc := range testCases {
			tc := tc

			suite.Run(tc.msg, func() {
				suite.SetupTest() // reset

				path = NewICAPath(suite.chainA, suite.chainB, encoding)
				suite.coordinator.SetupConnections(path)

				err := SetupICAPath(path, TestOwnerAddress)
				suite.Require().NoError(err)

				portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
				suite.Require().NoError(err)

				// Get the address of the interchain account stored in state during handshake step
				storedAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, portID)
				suite.Require().True(found)

				icaAddr, err := sdk.AccAddressFromBech32(storedAddr)
				suite.Require().NoError(err)

				// Check if account is created
				interchainAccount := suite.chainB.GetSimApp().AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
				suite.Require().Equal(interchainAccount.GetAddress().String(), storedAddr)

				suite.fundICAWallet(suite.chainB.GetContext(), path.EndpointA.ChannelConfig.PortID, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10000))))

				tc.malleate(encoding) // malleate mutates test data

				packet := channeltypes.NewPacket(
					packetData,
					suite.chainA.SenderAccount.GetSequence(),
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					path.EndpointB.ChannelConfig.PortID,
					path.EndpointB.ChannelID,
					suite.chainB.GetTimeoutHeight(),
					0,
				)

				txResponse, err := suite.chainB.GetSimApp().ICAHostKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)

				if tc.expPass {
					suite.Require().NoError(err)
					suite.Require().NotNil(txResponse)
				} else {
					suite.Require().Error(err)
					suite.Require().Nil(txResponse)
				}
			})
		}
	}
}

func (suite *KeeperTestSuite) TestJSONOnRecvPacket() {
	var (
		path       *ibctesting.Path
		packetData []byte
	)
	interchainAccountAddr := "cosmos15ulrf36d4wdtrtqzkgaan9ylwuhs7k7qz753uk"

	testCases := []struct {
		msg      string
		malleate func(icaAddress string)
		expPass  bool
	}{
		{
			"interchain account successfully executes an arbitrary message type using the * (allow all message types) param",
			func(icaAddress string) {
				proposal, err := govtypesv1.NewProposal([]sdk.Msg{getTestProposalMessage()}, govtypesv1.DefaultStartingProposalID, time.Now(), time.Now().Add(time.Hour), "test proposal", "title", "Description", sdk.AccAddress(interchainAccountAddr), false)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().GovKeeper.SetProposal(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)
				err = suite.chainB.GetSimApp().GovKeeper.ActivateVotingPeriod(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)

				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.gov.v1beta1.MsgVote",
							"voter": "` + icaAddress + `",
							"proposal_id": 1,
							"option": 1
						}
					]
				}`)
				// this is the way cosmwasm encodes byte arrays by default
				// golang doesn't use this encoding by default, but it can still deserialize:
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{"*"})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes banktypes.MsgSend",
			func(icaAddress string) {
				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.bank.v1beta1.MsgSend",
							"from_address": "` + icaAddress + `",
							"to_address": "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
							"amount": [{ "denom": "stake", "amount": "100" }]
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*banktypes.MsgSend)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes govtypes.MsgSubmitProposal",
			func(icaAddress string) {
				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.gov.v1.MsgSubmitProposal",
							"messages": [{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn","to_address":"cosmos1h03ugpwt928rxydxey3dll66zdvfs06t8dq5hl","amount":[{"denom":"stake","amount":"1000"}]}],
                            "initial_deposit":[{"denom":"stake","amount":"5000"}],
                            "proposer":"` + icaAddress + `",
                            "metadata":"",
                            "title":"test proposal",
                            "summary":"test proposal","expedited":false
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*govtypesv1.MsgSubmitProposal)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes govtypes.MsgVote",
			func(icaAddress string) {
				proposal, err := govtypesv1.NewProposal([]sdk.Msg{getTestProposalMessage()}, govtypesv1.DefaultStartingProposalID, suite.chainA.GetContext().BlockTime(), suite.chainA.GetContext().BlockTime(), "test proposal", "title", "Description", sdk.AccAddress(interchainAccountAddr), false)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().GovKeeper.SetProposal(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)
				err = suite.chainB.GetSimApp().GovKeeper.ActivateVotingPeriod(suite.chainB.GetContext(), proposal)
				suite.Require().NoError(err)

				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.gov.v1.MsgVote",
							"voter": "` + icaAddress + `",
							"proposal_id": 1,
							"option": 1
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*govtypesv1.MsgVote)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes govtypes.MsgSubmitProposal, govtypes.MsgDeposit, and then govtypes.MsgVote sequentially",
			func(icaAddress string) {
				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.gov.v1.MsgSubmitProposal",
							"messages": [{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn","to_address":"cosmos1h03ugpwt928rxydxey3dll66zdvfs06t8dq5hl","amount":[{"denom":"stake","amount":"1000"}]}],
                            "initial_deposit":[{"denom":"stake","amount":"5000"}],
                            "proposer":"` + icaAddress + `",
                            "metadata":"",
                            "title":"test proposal",
                            "summary":"test proposal","expedited":false
						},
						{
							"@type": "/cosmos.gov.v1.MsgDeposit",
							"proposal_id": 1,
							"depositor": "` + icaAddress + `",
							"amount": [{ "denom": "stake", "amount": "10000000" }]
						},
						{
							"@type": "/cosmos.gov.v1.MsgVote",
							"voter": "` + icaAddress + `",
							"proposal_id": 1,
							"option": 1
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*govtypesv1.MsgSubmitProposal)(nil)), sdk.MsgTypeURL((*govtypesv1.MsgDeposit)(nil)), sdk.MsgTypeURL((*govtypesv1.MsgVote)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"interchain account successfully executes transfertypes.MsgTransfer",
			func(icaAddress string) {
				transferPath := ibctesting.NewTransferPath(suite.chainB, suite.chainC)

				suite.coordinator.Setup(transferPath)

				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/ibc.applications.transfer.v1.MsgTransfer",
							"source_port": "transfer",
							"source_channel": "channel-1",
							"token": { "denom": "stake", "amount": "100" },
							"sender": "` + icaAddress + `",
							"receiver": "cosmos15ulrf36d4wdtrtqzkgaan9ylwuhs7k7qz753uk",
							"timeout_height": { "revision_number": 1, "revision_height": 100 },
							"timeout_timestamp": 0
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*transfertypes.MsgTransfer)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			true,
		},
		{
			"unregistered sdk.Msg",
			func(icaAddress string) {
				msgBytes := []byte(`{"messages":[{}]}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{"*"})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			false,
		},
		{
			"message type not allowed banktypes.MsgSend",
			func(icaAddress string) {
				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.bank.v1beta1.MsgSend",
							"from_address": "` + icaAddress + `",
							"to_address": "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
							"amount": [{ "denom": "stake", "amount": "100" }]
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*transfertypes.MsgTransfer)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			false,
		},
		{
			"unauthorised: signer address is not the interchain account associated with the controller portID",
			func(icaAddress string) {
				msgBytes := []byte(`{
					"messages": [
						{
							"@type": "/cosmos.bank.v1beta1.MsgSend",
							"from_address": "` + ibctesting.InvalidID + `",
							"to_address": "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
							"amount": [{ "denom": "stake", "amount": "100" }]
						}
					]
				}`)
				byteArrayString := strings.Join(strings.Fields(fmt.Sprint(msgBytes)), ",")

				packetData = []byte(`{
					"type": 1,
					"data":` + byteArrayString + `
				}`)

				params := types.NewParams(true, []string{sdk.MsgTypeURL((*banktypes.MsgSend)(nil))})
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset

			path = NewICAPath(suite.chainA, suite.chainB, icatypes.EncodingProto3JSON)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
			suite.Require().NoError(err)

			// Get the address of the interchain account stored in state during handshake step
			icaAddress, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, portID)
			suite.Require().True(found)

			suite.fundICAWallet(suite.chainB.GetContext(), path.EndpointA.ChannelConfig.PortID, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000))))

			tc.malleate(icaAddress) // malleate mutates test data

			packet := channeltypes.NewPacket(
				packetData,
				suite.chainA.SenderAccount.GetSequence(),
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				suite.chainB.GetTimeoutHeight(),
				0,
			)

			txResponse, err := suite.chainB.GetSimApp().ICAHostKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(txResponse)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(txResponse)
			}
		})
	}
}

func (suite *KeeperTestSuite) fundICAWallet(ctx sdk.Context, portID string, amount sdk.Coins) {
	interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(ctx, ibctesting.FirstConnectionID, portID)
	suite.Require().True(found)

	msgBankSend := &banktypes.MsgSend{
		FromAddress: suite.chainB.SenderAccount.GetAddress().String(),
		ToAddress:   interchainAccountAddr,
		Amount:      amount,
	}

	res, err := suite.chainB.SendMsgs(msgBankSend)
	suite.Require().NotEmpty(res)
	suite.Require().NoError(err)
}
