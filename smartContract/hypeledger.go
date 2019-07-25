//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"../model"
//	"github.com/hyperledger/fabric/core/chaincode/shim"
//	pb "github.com/hyperledger/fabric/protos/peer"
//)
//
//type PostChaincode struct{}
//
//
//func (c *PostChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
//	post:=model.Posts{
//		ID: "123456",
//	}
//	postAsJSONBytes, _ := json.Marshal(post)
//	//Add Tom to ledger
//	err := stub.PutState("post", postAsJSONBytes)
//	if err != nil {
//		return shim.Error("Failed to create asset " + post.ID)
//	}
//
//	return shim.Success([]byte("Assets created successfully."))
//}
//
//func (c *PostChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
//	fc, args := stub.GetFunctionAndParameters()
//	if fc == "TransferOwnership" {
//		return c.TransferOwnership(stub, args)
//	}
//	return shim.Error("Called function is not defined in the chaincode ")
//}
//
//// func (d *PostChaincode) getState(stub shim.ChaincodeStubInterface) pb.Response {
//
//// 	donation, cErr := stub.getState("post")
//// 	if cErr != nil {
//// 		return shim.Error(cErr.Error())
//// 	}
//// 	// r := response{CODEALLAOK, "OK", donation}
//// 	return shim.Success([]byte("Asset modified."))
//// }
//
//func (c *PostChaincode) TransferOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	// args[0]=> car serial no
//	// args[1]==> new owner national identity
//	// Read car asset
//	postAsBytes, _ := stub.GetState(args[0])
//	if postAsBytes == nil {
//		return shim.Error("post asset not found")
//	}
//
//	return shim.Success([]byte("post modified."))
//}
//
//func main() {
//	fmt.Println(shim.LogInfo)
//	// Start the chaincode process
//	err := shim.Start(new(PostChaincode))
//	if err != nil {
//		fmt.Println("Error starting PhantomChaincode - ", err.Error())
//	}
//}


/*
 * Smart Contract de ejemplo para HF Tutorial
 *
 * Autor: Antonio Paya Gonzalez
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define la estructura del SmartContract
type SmartContract struct {
}

// Define la estructura de Laptop
type Laptop struct {
	Marca       string `json:"marca"`
	Modelo      string `json:"modelo"`
	Color       string `json:"color"`
	Propietario string `json:"propietario"`
}

/*
 * El metodo Init se llama cuando el Smart Contract se instancia por la red blockchain
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * El metodo Invoke se llama como resultado de ejecutar el Smart Contract
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryLaptop" {
		return s.queryLaptop(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createLaptop" {
		return s.createLaptop(APIstub, args)
	} else if function == "deleteLaptop" {
		return s.deleteLaptop(APIstub, args)
	} else if function == "queryAllLaptops" {
		return s.queryAllLaptops(APIstub)
	} else if function == "cambiarPropietarioLaptop" {
		return s.cambiarPropietarioLaptop(APIstub, args)
	}

	return shim.Error("Nombre de funcion del SmartContract invalido o inexistente.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	laptops := []Laptop{
		Laptop{Marca: "HP", Modelo: "Omen", Color: "black", Propietario: "Microsoft"},
		Laptop{Marca: "Acer", Modelo: "Aspire", Color: "black", Propietario: "Microsoft"},
		Laptop{Marca: "Asus", Modelo: "N551J", Color: "silver", Propietario: "Apple"},
		Laptop{Marca: "Lenovo", Modelo: "80XL", Color: "white", Propietario: "Apple"},
	}

	i := 0
	for i < len(laptops) {
		fmt.Println("i is ", i)
		laptopAsBytes, _ := json.Marshal(laptops[i])
		APIstub.PutState("LAP"+strconv.Itoa(i), laptopAsBytes)
		fmt.Println("Added", laptops[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryLaptop(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Numero incorrecto de argumentos, se esperaba 1")
	}

	laptopAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(laptopAsBytes)
}

func (s *SmartContract) createLaptop(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Numero incorrecto de argumentos, se esperaban 5")
	}

	var laptop = Laptop{Marca: args[1], Modelo: args[2], Color: args[3], Propietario: args[4]}

	laptopAsBytes, _ := json.Marshal(laptop)
	APIstub.PutState(args[0], laptopAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) deleteLaptop(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Numero incorrecto de argumentos, se esperaba 1")
	}

	APIstub.DelState(args[0])
	return shim.Success(nil)
}

func (s *SmartContract) queryAllLaptops(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "LAP0"
	endKey := "LAP999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer es un array JSON con los resultados de la consulta
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Guardarlo como un objeto JSON
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllLaptops:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) cambiarPropietarioLaptop(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Numero incorrecto de argumentos, se esperaban 2")
	}

	laptopAsBytes, _ := APIstub.GetState(args[0])
	laptop := Laptop{}

	json.Unmarshal(laptopAsBytes, &laptop)
	laptop.Propietario = args[1]

	laptopAsBytes, _ = json.Marshal(laptop)
	APIstub.PutState(args[0], laptopAsBytes)

	return shim.Success(nil)
}

// Esta funcion solo es relevante para pruebas unitarias.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error al crear el Smart Contract: %s", err)
	}
}