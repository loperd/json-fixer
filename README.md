# JSON-Fixer

A simple proxy-server which gets content by passed url & removes extra-commas from content.

### Usage

To run http-server use following commands:
```bash 
make build && make up
```

And then you can get any json with fixing from extra-commas, using this example:
```
http://localhost:9900/?target=https%3A%2F%2Fexample.com%2Fsome-incorrect-json-file.json
```