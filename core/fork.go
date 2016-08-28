package core

import (
	"bytes"
	"math/big"

	"github.com/ethereumproject/go-ethereum/core/types"
)

type Fork struct {
	Name         string
	Support      bool
	NetworkSplit bool
	// Block is the block number where the hard-fork commences on
	// the Ethereum network.
	Block *big.Int
	// SplitHash to derive BadHashes to assist in avoiding sync issues
	// after network split.
	// ETC Block+1 "94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f"
	SplitHash string
	// ETF Block+1 "4985f5ca3d2afbec36529aa96f74de3cc10a2a4a6c44f2157a57d2c6059a11bb"
	ForkSplitHash string
	// ForkBlockExtra is the block header extra-data field to set for a fork
	// point and a number of consecutive blocks to allow fast/light syncers to correctly
	BlockExtra []byte
	// ForkExtraRange is the number of consecutive blocks from the fork point
	// to override the extra-data in to prevent no-fork attacks.
	ExtraRange *big.Int
	// TODO Derive Oracle contracts from fork struct (Version, Registrar, Release)
}

func (fork *Fork) ValidateForkHeaderExtraData(header *types.Header) error {
	// Short circuit validation if the node doesn't care about the DAO fork
	if fork.Block == nil {
		return nil
	}
	// Make sure the block is within the fork's modified extra-data range
	limit := new(big.Int).Add(fork.Block, fork.ExtraRange)
	if header.Number.Cmp(fork.Block) < 0 || header.Number.Cmp(limit) >= 0 {
		return nil
	}
	// Depending whether we support or oppose the fork, validate the extra-data contents
	if fork.Support {
		if bytes.Compare(header.Extra, fork.BlockExtra) != 0 {
			return ValidationError("Fork bad block extra-data: 0x%x", header.Extra)
		}
	} else {
		if bytes.Compare(header.Extra, fork.BlockExtra) == 0 {
			return ValidationError("No-fork bad block extra-data: 0x%x", header.Extra)
		}
	}
	// All ok, header has the same extra-data we expect
	return nil
}

// TODO Migrate hardcoded fork config into a json file
func LoadForks() []*Fork {
	return []*Fork{
		&Fork{
			Name:         "Homestead",
			NetworkSplit: false,
			Support:      true,
			Block:        big.NewInt(1150000),
		},
		&Fork{
			Name:         "ETF",
			NetworkSplit: true,
			Support:      false,
			Block:        big.NewInt(1920000),
			BlockExtra:   common.FromHex("0x64616f2d686172642d666f726b"),
			ExtraRange:   big.NewInt(10),
		},
	}
}