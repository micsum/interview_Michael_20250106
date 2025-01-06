# OXO GAME API

## Description

This project requires a `.env` file to configure the database connection for testing purposes. Please follow the instructions below to create this file.

## Creating the `.env` File

1. **Create a new file** in the root directory of the project and name it `.env`.

2. **Add the following environment variables** to the `.env` file:

   ```env
   TEST_DB_HOST=db
   TEST_DB_USER=<YOUR DB USER>
   TEST_DB_PASS=<YOUR PASSWORD>
   TEST_DB_NAME=<DATABASE NAME>>
   TEST_DB_PORT=<PORT NO.>
   ```

3. **Explanation of the Variables**:
   - `TEST_DB_HOST`: The hostname of the database server. For local testing, this is usually set to `db` due to docker container.
   - `TEST_DB_USER`: The username used to connect to the database. In this case, it is set to `postgres`.
   - `TEST_DB_PASS`: The password for the database user. Ensure this matches the password set for your PostgreSQL user.
   - `TEST_DB_NAME`: The name of the database to connect to for testing. This should be set to `test_oxo_db`.
   - `TEST_DB_PORT`: The port on which the database server is running. The default PostgreSQL port is `5432`.

4. **Save the File**: After adding the above variables, save the `.env` file.

## Running the Tests

Once you have created the `.env` file, you can run the tests for the project. Make sure your database server is running and accessible.

## Accessing the API Documentation

Once the server is running, you can access the API documentation via Swagger at the following URL:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Kick start the server with Docker, and you can read the documentation using the link above.

## Additional Notes

- Ensure that the values in the `.env` file match your local database configuration.
- If you have any questions or need further assistance, feel free to reach out.
