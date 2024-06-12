# Banking Application

## Introduction

This document serves as a guide to understand, deploy, and test the banking application developed in Go. It includes detailed explanations of the server and client applications, database setup, and test cases to ensure the reliability and functionality of the system. Whether you're a developer looking to understand the codebase or a tester validating the application's behavior, this documentation provides comprehensive insights into the banking application's workings.

## Database Connection

### Instructions for Setting Up a MySQL Database

To set up the MySQL database and create the necessary tables (`users` and `account`) to store user credentials and account balances, follow these steps:

1. Install the MySQL driver for Go:
    ```sh
    go get -u github.com/go-sql-driver/mysql
    ```

2. Create the necessary tables with the following SQL queries:
    ```sql
    CREATE TABLE users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        name VARCHAR(255)
    );

    CREATE TABLE account (
        username VARCHAR(255) PRIMARY KEY,
        balance DECIMAL(10, 2) NOT NULL DEFAULT 0
    );
    ```

## Test Cases

Various test cases are listed to verify the functionality of the banking application, including registration, login, deposit, withdrawal, and transfer operations. These test cases cover scenarios such as empty fields, invalid inputs, existing usernames, insufficient balances, and successful transactions.

### Register

- **Register username exists**:
    - Ensure the application correctly handles attempts to register with an existing username.

- **Empty fields**:
    - Verify that the application prompts the user to fill in all required fields.

- **Password empty**:
    - Check the application’s response when the password field is left empty.

### Invalid Login

- **Login**:
    - Test the login functionality with correct and incorrect credentials.

- **Wrong input**:
    - Validate the system’s behavior when incorrect login details are provided.

### Deposit

- **Deposit**:
    - Confirm that deposits are correctly added to the user's account balance.

### Withdrawal

- **Withdraw more amount**:
    - Ensure the system prevents withdrawals that exceed the current account balance.

- **Withdraw**:
    - Validate the correct deduction of funds from the user's account.

### Transactions

- **Invalid username selected for transaction**:
    - Verify the application handles transactions involving non-existent usernames.

- **Transferring more money**:
    - Test the application’s response when attempting to transfer more money than available in the sender's account.

- **Transaction Successful**:
    - Ensure the correct processing of a successful transaction.

    - **Receiver account initial balance**:
        - Check the initial balance of the receiver’s account before the transaction.

    - **Transaction from sender to receiver**:
        - Confirm that the funds are correctly transferred from the sender’s to the receiver’s account.

    - **Receiver account after transaction**:
        - Validate the updated balance of the receiver’s account post-transaction.

## Summary

In summary, this document provides an overview of a simple banking application developed in Go. It includes server-side and client-side implementations, database setup, and a set of test cases. The application allows users to register, log in, and perform banking operations such as deposit, withdrawal, and fund transfer. The test cases ensure the application's functionality and reliability. Overall, it offers a concise guide to understanding and testing the banking system.

