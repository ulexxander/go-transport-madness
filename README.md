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

![image](https://user-images.githubusercontent.com/55105865/127398015-2d1581f5-5875-4324-b7d9-e839af2b21dc.png)
