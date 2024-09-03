/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	mspID         = "Org1MSP"
	cryptoPath    = "../organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/Admin@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "basic"
)

// var now = time.Now()
// var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// mainProcess()
	fmt.Println("Servers running on port 8082:")
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/mint", minter)
	http.HandleFunc("/balance", balancer)
	http.HandleFunc("/transfer", transferer)
	http.HandleFunc("/accountid", clientAccountIDer)
	http.HandleFunc("/initializer", initializer)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Success"))
	return
}

type registerRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponseBody struct {
	Message string `json:"message"`
}

func register(w http.ResponseWriter, req *http.Request) {
	//Parse Request's body
	body, _ := io.ReadAll(req.Body)
	var reqData registerRequestBody
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not parse the request body!"))
		return
	}
	cmd := "./../registerIdentityOrg1.sh"
	username := reqData.Username
	password := reqData.Password
	cmdOut := exec.Command("/bin/sh", cmd, username, password)
	_, err0 := cmdOut.Output()
	if err0 != nil {
		fmt.Println(err0.Error())
	}
	var now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s has been registered.", now, username)

	var payload registerResponseBody
	payload.Message = "ok"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

type mintRequestBody struct {
	Username string `json:"username"`
	Value    int    `json:"value"`
}

type mintResponseBody struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Value    string `json:"value"`
}

func minter(w http.ResponseWriter, req *http.Request) {
	// Parse request's body
	body, _ := io.ReadAll(req.Body)
	var reqData mintRequestBody
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not parse the request body!"))
		return
	}
	var username string = reqData.Username
	var value int = reqData.Value

	certPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/signcerts/cert.pem", username)
	keyPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/keystore/", username)

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	var now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s Called for mint.", now, username)
	mint(contract, value)
	now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s minted %d token.", now, username, value)

	var payload mintResponseBody
	payload.Message = "ok"
	payload.Username = username
	payload.Value = strconv.Itoa(value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

type balanceRequestBody struct {
	Username string `json:"username"`
}

type balanceResponse struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}

func balancer(w http.ResponseWriter, req *http.Request) {
	// Parse request's body
	body, _ := io.ReadAll(req.Body)
	var reqData balanceRequestBody
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not parse the request body!"))
		return
	}
	var username string = reqData.Username

	certPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/signcerts/cert.pem", username)
	keyPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/keystore/", username)

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	var now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s Called for balance.", now, username)
	balance := clientAccountBalance(contract)

	// var msg string = fmt.Sprintf("200 - Account %s's balance is %s.", username, balance)
	// w.Write([]byte(msg))

	var payload balanceResponse
	payload.Message = "ok"
	payload.Value = balance
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

type transferRequestBody struct {
	Username string `json:"username"`
	Receiver string `json:"receiver"`
	Value    int    `json:"value"`
}

type transferResponse struct {
	Message  string `json:"message"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Value    string `json:"value"`
}

func transferer(w http.ResponseWriter, req *http.Request) {
	// Parse request's body
	body, _ := io.ReadAll(req.Body)
	var reqData transferRequestBody
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not parse the request body!"))
		return
	}
	var username string = reqData.Username
	var receiver string = reqData.Receiver
	var value int = reqData.Value

	certPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/signcerts/cert.pem", username)
	keyPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/keystore/", username)

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	var now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s Called for transfer.", now, username)
	result := transfer(contract, receiver, value)
	if result {
		now = time.Now().Format("20060102150405")
		log.Printf("%s - User %s transfered %d token to user %s.", now, username, value, receiver)
	}

	var payload transferResponse
	payload.Message = "ok"
	payload.Sender = username
	payload.Receiver = receiver
	payload.Value = strconv.Itoa(value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

type accountIDRequestBody struct {
	Username string `json:"username"`
}

type accountIDResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func clientAccountIDer(w http.ResponseWriter, req *http.Request) {
	// Parse request's body
	body, _ := io.ReadAll(req.Body)
	var reqData accountIDRequestBody
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not parse the request body!"))
		return
	}
	var username string = reqData.Username

	fmt.Println(username)
	certPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/signcerts/cert.pem", username)
	keyPath = fmt.Sprintf(cryptoPath+"/users/%s@org1.example.com/msp/keystore/", username)

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	var now = time.Now().Format("20060102150405")
	log.Printf("%s - User %s Called for account ID.", now, username)
	clientID := clientAccountID(contract)

	var payload accountIDResponse
	payload.Message = "ok"
	payload.ID = clientID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

// Finish
func initializer(w http.ResponseWriter, req *http.Request) {

	certPath = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts/cert.pem"
	keyPath = cryptoPath + "/users/Admin@org1.example.com/msp/keystore/"
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	var now = time.Now().Format("20060102150405")
	log.Printf("%s - Service initializer called.", now)
	initialize(contract)

	var msg string = fmt.Sprintf("200 - Service Initilized. Time: %s", now)
	w.Write([]byte(msg))
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func initialize(contract *client.Contract) {
	fmt.Printf("Submit Transaction: Initialize, function creates the initial set of varriables on the ledger \n")

	_, err := contract.SubmitTransaction("Initialize", "MeetCoin", "MTCN", "0")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func mint(contract *client.Contract, value int) {
	_, err := contract.SubmitTransaction("Mint", strconv.Itoa(value))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
}

func transfer(contract *client.Contract, receiver string, value int) bool {
	_, err := contract.SubmitTransaction("Transfer", receiver, strconv.Itoa(value))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	return true
}

func clientAccountBalance(contract *client.Contract) string {
	fmt.Printf("Evaluate Transaction: Account balance \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountBalance")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	return result
}

func clientAccountID(contract *client.Contract) string {
	fmt.Printf("Evaluate Transaction: Account ID \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountID")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	//result := formatJSON(evaluateResult)

	// fmt.Printf("*** Result:%s\n", evaluateResult)
	return string(evaluateResult)
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
