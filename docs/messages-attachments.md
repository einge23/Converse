# Message Flow Documentation

## Architecture Overview

**Client (React)** → **REST API (Go)** → **Socket.io Server (Node.js)** → **Message Consumer (Go)** → **Database (MySQL)**

## Message Types

### 1. Text-Only Messages

-   `content_type`: `"text"`
-   `content`: Message text
-   `attachments`: `[]` (empty array)

### 2. Image/File-Only Messages

-   `content_type`: `"image_url"` or `"file_url"`
-   `content`: Primary file URL or empty string
-   `attachments`: Array of file objects

### 3. Text + Attachments

-   `content_type`: `"text"`
-   `content`: Message text
-   `attachments`: Array of file objects

## Client Flow

### Text-Only Message

```typescript
// Direct socket emission
socket.emit("send_message", {
    room_id: "room-123",
    content: "Hello world!",
    content_type: "text",
    attachments: [],
});
```

### Image/File Upload + Message

```typescript
// Step 1: Upload files via REST API
const uploadFile = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("type", "messages");

    const response = await fetch("/api/v1/upload/message-attachment", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
    });

    return response.json();
};

// Step 2: Send message with attachment URLs
const attachments = await Promise.all(files.map(uploadFile));

socket.emit("send_message", {
    room_id: "room-123",
    content: "Check this out!", // Optional for image-only
    content_type: "text", // or 'image_url' for image-only
    attachments: attachments.map((upload) => ({
        id: generateUUID(),
        file_url: upload.url,
        file_key: upload.key,
        file_type: upload.content_type,
        file_size: upload.file_size,
        filename: upload.filename,
    })),
});
```

## Socket.io Server Processing

### Message Validation

```javascript
socket.on("send_message", (data) => {
    const message = {
        message_id: generateUUID(),
        room_id: data.room_id,
        thread_id: data.thread_id,
        sender_id: socket.userId,
        content_type: data.content_type || "text",
        content: data.content || "",
        attachments: data.attachments || [],
        created_at: new Date().toISOString(),
    };

    // Validate message structure
    if (!message.room_id && !message.thread_id) {
        socket.emit("error", { message: "Invalid message target" });
        return;
    }

    // Publish to RabbitMQ
    publishMessage(message);

    // Broadcast to room participants
    socket.to(data.room_id).emit("new_message", message);
});
```

## Backend Consumer Processing

### Message Storage Logic

```go
func storeMessage(db *sql.DB, message *Message) error {
    // Validation rules
    if message.ContentType == "text" && message.Content == "" {
        return errors.New("text messages require content")
    }

    if (message.ContentType == "image_url" || message.ContentType == "file_url") &&
       message.Content == "" && len(message.Attachments) == 0 {
        return errors.New("media messages require content or attachments")
    }

    // Transaction: Insert message + attachments
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert message
    _, err = tx.Exec(messageQuery, message.MessageID, message.RoomID, ...)
    if err != nil {
        return err
    }

    // Insert attachments
    for _, attachment := range message.Attachments {
        _, err = tx.Exec(attachmentQuery, attachment.ID, message.MessageID, ...)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

## Database Schema

### messages

```sql
CREATE TABLE messages (
    message_id CHAR(36) PRIMARY KEY,
    room_id CHAR(36),
    thread_id CHAR(36),
    sender_id CHAR(36),
    content_type ENUM('text','image_url','file_url','system_notification'),
    content TEXT NOT NULL,
    metadata JSON,
    created_at DATETIME(3) NOT NULL
);
```

### message_attachments

```sql
CREATE TABLE message_attachments (
    id CHAR(36) PRIMARY KEY,
    message_id CHAR(36) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_key VARCHAR(200) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INT,
    filename VARCHAR(255),
    FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE
);
```

## API Endpoints

### File Upload

```
POST /api/v1/upload/message-attachment
Authorization: Bearer <token>
Content-Type: multipart/form-data

Form Fields:
- file: File to upload
- room_id: Target room ID (optional)
- thread_id: Target thread ID (optional)

Response:
{
  "url": "https://bucket.s3.amazonaws.com/messages/rooms/123/user456/1642607890.jpg",
  "key": "messages/rooms/123/user456/1642607890.jpg",
  "filename": "image.jpg",
  "file_size": 2048576
}
```

## Error Handling

### Client-Side

-   File upload failures: Retry mechanism
-   Socket disconnection: Queue messages locally
-   Invalid file types: Show user-friendly errors

### Server-Side

-   Database failures: Requeue messages in RabbitMQ
-   S3 upload failures: Return HTTP 500
-   Invalid message format: Reject and log

## Performance Considerations

1. **File Upload**: Separate from real-time messaging
2. **Progress Tracking**: Use multipart upload for large files
3. **Caching**: Store frequently accessed attachments in CDN
4. **Optimization**: Compress images before upload
5. **Cleanup**: Remove orphaned S3 files periodically
