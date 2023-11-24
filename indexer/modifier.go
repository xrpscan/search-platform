package indexer

import (
	"github.com/xrpscan/platform/models"
	"github.com/xrpscan/xrpl-go"
)

// Raw transaction object represented as a map-string-interface
type MapStringInterface map[string]interface{}

// Modify transacion object to normalize Amount-like and other fields.
func ModifyTransaction(tx map[string]interface{}) (map[string]interface{}, error) {
	// Detect network from NetworkID field. Default is XRPL Mainnet
	network := xrpl.NetworkXrplMainnet
	networkId, ok := tx["NetworkID"].(int)
	if ok {
		network = xrpl.GetNetwork(networkId)
	}

	// Modify Fee from string to uint64
	fee, ok := tx["Fee"].(uint64)
	if ok {
		tx["Fee"] = fee
	}

	// Modify DestinationTag from integer to uint32
	destinationTag, ok := tx["DestinationTag"].(uint32)
	if ok {
		tx["DestinationTag"] = destinationTag
	}

	// Modify Amount-like fields listed in models.AmountFields
	for _, field := range models.AmountFields {
		ModifyAmount(tx, field.String(), network)
	}

	// Rename tx.metaData property to tx.meta
	metaDataField := "metaData"
	_, ok2 := tx[metaDataField]
	if ok2 {
		tx["meta"] = tx[metaDataField]
		delete(tx, metaDataField)
	}

	// Modify Amount-like fields in meta
	meta, ok := tx["meta"].(map[string]interface{})
	if ok {
		// For simplicity, AffectedNodes field is dropped. This field may indexed
		// in a future release after due consideration.
		delete(meta, "AffectedNodes")
		ModifyAmount(meta, models.DeliveredAmount.String(), network)
		ModifyAmount(meta, models.Delivered_Amount.String(), network)
		tx["meta"] = meta
	}

	return tx, nil
}

func ModifyAmount(tx MapStringInterface, field string, network xrpl.Network) error {
	value, ok := tx[field].(string)
	if ok {
		tx[field] = MapStringInterface{"currency": network.Asset(), "value": value}
	}
	return nil
}
