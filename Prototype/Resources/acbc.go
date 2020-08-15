package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/op/go-logging"
    "strconv"
    "strings"
    "reflect"
    "time"
)

type AcbcChaincode struct {
}

func main() {
    err := shim.Start(new(AcbcChaincode))
    if err != nil {
        logger.Debug(fmt.Sprintf("Error starting AcbcChaincode: %s", err))
    }
}

var logger = logging.MustGetLogger("acbc")


///////////////////////////////////////////////////
//                   FIELDS DEFINITIONS 
///////////////////////////////////////////////////

const (
    UserTable    = "User"
    ProductTable = "Product"
    TransactionTable = "Transaction"
    //DepositBoxTable = "DepositBoxTable"
)

var countryCodes  = []string{"BE", "ES", "FR", "CH", "NO"}

type UserObject struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

type ProductObject struct {
    Id          string `json:"id"`
    Name        string `json:"name"`
    Pseudo      string `json:"pseudo"`
    Title       string `json:"title"`
    Weight      string `json:"weight"`
    TrType      string `json:"trType"`
    DateIn      string `json:"dateIn"`
    DateOut     string `json:"dateOut"`
    PieceStatus string `json:"pieceStatus"`
    Comment     string `json:"comment"`
    Country     string `json:"country"`
    Bank        string `json:"bank"`
    DepositBox  string `json:"depositBox"`
    Album       string `json:"album"`
    PaidPrice   string `json:"paidPrice"`
    QuotedPrice string `json:"quotedPrice"`
    LastTrDate  string `json:"lastTrDate"`
}


type TransactionObject struct {
    Id          string `json:"id" key:"true"`
    TrDate      string `json:"trDate" key:"false"`
    Seller      string `json:"seller" key:"false"`
    TrReference string `json:"trReference" key:"false"`
    IdProduct   string `json:"idProduct" key:"false"`
    NameProduct string `json:"nameProduct" key:"false"`
    Buyer       string `json:"buyer" key:"false"`
    Title       string `json:"title" key:"false"`
    Weight      string `json:"weight" key:"false"`
    TrType      string `json:"trType" key:"false"`
    DateIn      string `json:"dateIn" key:"false"`
    DateOut     string `json:"dateOut" key:"false"`
    PieceStatus string `json:"pieceStatus" key:"false"`
    Comment     string `json:"comment" key:"false"`
    Country     string `json:"country" key:"false"`
    Bank        string `json:"bank" key:"false"`
    DepositBox  string `json:"depositBox" key:"false"`
    Album       string `json:"album" key:"false"`
    PaidPrice   string `json:"paidPrice" key:"false"`
    QuotedPrice string `json:"quotedPrice" key:"false"`
}

/*
type DepositBoxObject struct {
    DepositBox  string `json:"depositBox" key:"true"`
    Album       string `json:"album" key:"true"`
    IdProduct   string `json:"id" key:"true"` 
}*/


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                        INIT PART
//
// Init is Hyperledger fabric's function ( chaincode shim interface ) which called when we first deploy the chaincode. We use Init to create 
// the Hyperledger Fabric's tables used by ACBC API. 
// The Init function take in a 'stub', which used to read from and write to the ledger, a function name, and an array of strings.  
// Init function returning nothing.
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *AcbcChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with function %v and args %v", "Init", function, args))

    initUserTable(stub)
    initProductTable(stub)
    initTransactionTable(stub)
    //initDepositBoxTable(stub)

    return nil, nil
}

func initUserTable(stub shim.ChaincodeStubInterface){ 

    var object UserObject
	typeObject := reflect.TypeOf(object)
    var columnsTableDef []*shim.ColumnDefinition
    columnOneTableDef := shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true}
    columnsTableDef = append(columnsTableDef, &columnOneTableDef)
    for i := 1; i < typeObject.NumField(); i++ {
        columnTableDef := shim.ColumnDefinition{Name: typeObject.Field(i).Tag.Get("json") , Type: shim.ColumnDefinition_STRING, Key: false}
        columnsTableDef = append(columnsTableDef, &columnTableDef)
	}
	err := stub.CreateTable(UserTable, columnsTableDef)
    panicErr(err)
}

func initProductTable(stub shim.ChaincodeStubInterface){ 

    var object ProductObject
	typeObject := reflect.TypeOf(object)
	var columnsTableDef []*shim.ColumnDefinition
    columnOneTableDef := shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true}
    columnsTableDef = append(columnsTableDef, &columnOneTableDef)
	for i := 1; i < typeObject.NumField(); i++ {
        columnTableDef := shim.ColumnDefinition{Name: typeObject.Field(i).Tag.Get("json") , Type: shim.ColumnDefinition_STRING, Key: false}
		columnsTableDef = append(columnsTableDef, &columnTableDef)
	}
	err := stub.CreateTable(ProductTable, columnsTableDef)
    panicErr(err)
}

