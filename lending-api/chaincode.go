// package chaincode

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// )

// // SmartContract provides functions for managing the P2P lending platform
// type SmartContract struct {
// 	contractapi.Contract
// }

// // Account model
// type Account struct {
// 	Name         string   `json:"name"`
// 	Risk         int      `json:"risk"`
// 	Fund         int      `json:"fund"`   //money in is account
// 	Debt         int      `json:"debt"`   //loan taken by him
// 	Credit       int      `json:"credit"` //loan given by him
// 	Loan         int      `json:"loan"`   //loan he is willing to take
// 	Transactions []string `json:"transactions"`
// }

// type Transaction struct {
// 	LenderId   string `json:"lenderId"`
// 	BorrowerId string `json:"borrowerId"`
// 	Loan       int    `json:"loan"`
// 	AmountLeft int    `json:"amountLeft"`
// 	Paid       bool   `json:"paid"`
// }

// // Init initializes the ledger
// func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
// 	return nil
// }

// // Invoke routes function calls to the appropriate function
// // func (s *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) (*contractapi.TransactionResponse, error) {
// // 	function, args := ctx.GetStub().GetFunctionAndParameters()
// // 	fmt.Println("Invoke is running " + function)

// // 	switch function {
// // 	case "initLedger":
// // 		return s.initLedger(ctx, args)
// // 	case "borrow":
// // 		return s.borrow(ctx, args)
// // 	case "updateRisk":
// // 		return s.updateRisk(ctx, args)
// // 	case "query":
// // 		return s.query(ctx, args)
// // 	default:
// // 		return nil, fmt.Errorf("Received unknown function invocation: %s", function)
// // 	}
// // }

// // initLedger adds a base set of accounts to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	accounts := []Account{
// 		{Name: "Mahesh", Risk: 1, Fund: 1000, Debt: 0, Credit: 1000, Loan: 0, Transactions: []string{"TRANSACTION0", "TRANSACTION1"}},
// 		{Name: "Ramesh", Risk: 2, Fund: 0, Debt: 1000, Credit: 0, Loan: 0, Transactions: []string{"TRANSACTION1"}},
// 		{Name: "Suresh", Risk: 1, Fund: 200, Debt: 0, Credit: 0, Loan: 1000, Transactions: []string{"TRANSACTION0"}},
// 	}

// 	transactions := []Transaction{
// 		{LenderId: "ACCOUNT0", BorrowerId: "ACCOUNT2", Loan: 2300, AmountLeft: 2300, Paid: true},
// 		{LenderId: "ACCOUNT0", BorrowerId: "ACCOUNT1", Loan: 1000, AmountLeft: 1000, Paid: false},
// 	}
// 	for i, account := range accounts {
// 		accountJSON, err := json.Marshal(account)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState("ACCOUNT"+strconv.Itoa(i), accountJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %s", err.Error())
// 		}
// 	}

// 	for i, transaction := range transactions {
// 		transactionJSON, err := json.Marshal(transaction)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState("TRANSACTION"+strconv.Itoa(i), transactionJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %s", err.Error())
// 		}
// 	}
// 	return nil
// }

// // create account
// func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, name string, risk int, fund int) error {
// 	account := Account{
// 		Name:         name,
// 		Risk:         risk,
// 		Fund:         fund,
// 		Debt:         0,
// 		Credit:       0,
// 		Loan:         0,
// 		Transactions: []string{},
// 	}

// 	accountJSON, err := json.Marshal(account)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, accountJSON)
// }

// // add funds
// func (s *SmartContract) AddFunds(ctx contractapi.TransactionContextInterface, accountId string, fund int) error {
// 	accountJSON, err := ctx.GetStub().GetState(accountId)

// 	if err != nil {
// 		return fmt.Errorf("failed to read account from the world database: %v", err)
// 	}
// 	if accountJSON == nil {
// 		return fmt.Errorf("user does not exist")
// 	}

// 	var account Account
// 	err = json.Unmarshal(accountJSON, &account)
// 	if err != nil {
// 		return err
// 	}

// 	account.Fund += fund
// 	accountJSON, err = json.Marshal(account)

// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(accountId, accountJSON)
// }

// // borrow handles borrowing of funds
// func (s *SmartContract) Borrow(ctx contractapi.TransactionContextInterface, accountID string, amount int) error {
// 	borrowerJSON, err := ctx.GetStub().GetState(accountID)
// 	if err != nil {
// 		return fmt.Errorf("failed to read borrower from world state: %v", err)
// 	}
// 	if borrowerJSON == nil {
// 		return fmt.Errorf("borrower %s does not exist", accountID)
// 	}

// 	var borrower Account
// 	err = json.Unmarshal(borrowerJSON, &borrower)
// 	if err != nil {
// 		return err
// 	}

// 	borrower.Loan += amount
// 	borrowerJSON, err = json.Marshal(borrower)
// 	if err != nil {
// 		return err
// 	}

