package xPub

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/blockchain"
	"github.com/Primexz/Kraken-InvestMetrics/config"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/primexz/KrakenDCA/logger"
)

type XPub struct {
	address string
	account int
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("xpub")
}

func NewXPub() *XPub {
	addr := config.BitcoinAddress
	if !IsXPub() {
		panic("Invalid xPub address " + addr)
	}

	return &XPub{
		address: addr,
		account: config.BitcoinAccount,
	}
}

func (x *XPub) GetTotalSats() (float64, error) {
	bitcoinArray := make([]float64, 0)
	searchEnd := config.BitcoinAddressGapLimit

	i := 0
	for i <= searchEnd {
		address := x.getAddressFromIndex(i)
		log.Info("Computing bitcoin address ", i, address)

		data, err := blockchain.GetAddressInfo(address)
		if err != nil {
			log.Error("Failed to resolve address:", err)
			return 0, err
		}

		fundedSum := float64(data.ChainStats.FundedTxoSum)
		if fundedSum > 0 {
			bitcoinArray = append(bitcoinArray, fundedSum)
			searchEnd += 1
		}

		i++
		time.Sleep(250 * time.Millisecond)
	}

	return Sum(bitcoinArray), nil
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