func initTransactionTable(stub shim.ChaincodeStubInterface){ 

    var object TransactionObject
    typeObject := reflect.TypeOf(object)
    var columnsTableDef []*shim.ColumnDefinition
    columnOneTableDef := shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true}
    columnsTableDef = append(columnsTableDef, &columnOneTableDef)
    for i := 1; i < typeObject.NumField(); i++ {
        columnTableDef := shim.ColumnDefinition{Name: typeObject.Field(i).Tag.Get("json") , Type: shim.ColumnDefinition_STRING, Key: false}
        columnsTableDef = append(columnsTableDef, &columnTableDef)
    }  
    err := stub.CreateTable(TransactionTable, columnsTableDef)
    panicErr(err)
}

/*func initDepositBoxTable(stub shim.ChaincodeStubInterface){ 

    var object DepositBoxObject
    typeObject := reflect.TypeOf(object)
    var columnsTableDef []*shim.ColumnDefinition
    for i := 0; i < typeObject.NumField(); i++ {
        columnTableDef := shim.ColumnDefinition{Name: typeObject.Field(i).Tag.Get("json") , Type: shim.ColumnDefinition_STRING, Key: true}
        columnsTableDef = append(columnsTableDef, &columnTableDef)
    }
    err := stub.CreateTable(DepositBoxTable, columnsTableDef)
    panicErr(err)
}*/



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                           INVOKE PART
//
// Invoke is Hyperledger fabric's function ( chaincode shim interface ) which called when we want to call chaincode functions to do real work.
// Invocations will be captured as a transactions, which get grouped into blocks on the chain. When we need to update the ledger, we will
// do so by invoking our chaincode. 
// The structure of Invoke is simple. It take in a 'stub', a function and an array of arguments. Based on what function was passed in 
// through the function parameter in the invoke request, Invoke will  call custom function or return an error.
// 
// ACBC Api makes the following custom functions available that permit to modify the ledger state:
// - PostUser - used to register a user 
// - PostProduct - used to register a product, called by PostTransaction function
// - PostTransaction - used to register a transaction
// - TransferProduct - used to transfer a product (could be only used to force transfer)
// 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *AcbcChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) (response []byte, err error) {
    logger.Debug(fmt.Sprintf("%v called with function %v and args %v", "Invoke", function, args))

    defer func() {
        if r := recover(); r != nil {
            if errorFromCallee, ok := r.(error); ok {
                logger.Debug(errorFromCallee)
                err = errorFromCallee
            }
        }
    }()

    switch args[0] {
    case "PostUser":
        checkArgsNumberOrPanic(args, 3)
        return t.PostUser(stub, args[1:])
    case "PostTransaction":
        checkArgsNumberOrPanic(args, 20)
        return t.PostTransaction(stub, args[1:])
    case "TransferProduct":
        checkArgsNumberOrPanic(args, 3)
        return t.TransferProduct(stub, args[1:])
    }
    return nil, errors.New("Couldn't find requested method")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// PostUser - used to register a user 
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["PostUser", "3656469441", "MAXAVE"]}'
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) PostUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "PostUser", args))
 
    var userObject UserObject
    userObject = UserObject{args[0], args[1]}
    typeObject := reflect.TypeOf(userObject)

    var columns []*shim.Column  
    for i := 0; i < typeObject.NumField(); i++ {
        colNext := shim.Column{Value: &shim.Column_String_{String_: args[i]}}
        columns = append(columns, &colNext)
    }
    rowToInsert := shim.Row{Columns: columns}
    _, err := stub.InsertRow(UserTable, rowToInsert)
    panicErr(err)
    return userRowToJSON(rowToInsert), nil
}


