package infobot

import (
	"time"

	"samhofi.us/x/keybase/v2/types/chat1"
)

var (
	DefaultGreetingKey     = "greeting"
	DefaultGreetingChannel = ""
)

// Team holds information about a conversation
type Team struct {
	Name     string            `json:"name"`
	ConvIDs  []chat1.ConvIDStr `json:"conv_ids"`
	Settings TeamSettings      `json:"team_settings"`
}

// NewTeam returns a new Team struct
func NewTeam(Name string, ConvIDs ...chat1.ConvIDStr) *Team {
	var teamSettings = NewTeamSettings()

	return &Team{
		Name:     Name,
		ConvIDs:  ConvIDs,
		Settings: *teamSettings,
	}
}

// TeamSettings holds the settings for a team
type TeamSettings struct {
	NonAdminCreate  bool   `json:"non_admin_create"`
	NonAdminEdit    bool   `json:"non_admin_edit"`
	NonAdminDelete  bool   `json:"non_admin_delete"`
	GreetingEnabled bool   `json:"greeting_enabled"`
	GreetingChannel string `json:"greeting_channel"`
	GreetingKey     string `json:"greeting_key"`
}

// NewTeamSettings returns a new TeamSettings struct
func NewTeamSettings() *TeamSettings {
	return &TeamSettings{
		NonAdminCreate:  false,
		NonAdminEdit:    false,
		NonAdminDelete:  false,
		GreetingEnabled: false,
		GreetingChannel: DefaultGreetingChannel,
		GreetingKey:     DefaultGreetingKey,
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