// 	// lenderIterator, err := ctx.GetStub().GetStateByRange("ACCOUNT0", "ACCOUNT9999")
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer lenderIterator.Close()

// 	// remaining := amount

// 	// for lenderIterator.HasNext() {
// 	// 	queryResponse, err := lenderIterator.Next()
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}

// 	// 	var lender Account
// 	// 	err = json.Unmarshal(queryResponse.Value, &lender)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}

// 	// 	if lender.Type == "lender" && (lender.Risk == borrower.Risk || lender.Auto) {
// 	// 		if lender.Fund >= remaining {
// 	// 			lender.Fund -= remaining
// 	// 			borrower.Fund += remaining
// 	// 			lender.Loan += remaining
// 	// 			remaining = 0
// 	// 		} else {
// 	// 			borrower.Fund += lender.Fund
// 	// 			lender.Loan += lender.Fund
// 	// 			remaining -= lender.Fund
// 	// 			lender.Fund = 0
// 	// 		}

// 	// 		lenderJSON, err := json.Marshal(lender)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}

// 	// 		err = ctx.GetStub().PutState(queryResponse.Key, lenderJSON)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}

// 	// 		if remaining == 0 {
// 	// 			break
// 	// 		}
// 	// 	}
// 	// }

// 	// if remaining > 0 {
// 	// 	return fmt.Errorf("not enough funds available to fulfill the loan")
// 	// }
// 	return ctx.GetStub().PutState(accountID, borrowerJSON)
// }

// // to give loans
// func (s *SmartContract) FundGiving(ctx contractapi.TransactionContextInterface, lenderId string, borrowerId string, tid string) error {
// 	lenderJSON, err := ctx.GetStub().GetState(lenderId)

// 	if err != nil {
// 		return fmt.Errorf("unable to get lender details from the world database: %v", err)
// 	}
// 	if lenderJSON == nil {
// 		return fmt.Errorf("lender Id is invalid")
// 	}

// 	borrowerJSON, err := ctx.GetStub().GetState(borrowerId)
// 	if err != nil {
// 		return fmt.Errorf("unable to get borrower details from the world database: %v", err)
// 	}
// 	if borrowerJSON == nil {
// 		return fmt.Errorf("lender Id is invalid")
// 	}

// 	var lender Account
// 	err = json.Unmarshal(lenderJSON, &lender)
// 	if err != nil {
// 		return err
// 	}

// 	var borrower Account

// 	err = json.Unmarshal(borrowerJSON, &borrower)
// 	if err != nil {
// 		return err
// 	}

// 	if borrower.Loan == 0 || lender.Fund < borrower.Loan {
// 		return fmt.Errorf("you cannot give loan to this user")
// 	}
// 	fund := borrower.Loan
// 	lender.Fund -= borrower.Loan
// 	lender.Credit += borrower.Loan
// 	borrower.Debt += borrower.Loan
// 	borrower.Fund += fund
// 	lender.Transactions = append(lender.Transactions, tid)
// 	borrower.Transactions = append(borrower.Transactions, tid)

// 	borrower.Loan = 0

// 	lenderJSON, err = json.Marshal(lender)
// 	if err != nil {
// 		return err
// 	}
// 	err = ctx.GetStub().PutState(lenderId, lenderJSON)
// 	if err != nil {
// 		return err
// 	}

// 	borrowerJSON, err = json.Marshal(borrower)
// 	if err != nil {
// 		return err
// 	}
// 	err = ctx.GetStub().PutState(borrowerId, borrowerJSON)
// 	if err != nil {
// 		return err
// 	}

// 	transaction := Transaction{
// 		LenderId:   lenderId,
// 		BorrowerId: borrowerId,
// 		Loan:       fund,
// 		AmountLeft: fund,
// 		Paid:       false,
// 	}

// 	transactionJSON, err := json.Marshal(transaction)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(tid, transactionJSON)
// }

// // to repay loans
// func (s *SmartContract) LoanRepayment(ctx contractapi.TransactionContextInterface, lenderId string, borrowerId string, tid string, amount int) error {
// 	lenderJSON, err := ctx.GetStub().GetState(lenderId)

// 	if err != nil {
// 		return fmt.Errorf("unable to get lender details from the world database: %v", err)
// 	}
// 	if lenderJSON == nil {
// 		return fmt.Errorf("lender Id is invalid")
// 	}

// 	borrowerJSON, err := ctx.GetStub().GetState(borrowerId)
// 	if err != nil {
// 		return fmt.Errorf("unable to get borrower details from the world database: %v", err)
// 	}
// 	if borrowerJSON == nil {
// 		return fmt.Errorf("lender Id is invalid")
// 	}

// 	transactionJSON, err := ctx.GetStub().GetState(tid)
// 	if err != nil {
// 		return fmt.Errorf("unable to get transaction details from the world database: %v", err)
// 	}
// 	if transactionJSON == nil {
// 		return fmt.Errorf("transaction Id is invalid")
// 	}

