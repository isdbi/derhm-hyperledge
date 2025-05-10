package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type RoomAlertContract struct {
	contractapi.Contract
}

type RoomAlert struct {
	Timestamp time.Time `json:"timestamp"`
	LockID    int       `json:"lockid"`
	AlertType string    `json:"alertType"` // e.g., "FORCED_ENTRY", "SENSOR_FAILURE"
	Message   string    `json:"message"`
	Reporter  string    `json:"reporter"`  // optional: "sensor-x", "admin", etc.
}

// EmitAlert stores a new alert in the ledger
func (c *RoomAlertContract) EmitAlert(ctx contractapi.TransactionContextInterface, lockid int, alertType string, message string, reporter string) error {
	timestamp := time.Now()
	alert := RoomAlert{
		Timestamp: timestamp,
		LockID:    lockid,
		AlertType: alertType,
		Message:   message,
		Reporter:  reporter,
	}

	key := fmt.Sprintf("alert_%d_%d", lockid, timestamp.UnixNano())
	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to serialize alert: %v", err)
	}

	log.Printf("ALERT EMITTED: Room %d | Type: %s | Msg: %s", lockid, alertType, message)
	return ctx.GetStub().PutState(key, data)
}

// GetAllAlerts returns all alerts (simple scan — for demo/testing only)
func (c *RoomAlertContract) GetAllAlerts(ctx contractapi.TransactionContextInterface) ([]*RoomAlert, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("alert_", "alert_z")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var alerts []*RoomAlert
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		var alert RoomAlert
		err := json.Unmarshal(queryResponse.Value, &alert)
		if err == nil {
			alerts = append(alerts, &alert)
		}
	}

	return alerts, nil
}