////////////////////////////////////////////////////////////////////////////////////////////////////////
// PostProduct - used to register a product, called by PostTransaction function
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["PostProduct", "1288915394", "123", 
// "4977", "Vera Valor 1 once (LSP) 2012 - 6 langues PROOF 2012 LSP", "3656469441", "900.000%", "6.45 g",
// "Acceptation don cadeau", "1230768000", "1230768000", "message", "Belgique", "BELGIQUE-CBC BRUXELLES",
// "3003", "3543", "118.00", "118.80"]}'
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) PostProduct(stub shim.ChaincodeStubInterface, args []string) (string) {
    logger.Debug(fmt.Sprintf("%v called with %v", "PostProduct", args))

    idproduit:=args[2]
    var previousOwner string
    previousOwner=""
    var previousAlbum string
    previousAlbum=""

    // Map arguments with the shim.Column and create shim.Row
    rowToInsert := shim.Row{
        Columns: []*shim.Column{
            &shim.Column{Value: &shim.Column_String_{String_: args[2]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[3]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[4]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[5]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[6]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[7]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[8]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[9]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[10]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[11]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[12]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[13]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[14]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[15]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[16]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[17]}},
            &shim.Column{Value: &shim.Column_String_{String_: args[0]}}},
            }

    // Insert created row to the Product table        
    ok, err := stub.InsertRow(ProductTable, rowToInsert)

    // If a row with the given ID Product exists we start the update procedure 
    if !ok {
        logger.Debug(fmt.Sprintf("Product row already present"))

        // First we retrieve the row which contains the given ID from the Product table 
        var searchColumns []shim.Column
        searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: idproduit}})
        foundRow, err := stub.GetRow(ProductTable, searchColumns)
        panicErr(err)

        // Then we handle user's possesions:
        // We found the current Owner of the Item
        var currentOwner string
        currentOwner=foundRow.Columns[2].GetString_()
        previousOwner=currentOwner
        // We retrieve the string of users possessions (value) from the State for the current Owner (key)
        usersPossessions, err := stub.GetState(currentOwner)
        panicErr(err)
        // We delete the given ID Product from the string of users possessions 
        result:=strings.Replace(string(usersPossessions[:]), "/"+idproduit+"/", "/", -1)
        // We save the updated string (value) to the State
        // key/value State: key   = ID User, 
        //                  value = ID Product1/ID Product2/ID Product3
        err = stub.PutState(currentOwner, []byte(result))
        panicErr(err)

        // Thereafter, we update album content if necessary
        var currentAlbum string
        currentAlbum=foundRow.Columns[13].GetString_()
        if len(currentAlbum)!=0 {
            previousAlbum="album-"+currentAlbum
            fmt.Println(previousAlbum, "previousAlbum")
            // We retrieve the string of IDs Product (value) from the State "album-" for the Album (key)
            albumContent, err := stub.GetState(previousAlbum)
            panicErr(err)
            // We delete the given ID Product from the string of IDs 
            resultingContent:=strings.Replace(string(albumContent[:]), "/"+idproduit+"/", "/", -1)
            // We save the updated string (value) to the State
            // key/value State "album-": key   = ID Album, 
            //                           value = ID Product1/ID Product2/ID Product3
            err = stub.PutState(previousAlbum, []byte(resultingContent))
            panicErr(err)
        }

        // Finally, we update depositbox content
        // TODO update only if different from current box 
        // We retrieve the IDs of Deposit Box (value) from the State "albumToBoxPivot-" for the Album (key)
        currentBox, err := stub.GetState("albumToBoxPivot-"+currentAlbum)
        panicErr(err)
        if currentBox[:]!=nil {
            // We retrieve the string of IDs Albums (value) from the State "depositBox-" for the Deposit Box (key)
            previousBoxContent, err:=stub.GetState("depositBox-"+string(currentBox[:]))
            panicErr(err)
            // We delete the given ID Album from the string of IDs 
            newVal:=strings.Replace(string(previousBoxContent[:]), "/"+currentAlbum+"/", "/", -1)
            // We save the updated string (value) to the State
            // key/value State "depositBox-": key   = ID DepositBox, 
            //                                value = ID Album1/ID Album2/ID Album3
            err = stub.PutState("depositBox-"+string(currentBox[:]), []byte(newVal))
            panicErr(err)
        }

        _, err = stub.ReplaceRow(ProductTable, rowToInsert)
        panicErr(err)
    }

    panicErr(err)


    // Below we update the States independently from that if the Product row already present or not"
    // Update the product Owner
    var newOwner string
    newOwner=args[4]
    // We retrieve the string of users possessions (value) from the State for the new Owner (key)
    newUsersPossessions, err := stub.GetState(newOwner)
    panicErr(err)
    if newUsersPossessions != nil {
        // If the new Owner has other Product Item we add the given Id Product to the string  
        // and we save the updated string (value) to the State
        stub.PutState(newOwner, []byte(string(newUsersPossessions[:])+""+idproduit+"/"))
    } else {
        // If it's first Product Item for the new Owner we add the given Id Product (value) to the State
        stub.PutState(newOwner, []byte("/"+idproduit+"/"))
    }

    // Update the containing Album
    var newAlbum string
    newAlbum=args[15]
    // We retrieve the string of IDs Product (value) from the State "album-" for the new Album (key) 
    newAlbumContent, err := stub.GetState("album-"+newAlbum)
    panicErr(err)
    if len(newAlbum) != 0 {
        if newAlbumContent != nil {
            // If the Album has other Product Item we add the given Id Product to the string  
            // and we save the updated string (value) to the State
            stub.PutState("album-"+newAlbum, []byte(string(newAlbumContent[:])+""+idproduit+"/"))
        } else {
            // If it's first Product Item for the Album we add the given Id Product (value) to the State
            stub.PutState("album-"+newAlbum, []byte("/"+idproduit+"/"))
        }
    }

    // Update the containing Deposit Box
    var newDepositBox string
    newDepositBox=args[14]
    // We retrieve the string of IDs Albums (value) from the State "depositBox-" for the Deposit Box (key)
    newDepositBoxContent, err := stub.GetState("depositBox-"+newDepositBox)
    panicErr(err)
    if len(newDepositBox) != 0 {
        if newDepositBoxContent != nil {
            // If the Deposit Box has other Albums  we add the given Id Album to the string 
            // and we save the updated string (value) to the State
            stub.PutState("depositBox-"+newDepositBox, []byte(string(newDepositBoxContent[:])+""+newAlbum+"/"))
        } else {
            // If it's first Album for the Deposit Box we add the given ID Album (value) to the State
            stub.PutState("depositBox-"+newDepositBox, []byte("/"+newAlbum+"/"))
        }
        // Also we fill out the State which keep the pair Album/Deposit Box
        // key/value State "albumToBoxPivot-": key   = ID Album, 
        //                                value = ID Deposit Box
        stub.PutState("albumToBoxPivot-"+newAlbum, []byte(newDepositBox)) 
    }

    return previousOwner
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// PostTransaction - used to register a transaction 
// Example of call:
// peer chaincode invoke -n $NAME -c '{"Function": "invoke", "Args": ["PostTransaction", "100", "1288915394",
// "100", "4977", "Vera Valor 1 once (LSP) 2012 - 6 langues PROOF 2012 LSP", "MAXAVE", "900.000", "6.45 g",
// "Acceptation don cadeau", "1230768000", "1230768000", "4", "message", "Belgique", "BELGIQUE-CBC BRUXELLES",
// "3003", "3543", "118.00", "118.80"]}'
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) PostTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "PostTransaction", args))

    var trDateStr string
    var transactionId string
    var stateValue string
    var countryCode string
    var transactionKey string

    switch {

    case args[13]=="Belgique":

        countryCode = "BE"

    case args[13]=="Espagne":

        countryCode = "ES"

    case args[13]=="France":

        countryCode = "FR"

    case args[13]=="Suisse":

        countryCode = "CH"

    default:

        countryCode = "NO"
    }
    
    // Transaction date is converted into date format from timestamp
    trDateInt, err := strconv.ParseInt(args[1], 10, 64)
    panicErr(err)
    trDate := time.Unix(trDateInt, 0)
    const layout = "20060102" 
    trDateStr = trDate.Format(layout)

    // On the current block we update key/value State "transaction-"
    // Create the key - Transaction date + country code
    transactionKey="transaction-"+trDateStr+"/"+countryCode
    transactionId=args[0]
    fmt.Println(transactionKey, transactionId, "key, value, PostTransaction")
    // We retrieve the string of Transaction IDs (value) from the State "transaction-" for the given key
    transactionIds, err := stub.GetState(transactionKey)
    panicErr(err)
    if transactionIds != nil {
        // If there are other IDs for the given key we add the Id Transaction to the string 
        stateValue=string(transactionIds[:])+""+transactionId+"/"
    } else {
        // If not we assign the current Id Transaction to stateValue
        stateValue="/"+transactionId+"/"
    }
    fmt.Println(stateValue, "<= State value is")
    // And we save the updated/created string (value) to the State
    // key/value State "transaction-": key   = Transaction date + country code 
    //                                 value = ID Transaction1/ID Transaction2/ID Transaction3
    err = stub.PutState(transactionKey, []byte(stateValue))
    panicErr(err)


    // Deposit Box part
    // Current part is not used currently due to the volumetry issue (see comment to the GetDepositBox method )

    /*var searchColumnsDB []shim.Column
    searchColumnsDB = append(searchColumnsDB, shim.Column{Value: &shim.Column_String_{String_: args[2]}})
    foundRowDB, err := stub.GetRow(ProductTable, searchColumnsDB)
    if err!=nil {
        panic(err)
    }

    if len(foundRowDB.Columns[:]) != 0 {
        if (foundRowDB.Columns[12].GetString_()!=args[14])||(foundRowDB.Columns[13].GetString_()!=args[15]){
      
            var rowToDeleteDB []shim.Column
            col1 := shim.Column{Value: &shim.Column_String_{String_: foundRowDB.Columns[12].GetString_()}}
            rowToDeleteDB = append(rowToDeleteDB, col1)
            col2 := shim.Column{Value: &shim.Column_String_{String_: foundRowDB.Columns[13].GetString_()}}
            rowToDeleteDB = append(rowToDeleteDB, col2)
            col3 := shim.Column{Value: &shim.Column_String_{String_: foundRowDB.Columns[0].GetString_()}}
            rowToDeleteDB = append(rowToDeleteDB, col3)


            err := stub.DeleteRow(DepositBoxTable, rowToDeleteDB)
            if err!= nil {
                return nil, fmt.Errorf("deleteRowoperation failed. %s", err)
            } 
        }
    }

    rowToInsertDB := shim.Row{
    Columns: []*shim.Column{
        {Value: &shim.Column_String_{String_: args[14]}},
        {Value: &shim.Column_String_{String_: args[15]}},
        {Value: &shim.Column_String_{String_: args[2]}},
    },}

    ok, err := stub.InsertRow(DepositBoxTable, rowToInsertDB)
    panicErr(err)
    if !ok {
        logger.Debug(fmt.Sprintf("Deposit row already present"))
    } 

    if err!=nil {
        panic(err)
    }*/


    // Below we call PostProduct method which create Product object and insert it to the Product Table or update existing row
    // PostProduct method return Seller (previous Owner of Product item)
    seller := t.PostProduct(stub, args[1:])

    // Map arguments with the shim.Column and create shim.Row
    var transactionObject TransactionObject
    transactionObject = TransactionObject{args[0], args[1], seller, args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9],args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18]}
    typeObject := reflect.TypeOf(transactionObject)

    var columns []*shim.Column  
    var sellerFlag = 0
    for i := 0; i < typeObject.NumField()-1; i++ {
        if i == 2 && sellerFlag!= 1 {
            sellerFlag = 1 
            var colNext shim.Column         
            colNext = shim.Column{Value: &shim.Column_String_{String_: seller}}
            columns = append(columns, &colNext)
            i=i-1
        } else {
            var colNext shim.Column
            colNext = shim.Column{Value: &shim.Column_String_{String_: args[i]}}
            columns = append(columns, &colNext)
       }
    }
    rowToInsert := shim.Row{Columns: columns}

    // Insert created row to the Product table
     _, err = stub.InsertRow(TransactionTable, rowToInsert)   
    panicErr(err)

    return transactionRowToJSON(rowToInsert), nil 
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// TransferProduct - used to transfer a product (could be only used to force transfer)
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["TransferProduct", "UZAVAV"]}'
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) TransferProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "TransferProduct", args))

    var searchColumns []shim.Column
    searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: args[0]}})

    foundRow, err := stub.GetRow(ProductTable, searchColumns)
    panicErr(err)

    foundRow.Columns[2] = &shim.Column{Value: &shim.Column_String_{String_: args[1]}}

    _, err = stub.ReplaceRow(ProductTable, foundRow)
    panicErr(err)

    return productRowToJSON(foundRow), nil
}



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                           QUERY PART
//
// Query is Hyperledger fabric's function ( chaincode shim interface ) which called whenever we query the chaincode's state. 
// It take in a 'stub', a function and an array of arguments. Based on what function was passed in 
// through the function parameter in the query request, Query will  call custom function or return an error.
// 
// ACBC Api makes the following custom functions available that permit to modify the ledger state:
// - GetUser - retrieve User information by User ID
// - GetProduct - retrieve a product Item by Product ID
// - GetProductsByUser - retrieve all Product IDs of product Items held by the given User 
// - GetProductsByAlbum - retrieve all Product IDs of product Items which belongs to the given Album
// - GetProductsByDepositBox - retrieve all Product IDs of product Items which belongs to the given DepositBox 
// - GetAlbumsByDepositBox - retrieve all Album (IDs) which belongs to the given DepositBox
// - GetDepositBox - retrieve all Album and Product IDs which belongs to the given DepositBow (method is not used currently due to the volumetry issue) 
// - GetTransactionById - retrieve Transaction information by Transaction ID
// - GetTransactionsBetweenDates - retrieve all Transaction IDs between given dates (date of Transaction) 
// - GetTransactionsBetweenDatesByCountry - retrieve all Transaction IDs between given dates which correspond to Items 
// - GetTransactionsBetweenDatesForUser - retrieve all Transaction IDs between given dates which correspond to the ownership of given User
// 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *AcbcChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) (response []byte, err error) {
    logger.Debug(fmt.Sprintf("%v called with function %v and args %v", "Query", function, args))

    defer func() {
        if r := recover(); r != nil {
            if errorFromCallee, ok := r.(error); ok {
                logger.Debug(errorFromCallee)
                err = errorFromCallee
            }
        }
    }()

    switch args[0] {
    case "GetUser":
        checkArgsNumberOrPanic(args, 2)
        return t.GetUser(stub, args[1:])
    case "GetProduct":
        checkArgsNumberOrPanic(args, 2)
        return t.GetProduct(stub, args[1:])
    case "GetProductsByUser":
        checkArgsNumberOrPanic(args, 2)
        return t.GetProductsByUser(stub, args[1:])
    case "GetProductsByAlbum":
        checkArgsNumberOrPanic(args, 2)
        return t.GetProductsByAlbum(stub, args[1:])
    case "GetProductsByDepositBox":
        checkArgsNumberOrPanic(args, 2)
        return t.GetProductsByDepositBox(stub, args[1:])
    /*case "GetDepositBox":
        checkArgsNumberOrPanic(args, 2)
        return t.GetDepositBox(stub, args[1:])*/
    case "GetAlbumsByDepositBox":
        checkArgsNumberOrPanic(args, 2)
        return t.GetAlbumsByDepositBox(stub, args[1:])
    case "GetTransactionById":
        checkArgsNumberOrPanic(args, 2)
        return t.GetTransactionById(stub, args[1:])
    case "GetTransactionsBetweenDates":
        checkArgsNumberOrPanic(args, 3)
        return t.GetTransactionsBetweenDates(stub, args[1:]) 
    case "GetTransactionsBetweenDatesByCountry":
        checkArgsNumberOrPanic(args, 4)
        return t.GetTransactionsBetweenDatesByCountry(stub, args[1:])
    case "GetTransactionsBetweenDatesForUser":
        checkArgsNumberOrPanic(args, 4)
        return t.GetTransactionsBetweenDatesForUser(stub, args[1:])
    }
    return nil, errors.New("Couldn't find requested method")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetUser - retrieve User information by User ID
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetUser", "3656469441"]}'
// Return User Object
// Example of return: 
// ["id: 3656469441", "name: MAXAVE"]
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetUser", args))

    var searchColumns []shim.Column
    searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: args[0]}})

    foundRow, err := stub.GetRow(UserTable, searchColumns)
    panicErr(err)
    return userRowToJSON(foundRow), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetProduct - retrieve a product Item by Product ID
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetProduct", "4977"]}'
// Return Product Object
// Example of return: 
// ["id: 4977","name: Vera Valor 1 once (LSP) 2012 - 6 langues PROOF 2012 LSP","pseudo: MAXAVE",
// "title: 900.000","weight: 6.45 g","trType: Acceptation don cadeau","dateIn: 1230768000",
// "dateOut: 1230768000","pieceStatus: 4","comment: message","country: Belgique",
// "bank: BELGIQUE-CBC BRUXELLES","depositBox: 3003","album: 3543","paidPrice: 118.00",
// "quotedPrice: 118.80", "lastTrDate: 1289990000"]
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetProduct", args))

    var searchColumns []shim.Column
    searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: args[0]}})

    foundRow, err := stub.GetRow(ProductTable, searchColumns)
    panicErr(err)
    return productRowToJSON(foundRow), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetProductsByUser - retrieve all Product IDs of product Items held by the given User 
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetProductByUser", "MAXAVE"]}'
// Return IDs Product separated by '/'
// Example of return: 
// /4977/4993/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetProductsByUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetProductsByUser", args))

    pseudo:=args[0]
    // We retrieve the string of users possessions (value) from the State for the given User (key)
    usersPossessions, err := stub.GetState(pseudo)
    panicErr(err)

    if usersPossessions!=nil {
        return []byte(usersPossessions), nil
    }
    return []byte("/"), nil
} 

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetProductsByAlbum - retrieve all Product IDs of product Items which belongs to the given Album 
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetProductsByAlbum", "3543"]}'
// Return IDs Product separated by '/'
// Example of return: 
// /4977/4993/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetProductsByAlbum(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetProductsByAlbum", args))

    album:=args[0]
    // We retrieve the string of IDs Product (value) from the State "album-" for the given Album ID (key)
    productsByAlbum, err := stub.GetState("album-"+album)
    panicErr(err)

    if productsByAlbum!=nil {
        return productsByAlbum, nil
    }
    return []byte("/"), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetProductsByDepositBox - retrieve all Product IDs of product Items which belongs to the given DepositBox
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetProductsByDepositBox", "4003"]}'
// Return IDs Product separated by '/'
// Example of return: 
// /4977/4993/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetProductsByDepositBox(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetProductsByDepositBox", args))

    depositBox:=args[0]
    // We retrieve the string of IDs Albums (value) from the State "depositBox-" for the Deposit Box (key)
    depositBoxAlbums, err := stub.GetState("depositBox-"+depositBox)
    panicErr(err)

    var resultingIds string
    resultingIds="/"

    if len(depositBoxAlbums) != 0 {
        result := strings.Split(string(depositBoxAlbums[:]), "/")
        // Loop for every ID Album
        for i := range result {
            var albumId string
            albumId=result[i]
            if len(result[i]) != 0 {
                // We retrieve the string of IDs Product (value) from the State "album-" for the given Album ID (key)
                productsByAlbum, err := stub.GetState("album-"+albumId)
                panicErr(err)
                // Concatenate all IDs Product from every iteration 
                resultingIds=resultingIds+""+string(productsByAlbum[:])
                fmt.Println(albumId, string(productsByAlbum[:]), "albumId, string(productsByAlbum[:])")
                resultingIds=strings.Replace(resultingIds, "//", "/", -1)
            }
        }
    }
    return []byte(resultingIds), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetAlbumsByDepositBox - retrieve all Album (IDs) which belongs to the given DepositBox
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetAlbumsByDepositBox", "4003"]}'
// Return IDs Album separated by '/'
// Example of return: 
// /4543/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetAlbumsByDepositBox(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetAlbumsByDepositBox", args))

    depositBox:=args[0]
    // We retrieve the string of IDs Albums (value) from the State "depositBox-" for the Deposit Box (key)
    depositBoxAlbums, err := stub.GetState("depositBox-"+depositBox)
    panicErr(err)

    if depositBoxAlbums!=nil {
        return depositBoxAlbums, nil
    }
    return []byte("/"), nil
}


