/*
Assignment 4
Part A

Name: Rajat Rathi
Roll No.: 19IE10041

*/

package main
import (
"fmt"
"log"
"encoding/json"
"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract - provides functions for storing and
// retrieving keys and values from the world state
//
type SmartContract struct {
	contractapi.Contract
}

type OurStruct struct {
	Roll string
	Name string
}

func (s *SmartContract) StudentExists(ctx contractapi.TransactionContextInterface, roll string)(bool, error) {
	assetJSON, err := ctx.GetStub().GetState(roll)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return false, nil
	}

	return true, nil
}

func (s *SmartContract) CreateStudent(ctx contractapi.TransactionContextInterface, roll string, name string) error {
	asset := OurStruct {
		Roll: roll,
		Name: name,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	val, mess := s.StudentExists(ctx, roll)
	if mess != nil {
		return mess
	}

	if val == true{
		return fmt.Errorf("The student %s already exists", roll)
	} else {
		ctx.GetStub().PutState(roll, assetJSON)
	}

	return nil
}

func (s *SmartContract) ReadStudent(ctx contractapi.TransactionContextInterface, roll string)(string, error) {
	assetJSON, err := ctx.GetStub().GetState(roll)
	if err != nil {
		return "", fmt.Errorf("Failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return "", fmt.Errorf("The asset %s does not exist", roll)
	}

	var student OurStruct
	err = json.Unmarshal(assetJSON, &student)
	if err != nil {
		return "", err
	}

	return student.Name, nil
}

func (s *SmartContract) ReadAllStudents(ctx contractapi.TransactionContextInterface)(string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var allstudents string = ""
  	for resultsIterator.HasNext() {
    	queryResponse, err := resultsIterator.Next()
    	if err != nil {
      	return "", err
    	}

		var student OurStruct
		err = json.Unmarshal(queryResponse.Value, &student)
		if err != nil {
		return "", err
		}
		allstudents = allstudents + student.Roll + " " + student.Name + "       "
  	}

  	return allstudents, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating test-network chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting test-network chaincode: %v", err)
	}
}