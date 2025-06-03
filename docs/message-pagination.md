# Message Pagination API Documentation

This document outlines how to use the message pagination features in the Converse API.

## Endpoints

### Get Messages by Room ID

Retrieves paginated messages for a specific room.

```
GET /api/v1/messages/rooms/:room_id
```

#### URL Parameters

-   `room_id`: The ID of the room to get messages from

#### Query Parameters

-   `page` (optional): The page number to retrieve (default: 1)
-   `page_size` (optional): Number of messages per page (default: 20, max: 100)

#### Response

```json
{
    "messages": [
        {
            "message_id": "123e4567-e89b-12d3-a456-426614174000",
            "room_id": "123e4567-e89b-12d3-a456-426614174001",
            "thread_id": null,
            "sender_id": "123e4567-e89b-12d3-a456-426614174002",
            "content_type": "text",
            "content": "Hello world!",
            "metadata": null,
            "created_at": "2023-06-01T12:00:00Z",
            "updated_at": null,
            "deleted_at": null
        }
        // ... more messages
    ],
    "current_page": 1,
    "page_size": 20,
    "has_more": true
}
```

### Get Messages by Thread ID

Retrieves paginated messages for a specific thread.

```
GET /api/v1/messages/threads/:thread_id
```

#### URL Parameters

-   `thread_id`: The ID of the thread to get messages from

#### Query Parameters

-   `page` (optional): The page number to retrieve (default: 1)
-   `page_size` (optional): Number of messages per page (default: 20, max: 100)

#### Response

```json
{
    "messages": [
        {
            "message_id": "123e4567-e89b-12d3-a456-426614174000",
            "room_id": null,
            "thread_id": "123e4567-e89b-12d3-a456-426614174003",
            "sender_id": "123e4567-e89b-12d3-a456-426614174002",
            "content_type": "text",
            "content": "Hello world!",
            "metadata": null,
            "created_at": "2023-06-01T12:00:00Z",
            "updated_at": null,
            "deleted_at": null
        }
        // ... more messages
    ],
    "current_page": 1,
    "page_size": 20,
    "has_more": true
}
```

## Pagination

The API uses page-based pagination for retrieving messages:

-   `current_page`: The current page number
-   `page_size`: Number of messages per page
-   `has_more`: Boolean indicating if there are more pages available

To get the next page, increment the `page` parameter. For example, if you're on page 1, request page 2 by using `?page=2`.

## Example Usage

### Retrieving the first page of messages for a room

```
GET /api/v1/messages/rooms/123e4567-e89b-12d3-a456-426614174001?page=1&page_size=20
```

### Retrieving the second page of messages for a thread

```
GET /api/v1/messages/threads/123e4567-e89b-12d3-a456-426614174003?page=2&page_size=20
```

### Retrieving messages with a smaller page size

```
GET /api/v1/messages/rooms/123e4567-e89b-12d3-a456-426614174001?page=1&page_size=10
```
