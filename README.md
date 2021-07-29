# Transport Madness o_0

## Structure

- `models/` - shared structs between services and transport layers (User, Message)
- `services/` - business logic and tests for it (in this case very simple - UsersService, MessagesService)
- `transport/` - package itself defines **Publisher** interface that must be implemented by event-driven transports (websocket, nats)
- `transport/*` - different transport layers, some of them define only Responder (that depends on services), others Publisher implementations

## Implemented Transports

- `graphql` - only Responder (no events)
- `http` - only Responder (no events)
- `nats` - Responder and Publisher
- `websocket` - only Publisher (no requests)

## Preview

### Websocket, Nats, GraphQL

![image](https://user-images.githubusercontent.com/55105865/127398015-2d1581f5-5875-4324-b7d9-e839af2b21dc.png)

### GRPC

```
~/go/src/github.com/ulexxander/transport-madness grpcurl -proto transport/grpc/services.proto -plaintext -d '{"Username": "alex grpc!"}' localhost:4008 services.Users/Create{
  "User": {
    "Username": "alex grpc!",
    "CreatedAt": "2021-07-29 23:08:48.830969528 +0200 CEST m=+83.909324685"
  }
}
~/go/src/github.com/ulexxander/transport-madness grpcurl -proto transport/grpc/services.proto -plaintext -d '{"SenderUsername": "alex grpc!", "Content": "yaaaaaa"}' localhost:4008 services.Messages/Create
{
  "Message": {
    "SenderUsername": "alex grpc!",
    "Content": "yaaaaaa",
    "CreatedAt": "2021-07-29 23:09:05.212179931 +0200 CEST m=+100.290535198"
  }
}
~/go/src/github.com/ulexxander/transport-madness grpcurl -proto transport/grpc/services.proto -plaintext -d '{"Page": 0, "PageSize": 1}' localhost:4008 services.Messages/Pagination{
  "Messages": [
    {
      "SenderUsername": "alex grpc!",
      "Content": "yaaaaaa",
      "CreatedAt": "2021-07-29 23:09:05.212179931 +0200 CEST m=+100.290535198"
    }
  ]
}
~/go/src/github.com/ulexxander/transport-madness
```
