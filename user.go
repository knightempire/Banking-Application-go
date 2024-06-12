package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Prompt user to choose between login and register
	fmt.Println("Choose an option:")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		// Login
		fmt.Println("Enter username:")
		username, _ := reader.ReadString('\n')
		fmt.Println("Enter password:")
		password, _ := reader.ReadString('\n')

		// Send option and login details to server
		conn.Write([]byte("1\n")) // Option 1 for login
		conn.Write([]byte(username))
		conn.Write([]byte(password))

		// Read response from server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println(response)

		// If login successful, show options
		if strings.Contains(response, "Welcome") {
			showOptions(conn, reader)
		}
	case "2":
		// Register
		// Prompt user to enter registration details
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		fmt.Print("Enter name: ")
		name, _ := reader.ReadString('\n')
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')

		// Send option and registration details to server
		conn.Write([]byte("2\n")) // Option 2 for register
		conn.Write([]byte(username))
		conn.Write([]byte(name))
		conn.Write([]byte(password))

		// Read response from server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println(response)
	default:
		fmt.Println("Invalid option")
	}
}

func showOptions(conn net.Conn, reader *bufio.Reader) {
	fmt.Println("Choose an option:")
	fmt.Println("1. Deposit")
	fmt.Println("2. Withdraw")
	fmt.Println("3. Transfer")
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		fmt.Println("Deposit now")
		conn.Write([]byte("1\n")) // Send deposit option

		// Enter deposit amount
		fmt.Println("Enter deposit amount:")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)

		// Send deposit amount to server
		conn.Write([]byte(amountStr + "\n"))

		// Read response from server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println("Received response:", response) // Log the received response
		fmt.Println(response)

	case "2":
		fmt.Println("Withdraw option selected")

		// Send the withdraw option to the server
		conn.Write([]byte("2\n"))

		// Enter withdraw amount
		fmt.Println("Enter withdraw amount:")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)

		// Send withdraw amount to server
		conn.Write([]byte(amountStr + "\n"))

		// Read response from server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println(response)

	case "3":
		fmt.Println("Transfer option selected")

		// Send the transfer option to the server
		conn.Write([]byte("3\n"))

		// Enter transfer details: recipient username and amount
		fmt.Println("Enter recipient username:")
		recipientUsername, _ := reader.ReadString('\n')
		recipientUsername = strings.TrimSpace(recipientUsername)
		conn.Write([]byte(recipientUsername + "\n"))

		fmt.Println("Enter transfer amount:")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)
		conn.Write([]byte(amountStr + "\n"))

		// Read response from server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println(response)

	default:
		fmt.Println("Invalid option")
	}
}
