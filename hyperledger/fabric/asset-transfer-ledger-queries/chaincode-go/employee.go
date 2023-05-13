package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Employee struct {
	DocType string `json:"docType"`
	ID      string `json:"id"`
	DID     string `json:"did"`
}

type EmployeeDID struct {
	ID        string      `json:"id"`
	PublicKey []PublicKey `json:"publicKey"`
}

type PublicKey struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	PublicKeyHex string `json:"publicKeyHex"`
}

type DIDVerificationResult struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

type DIDChaincode struct {
	contractapi.Contract
}

func (dcc *DIDChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	employees := []Employee{
		{DocType: "employee", ID: "emp1", DID: ""},
		{DocType: "employee", ID: "emp2", DID: ""},
	}

	for _, emp := range employees {
		did := generateDID(emp.ID)
		emp.DID = did

		// Create EmployeeDID
		empDID, err := createEmployeeDID(did)
		if err != nil {
			return fmt.Errorf("failed to create employee DID: %v", err)
		}

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

func generateRandomEmployeeID() string {
	employeeIDBytes := make([]byte, 16)
	rand.Read(employeeIDBytes)
	return hex.EncodeToString(employeeIDBytes)
}

func (dcc *DIDChaincode) CreateEmployee(ctx contractapi.TransactionContextInterface, employee *Employee) error {
	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee JSON: %v", err)
	}

	err = ctx.GetStub().PutState(employee.ID, employeeJSON)
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

func (dcc *DIDChaincode) GenerateRandomEmployee(ctx contractapi.TransactionContextInterface) (*Employee, error) {
	employeeID := generateRandomEmployeeID()
	employee := &Employee{
		DocType: "employee",
		ID:      employeeID,
		DID:     generateDID(employeeID),
	}
	return employee, nil
}


func generateDID(id string) string {
	return "did:ipid:" + hashString(id)
}

func hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func createEmployeeDID(did string) (EmployeeDID, error) {
	// Generate public key
	publicKeyBytes := make([]byte, 32)
	_, err := rand.Read(publicKeyBytes)
	if err != nil {
		return EmployeeDID{}, fmt.Errorf("failed to generate public key: %v", err)
	}
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	return EmployeeDID{
		ID: did,
		PublicKey: []PublicKey{
			{
				ID:           did + "#keys-1",
				Type:         "Ed25519VerificationKey2018",
				PublicKeyHex: publicKeyHex,
			},
		},
	}, nil
}

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

	didDocument, err := createEmployeeDID(employee.DID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee DID: %v", err)
	}

	return &didDocument, nil
}

// 사원정보
func (dcc *DIDChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Employee, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var employees []*Employee
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate query result: %v", err)
		}

		employee := new(Employee)
		err = json.Unmarshal(result.Value, employee)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal employee JSON: %v", err)
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

// DID 정보
func (dcc *DIDChaincode) GetDID(ctx contractapi.TransactionContextInterface, id string) (*EmployeeDID, error) {
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

	didDocument, err := createEmployeeDID(employee.DID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee DID: %v", err)
	}

	return &didDocument, nil
}

func (dcc *DIDChaincode) VerifyEmployee(ctx contractapi.TransactionContextInterface, id string) (*DIDVerificationResult, error) {
	// 사원 did document
	employeeDID, err := dcc.GetDIDDocument(ctx, id)
	fmt.Println(employeeDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee did: %v", err)
	}

	// 추가적인 검증 로직 수행
	// ...

	// 검증 결과 생성
	verified := true // 임시로 검증 결과를 true로 설정
	result := &DIDVerificationResult{
		Verified: verified,
		Message: "employee verification result",
	}

	return result, nil
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