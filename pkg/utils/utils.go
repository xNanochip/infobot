package utils

import (
	"encoding/json"
	"strings"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// ChannelToString takes a Keybase ChatChannel and returns a string representation of it
func ChannelToString(channel chat1.ChatChannel) string {
	return ""
}

// HasMinRole checks whether a user has a role that is at least as high as the provided role
func HasMinRole(kb *keybase.Keybase, role string, user string, conv chat1.ConvIDStr) bool {
	conversation, err := kb.ListMembersOfConversation(conv)
	if err != nil {
		return false
	}

	memberTypes := make(map[string]struct{})
	memberTypes["owner"] = struct{}{}
	memberTypes["admin"] = struct{}{}
	memberTypes["writer"] = struct{}{}
	memberTypes["reader"] = struct{}{}

	if _, ok := memberTypes[strings.ToLower(role)]; !ok {
		// invalid role
		return false
	}

	for _, member := range conversation.Owners {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "owner" {
		return false
	}

	for _, member := range conversation.Admins {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "admin" {
		return false
	}

	for _, member := range conversation.Writers {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "writer" {
		return false
	}

	for _, member := range conversation.Readers {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "reader" {
		return false
	}

	return false
}

// ToJson marshals an object into a json string
func ToJson(b interface{}) string {
	s, _ := json.Marshal(b)
	return string(s)
}

// ToJson marshals an object into a json string with indenting
func ToJsonPretty(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "  ")
	return string(s)
}
