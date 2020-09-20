package infobot

import (
	"time"
)

var (
	DefaultGreetingKey     = "greeting"
	DefaultGreetingChannel = ""
)

// TeamSettings holds the settings for a team
type TeamSettings struct {
	NonAdminCreate bool `json:"non_admin_create"`
	NonAdminEdit   bool `json:"non_admin_edit"`
	NonAdminDelete bool `json:"non_admin_delete"`
	//GreetingEnabled bool   `json:"greeting_enabled"`
	//GreetingChannel string `json:"greeting_channel"`
	//GreetingKey     string `json:"greeting_key"`
	revision int
}

// NewTeamSettings returns a new TeamSettings struct
func NewTeamSettings() *TeamSettings {
	return &TeamSettings{
		NonAdminCreate: false,
		NonAdminEdit:   false,
		NonAdminDelete: false,
		//GreetingEnabled: false,
		//GreetingChannel: DefaultGreetingChannel,
		//GreetingKey:     DefaultGreetingKey,
		revision: 0,
	}
}

// Info holds the information related to a particular key/value pair
type Info struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	CreatedBy   string `json:"created_by"`
	CreatedTime int64  `json:"created_time"`
	Edits       []Edit `json:"edits"`
	revision    int
}

// NewInfo returns a new Info struct
func NewInfo(Key, Value, CreatedBy string) *Info {
	var CreatedTime = time.Now().UTC().Unix()

	return &Info{
		Key:         Key,
		Value:       Value,
		CreatedBy:   CreatedBy,
		CreatedTime: CreatedTime,
		Edits: []Edit{
			{
				EditedBy:  CreatedBy,
				Timestamp: CreatedTime,
				NewValue:  Value,
			},
		},
		revision: 0,
	}
}

// Edit holds one edit record for an Info
type Edit struct {
	EditedBy  string `json:"edited_by"`
	Timestamp int64  `json:"timestamp"`
	NewValue  string `json:"new_value"`
}

// NewEdit returns a new Edit struct
func NewEdit(EditedBy, NewValue string) *Edit {
	var Timestamp = time.Now().UTC().Unix()

	return &Edit{
		EditedBy:  EditedBy,
		Timestamp: Timestamp,
		NewValue:  NewValue,
	}
}
