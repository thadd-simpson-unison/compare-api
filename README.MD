# Project Overview

Create a simple server that can:

1. Ingest a CSV list of emails
2. Call two external APIs passing the emails in different ways
3. Return the results in a comparable way
4. Be triggered using Swagger and no Auth
5. Pull external api credentials from Config

## Notes

- BLUE API is the service that takes in a bulk
- RED API is the service that takes in requests one at a time
