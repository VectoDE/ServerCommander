# ServerCommander API Documentation

## Overview

The **ServerCommander API** provides a RESTful interface for managing remote servers via SSH, FTP, and SFTP. This API enables automated interactions with remote systems, including process management, file transfers, and system monitoring.

### Base URL

```bash
http://localhost:8080/api/v1
```

### Authentication

Authentication is done via API keys. Include your API key in the request headers:

```bash
Authorization: Bearer YOUR_API_KEY
```

## Endpoints

### 1. Server Management

#### 1.1 Get Server List

```bash
GET /servers
```

**Response**:

```bash
[
  {
    "id": "server1",
    "hostname": "example.com",
    "ip": "192.168.1.100",
    "status": "connected"
  }
]
```

#### 1.2 Connect to Server

```bash
POST /servers/connect
```

**Request Body**:

```bash
{
  "hostname": "example.com",
  "port": 22,
  "username": "admin",
  "password": "securepassword"
}
```

**Response**:

```bash
{
  "message": "Connection successful",
  "server_id": "server1"
}
```

#### 1.3 Disconnect from Server

```bash
POST /servers/disconnect
```

**Request Body**:

```bash
{
  "server_id": "server1"
}
```

**Response**:

```bash
{
  "message": "Disconnected successfully"
}
```

### 2. Process Management

#### 2.1 List Running Processes

```bash
GET /servers/{server_id}/processes
```

**Response**:

```bash
[
  {
    "pid": 1234,
    "name": "nginx",
    "cpu": "2.5%",
    "memory": "50MB"
  }
]
```

#### 2.2 Kill a Process

```bash
POST /servers/{server_id}/processes/kill
```

**Request Body**:

```bash
{
  "pid": 1234
}
```

**Response**:

```bash
{
  "message": "Process killed successfully"
}
```

### 3. File Transfer

#### 3.1 Upload File

```bash
POST /servers/{server_id}/upload
```

**Request Body**:

```bash
{
  "file_path": "/local/path/file.txt",
  "destination": "/remote/path/"
}
```

**Response**:

```bash
{
  "message": "File uploaded successfully"
}
```

#### 3.2 Download File

```bash
GET /servers/{server_id}/download?file_path=/remote/path/file.txt
```

**Response**:

The requested file is returned as a binary stream.

### 4. Server Monitoring

#### 4.1 Get Server Status

```bash
GET /servers/{server_id}/status
```

**Response**:

```bash
{
  "cpu_usage": "10%",
  "memory_usage": "2GB/8GB",
  "disk_usage": "50GB/100GB"
}
```

#### 4.2 Get Logs

```bash
GET /servers/{server_id}/logs
```

**Response**:

```bash
{
  "logs": [
    "[INFO] Server started",
    "[ERROR] Failed login attempt"
  ]
}
```

### Error Handling

Errors are returned in the following format:

```bash
{
  "error": "Invalid credentials",
  "code": 401
}
```

## Contribution

If you want to contribute, check out [CONTRIBUTING](CONTRIBUTING.md).

## License

This project is licensed under the MIT License. See [LICENSE](https://github.com/VectoDE/ServerCommander/blob/main/LICENSE) for details.
