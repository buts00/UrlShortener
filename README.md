### URL Shortener Application

This program is a simple REST API server designed for shortening URLs. It interacts with an SQLite database to store mappings of shortened aliases to URLs.

#### Configuration

The program's configuration is stored in the `config/local.yaml` file. You can modify server settings and other configurations within this file.

#### Usage

1. **POST Request to Shorten URL:**

   To shorten a URL, send a POST request to `localhost:8080/url/` with the following JSON body:

   ```json
   {
       "url": "your_long_url_here",
       "alias": "optional_alias_here"
   }
   ```

   - `url` (required): The long URL that you want to shorten.
   - `alias` (optional): Custom alias to be used for the shortened URL.

2. **GET Request to Retrieve URL:**

   To retrieve a shortened URL, send a GET request to `localhost:8080/{alias}`, where `alias` is the shortened name. The server will check if the alias exists in the database and redirect you to the corresponding URL if found.

#### Response

- For successful requests, the response will be in the following format:

   ```json
   {
       "status": "success",
       "alias": "shortened_alias",
       "url": "original_long_url"
   }
   ```

- If there's an error, the response will contain:

   ```json
   {
       "status": "error",
       "error": "error_message_here"
   }
   ```

#### Environment

The path to the configuration file (`config/local.yaml`) is specified in the `.env` file.


