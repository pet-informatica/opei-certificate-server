# opei-certificates-server

golang web server that generates and downloads students OPEI certificates

## Getting Started

This project runs in go 1.10.x ([install on ubuntu](https://gist.github.com/ndaidong/4c0e9fbae8d3729510b1c04eb42d2a80)).

### Installing

At the project root folder, run the following:

```
go get github.com/jung-kurt/gofpdf
go get github.com/mattn/go-sqlite3
go build
```

### Running

Copy the `opei.db` file into to the project root folder.

```
./opei-certificate-server
```

### Testing

After starting the web server, it will respond to:

* POST requests
```
curl -H "Content-type: application/json" -X POST -d '{"cpf":"<insert CPF here>"}' http://localhost:8080/
```

* GET requests
  
    At your browser, access: `http://localhost:8080/?cpf=<insert CPF here>`

Both methods will serve the OPEI certificate with name following

