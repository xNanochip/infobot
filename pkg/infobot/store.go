package infobot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"samhofi.us/x/keybase/v2"
)

var (
	namespacePrefix = "infobot_"
)

// GetKeys returns a slice of all keys for a team
func GetKeys(kb *keybase.Keybase, teamName string) ([]string, error) {
	keys, err := kb.KVListKeys(&teamName, namespacePrefix+"keys")
	if err != nil {
		return []string{}, err
	}

	res := make([]string, len(keys.EntryKeys))
	for i, k := range keys.EntryKeys {
		// TODO: probably shouldn't ignore the error here but it shouldn't be fatal
		// and i'm not too sure how i'd want to handle it just yet... that's a problem
		// for future self i guess
		key, _ := base64.StdEncoding.DecodeString(k.EntryKey)
		res[i] = string(key)
	}

	return res, nil
}

// StringToInfo takes base64 encoded JSON and returns an unmarshaled Info
func StringToInfo(s string) (Info, error) {
	var res Info

	jsonBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// InfoToString returns a marshaled Info
func InfoToString(i Info) (string, error) {
	res, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

// FetchKey fetches a key from the KVStore and returns it as an Info
func FetchKey(kb *keybase.Keybase, teamName, key string) (Info, error) {
	var res Info

	b64Key := base64.StdEncoding.EncodeToString([]byte(strings.ToLower(key)))
	kv, err := kb.KVGet(&teamName, namespacePrefix+"keys", b64Key)
	if err != nil {
		return res, err
	}

	if kv.EntryValue == "" {
		res.revision = kv.Revision
		return res, fmt.Errorf("key does not exist")
	}

	res, err = StringToInfo(kv.EntryValue)
	if err != nil {
		return res, err
	}
	res.revision = kv.Revision

	return res, nil
}

// WriteInfo writes an Info to the KVStore
func WriteInfo(kb *keybase.Keybase, teamName string, info Info) error {
	var err error

	b64Key := base64.StdEncoding.EncodeToString([]byte(strings.ToLower(info.Key)))
	infoString, err := InfoToString(info)
	if err != nil {
		return err
	}

	if info.revision == 0 {
		_, err = kb.KVPut(&teamName, namespacePrefix+"keys", b64Key, infoString)
		return err
	}
	_, err = kb.KVPutWithRevision(&teamName, namespacePrefix+"keys", b64Key, infoString, info.revision+1)
	return err
}

// DeleteKey deletes a key from the KVStore
func DeleteKey(kb *keybase.Keybase, teamName string, info Info) error {
	var err error

	b64Key := base64.StdEncoding.EncodeToString([]byte(strings.ToLower(info.Key)))
	_, err = kb.KVDeleteWithRevision(&teamName, namespacePrefix+"keys", b64Key, info.revision+1)
	return err
}

// WriteNewKey creates a new key and saves it to a team
func WriteNewKey(kb *keybase.Keybase, teamName, key, value, createdBy string) error {
	info := NewInfo(key, value, createdBy)
	return WriteInfo(kb, teamName, *info)
}

// EditKey edits an existing key and saves it to a team
func EditKey(kb *keybase.Keybase, teamName, key, editedBy, newValue string) error {
	edit := NewEdit(editedBy, newValue)
	info, err := FetchKey(kb, teamName, key)
	if err != nil {
		return err
	}

	info.Edits = append(info.Edits, *edit)
	info.Value = newValue

	return WriteInfo(kb, teamName, info)
}

func cleanPercents(s string) string {
	return strings.Replace(s, "%", "%%", -1)
}
