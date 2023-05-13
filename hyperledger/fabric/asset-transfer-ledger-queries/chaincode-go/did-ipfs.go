package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	shell "github.com/ipfs/go-ipfs-api"
)

// Employee 모델 정의
type Employee struct {
	DocType     string `json:"docType"`
	ID          string `json:"id"`
	KoreanName  string `json:"koreanName"`
	EnglishName string `json:"englishName"`
	Email       string `json:"email"`
	Designation string `json:"designation"`
	DID         string `json:"did"`
}

// DID Document 모델 정의
type DIDDocument struct {
	Context   []string    `json:"@context"`
	ID        string      `json:"id"`
	PublicKey []PublicKey `json:"publicKey"`
	Service   []Service   `json:"service"`
}

// PublicKey 모델 정의
type PublicKey struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	PublicKeyHex string `json:"publicKeyHex"`
}

// Service 모델 정의
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// DIDChaincode 정의
type DIDChaincode struct {
	contractapi.Contract
}

const didRegistryChaincodeName = "didregistry"

func (dcc *DIDChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	employees := []Employee{
		{DocType: "employee", ID: "emp1", KoreanName: "John", EnglishName: "Doe", Email: "john.doe@example.com", Designation: "Software Engineer", DID: ""},
		{DocType: "employee", ID: "emp2", KoreanName: "Jane", EnglishName: "Doe", Email: "jane.doe@example.com", Designation: "Project Manager", DID: ""},
	}

	shell := shell.NewShell("localhost:5001")
	for _, emp := range employees {
		// generate DID
		// DID 형식은 did:ipid:CID
		res, err := shell.Add(strings.NewReader(emp.KoreanName))
		if err != nil {
			return fmt.Errorf("failed to generate DID: %v", err)
		}
		cid := res
		did := fmt.Sprintf("did:ipid:%s", cid)

		// save DID document to IPFS
		// DID Document는 IPFS에 저장
		didDoc := EmployeeDID{
			Context: []string{"https://www.w3.org/ns/did/v1", "https://www.w3.org/2018/credentials/v1"},
			ID:      did,
			PublicKey: []PublicKey{
				{
					ID:           fmt.Sprintf("%s#keys-1", did),
					Type:         "Ed25519VerificationKey2018",
					PublicKeyHex: "dummy public key",
				},
			},
			Service: []Service{
				{
					ID:              fmt.Sprintf("%s#vcs", did),
					Type:            "VerifiableCredentialService",
					ServiceEndpoint: "https://dummy-service-endpoint.com/",
				},
			},
		}
		didDocJSON, err := json.Marshal(didDoc)
		if err != nil {
			return fmt.Errorf("failed to marshal DID document JSON: %v", err)
		}
		didDocCID, err := shell.Add(strings.NewReader(string(didDocJSON)))
		if err != nil {
			return fmt.Errorf("failed to save DID document to IPFS: %v", err)
		}

		// save employee data with DID to the ledger
		emp.DID = fmt.Sprintf("%s#keys-1", did)
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

func (dcc *DIDChaincode) CreateEmployee(ctx contractapi.TransactionContextInterface, docType string, id string, koreanName string, englishName string, email string, designation string) error {
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
		KoreanName:  koreanName,
		EnglishName: englishName,
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

func (dcc *DIDChaincode) UpdateEmployee(ctx contractapi.TransactionContextInterface, docType string, id string, koreanName string, englishName string, email string, designation string) error {
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
	employee.KoreanName = koreanName
	employee.EnglishName = englishName
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
		ID:      "did:example:" + id,
		PublicKey: []PublicKey{
			{
				ID:           "did:example:" + id + "#keys-1",
				Type:         "Ed25519VerificationKey2018",
				PublicKeyHex: "abc123...",
			},
		},
		Service: []Service{
			{
				ID:              "did:example:" + id + "#vcs",
				Type:            "VerifiableCredentialService",
				ServiceEndpoint: "https://example.com/vc/",
			},
		},
	}

	didJSON, err := json.Marshal(didDocument)
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

func (dcc *DIDChaincode) VerifyEmployeeDID(ctx contractapi.TransactionContextInterface, id string, did string) (*DIDVerificationResult, error) {
	employeeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if employeeJSON == nil {
		return &DIDVerificationResult{
			Verified: false,
			Message:  fmt.Sprintf("the employee %s does not exist", id),
		}, nil
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal employee JSON: %v", err)
	}

	if employee.DID != did {
		return &DIDVerificationResult{
			Verified: false,
			Message:  "DID mismatch",
		}, nil
	}

	// TODO: 실제 DID 검증 로직 작성
	// 여기에서는 임시로 true를 반환하는 방식으로 작성합니다.
	return &DIDVerificationResult{
		Verified: true,
		Message:  "DID verified",
	}, nil
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
