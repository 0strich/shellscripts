package main

import (
	"encoding/json"
	"fmt"
	"log"
	"bytes"
	"context"
	"io/ioutil"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	shell "github.com/ipfs/go-ipfs-api"
)

type Employee struct {
	DocType 	  string `json:"docType"`
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Designation string `json:"designation"`
	DID         string `json:"did"`
}

// Issuer

type EmployeeDID struct {
	Context    []string    `json:"@context"`
	ID         string      `json:"id"`
	PublicKey  []PublicKey `json:"publicKey"`
	Service    []Service   `json:"service"`
}

type PublicKey struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	PublicKeyHex string `json:"publicKeyHex"`
}

// Verifier

type DIDVerificationResult struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// DIDChaincode defines chaincode methods
type DIDChaincode struct {
	contractapi.Contract
}

const didRegistryChaincodeName = "didregistry"

func (dcc *DIDChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	employees := []Employee{
		{DocType: "employee", ID: "emp1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Designation: "Software Engineer", DID: ""},
		{DocType: "employee", ID: "emp2", FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Designation: "Project Manager", DID: ""},
	}

	for _, emp := range employees {
		empJSON, err := json.Marshal(emp)
		if err != nil {
			return fmt.Errorf("failed to marshal employee JSON: %v", err)
		}

		err = ctx.GetStub().PutState(emp.ID, empJSON)
		if err != nil {
			return fmt.Errorf("failed to put employee data: %v", err)
		}
	}

	return nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Employee, error) {
	var employees []*Employee
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var employee Employee
		err = json.Unmarshal(queryResult.Value, &employee)
		if err != nil {
			return nil, err
		}
		employees = append(employees, &employee)
	}

	return employees, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Employee, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (dcc *DIDChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Employee, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func (dcc *DIDChaincode) CreateEmployee(ctx contractapi.TransactionContextInterface, docType string, id string, firstName string, lastName string, email string, designation string) error {
	existingData, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingData != nil {
		return fmt.Errorf("the employee %s already exists", id)
	}

	employee := Employee{
		DocType:     docType,
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Designation: designation,
		DID:         "",
	}
	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee JSON: %v", err)
	}

	err = ctx.GetStub().PutState(id, employeeJSON)
	if err != nil {
		return fmt.Errorf("failed to put employee data: %v", err)
	}

	return nil
}

func (dcc *DIDChaincode) UpdateEmployee(ctx contractapi.TransactionContextInterface, docType string,id string, firstName string, lastName string, email string, designation string) error {
	existingData, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingData == nil {
		return fmt.Errorf("the employee %s does not exist", id)
	}

	employee := Employee{}
	err = json.Unmarshal(existingData, &employee)
	if err != nil {
		return fmt.Errorf("failed to unmarshal employee JSON: %v", err)
	}

	employee.DocType = docType
	employee.FirstName = firstName
	employee.LastName = lastName
	employee.Email = email
	employee.Designation = designation

	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee JSON: %v", err)
	}

	err = ctx.GetStub().PutState(id, employeeJSON)
	if err != nil {
		return fmt.Errorf("failed to put employee data: %v", err)
	}

	return nil
}

func (dcc *DIDChaincode) GetEmployee(ctx contractapi.TransactionContextInterface, id string) (*Employee, error) {
	employeeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if employeeJSON == nil {
		return nil, fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal employee JSON: %v", err)
	}

	return employee, nil
}

func (dcc *DIDChaincode) DeleteEmployee(ctx contractapi.TransactionContextInterface, id string) error {
	existingData, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingData == nil {
		return fmt.Errorf("the employee %s does not exist", id)
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		return fmt.Errorf("failed to delete employee data: %v", err)
	}

	return nil
}


