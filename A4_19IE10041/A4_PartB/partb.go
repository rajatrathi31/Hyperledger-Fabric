/*
Assignment 4
Part B

Name: Rajat Rathi
Roll No.: 19IE10041

*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type BinaryNode struct {
	Data  int	`json: "Data"`
    Left  *BinaryNode	`json: "Left"`
    Right *BinaryNode	`json: "Right"`
}
 
type BinaryTree struct {
	PrimaryKey string	`json: "PrimaryKey"`
    Root *BinaryNode	`json: "Root"`
}

func ReadMyBST(ctx contractapi.TransactionContextInterface) (*BinaryTree, error) {
	bstJSON, err := ctx.GetStub().GetState("bst")
	if err != nil {
		return nil,fmt.Errorf("Failed to read from world state: %v", err)
	}
	if bstJSON == nil {
		return nil,nil
	}

	var bst *BinaryTree
	var asset BinaryTree
	err = json.Unmarshal(bstJSON, &asset)
	bst = &asset
	if err != nil {
		return nil, err
	}
	
	return bst, nil
}

func (s* SmartContract) Insert(ctx contractapi.TransactionContextInterface, val int) error {
	bst, err := ReadMyBST(ctx)
	if err != nil {
		return err
	}

	if bst == nil{
		bst = new(BinaryTree)
		bst.PrimaryKey = "bst"
		bst.Root = new(BinaryNode)
		bst.Root.Data = val
		bst.Root.Left = nil
		bst.Root.Right = nil

		bstJSON, err := json.Marshal(bst)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState("bst", bstJSON)
	} else if bst.Root == nil {
		bst.Root = new(BinaryNode)
		bst.Root.Data = val
		bst.Root.Left = nil
		bst.Root.Right = nil

		bstJSON, err := json.Marshal(bst)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState("bst", bstJSON)
	} else {
		err := UpdateMyBST(ctx, val, bst, 0)
		return err
	}
	return nil
}

func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, val int) error {
	bst, err := ReadMyBST(ctx)
	if err != nil {
		return err
	}

	if bst == nil {
		return fmt.Errorf("There is no BST in the ledger")
	}
	err = UpdateMyBST(ctx, val, bst, 1)
	return err
}

func (s *SmartContract) Preorder(ctx contractapi.TransactionContextInterface) (string, error) {
	bst, err := ReadMyBST(ctx)
	if err != nil {
		return "", err
	}

	if bst == nil {
		return "", fmt.Errorf("There is no BST in the ledger")
	}

	output := (bst.Root).preorderTraversal()
	var traversal string = ""
	for i := 0; i < len(output); i++ {
        traversal = traversal + strconv.Itoa(output[i])
		if i + 1 < len(output) {
			traversal = traversal + ","
		}
    }
	return traversal, err
	// return (bst.Root).preorderTraversal(), nil
}

func (s *SmartContract) Inorder(ctx contractapi.TransactionContextInterface) (string, error) {
	bst, err := ReadMyBST(ctx)
	if err != nil {
		return "", err
	}

	if bst == nil {
		return "", fmt.Errorf("There is no BST in the ledger")
	}

	output := (bst.Root).inorderTraversal()
	var traversal string = ""
	for i := 0; i < len(output); i++ {
        traversal = traversal + strconv.Itoa(output[i])
		if i + 1 < len(output) {
			traversal = traversal + ","
		}
    }
	return traversal, err
}

func (s *SmartContract) TreeHeight(ctx contractapi.TransactionContextInterface) (string, error) {
	bst, err := ReadMyBST(ctx)
	if err != nil {
		return "", err
	}

	if bst == nil {
		return "0", fmt.Errorf("There is no BST in the ledger")
	}
	var result string = strconv.Itoa((bst.Root).heightOfTree())
	return result, nil
}

func UpdateMyBST(ctx contractapi.TransactionContextInterface, val int, bst *BinaryTree, operation int) error {
	if operation == 0 {
		(bst.Root).InsertValue(val)
		bstJSON, err := json.Marshal(bst)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState("bst", bstJSON)
		return nil
	} else {
		var errdel error = nil
		var newbst *BinaryNode
		newbst, errdel = (bst.Root).DeleteValue(val)
		bst.Root = newbst
		bstJSON, err := json.Marshal(bst)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState("bst", bstJSON)
		return errdel
	}
	return nil
}
 
func (n *BinaryNode) InsertValue(val int) {
    if n == nil {
        return
    } else if val < n.Data {
        if n.Left == nil {
            n.Left = &BinaryNode{Data: val, Left: nil, Right: nil}
        } else {
            n.Left.InsertValue(val)
        }
    } else if val > n.Data {
        if n.Right == nil {
            n.Right = &BinaryNode{Data: val, Left: nil, Right: nil}
        } else {
            n.Right.InsertValue(val)
        }
    } else{ 
		return 
	}  
}

func (n *BinaryNode) DeleteValue(val int) (*BinaryNode, error) {
	if n == nil {
		return n, fmt.Errorf("Node not found in the Binary Tree")
	}

	var err error = nil
	if n.Data > val {
		n.Left, err = n.Left.DeleteValue(val)
	}
	if n.Data < val {
		n.Right, err = n.Right.DeleteValue(val)
	}
	if n.Data == val {
		if n.Left == nil && n.Right == nil {
			n = nil
			return n, err
		}
		if n.Left == nil && n.Right != nil {
			temp := n.Right
			n = nil
			n = temp
			return n, err
		}
		if n.Left != nil && n.Right == nil {
			temp := n.Left
			n = nil
			n = temp
			return n, err
		}

		tempNode := minValued(n.Right)
		n.Data = tempNode.Data
		n.Right, err = n.Right.DeleteValue(tempNode.Data)
	}
	return n, err
}

func minValued(Root *BinaryNode) *BinaryNode {
	temp := Root
	for nil != temp && temp.Left != nil {
		temp = temp.Left
	}
	return temp
}

func (n *BinaryNode) inorderTraversal() []int {
	if n == nil {
		return nil
	}

	left := n.Left.inorderTraversal()
	right := n.Right.inorderTraversal()

	result := make([]int, 0)

	result = append(result, left...)
	result = append(result, n.Data)
	result = append(result, right...)
	return result
	
	// var result string = n.Left.inorderTraversal() + "," + strconv.Itoa(n.Data) + "," + n.Right.inorderTraversal()
	// return result
}

func (n *BinaryNode) preorderTraversal() []int {
	if n == nil {
		return nil
	}
	
	left := n.Left.preorderTraversal()
	right := n.Right.preorderTraversal()

	result := make([]int, 0)

	result = append(result, n.Data)
	result = append(result, left...)
	result = append(result, right...)
	return result

	// var result string = strconv.Itoa(n.Data) + "," + n.Left.preorderTraversal() + "," + n.Right.preorderTraversal()
	// return result
}

func (n *BinaryNode) heightOfTree() int {
	var height int = 0
	if n == nil {
		return 0
	}

	var leftheight int = (n.Left).heightOfTree()
	var rightheight int = (n.Right).heightOfTree()
	if leftheight > rightheight {
		height = 1 + leftheight
	} else {
		height = 1 + rightheight
	}
	return height
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