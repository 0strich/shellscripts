package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

// 사원증, 사원증 DID 초기화
func (dcc *DIDChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	employees := []Employee{
		{DocType: "employee", ID: "emp1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Designation: "Software Engineer", DID: ""},
		{DocType: "employee", ID: "emp2", FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Designation: "Project Manager", DID: ""},
	}

	for _, emp := range employees {
		did := generateDID(emp.ID)
		emp.DID = did

		// Create EmployeeDID
		empDID := createEmployeeDID(did)

		empDIDJSON, err := json.Marshal(empDID)
		if err != nil {
			return fmt.Errorf("failed to marshal employee DID JSON: %v", err)
		}

		err = ctx.GetStub().PutState(emp.ID, empDIDJSON)
		if err != nil {
			return fmt.Errorf("failed to put employee DID data: %v", err)
		}

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

// DID String 생성
func generateDID(id string) string {
	return "did:ipid:" + hashString(id)
}

// DID Hash 생성
func hashString(str string) string {
	// TODO: 실제 해싱 알고리즘을 사용하여 문자열 해싱 로직을 작성하세요.
	// 여기에서는 단순히 입력된 문자열을 반환하는 방식으로 대체합니다.
	return str
}

// DID Document 생성
func createEmployeeDID(did string) EmployeeDID {
	return EmployeeDID{
		Context: []string{"https://w3id.org/did/v1"},
		ID:      did,
		PublicKey: []PublicKey{
			{
				ID:           did + "#keys-1",
				Type:         "Ed25519VerificationKey2018",
				PublicKeyHex: "publicKeyHex",
			},
		},
		Service: []Service{
			{
				ID:              did + "#vcr",
				Type:            "VerifiableCredentialRegistry",
				ServiceEndpoint: "https://example.com/vcr",
			},
		},
	}
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


// 사원 정보 조회
func (dcc *DIDChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Employee, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// 사원 DID Document 조회
func (dcc *DIDChaincode) GetDIDDocument(ctx contractapi.TransactionContextInterface, id string) (*EmployeeDID, error) {
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

	didDocument := createEmployeeDID(employee.DID)

	return &didDocument, nil
}

// 사원정보 생성
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

// 사원정보 수정
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

// 사원정보 조회
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

// 사원정보 삭제
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

func main() {
	didChaincode, err := contractapi.NewChaincode(&DIDChaincode{})
	if err != nil {
		log.Fatalf("Error creating DID chaincode: %v", err)
	}

	if err := didChaincode.Start(); err != nil {
		log.Fatalf("Error starting DID chaincode: %v", err)
	}
}