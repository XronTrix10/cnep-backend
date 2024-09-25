# CNEP Backend

This is a Go backend for the CNEP project. It is built using the Fiber framework and uses PostgreSQL as the database.

## Folder Structure

```bash
cnep-backend/
├── cmd/
│   └── server/
│       └── main.go
├── source/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── users.go
│   │   ├── posts.go
│   │   ├── comments.go
│   │   ├── connections.go
│   │   ├── conversations.go
│   │   └── notifications.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logger.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── comment.go
│   │   ├── connection.go
│   │   ├── conversation.go
│   │   └── notification.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   ├── connection_repository.go
│   │   ├── conversation_repository.go
│   │   └── notification_repository.go
│   ├── routes/
│   │   └── routes.go
│   └── websocket/
│       └── hub.go
├── pkg/
│   └── utils/
│       ├── jwt.go
│       └── password.go
├── docs/
│   └── PostgreSQl.md
├── LICENSE
├── README.md
├── .env
├── go.mod
└── go.sum
```
