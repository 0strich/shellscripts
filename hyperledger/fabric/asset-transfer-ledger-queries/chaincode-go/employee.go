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

// 에러 핸들
func checkError(err error) {
	if err != nil {
		fmt.Errorf("error occurred %v", err)
	}
}

// 원장 초기화
func (dcc *DIDChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	employees := []Employee{
		{DocType: "employee", ID: "olive", DID: ""},
		{DocType: "employee", ID: "austin", DID: ""},
	}

	for _, employee := range employees {
		// DID 생성
		did := generateDID(employee.ID)
		employee.DID = did

		// DID Document 생성 및 저장
		employeeDIDDocument, err := createEmployeeDIDDocument(did)
		checkError(err)

		employeeDIDDocumentJSON, err := json.Marshal(employeeDIDDocument)
		checkError(err)

		err = ctx.GetStub().PutState(employee.ID, employeeDIDDocumentJSON)
		checkError(err)

		// 사원정보 저장
		employeeJSON, err := json.Marshal(employee)
		checkError(err)

		err = ctx.GetStub().PutState(employee.ID, employeeJSON)
		checkError(err)
	}

	return nil
}

func generateRandomEmployeeID() string {
	employeeIDBytes := make([]byte, 16)
	rand.Read(employeeIDBytes)
	return hex.EncodeToString(employeeIDBytes)
}

func (dcc *DIDChaincode) CreateEmployee(ctx contractapi.TransactionContextInterface, docType string, id string) error {
	// 존재 유무 체크
	existingData, err := ctx.GetStub().GetState(id)
	checkError(err)
	if existingData != nil {
		return fmt.Errorf("the employee %s already exists", id)
	}

	// 사원 정보
	employee := Employee{
		DocType: docType,
		ID:      id,
		DID:     "",
	}

	// DID 생성
	did := generateDID(employee.ID)
	employee.DID = did

	// DID Document 생성 및 저장
	employeeDIDDocument, err := createEmployeeDIDDocument(did)
	checkError(err)

	employeeDIDDocumentJSON, err := json.Marshal(employeeDIDDocument)
	checkError(err)

	err = ctx.GetStub().PutState(employee.ID, employeeDIDDocumentJSON)
	checkError(err)

	// 사원정보 저장
	employeeJSON, err := json.Marshal(employee)
	checkError(err)

	err = ctx.GetStub().PutState(employee.ID, employeeJSON)
	checkError(err)

	return nil
}

func (dcc *DIDChaincode) UpdateEmployee(ctx contractapi.TransactionContextInterface, docType string, id string) error {
	existingData, err := ctx.GetStub().GetState(id)
	checkError(err)
	if existingData == nil {
		return fmt.Errorf("the employee %s does not exist", id)
	}

	employee := Employee{}
	err = json.Unmarshal(existingData, &employee)
	checkError(err)

	employee.DocType = docType

	employeeJSON, err := json.Marshal(employee)
	checkError(err)

	err = ctx.GetStub().PutState(id, employeeJSON)
	checkError(err)

	return nil
}

func (dcc *DIDChaincode) GetEmployee(ctx contractapi.TransactionContextInterface, id string) (*Employee, error) {
	employeeJSON, err := ctx.GetStub().GetState(id)
	checkError(err)
	if employeeJSON == nil {
		return nil, fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	checkError(err)

	return employee, nil
}

func (dcc *DIDChaincode) DeleteEmployee(ctx contractapi.TransactionContextInterface, id string) error {
	existingData, err := ctx.GetStub().GetState(id)
	checkError(err)
	if existingData == nil {
		return fmt.Errorf("the employee %s does not exist", id)
	}

	err = ctx.GetStub().DelState(id)
	checkError(err)

	return nil
}

// 랜덤 사원
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
	hash := sha256.Sum256([]byte(id))
	return "did:ipid:" + hex.EncodeToString(hash[:])
}

func createEmployeeDIDDocument(did string) (EmployeeDID, error) {
	// Generate public key
	publicKeyBytes := make([]byte, 32)
	_, err := rand.Read(publicKeyBytes)
	checkError(err)
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
	checkError(err)
	if employeeJSON == nil {
		return nil, fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	checkError(err)

	didDocument, err := createEmployeeDIDDocument(employee.DID)
	checkError(err)

	return &didDocument, nil
}

// 사원정보 조회
func (dcc *DIDChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Employee, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	checkError(err)
	defer resultsIterator.Close()

	var employees []*Employee
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		checkError(err)

		employee := new(Employee)
		err = json.Unmarshal(result.Value, employee)
		checkError(err)

		employees = append(employees, employee)
	}

	return employees, nil
}

// DID 정보
func (dcc *DIDChaincode) GetDID(ctx contractapi.TransactionContextInterface, id string) (*EmployeeDID, error) {
	employeeJSON, err := ctx.GetStub().GetState(id)
	checkError(err)
	if employeeJSON == nil {
		return nil, fmt.Errorf("the employee %s does not exist", id)
	}

	employee := new(Employee)
	err = json.Unmarshal(employeeJSON, employee)
	checkError(err)

	didDocument, err := createEmployeeDIDDocument(employee.DID)
	checkError(err)

	return &didDocument, nil
}

func (dcc *DIDChaincode) VerifyEmployee(ctx contractapi.TransactionContextInterface, id string) (*DIDVerificationResult, error) {
	// 사원 did document
	employeeDID, err := dcc.GetDIDDocument(ctx, id)
	checkError(err)

	// 추가적인 검증 로직 수행
	if len(employeeDID.PublicKey) == 0 {
		return &DIDVerificationResult{
			Verified: false,
			Message:  "Employee DID does not have a public key",
		}, nil
	}

	for _, publicKey := range employeeDID.PublicKey {
		if publicKey.PublicKeyHex == "" {
			return &DIDVerificationResult{
				Verified: false,
				Message:  "Employee DID has an invalid public key",
			}, nil
		}
	}

	// 검증 결과 
	verified := true // 임시로 검증 결과를 true로 설정
	result := &DIDVerificationResult{
		Verified: verified,
		Message:  "employee verification result",
	}

	return result, nil
}

func main() {
	didChaincode, err := contractapi.NewChaincode(&DIDChaincode{})
	checkError(err)

	if err := didChaincode.Start(); err != nil {
		log.Fatalf("Error starting DID chaincode: %v", err)
	}
}
