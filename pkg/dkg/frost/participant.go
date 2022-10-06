//
// Copyright Coinbase, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

// Package frost is an implementation of the DKG part of  https://eprint.iacr.org/2020/852.pdf
package frost

import (
	"strconv"

	"github.com/coinbase/kryptology/internal"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/sharing"
)

type DkgParticipant struct {
	round             int
	Curve             *curves.Curve
	participantShares map[uint32]*dkgParticipantData
	Id                uint32
	SkShare           curves.Scalar
	VerificationKey   curves.Point
	VkShare           curves.Point
	feldman           *sharing.Feldman
	verifiers         *sharing.FeldmanVerifier
	secretShares      []*sharing.ShamirShare
	ctx               byte
}

type dkgParticipantData struct {
	Id        uint32
	Share     *sharing.ShamirShare
	Verifiers *sharing.FeldmanVerifier
}

func NewDkgParticipant(id, threshold uint32, ctx string, curve *curves.Curve, committee ...uint32) (*DkgParticipant, error) {
	if curve == nil || len(committee) < 1 {
		return nil, internal.ErrNilArguments
	}
	limit := uint32(len(committee))
	feldman, err := sharing.NewFeldman(threshold, limit, curve)
	if err != nil {
		return nil, err
	}
	participantShares := make(map[uint32]*dkgParticipantData, len(committee))
	for _, id := range committee {
		participantShares[id] = &dkgParticipantData{
			Id: id,
		}
	}

	// SetBigInt the common fixed string
	ctxV, _ := strconv.Atoi(ctx)

	return &DkgParticipant{
		Id:                id,
		round:             1,
		Curve:             curve,
		feldman:           feldman,
		participantShares: participantShares,
		ctx:               byte(ctxV),
	}, nil
}
