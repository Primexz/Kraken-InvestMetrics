package xPub

import (
	"fmt"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/modules/blockchain"
	"github.com/Primexz/Kraken-InvestMetrics/util"
	"github.com/checksum0/go-electrum/electrum"
	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/btcec/v2"
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

	xpub, err := hdkeychain.NewKeyFromString(config.C.BitcoinAddress)
	if err != nil {
		log.Fatal("Error parsing xpub:", err)
	}

	externalGapLimit := config.C.BitcoinAddressGapLimit // Path 0/* (receiving addresses)
	internalGapLimit := config.C.BitcoinAddressGapLimit // Path 1/* (change addresses)

	scanDerivationPath := func(pathIndex uint32, gapLimit *int) error {
		consecutiveEmptyAddresses := 0
		addressIndex := uint32(0)

		for {
			derivation := []uint32{pathIndex, addressIndex}

			log.WithFields(logrus.Fields{
				"derivation": derivation,
				"path":       pathIndex,
				"index":      addressIndex,
			}).Debug("Computing bitcoin address")

			// Derive address
			addr, err := x.deriveAddress(xpub, derivation, &chaincfg.MainNetParams)
			if err != nil {
				log.Warn("Failed to derive address:", err)
				addressIndex++
				continue
			}

			// Convert to script hash for Electrum
			scriptHash, err := electrum.AddressToElectrumScriptHash(addr)
			if err != nil {
				log.Error("Failed to convert address to script hash:", err)
				return err
			}

			// Get balance
			balance, err := blockchain.Client.GetBalance(blockchain.Ctx, scriptHash)
			if err != nil {
				log.Error("Failed to resolve address:", err)
				return err
			}

			// Get history to check if address was ever used
			history, err := blockchain.Client.GetHistory(blockchain.Ctx, scriptHash)
			if err != nil {
				log.Error("Failed to get history:", err)
				return err
			}

			confirmed := balance.Confirmed

			// Check if address has been used (has history or balance)
			addressUsed := len(history) > 0 || confirmed > 0

			if addressUsed {
				// Address has been used, add to map and reset gap counter
				bitcoinMap[addr] = confirmed
				consecutiveEmptyAddresses = 0

				// Extend gap limit for this specific path since we found a used address
				*gapLimit = int(addressIndex) + config.C.BitcoinAddressGapLimit + 1

				log.WithFields(logrus.Fields{
					"address":       addr,
					"balance":       confirmed,
					"path":          pathIndex,
					"index":         addressIndex,
					"new_gap_limit": *gapLimit,
				}).Info("Found used address, extending gap limit")
			} else {
				// Address unused, increment gap counter
				consecutiveEmptyAddresses++
			}

			fmt.Printf("Path %d/%d: %s (Balance: %f, Used: %t)\n",
				pathIndex, addressIndex, addr, confirmed, addressUsed)

			// Check if we've reached the gap limit
			if consecutiveEmptyAddresses >= config.C.BitcoinAddressGapLimit {
				log.WithFields(logrus.Fields{
					"path":        pathIndex,
					"final_index": addressIndex,
					"gap_limit":   *gapLimit,
				}).Info("Reached gap limit for derivation path")
				break
			}

			// Safety check: don't scan beyond the extended gap limit
			if int(addressIndex) >= *gapLimit {
				log.WithFields(logrus.Fields{
					"path":        pathIndex,
					"final_index": addressIndex,
					"gap_limit":   *gapLimit,
				}).Info("Reached extended gap limit for derivation path")
				break
			}

			addressIndex++
		}
		return nil
	}

	// Scan external addresses (path 0/*)
	log.Info("Starting scan of external addresses (path 0/*)")
	if err := scanDerivationPath(0, &externalGapLimit); err != nil {
		return nil, fmt.Errorf("error scanning external addresses: %w", err)
	}

	// Scan internal/change addresses (path 1/*)
	log.Info("Starting scan of internal addresses (path 1/*)")
	if err := scanDerivationPath(1, &internalGapLimit); err != nil {
		return nil, fmt.Errorf("error scanning internal addresses: %w", err)
	}

	log.WithFields(logrus.Fields{
		"total_addresses":    len(bitcoinMap),
		"external_gap_limit": externalGapLimit,
		"internal_gap_limit": internalGapLimit,
	}).Info("Address scan completed")

	return bitcoinMap, nil
}

func (x *XPub) deriveAddress(xpub *hdkeychain.ExtendedKey, path []uint32, netParams *chaincfg.Params) (string, error) {
	derivedKey := xpub

	// Derive through each level of the path
	for _, index := range path {
		// For public key derivation, we can only use non-hardened derivation
		// (hardened derivation requires private keys)
		if index >= hdkeychain.HardenedKeyStart {
			return "", fmt.Errorf("cannot derive hardened path from public key: index %d", index)
		}

		var err error
		derivedKey, err = derivedKey.Derive(index)
		if err != nil {
			return "", fmt.Errorf("failed to derive key at index %d: %v", index, err)
		}
	}

	// Get the public key
	pubKey, err := derivedKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %v", err)
	}

	return x.createAddresses(pubKey, netParams)
}

func (x *XPub) createAddresses(pubKey *btcec.PublicKey, netParams *chaincfg.Params) (string, error) {
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal("Failed to create address: ", err)
	}

	return addressWitnessPubKeyHash.EncodeAddress(), nil
}