////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetDepositBox - retrieve all Album and Product IDs which belongs to the given Deposit Box,
// Current method is not used currently due to the volumetry issue - method works fine with several
// lines but implemented shim method 'GetRows' doesn't work from certain number of lines. What why
// it was replaced by three other methods which use the key/value states.  
////////////////////////////////////////////////////////////////////////////////////////////////////////
/*func (t *AcbcChaincode) GetDepositBox(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetDepositBox", args))

    var searchColumns []shim.Column  
    searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: args[0]}})
    rowChannel, err := stub.GetRows(DepositBoxTable, searchColumns)
    if err!=nil {
        panic(err)
    }

    var returnString string
    returnString=""

        for {
            select {
            case row, ok := <-rowChannel:
                if !ok {
                    rowChannel = nil
                } else {
                    returnString=returnString+""+string(depositBoxRowToJSON(row))
                }
            }
            if rowChannel == nil {
                break
            }
        }
    return []byte(returnString), nil
}*/

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetTransactionById - retrieve Transaction information by Transaction ID
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetTransactionById", "100"]}'
// Return Transaction object
// Example of return: 
// ["id: 100","trDate: 1288915394","seller: MAXAVE","trReference: 100","idProduct: 4977",
// "nameProduct: Vera Valor 1 once (LSP) 2012 - 6 langues PROOF 2012 LSP","buyer: MAXAVE","title: 900.000",
// "weight: 6.45 g","trType: Acceptation don cadeau","dateIn: 1230768000","dateOut: 1230768000",
// "pieceStatus: 4","comment: message","country: Belgique","bank: BELGIQUE-CBC BRUXELLES","depositBox: 3003",
// "album: 3543","paidPrice: 118.00","quotedPrice: 118.80"]
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetTransactionById(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetTransactionById", args))

    var searchColumns []shim.Column
    searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: args[0]}})
    foundRow, err := stub.GetRow(TransactionTable, searchColumns)
    panicErr(err)
    return transactionRowToJSON(foundRow), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetTransactionsBetweenDates - retrieve all Transaction IDs between given dates (date of Transaction) 
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetTransactionsBetweenDates", "20100511", "20110313"]}'
// Return IDs Transaction separated by '/'
// Example of return: 
// /100///101/102///103/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetTransactionsBetweenDates(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetTransactionsBetweenDates", args))

    startDate := args[0]
    endDate := args[1]
    numStart,err := strconv.Atoi(startDate)
    numEnd,err := strconv.Atoi(endDate)
    panicErr(err)

    fmt.Println(numStart, "<= numStart")
    fmt.Println(numEnd, "<= numEnd")

    var transactionIds string
    transactionIds=""
    var queryKey string

    for _, countryCode := range countryCodes {

        for i := numStart; i <= numEnd; i++ {
            queryKey="transaction-"+strconv.Itoa(i)+"/"+countryCode
            fmt.Println(queryKey, "queryKey")
            value, err := stub.GetState(queryKey)
            panicErr(err)

            if value != nil {
                if transactionIds != "" {
                    transactionIds=transactionIds+"/"+string(value[:])
                } else {
                    transactionIds=string(value[:])
                }
            }
        }

    }

    fmt.Println(transactionIds, "<= resulting string")

    return []byte(transactionIds), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetTransactionsBetweenDatesByCountry - retrieve all Transaction IDs between given dates which correspond to Items 
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetTransactionsBetweenDatesByCountry", "20100511", "20101711","France"]}
// Return IDs Transaction separated by '/'
// Example of return: 
// /100/101/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetTransactionsBetweenDatesByCountry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    logger.Debug(fmt.Sprintf("%v called with %v", "GetTransactionsBetweenDatesByCountry", args))

    startDate := args[0]
    endDate := args[1]

    var countryCode string
    switch {

    case args[2]=="Belgique":

        countryCode = "BE"

    case args[2]=="Espagne":

        countryCode = "ES"

    case args[2]=="France":

        countryCode = "FR"

    case args[2]=="Suisse":

        countryCode = "CH"

    default:

        countryCode = "NO"
    }


    numStart,err := strconv.Atoi(startDate)
    numEnd,err := strconv.Atoi(endDate)
    panicErr(err)

    fmt.Println(numStart, "<= numStart")
    fmt.Println(numEnd, "<= numEnd")

    var transactionIds string
    transactionIds=""
    var queryKey string

    for i := numStart; i <= numEnd; i++ {
        queryKey="transaction-"+strconv.Itoa(i)+"/"+countryCode
        fmt.Println(queryKey, "queryKey")
        value, err := stub.GetState(queryKey)
        panicErr(err)

        if value != nil {
            if transactionIds != "" {
                transactionIds=transactionIds+"/"+string(value[:])
            } else {
                transactionIds=string(value[:])
            }
        }
    }
    fmt.Println(transactionIds, "<= resulting string")

    return []byte(transactionIds), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetTransactionsBetweenDatesForUser - retrieve all Transaction IDs between given dates which correspond to the ownership of given User
