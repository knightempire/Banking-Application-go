package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db              *sql.DB
	activeUsers     = make(map[string]net.Conn) // Map to store active users and their connections
	activeUsersLock sync.Mutex                  // Mutex to synchronize access to activeUsers map
)

func main() {
	// Connect to MySQL database
	var err error
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/go")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if the database connection is successful
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging database:", err)
		os.Exit(1)
	}
	fmt.Println("Database connected successfully.")

	// Start server
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on", port)

	// Accept connections indefinitely
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected.")

		// Handle client connection in a new goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Minute))
	var username string // Define username variable outside the switch statement

	// Read option from client
	reader := bufio.NewReader(conn)
	option, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading option:", err)
		return
	}
	option = strings.TrimSpace(option)

	switch option {
	case "1": // Login
		username = handleLogin(conn, reader) // Store the username returned by handleLogin
	case "2": // Register
		handleRegistration(conn, reader)
	default:
		fmt.Fprintln(conn, "Invalid option")
		return
	}

	// After login, handle deposit, withdrawal, transfer options
	for {

		// Read option from client
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading option:", err)
			return
		}
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			handleDeposit(conn, reader, username)

		case "2":
			handleWithdraw(conn, reader, username) // Uncomment this line when implementing withdrawal

		case "3":
			handleTransfer(conn, reader, username) // Uncomment this line when implementing transfer

		default:
			fmt.Fprintln(conn, "Invalid option")
		}

		// After handling each option selection and sending "Option selection successful" message
		fmt.Fprintln(conn, "Option selection received")

	}
}

func handleLogin(conn net.Conn, reader *bufio.Reader) string {
	// Read username from client
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		conn.Write([]byte("Internal server error\n"))
		return ""
	}
	username = strings.TrimSpace(username)

	// Check if username is empty
	if username == "" {
		conn.Write([]byte("Username is empty\n"))
		return ""
	}

	// Read password from client
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		conn.Write([]byte("Internal server error\n"))
		return ""
	}
	password = strings.TrimSpace(password)

	// Check if password is empty
	if password == "" {
		conn.Write([]byte("Password is empty\n"))
		return ""
	}

	// Perform authentication (check username and password against the database)
	validUser, err := dbQueryUserAndPassword(username, password)
	if err != nil {
		fmt.Println("Error querying database:", err)
		conn.Write([]byte("Internal server error\n"))
		return ""
	}
	if !validUser {
		conn.Write([]byte("Invalid username or password\n"))
		return ""
	}

	// Get the user's current balance from the database
	currentBalance, err := getCurrentBalance(username)
	if err != nil {
		fmt.Println("Error getting current balance:", err)
		conn.Write([]byte("Error getting current balance\n"))
		return ""
	}

	// Send the current balance to the client
	conn.Write([]byte(fmt.Sprintf("Welcome %s. | Your current balance is: %.2f\n", username, currentBalance)))

	// Mark the user as active
	activeUsers[username] = conn

	return username
}

// Two-Phase Commit Protocol

type Transaction struct {
	ID        int
	Operation string
	Data      string
}

type Coordinator struct {
	Transactions map[int]Transaction
	Lock         sync.Mutex
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Transactions: make(map[int]Transaction),
	}
}

func (c *Coordinator) Prepare() bool {
	// Simulate preparation phase
	// In a real implementation, this function would coordinate with the participating nodes
	return true
}

func (c *Coordinator) Commit() bool {
	// Simulate commit phase
	// In a real implementation, this function would coordinate with the participating nodes
	return true
}

func (c *Coordinator) Rollback() bool {
	// Simulate rollback phase
	// In a real implementation, this function would coordinate with the participating nodes
	return true
}

func (c *Coordinator) AddTransaction(id int, operation, data string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Transactions[id] = Transaction{ID: id, Operation: operation, Data: data}
}

func (c *Coordinator) RemoveTransaction(id int) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	delete(c.Transactions, id)
}

func (c *Coordinator) GetTransaction(id int) (Transaction, bool) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	tx, ok := c.Transactions[id]
	return tx, ok
}

var coordinator = NewCoordinator()

func handleDeposit(conn net.Conn, reader *bufio.Reader, username string) {
	// Read deposit amount from client
	amountStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading deposit amount:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	amountStr = strings.TrimSpace(amountStr)

	// Parse the deposit amount
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error parsing deposit amount:", err)
		conn.Write([]byte("Invalid deposit amount\n"))
		return
	}

	// Perform the deposit operation
	err = depositAmountWithTwoPhaseCommit(username, amount)
	if err != nil {
		fmt.Println("Error depositing amount:", err)
		conn.Write([]byte("Error depositing amount\n"))
		return
	}

	// Get the current balance after the deposit
	currentBalance, err := getCurrentBalance(username)
	if err != nil {
		fmt.Println("Error getting current balance:", err)
		conn.Write([]byte("Error getting current balance\n"))
		return
	}

	// Notify the client about the successful deposit and include the current balance
	message := fmt.Sprintf("Deposit of %.2f successful. Your current balance is %.2f\n", amount, currentBalance)

	conn.Write([]byte(message))
}

