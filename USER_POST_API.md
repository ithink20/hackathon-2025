# User Post API Documentation

## Overview
This API provides CRUD operations for user posts with support for comments and tags stored in JSONB metadata.

## Base URL
`http://localhost:8080/user/post`

## Operations

### 1. Create Post
**Endpoint:** `POST /user/post?op_type=create`

**Request Body:**
```json
{
  "title": "Optimizing React Performance in Large Applications",
  "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
  "authorName": "John Doe",
  "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
  "authorId": "1",
  "tags": ["React", "Performance", "Frontend"]
}
```

**Response:**
```json
{
  "post": {
    "post_id": "1703123456789",
    "title": "Optimizing React Performance in Large Applications",
    "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
    "authorName": "John Doe",
    "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
    "authorId": "1",
    "timestamp": 1703123456,
    "metaData": {
      "tags": ["React", "Performance", "Frontend"],
      "comments": []
    },
    "likes": 0,
    "created_at": "2023-12-21T10:30:56Z",
    "updated_at": "2023-12-21T10:30:56Z"
  },
  "message": "Post created successfully",
  "timestamp": "2023-12-21T10:30:56Z",
  "status": "success"
}
```

### 2. Read Post
**Endpoint:** `GET /user/post?op_type=read&post_id=<post_id>`

**Response:**
```json
{
  "post": {
    "post_id": "1703123456789",
    "title": "Optimizing React Performance in Large Applications",
    "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
    "authorName": "John Doe",
    "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
    "authorId": "1",
    "timestamp": 1703123456,
    "metaData": {
      "tags": ["React", "Performance", "Frontend"],
      "comments": [
        {
          "id": "2",
          "authorName": "Jane Doe",
          "authorId": "2",
          "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
          "content": "Great tips! I found code splitting especially useful in my projects.",
          "timestamp": 1715405800,
          "likes": 12
        }
      ]
    },
    "likes": 56,
    "created_at": "2023-12-21T10:30:56Z",
    "updated_at": "2023-12-21T10:30:56Z"
  },
  "message": "Post retrieved successfully",
  "timestamp": "2023-12-21T10:30:56Z",
  "status": "success"
}
```

### 3. Update Post
**Endpoint:** `PUT /user/post?op_type=update&post_id=<post_id>`

**Request Body:**
```json
{
  "title": "Optimizing React Performance in Large Applications",
  "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
  "authorName": "John Doe",
  "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
  "authorId": "1",
  "tags": ["React", "Performance", "Frontend"]
}
```

**Response:**
```json
{
  "post": {
    "post_id": "1703123456789",
    "title": "Optimizing React Performance in Large Applications",
    "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
    "authorName": "John Doe",
    "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
    "authorId": "1",
    "timestamp": 1703123456,
    "metaData": {
      "tags": ["React", "Performance", "Frontend"],
      "comments": []
    },
    "likes": 0,
    "created_at": "2023-12-21T10:30:56Z",
    "updated_at": "2023-12-21T10:30:56Z"
  },
  "message": "Post updated successfully",
  "timestamp": "2023-12-21T10:30:56Z",
  "status": "success"
}
```

### 4. Delete Post
**Endpoint:** `DELETE /user/post?op_type=delete&post_id=<post_id>`

**Response:**
```json
{
  "message": "Post deleted successfully",
  "timestamp": "2023-12-21T10:40:00Z",
  "status": "success"
}
```

### 5. List Posts
**Endpoint:** `GET /user/post?op_type=list`

**Query Parameters:**
- `limit` (optional): Number of posts to return (default: 10)
- `offset` (optional): Number of posts to skip (default: 0)
- `author_id` (optional): Filter posts by author ID

**Example:** `GET /user/post?op_type=list&limit=5&offset=0&author_id=1`

**Response:**
```json
{
  "posts": [
    {
      "post_id": "1703123456789",
      "title": "Optimizing React Performance in Large Applications",
      "content": "In this post, I share strategies and tools for improving the performance of large-scale React apps, including memoization, code splitting, and virtualization.",
      "authorName": "John Doe",
      "authorImg": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face",
      "authorId": "1",
      "timestamp": 1703123456,
      "metaData": {
        "tags": ["React", "Performance", "Frontend"],
        "comments": []
      },
      "likes": 56,
      "created_at": "2023-12-21T10:30:56Z",
      "updated_at": "2023-12-21T10:30:56Z"
    }
  ],
  "message": "Retrieved 1 posts",
  "timestamp": "2023-12-21T10:45:00Z",
  "status": "success"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid op_type. Must be one of: create, read, update, delete, list"
}
```

### 404 Not Found
```json
{
  "error": "Post not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create post"
}
```

## Database Schema

The `user_post` table has the following structure:

```sql
CREATE TABLE user_post (
    post_id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    author_image TEXT,
    author_id VARCHAR(255),
    timestamp BIGINT NOT NULL,
    metaData JSONB,
    likes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## Notes

1. **Post ID**: Automatically generated using timestamp for uniqueness
2. **Metadata**: Stored as JSONB in PostgreSQL, contains tags and comments
3. **Soft Delete**: Posts are soft deleted (deleted_at field) rather than hard deleted
4. **Timestamps**: All timestamps are in Unix format for consistency
5. **Comments**: Stored in the metadata JSONB field as an array
6. **Tags**: Stored in the metadata JSONB field as an array 