// 	var lender Account
// 	err = json.Unmarshal(lenderJSON, &lender)
// 	if err != nil {
// 		return err
// 	}

// 	var borrower Account
// 	err = json.Unmarshal(borrowerJSON, &borrower)
// 	if err != nil {
// 		return err
// 	}

// 	var transaction Transaction
// 	err = json.Unmarshal(transactionJSON, &transaction)
// 	if err != nil {
// 		return err
// 	}

// 	if transaction.LenderId != lenderId || transaction.BorrowerId != borrowerId {
// 		return fmt.Errorf("you have not taken any loan from this user")
// 	} else if transaction.Paid {
// 		return fmt.Errorf("you have already repaid the loan you took from this user")
// 	}

// 	if amount > lender.Fund {
// 		return fmt.Errorf("you donot have much funds left in your wallet to return the amount %v you selected", amount)
// 	}

// 	payableAmount := amount
// 	if amount >= transaction.AmountLeft {
// 		payableAmount = transaction.AmountLeft
// 		transaction.Paid = true
// 	}

// 	lender.Credit -= payableAmount
// 	lender.Fund += payableAmount
// 	borrower.Debt -= payableAmount
// 	borrower.Fund -= payableAmount

// 	transaction.AmountLeft -= payableAmount

// 	lenderJSON, err = json.Marshal(lender)
// 	if err != nil {
// 		return err
// 	}
// 	err = ctx.GetStub().PutState(lenderId, lenderJSON)
// 	if err != nil {
// 		return err
// 	}

// 	borrowerJSON, err = json.Marshal(borrower)
// 	if err != nil {
// 		return err
// 	}
// 	err = ctx.GetStub().PutState(borrowerId, borrowerJSON)
// 	if err != nil {
// 		return err
// 	}

// 	transactionJSON, err = json.Marshal(transaction)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(tid, transactionJSON)
// }

// // updateRisk updates the risk and auto status of an account
// func (s *SmartContract) UpdateRisk(ctx contractapi.TransactionContextInterface, accountID string, newRisk int) error {
// 	accountJSON, err := ctx.GetStub().GetState(accountID)
// 	if err != nil {
// 		return fmt.Errorf("failed to read account from world state: %v", err)
// 	}
// 	if accountJSON == nil {
// 		return fmt.Errorf("account %s does not exist", accountID)
// 	}

// 	var account Account
// 	err = json.Unmarshal(accountJSON, &account)
// 	if err != nil {
// 		return err
// 	}

// 	account.Risk = newRisk

// 	accountJSON, err = json.Marshal(account)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(accountID, accountJSON)
// }

// // returns the state of the account with the given key
// func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, borrowerId string, lenderId string) (*Account, error) {
// 	borrowerJSON, err := ctx.GetStub().GetState(borrowerId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read account from world state: %v", err)
// 	}
// 	if borrowerJSON == nil {
// 		return nil, fmt.Errorf("account %s does not exist", borrowerId)
// 	}

// 	lenderJSON, err := ctx.GetStub().GetState(lenderId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read account from world state: %v", err)
// 	}
// 	if lenderJSON == nil {
// 		return nil, fmt.Errorf("user %s does not exist", borrowerId)
// 	}
// 	var borrower Account
// 	err = json.Unmarshal(borrowerJSON, &borrower)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var lender Account
// 	err = json.Unmarshal(lenderJSON, &lender)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if borrower.Loan == 0 || lender.Loan != 0 || lender.Fund < borrower.Loan {
// 		return nil, fmt.Errorf("user %s is not authorised to check user %s", lender.Name, borrower.Name)
// 	}
// 	return &borrower, nil
// }

// // returns the state of all accounts
// func (s *SmartContract) GetAllAccounts(ctx contractapi.TransactionContextInterface) ([]*Account, error) {
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, fmt.Errorf("stub error %s", err)
// 	}
// 	defer resultsIterator.Close()

// 	var accounts []*Account
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, fmt.Errorf("query error %s", err)
// 		}

// 		// Check if the key starts with "ACCOUNT_"
// 		if strings.HasPrefix(queryResponse.Key, "ACCOUNT") {
// 			var account Account
// 			err = json.Unmarshal(queryResponse.Value, &account)
// 			if err != nil {
// 				return nil, fmt.Errorf("marshal error %s", err)
// 			}
// 			accounts = append(accounts, &account)
// 		}
// 	}

// 	return accounts, nil
// }

// // returns the state of all transactions
// func (s *SmartContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]*Transaction, error) {
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var transactions []*Transaction
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		if strings.HasPrefix(queryResponse.Key, "TRANSACTION") {
// 			var transaction Transaction
// 			err = json.Unmarshal(queryResponse.Value, &transaction)
// 			if err != nil {
// 				return nil, err
// 			}
// 			transactions = append(transactions, &transaction)
// 		}
// 	}

// 	return transactions, nil
// }