func depositAmountWithTwoPhaseCommit(username string, amount float64) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the transaction
	if !coordinator.Prepare() {
		// Rollback if preparation fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction preparation failed")
	}

	// Perform the deposit operation
	_, err = tx.Exec("UPDATE account SET balance = balance + ? WHERE username = ?", amount, username)
	if err != nil {
		// Rollback if deposit fails
		_ = tx.Rollback()
		return err
	}

	// Add the transaction to the coordinator
	coordinator.AddTransaction(txID, "deposit", fmt.Sprintf("%s deposited %.2f", username, amount))

	// Commit the transaction
	if !coordinator.Commit() {
		// Rollback if commit fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction commit failed")
	}

	return nil
}

func handleWithdraw(conn net.Conn, reader *bufio.Reader, username string) {
	// Read withdraw amount from client
	amountStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading withdraw amount:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	amountStr = strings.TrimSpace(amountStr)

	// Parse the withdraw amount
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error parsing withdraw amount:", err)
		conn.Write([]byte("Invalid withdraw amount\n"))
		return
	}

	// Perform the withdraw operation
	err = withdrawAmountWithTwoPhaseCommit(username, amount)
	if err != nil {
		fmt.Println("Error withdrawing amount:", err)
		conn.Write([]byte("Error withdrawing amount\n"))
		return
	}

	// Get the current balance after the withdrawal
	currentBalance, err := getCurrentBalance(username)
	if err != nil {
		fmt.Println("Error getting current balance:", err)
		conn.Write([]byte("Error getting current balance\n"))
		return
	}

	// Notify the client about the successful withdrawal and include the current balance
	message := fmt.Sprintf("Withdrawal of %.2f successful. Your current balance is %.2f\n", amount, currentBalance)
	conn.Write([]byte(message))
}

func withdrawAmountWithTwoPhaseCommit(username string, amount float64) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the transaction
	if !coordinator.Prepare() {
		// Rollback if preparation fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction preparation failed")
	}

	// Perform the withdraw operation
	_, err = tx.Exec("UPDATE account SET balance = balance - ? WHERE username = ?", amount, username)
	if err != nil {
		// Rollback if withdrawal fails
		_ = tx.Rollback()
		return err
	}

	// Add the transaction to the coordinator
	coordinator.AddTransaction(tx, "withdraw", fmt.Sprintf("%s withdrew %.2f", username, amount))

	// Commit the transaction
	if !coordinator.Commit() {
		// Rollback if commit fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction commit failed")
	}

	return nil
}

func handleTransfer(conn net.Conn, reader *bufio.Reader, username string) {
	// Read recipient username from the client
	recipientUsername, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading recipient username:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	recipientUsername = strings.TrimSpace(recipientUsername)

	// Validate recipient username
	if recipientUsername == username {
		fmt.Println("Self-transfer not allowed.")
		conn.Write([]byte("Self-transfer not allowed.\n"))
		return
	}

	// Read transfer amount from the client
	amountStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading transfer amount:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	amountStr = strings.TrimSpace(amountStr)

	// Parse the transfer amount
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error parsing transfer amount:", err)
		conn.Write([]byte("Invalid transfer amount\n"))
		return
	}

	// Perform the transfer operation
	err = transferAmountWithTwoPhaseCommit(username, recipientUsername, amount)
	if err != nil {
		fmt.Println("Error transferring amount:", err)
		conn.Write([]byte("Error transferring amount\n"))
		return
	}

	// Get the current balance after the transfer
	senderCurrentBalance, err := getCurrentBalance(username)
	if err != nil {
		fmt.Println("Error getting sender's current balance:", err)
		conn.Write([]byte("Error getting sender's current balance\n"))
		return
	}

	// Notify the client about the successful transfer including the current balance
	message := fmt.Sprintf("Transfer of %.2f to %s successful. Your current balance is %.2f\n", amount, recipientUsername, senderCurrentBalance)
	conn.Write([]byte(message))
}

