package xPub

import (
	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/modules/blockchain"
	"github.com/Primexz/Kraken-InvestMetrics/util"
	"github.com/checksum0/go-electrum/electrum"
	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

type XPub struct {
	address string
	account int
}

var (
	log = util.LoggerWithPrefix("xpub")
)

func NewXPub() *XPub {
	addr := config.C.BitcoinAddress
	if !util.IsXPub() {
		panic("Invalid xPub address " + addr)
	}

	return &XPub{
		address: addr,
		account: config.C.BitcoinAccount,
	}
}

func (x *XPub) GetTotalSats() (float64, error) {
	addressMap, err := x.GetAddressSatMap()
	if err != nil {
		log.Error("Failed to get address map:", err)
		return 0, err
	}

	var total float64
	for _, sat := range addressMap {
		total += sat
	}

	return total, nil
}

func (x *XPub) GetAddressSatMap() (map[string]float64, error) {
	bitcoinMap := make(map[string]float64)
	searchEnd := config.C.BitcoinAddressGapLimit

	i := 0
	for i <= searchEnd {
		address := x.getAddressFromIndex(i)
		log.WithFields(logrus.Fields{
			"index":   i,
			"address": address,
		}).Debug("Computing bitcoin address ")

		scriptHash, err := electrum.AddressToElectrumScriptHash(address)
		if err != nil {
			log.Error("Failed to convert address to script hash:", err)
			return nil, err
		}

		balance, err := blockchain.Client.GetBalance(blockchain.Ctx, scriptHash)
		if err != nil {
			log.Error("Failed to resolve address:", err)
			return nil, err
		}

		history, err := blockchain.Client.GetHistory(blockchain.Ctx, scriptHash)
		if err != nil {
			log.Error("Failed to get history:", err)
			return nil, err
		}

		confirmed := balance.Confirmed

		if len(history) > 0 {
			bitcoinMap[address] = confirmed
			searchEnd += 1
		}

		i++
	}

	return bitcoinMap, nil
}

func (x *XPub) getAddressFromIndex(index int) string {
	child, _ := x.getAccountKey().Derive(uint32(index))
	pubKey, _ := child.ECPubKey()

	witnessProg := btcutil.Hash160(pubKey.SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal("Failed to create address: ", err)
	}

	return addressWitnessPubKeyHash.EncodeAddress()
}

func (x *XPub) getAccountKey() *hdkeychain.ExtendedKey {
	masterKey, err := hdkeychain.NewKeyFromString(x.address)
	if err != nil {
		log.Fatal("Failed to create master key: ", err)
	}

	accountKey, err := masterKey.Derive(uint32(x.account))
	if err != nil {
		log.Fatal("Failed to derive account key: ", err)
	}

	return accountKey
}
