package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type HashInfo struct {
	Hash        string `json:"hash,omitempty"`
	Name        string `json:"name,omitempty"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
}

type Setting struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type TemplateS struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "upload_hash":
		return t.upload_hash(stub, args)
	case "query_hash":
		return t.query_hash(stub, args)
	case "query_user_hash":
		return t.query_user_hash(stub, args)
	case "upload_setting":
		return t.upload_setting(stub, args)
	case "query_setting":
		return t.query_setting(stub, args)
	case "query_all_setting":
		return t.query_all_setting(stub, args)
	case "upload_temp":
		return t.upload_temp(stub, args)
	case "upload_selected_temp":
		return t.upload_selected_temp(stub, args)
	case "update_temp":
		return t.update_temp(stub, args)
	case "delete_temp":
		return t.delete_temp(stub, args)
	case "query_temp":
		return t.query_temp(stub, args)
	case "query_selected_temp":
		return t.query_selected_temp(stub, args)
	}

	return shim.Error("Invalid invoke function Name. " + function)
}

func (t *SimpleChaincode) upload_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	info := HashInfo{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	value, err := stub.GetState("Hash#0#" + info.Hash)
	if err != nil {
		return shim.Error(err.Error())
	}

	if value != nil {
		return shim.Error("hash existed")
	}

	err = stub.PutState("Hash#0#"+info.Hash, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("User#"+info.Name+"#0#"+info.Hash, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) upload_setting(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	info := Setting{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("setting#0#"+info.Name, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) query_user_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result, err := stub.GetStateByRange("User#"+args[0]+"#0", "User#"+args[0]+"#F")
	if err != nil {
		return shim.Error(err.Error())
	}
	var hashInfoList []HashInfo
	if result != nil {
		if result != nil {
			defer func() {
				result.Close()
			}()

			for result.HasNext() {
				record, err := result.Next()
				if err != nil {
					return shim.Error(err.Error())
				}

				g := HashInfo{}
				err = json.Unmarshal(record.Value, &g)
				if err != nil {
					return shim.Error(err.Error())
				}

				hashInfoList = append(hashInfoList, g)
			}
		}
	}

	m, _ := json.Marshal(hashInfoList)

	return shim.Success(m)
}

func (t *SimpleChaincode) query_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var hashInfo *HashInfo
	result, err := stub.GetState("Hash#0#" + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if result != nil {
		err = json.Unmarshal(result, &hashInfo)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	m, _ := json.Marshal(hashInfo)

	return shim.Success(m)
}

func (t *SimpleChaincode) query_setting(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var setting *Setting
	result, err := stub.GetState("setting#0#" + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if result != nil {
		err = json.Unmarshal(result, &setting)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	m, _ := json.Marshal(setting)

	return shim.Success(m)
}

func (t *SimpleChaincode) query_all_setting(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result, err := stub.GetStateByRange("setting#0", "setting#F")
	if err != nil {
		return shim.Error(err.Error())
	}
	var settingList []Setting
	if result != nil {
		if result != nil {
			defer func() {
				result.Close()
			}()

			for result.HasNext() {
				record, err := result.Next()
				if err != nil {
					return shim.Error(err.Error())
				}

				g := Setting{}
				err = json.Unmarshal(record.Value, &g)
				if err != nil {
					return shim.Error(err.Error())
				}

				settingList = append(settingList, g)
			}
		}
	}

	m, _ := json.Marshal(settingList)

	return shim.Success(m)
}

func (t *SimpleChaincode) upload_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	info := TemplateS{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	value, err := stub.GetState("temp#0#" + info.Name)
	if err != nil {
		return shim.Error(err.Error())
	}

	if value != nil {
		return shim.Error("template existed")
	}

	err = stub.PutState("temp#0#"+info.Name, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) upload_selected_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	info := TemplateS{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("selectedtemp", argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) update_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	info := TemplateS{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("temp#0#"+info.Name, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) delete_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.DelState("temp#0#" + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) query_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result, err := stub.GetStateByRange("temp#0", "temp#F")
	if err != nil {
		return shim.Error(err.Error())
	}
	var templateSList []TemplateS
	if result != nil {
		if result != nil {
			defer func() {
				result.Close()
			}()

			for result.HasNext() {
				record, err := result.Next()
				if err != nil {
					return shim.Error(err.Error())
				}

				g := TemplateS{}
				err = json.Unmarshal(record.Value, &g)
				if err != nil {
					return shim.Error(err.Error())
				}

				templateSList = append(templateSList, g)
			}
		}
	}

	m, _ := json.Marshal(templateSList)

	return shim.Success(m)
}

func (t *SimpleChaincode) query_selected_temp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	value, err := stub.GetState("selectedtemp")
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
