# My Toy DNS

🚀 Welcome to My Toy DNS!

A simple, fun, and educational DNS client built with Go. Send DNS queries effortlessly and see the magic happen!

🎯 Features

- Craft DNS A record queries with ease
- Communicate over UDP for quick responses
- View detailed DNS responses

🛠 Requirements

- Go 1.24.2 or newer

🔧 Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/my_toy_dns.git
   cd my_toy_dns
   ```

2. Build the project:

   ```bash
   go build
   ```

🚀 Usage

Run the program with a domain name:

```bash
./my_toy_dns example.com
```

It will query Google's DNS server (8.8.8.8) for the A record of `example.com` and display the response.

🧪 Testing

Run tests with:

```bash
go test
```

📁 Project Structure

- `main.go`: Core DNS query logic
- `main_test.go`: Tests for DNS query functions
- `go.mod`: Go module dependencies

📜 License

Licensed under MIT. See [LICENSE](LICENSE).

🤝 Contributing

Contributions are welcome! Open issues or pull requests for improvements or bug fixes.

This project is a simple DNS client written in Go. It constructs and sends DNS queries to a specified DNS server and prints the response.

## Features

- Constructs DNS queries for A records.
- Sends queries to a DNS server over UDP.
- Receives and displays the DNS response.

## Requirements

- Go 1.24.2 or later

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/my_toy_dns.git
   cd my_toy_dns
   ```

2. Build the project:

   ```bash
   go build
   ```

## Usage

Run the program with a domain name as an argument:

```bash
./my_toy_dns example.com
```

This will send a DNS query for the A record of `example.com` to Google's public DNS server (8.8.8.8) and print the response.

## Testing

To run the tests, use the following command:

```bash
go test
```

## Project Structure

- `main.go`: Contains the main logic for constructing and sending DNS queries.
- `main_test.go`: Contains unit tests for the DNS query construction functions.
- `go.mod`: Go module file.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
