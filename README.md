# Terraform Provider Netorca

A Terraform provider for Netorca. Supports all required methods for integrating NetOrca into terraform.  
See more in [Netorca Docs](https://netautomate.gitlab.io/netorca_docs/) and [Contact us to book a demo call](https://netorca.io/contact-us/).  

## Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/)
- Terraform (to use this provider)

## Installation
Clone the repository and use the provided Makefile targets to build and install the provider:

```bash
git clone https://github.com/yourusername/terraform-provider-netorca.git
cd terraform-provider-netorca
make install
```

The `make install` command will:
- **Build:** Compile the provider binary.
- **Format & Lint:** Ensure the codebase adheres to standard coding practices.
- **Generate:** Execute any necessary code generation.

## Development

The repository includes a Makefile with several helpful targets to streamline development:

- **Default:** Runs the full suite of code formatting, linting, installation, and code generation.
  
  ```bash
  make
  ```

- **Build:** Compile the provider with verbose output.
  
  ```bash
  make build
  ```

- **Install:** Build and install the provider.
  
  ```bash
  make install
  ```

- **Lint:** Run static analysis using `golangci-lint`.
  
  ```bash
  make lint
  ```

- **Generate:** Execute code generation commands located in the `tools` directory.
  
  ```bash
  make generate
  ```

- **Format:** Format the source code using `gofmt`.
  
  ```bash
  make fmt
  ```

- **Test:** Run unit tests with code coverage analysis.
  
  ```bash
  make test
  ```

- **Test Acc:** Run acceptance tests (ensure proper environment variables are set, e.g., `TF_ACC=1`).
  
  ```bash
  make testacc
  ```

The default Makefile target is defined as:

```makefile
default: fmt lint install generate
```

This target ensures that the code is formatted, linted, built, installed, and that any necessary code generation is performed.

## Testing

To verify that your changes do not break the provider, run the tests:

- **Unit Tests:**

  ```bash
  make test
  ```

- **Acceptance Tests:**

  ```bash
  make testacc
  ```

## License

This project is licensed under the [MIT License](LICENSE).
