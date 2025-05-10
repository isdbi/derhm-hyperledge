/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type room struct {
	lockid int
	ids    []int
}

type rooms struct {
	rooms []room
}

type Userreturn struct {
	lockid     int
	id         int
	startstime time.Time
	endtime    time.Time
}

// Global room manager
var roomsobj = rooms{
	rooms: []room{
		{ids: []int{1, 2, 3}, lockid: 1},
		{ids: []int{4, 5, 6}, lockid: 2},
		{ids: []int{7, 8, 9}, lockid: 3},
		{ids: []int{10, 11, 12}, lockid: 4},
	},
}

// Global user tracking
var userSessions = make(map[int]*Userreturn) // active sessions
var sessionLogs []Userreturn                 // completed sessions

func (r *rooms) addroom(ids []int, lockid int) {
	r.rooms = append(r.rooms, room{ids: ids, lockid: lockid})
}

// Main logic function
func (S *SmartContract) addOrEditUser(id int, lockid int) {
	for i := 0; i < len(roomsobj.rooms); i++ {
		if roomsobj.rooms[i].lockid == lockid {
			index := indexOf(roomsobj.rooms[i].ids, id)

			if index == -1 {
				// User is entering
				roomsobj.rooms[i].ids = append(roomsobj.rooms[i].ids, id)
				userSessions[id] = &Userreturn{
					lockid:     lockid,
					id:         id,
					startstime: time.Now(),
				}
				log.Printf("User %d ENTERED room %d at %v", id, lockid, userSessions[id].startstime)
			} else {
				// User is exiting
				roomsobj.rooms[i].ids = removeAtIndex(roomsobj.rooms[i].ids, index)

				if session, exists := userSessions[id]; exists {
					session.endtime = time.Now()
					sessionLogs = append(sessionLogs, *session)
					log.Printf("User %d EXITED room %d at %v", id, lockid, session.endtime)
					delete(userSessions, id)
				} else {
					log.Printf("Warning: User %d was in room %d but no session found", id, lockid)
				}
			}
			return
		}
	}

	// Room not found – create new room and treat as entry
	roomsobj.addroom([]int{id}, lockid)
	userSessions[id] = &Userreturn{
		lockid:     lockid,
		id:         id,
		startstime: time.Now(),
	}
	log.Printf("User %d ENTERED new room %d at %v", id, lockid, userSessions[id].startstime)
}

// Helpers
func indexOf(s []int, val int) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func removeAtIndex(s []int, index int) []int {
	if index < 0 || index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}

// Optional: view all session logs
func (S *SmartContract) GetSessionLogs() []Userreturn {
	return sessionLogs
}

// Main entrypoint
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
