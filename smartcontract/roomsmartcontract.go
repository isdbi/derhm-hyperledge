package main

import (
	"encoding/json"
	"fmt"
	
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type RoomStatsContract struct {
	contractapi.Contract
}

// RoomAccessStats holds number of accesses per room
type RoomAccessStats struct {
	LockID    int `json:"lockid"`
	EntryCount int `json:"entryCount"`
}

// AddAccess increments room access counter
func (c *RoomStatsContract) AddAccess(ctx contractapi.TransactionContextInterface, lockid int) error {
	key := "room_" + strconv.Itoa(lockid)
	dataBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	var stats RoomAccessStats
	if dataBytes != nil {
		err = json.Unmarshal(dataBytes, &stats)
		if err != nil {
			return fmt.Errorf("failed to parse room data: %v", err)
		}
		stats.EntryCount += 1
	} else {
		stats = RoomAccessStats{
			LockID:     lockid,
			EntryCount: 1,
		}
	}

	newBytes, _ := json.Marshal(stats)
	return ctx.GetStub().PutState(key, newBytes)
}

// GetAccessStats returns the current access count for a room
func (c *RoomStatsContract) GetAccessStats(ctx contractapi.TransactionContextInterface, lockid int) (*RoomAccessStats, error) {
	key := "room_" + strconv.Itoa(lockid)
	dataBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataBytes == nil {
		return nil, fmt.Errorf("room %d not found", lockid)
	}

	var stats RoomAccessStats
	err = json.Unmarshal(dataBytes, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse room data: %v", err)
	}

	return &stats, nil
}