func (dcc *DIDChaincode) IssueEmployeeDID(ctx contractapi.TransactionContextInterface, id string) error {
	employeeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if employeeJSON == nil {
		return fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	if err != nil {
		return fmt.Errorf("failed to unmarshal employee JSON: %v", err)
	}
	if employee.DID != "" {
		return fmt.Errorf("the employee %s already has a DID", id)
	}

	// DID Document 생성
	didDocument := EmployeeDID{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		ID:      "did:ipfs:" + id,
		PublicKey: []PublicKey{
			{
				ID:           "did:ipfs:" + id + "#keys-1",
				Type:         "Ed25519VerificationKey2018",
				PublicKeyHex: "abc123...",
			},
		},
		Service: []Service{
			{
				ID:              "did:ipfs:" + id + "#vcs",
				Type:            "VerifiableCredentialService",
				ServiceEndpoint: "https://example.com/vc/",
			},
		},
	}

	didJSON, err := json.Marshal(didDocument)
	if err != nil {
		return fmt.Errorf("failed to marshal employee DID JSON: %v", err)
	}

	// IPFS에 DID Document 저장
	sh := shell.NewShell("localhost:5001")
	r := bytes.NewReader(didJSON)
	cid, err := sh.Add(context.Background(), shell.NewReader(r))
	if err != nil {
		return fmt.Errorf("failed to add DID Document to IPFS: %v", err)
	}

	// DID Document ID를 IPFS 해시값으로 변경
	didDocument.ID = "did:ipfs:" + cid

	didJSON, err = json.Marshal(didDocument)
	if err != nil {
		return fmt.Errorf("failed to marshal employee DID JSON: %v", err)
	}

	employee.DID = string(didJSON)

	employeeJSON, err = json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee JSON: %v", err)
	}

	err = ctx.GetStub().PutState(id, employeeJSON)
	if err != nil {
		return fmt.Errorf("failed to put employee data: %v", err)
	}

	return nil
}

func (v *Verifier) VerifyEmployeeDID(did string) (bool, error) {
	// IPFS 연결
	sh := shell.NewShell("localhost:5001")

	// DID Document 가져오기
	didBytes, err := sh.Cat(did)
	if err != nil {
		return false, fmt.Errorf("failed to get DID Document: %v", err)
	}

	// DID Document 파싱
	var didDoc map[string]interface{}
	err = json.Unmarshal(didBytes, &didDoc)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal DID Document: %v", err)
	}

	// 공개키 가져오기
	publicKeys, ok := didDoc["publicKey"].([]interface{})
	if !ok {
		return false, fmt.Errorf("failed to get public keys from DID Document")
	}

	// 공개키로 DID 검증
	for _, publicKey := range publicKeys {
		keyMap, ok := publicKey.(map[string]interface{})
		if !ok {
			continue
		}

		// 공개키 타입이 RSA인 경우만 검증
		keyType, ok := keyMap["type"].(string)
		if !ok || keyType != "RsaVerificationKey2018" {
			continue
		}

		// 공개키 가져오기
		publicKeyStr, ok := keyMap["publicKeyPem"].(string)
		if !ok {
			continue
		}

		// 공개키로 DID 검증
		// TODO: 실제 DID 검증 로직 작성
		// 여기에서는 임시로 true를 반환하는 방식으로 작성합니다.
		return true, nil
	}

	return false, nil
}

func (dcc *DIDChaincode) GetEmployeeDID(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	employeeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if employeeJSON == nil {
		return "", fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal employee JSON: %v", err)
	}

	return employee.DID, nil
}

// func (dcc *DIDChaincode) getPublicKeyFromDIDRegistry(ctx contractapi.TransactionContextInterface, did string) (string) {
// 	result := ctx.GetStub().InvokeChaincode(didRegistryChaincodeName, [][]byte{[]byte("getPublicKey"), []byte(did)}, "")
// 	// if result.err != nil {
// 	// 	return "", fmt.Errorf("failed to invoke DID registry chaincode: %v", err)
// 	// }

// 	// if result.payload == nil {
// 	// 	return "", fmt.Errorf("the public key for DID %s was not found in the DID registry", did)
// 	// }

// 	return string(result)
// }

// func (dcc *DIDChaincode) updatePublicKeyInDIDRegistry(ctx contractapi.TransactionContextInterface, did string, publicKey string) error {
// 	result := ctx.GetStub().InvokeChaincode(didRegistryChaincodeName, [][]byte{[]byte("updatePublicKey"), []byte(did), []byte(publicKey)}, "")
// 	// if result.err != nil {
// 	// 	return fmt.Errorf("failed to invoke DID registry chaincode: %v", err)
// 	// }

// 	// if result.payload == nil {
// 	// 	return fmt.Errorf("failed to update the public key for DID %s in the DID registry", did)
// 	// }

// 	return nil
// }

func main() {
	didChaincode, err := contractapi.NewChaincode(&DIDChaincode{})
	if err != nil {
		log.Fatalf("Error creating DID chaincode: %v", err)
	}

	if err := didChaincode.Start(); err != nil {
		log.Fatalf("Error starting DID chaincode: %v", err)
	}
}