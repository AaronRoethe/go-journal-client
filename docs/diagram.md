<!-- ```mermaid
sequenceDiagram
    participant User
    participant AppEndpoint
    participant Messenger
    participant BlobStorage

    User->>AppEndpoint: Sends POST request with journal entry
    AppEndpoint->>Messenger: Sends journal entry to Messenger
    Messenger->>BlobStorage: Appends journal entry to BlobStorage

    AppEndpoint->>User: Returns success/failure response
``` -->

```mermaid
sequenceDiagram
    participant User
    participant AppEndpoint
    participant AuthServer
    participant ValidationServer
    participant Messenger
    participant BlobStorage

    User ->> AppEndpoint: Send POST request with journal entry and auth token
    AppEndpoint ->> AuthServer: Send auth token for validation
    AuthServer -->> AppEndpoint: Returns validation result
    AppEndpoint ->> ValidationServer: Send journal entry for validation
    ValidationServer -->> AppEndpoint: Returns validation result
    AppEndpoint ->> Messenger: Send journal entry to Messenger
    Messenger ->> BlobStorage: Append journal entry to BlobStorage
    AppEndpoint -->> User: Returns success/failure response
```
