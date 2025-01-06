# OXO GAME API

## Description

This project requires a `.env.test` file to configure the database connection for testing purposes. Please follow the instructions below to create this file.

## Creating the `.env.test` File

1. **Create a new file** in the root directory of the project and name it `.env.test`.

2. **Add the following environment variables** to the `.env.test` file:

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

4. **Save the File**: After adding the above variables, save the `.env.test` file.

## Running the Tests

Once you have created the `.env.test` file, you can run the tests for the project. Make sure your database server is running and accessible.

## Additional Notes

- Ensure that the values in the `.env.test` file match your local database configuration.
- If you have any questions or need further assistance, feel free to reach out.