// Example of call:
// peer chaincode query -n $NAME -c '{"Function": "query", "Args": ["GetTransactionsBetweenDatesForUser", "20100511", "20110313","MAXAVE"]}'
// Return IDs Product separated by '/'
// Example of return: 
// /100/101/
////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *AcbcChaincode) GetTransactionsBetweenDatesForUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

    transacBetweenDates,err:=t.GetTransactionsBetweenDates(stub, args[0:2])
    panicErr(err)

    var allTransacIds string
    allTransacIds=string(transacBetweenDates[:])
    var resultingIds string
    resultingIds="/"

    if len(allTransacIds) != 0 {
        result := strings.Split(allTransacIds, "/")
        for i := range result {
            if len(result[i]) != 0 {
                var searchColumns []shim.Column
                searchColumns = append(searchColumns, shim.Column{Value: &shim.Column_String_{String_: result[i]}})
                foundRow, err := stub.GetRow(TransactionTable, searchColumns)
                panicErr(err)

                if len(foundRow.Columns[:]) != 0 {
                    // Column 5 is Buyer and Colum 2 is seller
                    if (foundRow.Columns[2].GetString_()==args[2])||(foundRow.Columns[6].GetString_()==args[2]) {
                        resultingIds=resultingIds+""+result[i]+"/"
                    }
                }
            }
        }

    }
    return []byte(resultingIds), nil
}



