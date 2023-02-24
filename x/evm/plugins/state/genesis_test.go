// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package state

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pkg.berachain.dev/stargazer/eth/common"
	"pkg.berachain.dev/stargazer/testutil"
	"pkg.berachain.dev/stargazer/x/evm/types"
)

var (
	alice = testutil.Alice
)

var _ = Describe("Genesis", func() {
	var (
		ctx         sdk.Context
		sp          Plugin
		contract    *types.Contract
		slotToValue map[string]string
		atc         map[string]*types.Contract
		htc         map[string]string
	)

	BeforeEach(func() {
		var ak AccountKeeper
		var bk BankKeeper
		ctx, ak, bk, _ = testutil.SetupMinimalKeepers()
		sp = NewPlugin(ak, bk, testutil.EvmKey, "abera", nil)

		// New Contract.
		codeHash := common.HexToHash("0x123")
		code := []byte("code")
		slotToValue = make(map[string]string)
		slotToValue[common.HexToHash("0x456").Hex()] = common.HexToHash("0x789").Hex()
		contract = types.NewContract(codeHash, code, slotToValue)

		// New Address to Contract.
		atc = make(map[string]*types.Contract)
		atc[alice.Hex()] = contract

		// New Hash to Code.
		htc = make(map[string]string)
		htc[codeHash.Hex()] = string(code)

		// Init Genesis.
		genesis := types.NewGenesisState(
			*types.DefaultParams(),
			atc,
			htc,
		)
		sp.InitGenesis(ctx, genesis)
	})

	It("should export current state", func() {
		sp.Reset(ctx)
		sp.SetState(alice, common.HexToHash("0x456"), common.HexToHash("0x789"))
		sp.Finalize()
		var gs types.GenesisState
		sp.ExportGenesis(ctx, &gs)

		Expect(gs.AddressToContract).To(HaveLen(1))
	})
})