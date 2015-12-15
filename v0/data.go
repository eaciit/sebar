package sebar

import (
	"github.com/eaciit/errorlib"
	//"reflect"
	"strings"
	"time"
)

/*
Data is struct of data value to be saved on Sebar
*/
type Data struct {
	Key    string
	Type   string // string, int, float, date, time, datetime, blob
	Expiry time.Time
	Value  interface{}
}

/*
ExtractKey extract Data.Key attribute into owner, cluster and key
*/
func (d *Data) ExtractKey() (string, string, string, error) {
	keys := strings.Split(d.Key, ".")
	if len(keys) < 3 {
		return "", "", "", errorlib.Error(packageName, modData, "ExtractKey", "Invalid key to be extracted")
	}

	keylength := len(keys)
	owner := keys[0]
	key := keys[keylength-1]
	cluster := ""
	for _, k := range keys[1:keylength] {
		if cluster == "" {
			cluster = k
		} else {
			cluster += "." + k
		}
	}
	return owner, cluster, key, nil
}
