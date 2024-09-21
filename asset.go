package passet

import (
	"context"
	"sync"

	"github.com/KarpelesLab/rest"
	"github.com/KarpelesLab/xuid"
)

type Asset struct {
	Id              xuid.XUID `json:"Unit__"`
	Key             string    // USD, etc
	Name            string
	Symbol          *string // can be null if no symbol available
	SymbolPosition  string  `json:"Symbol_Position"` // before | after
	Decimals        int     `json:",string"`         // force decimals to this value
	DisplayDecimals int     `json:"Display_Decimals,string"`
	Type            string  // currency | crypto_token
}

var (
	assetMap   = make(map[string]*Asset)
	assetMapLk sync.RWMutex
)

// GetAsset returns the informations for a given asset
func GetAsset(k string) *Asset {
	assetMapLk.RLock()
	v, ok := assetMap[k]
	assetMapLk.RUnlock()

	if ok {
		return v
	}

	err := rest.Apply(context.Background(), "Unit/"+k, "GET", nil, &v)
	if err != nil {
		return nil
	}

	if v.Key != k {
		return nil
	}

	assetMapLk.Lock()
	assetMap[k] = v
	assetMapLk.Unlock()

	return v
}
