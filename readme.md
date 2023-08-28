# receipt processor challenge

## structure
This repository follows the clean arquitecture structure, where the software components are 
located into the following directories:

- **internal/domain/receipt/**: Entities and enterprice business rules.
- **internal/interfaceadapters/storage/**: In memory storage.
- **internal/inputports/http/**: exposed endpoints.
- **internal/app/**: application business.

## Makefile

This repository has a Makefile with the following targets to help us with some repetitive tasks:

- **build**: build the receipt-processor-challenge inside build directory.
- **run**: build the receipt-processor-challenge and ejecute it.
- **test**: execute all tests inside the project.
- **test-report**: shows the coverage for all executed tests.
- **lint**: check the code using golangci-lint.
- **clean**: delete the binary generated file.
- **docker**: build the image for the project.

## HOW TO RUN

There are different ways to run the project, following the description for all of them:

1.- **Building it manually:**

go to root directory and ejecute the following commands:

```bash
go build -o  build/receipt-processor-challenge cmd/receipt-processor-challenge/main.go;
./build/receipt-processor-challenge;
```

2.- **Use the make command to build:**. 

```bash
make build;
./build/receipt-processor-challenge;
```

3.- **Use the make command to run:**
```bash
make run;
```

4.- **Using docker:**
```bash
make docker;
docker run -p 8080:8080 receipt-processor-challenge:1.0;
```

**Note**: the project opens the 8080 port in localhost.

