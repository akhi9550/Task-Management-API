# Task-Management-API

## Environment Variables

Before running the project, you need to set the following environment variables with your corresponding values:


- `DB_URL`: Db url
- `DB_NAME`: Database name
- `JWT_SECRET_KEY`: Secret key

Make sure to provide the appropriate values for these environment variables to configure the project correctly.

### Steps

1. **Clone the repository:**

    ```bash
    git clone https://github.com/akhi9550/Task-Management-API.git
    cd TaskManagementAPI
    ```

2. **Setup environment variables:**

    Create a `.env` file in the root directory and add the necessary environment variables.

3. **Install Dependencies:**

    ```bash
    go mode tidy
    ```

4. **Run the application:**

    ```bash
    go run cmd/main.go
    ```

To run unit tests, execute the following command:

```bash
go test ./...
```
## Postman Api documentation

[API Documentation](https://documenter.getpostman.com/view/29514478/2sAXxPBYoW)
