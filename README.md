# POP3 Client

POP3 Client written in Golang in accordance to [RFC1939](https://www.ietf.org/rfc/rfc1939.txt).

## Usage

### Initialize client
```go
// Create a connection to the server
c, err := pop3.DialTLS("REPLACE_THIS_SERVER_ADDRESS:993")
if err != nil {
    log.Fatal(err)
}
defer c.Quit()

// Authenticate with the server
if err = c.Authorization("REPLACE_THIS_USERNAME", "REPLACE_THIS_PASSWORD"); err != nil {
    log.Fatal(err)
}
```

### Commands

```go
// Check if there are any messages to retrieved.
count, _, err := pc.Stat()
if err != nil {
    log.Fatal(err)
}

message, err := pc.Retr(1)
if err != nil {
    log.Fatal(err)
}

log.Println(message.Text)

if err := pc.Dele(1); err != nil {
    log.Fatal(err)
}
```
## Testing

To run tests, run following command:
```bash
go test -race ./...
```

## Contributions & Issues
Contributions are welcome. Please clearly explain the purpose of the PR and follow the current style.

Issues can be resolved quickest if they are descriptive and include both a reduced test case, and a set of steps to reproduce.

## Licence
The `genert/pop3` library is copyrighted © by [Genert Org](http://genert.org) and licensed for use under the MIT License (MIT).

Please see [MIT License](LICENSE) for more information.