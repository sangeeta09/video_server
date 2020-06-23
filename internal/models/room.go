package models

import "github.com/google/uuid"

//Room struct declaration
type Room struct {
	GUID            uuid.UUID `json:"guid"`
	Name            string    `json:"name"`
	HostUserID      uint64    `json:"host_id"`
	ParticipantList []uint64  `json:"participants"`
	Capacity        int       `json:"capacity"`
}

var UserRoomMap map[uint64][]uuid.UUID

var GuidRoomDetailMap map[uuid.UUID]*Room
