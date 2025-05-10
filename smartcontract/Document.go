package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type IslamicDocumentContract struct {
	contractapi.Contract
}

type DocumentStatus string

const (
	StatusDraft     DocumentStatus = "DRAFT"
	StatusApproved  DocumentStatus = "APPROVED_BY_SHARIA_COUNCIL"
	StatusExecuted  DocumentStatus = "EXECUTED"
	StatusArbitration DocumentStatus = "UNDER_ARBITRATION"
)

type IslamicDocument struct {
	DocumentID      string        `json:"documentId"`
	Timestamp       time.Time     `json:"timestamp"`
	UserID          string        `json:"userId"`       // Customer or Account ID
	DocumentType    string        `json:"documentType"` // e.g., MURABAHA, IJARA, SUKUK, MUDARABA
	Objective       string        `json:"objective"`    // Purpose of the contract
	Parties         []string      `json:"parties"`     // Involved parties (Customer, Bank, Third Party)
	ContractDate    time.Time     `json:"contractDate"`
	ExpiryDate      time.Time     `json:"expiryDate"`
	Amount          float64       `json:"amount"`       // In Halal currency
	Status          DocumentStatus `json:"status"`
	ShariaComplianceCert string    `json:"shariaCert"`  // Reference to compliance certification
	AuthorizedBy    string        `json:"authorizedBy"` // Sharia board member ID
}

// CreateDocument stores a new Islamic financial agreement
func (c *IslamicDocumentContract) CreateDocument(
	ctx contractapi.TransactionContextInterface,
	userId string,
	docType string,
	objective string,
	parties []string,
	contractDate time.Time,
	expiryDate time.Time,
	amount float64,
	shariaCert string,
	authorizedBy string,
) error {
	
	timestamp := time.Now()
	docID := fmt.Sprintf("doc_%s_%s_%d", userId, docType, timestamp.UnixNano())

	document := IslamicDocument{
		DocumentID:      docID,
		Timestamp:       timestamp,
		UserID:          userId,
		DocumentType:    docType,
		Objective:       objective,
		Parties:         parties,
		ContractDate:    contractDate,
		ExpiryDate:      expiryDate,
		Amount:          amount,
		Status:          StatusDraft,
		ShariaComplianceCert: shariaCert,
		AuthorizedBy:    authorizedBy,
	}

	data, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to serialize Islamic document: %v", err)
	}

	log.Printf("ISLAMIC DOCUMENT CREATED: %s | Type: %s | User: %s | Status: %s", 
		docID, docType, userId, StatusDraft)
		
	return ctx.GetStub().PutState(docID, data)
}

// GetAllDocuments returns all financial documents
func (c *IslamicDocumentContract) GetAllDocuments(ctx contractapi.TransactionContextInterface) ([]*IslamicDocument, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("doc_", "doc_z")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var documents []*IslamicDocument
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		var doc IslamicDocument
		err := json.Unmarshal(queryResponse.Value, &doc)
		if err == nil {
			documents = append(documents, &doc)
		}
	}

	return documents, nil
}

// Additional functions could include:
// - UpdateDocumentStatus
// - GetDocumentsByUser
// - VerifyShariaCompliance
// - CalculateProfitDistribution
// - CheckExpiryDates