///////////////////////////////////////////////////
//               "technical" methods       
///////////////////////////////////////////////////

func userRowToJSON(row shim.Row) []byte {
 
	var object UserObject
    typeObject := reflect.TypeOf(object)
    var jsonResult []string

	for i := 0; i < typeObject.NumField(); i++ {
	    field := typeObject.Field(i)
        jsonType := field.Tag.Get("json")
	    result := jsonType + ":"+ " "+row.Columns[i].GetString_()

	jsonResult = append(jsonResult, result)
    }
    productJson, _ := json.Marshal(jsonResult)
    return productJson
}

func productRowToJSON(row shim.Row) []byte {
 
	var object ProductObject
    typeObject := reflect.TypeOf(object)
    var jsonResult []string

	for i := 0; i < typeObject.NumField(); i++ {
	    field := typeObject.Field(i)
        jsonType := field.Tag.Get("json")
	    result := jsonType + ":"+ " "+row.Columns[i].GetString_()

	jsonResult = append(jsonResult, result)
    }
    productJson, _ := json.Marshal(jsonResult)
    return productJson
}

func transactionRowToJSON(row shim.Row) []byte {
 
    var object TransactionObject
    typeObject := reflect.TypeOf(object)
    var jsonResult []string

    for i := 0; i < typeObject.NumField(); i++ {
        field := typeObject.Field(i)
        jsonType := field.Tag.Get("json")
        result := jsonType + ":"+ " "+row.Columns[i].GetString_()

    jsonResult = append(jsonResult, result)
    }
    transactionJson, _ := json.Marshal(jsonResult)
    return transactionJson
}

func depositBoxRowToJSON(row shim.Row) []byte {
 
    var object DepositBoxObject
    typeObject := reflect.TypeOf(object)
    var jsonResult []string

    for i := 0; i < typeObject.NumField(); i++ {
        field := typeObject.Field(i)
        jsonType := field.Tag.Get("json")
        result := jsonType + ":"+ " "+row.Columns[i].GetString_()

    jsonResult = append(jsonResult, result)
    }
    depositBoxJson, _ := json.Marshal(jsonResult)
    return depositBoxJson
}

///////////////////////////////////////////////////
//               ERRORS         
///////////////////////////////////////////////////

func checkArgsNumberOrPanic(args []string, argNumber int) {
    if len(args) != argNumber {
        panic(errors.New("Incorrect number of arguments."))
    }
}

func panicErr(err error) {
    if err != nil {
        panic(err)
    }
}