func transferAmountWithTwoPhaseCommit(sender, recipient string, amount float64) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the transaction
	if !coordinator.Prepare() {
		// Rollback if preparation fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction preparation failed")
	}

	// Deduct the transfer amount from the sender's balance
	_, err = tx.Exec("UPDATE account SET balance = balance - ? WHERE username = ?", amount, sender)
	if err != nil {
		// Rollback if deduction fails
		_ = tx.Rollback()
		return err
	}

	// Add the transaction to the coordinator
	coordinator.AddTransaction(tx, "transfer", fmt.Sprintf("%s transferred %.2f to %s", sender, amount, recipient))

	// Add the transfer amount to the recipient's balance
	_, err = tx.Exec("UPDATE account SET balance = balance + ? WHERE username = ?", amount, recipient)
	if err != nil {
		// Rollback if addition fails
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction
	if !coordinator.Commit() {
		// Rollback if commit fails
		_ = tx.Rollback()
		return fmt.Errorf("transaction commit failed")
	}

	return nil
}

func handleRegistration(conn net.Conn, reader *bufio.Reader) {
	// Read username, name, and password from client
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	username = strings.TrimSpace(username)

	// Read name from client
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	name = strings.TrimSpace(name)

	// Read password from client
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	password = strings.TrimSpace(password)

	// Check if any field is empty
	if username == "" || name == "" || password == "" {
		conn.Write([]byte("All fields are required\n"))
		return
	}

	// Check if username already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		fmt.Println("Error querying database:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}
	if count > 0 {
		conn.Write([]byte("Username is already taken\n"))
		return
	}

	// Insert new user into the database
	_, err = db.Exec("INSERT INTO users (username, name, password) VALUES (?, ?, ?)", username, name, password)
	if err != nil {
		fmt.Println("Error inserting user into database:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}

	// Insert new user into the account table with an initial balance of 0
	_, err = db.Exec("INSERT INTO account (username, balance) VALUES (?, ?)", username, 0)
	if err != nil {
		fmt.Println("Error inserting user into account table:", err)
		conn.Write([]byte("Internal server error\n"))
		return
	}

	// Registration successful
	conn.Write([]byte("Registration successful\n"))
}

func dbQueryUserAndPassword(username, password string) (bool, error) {
	// Query the database to check if the username and password match
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		return false, err
	}
	// If count is 1, it means there is a match
	return count == 1, nil
}

func getCurrentBalance(username string) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM account WHERE username = ?", username).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func depositAmount(username string, amount float64) error {
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if there's an error
			fmt.Println("Rolling back transaction due to error:", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Println("Error rolling back transaction:", rollbackErr)
			}
		}
	}()

	// Get the current balance
	currentBalance, err := getCurrentBalance(username)
	if err != nil {
		return err
	}

	// Calculate the new balance after deposit
	newBalance := currentBalance + amount

	// Update the balance in the database
	_, err = tx.Exec("UPDATE account SET balance = ? WHERE username = ?", newBalance, username)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func withdrawAmount(username string, amount float64) error {
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if there's an error
			fmt.Println("Rolling back transaction due to error:", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Println("Error rolling back transaction:", rollbackErr)
			}
		}
	}()

	// Get the current balance
	currentBalance, err := getCurrentBalance(username)
	if err != nil {
		return err
	}

	// Check if there are sufficient funds for withdrawal
	if currentBalance < amount {
		return fmt.Errorf("Insufficient balance")
	}

	// Calculate the new balance after withdrawal
	newBalance := currentBalance - amount

	// Update the balance in the database
	_, err = tx.Exec("UPDATE account SET balance = ? WHERE username = ?", newBalance, username)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func transferAmount(sender, recipient string, amount float64) error {
	// Check if the sender exists
	senderExists, err := userExists(sender)
	if err != nil {
		return fmt.Errorf("error checking sender: %v", err)
	}
	if !senderExists {
		return fmt.Errorf("sender '%s' not found", sender)
	}

	// Check if the recipient exists
	recipientExists, err := userExists(recipient)
	if err != nil {
		return fmt.Errorf("error checking recipient: %v", err)
	}
	if !recipientExists {
		return fmt.Errorf("recipient '%s' not found", recipient)
	}

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if there's an error
			fmt.Println("Rolling back transaction due to error:", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Println("Error rolling back transaction:", rollbackErr)
			}
		}
	}()

	// Get the current balance of the sender
	senderBalance, err := getCurrentBalance(sender)
	if err != nil {
		return err
	}

	// Check if the sender has sufficient funds for the transfer
	if senderBalance < amount {
		return fmt.Errorf("sender '%s' has insufficient balance for transfer", sender)
	}

	// Calculate the new balances
	senderNewBalance := senderBalance - amount
	recipientNewBalance, err := getCurrentBalance(recipient)
	if err != nil {
		return err
	}
	recipientNewBalance += amount

	// Update the balances in the database
	_, err = tx.Exec("UPDATE account SET balance = ? WHERE username = ?", senderNewBalance, sender)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE account SET balance = ? WHERE username = ?", recipientNewBalance, recipient)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func userExists(username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
