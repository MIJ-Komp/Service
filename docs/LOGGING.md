# Logging Documentation

## Overview

Sistem logging telah diperbaiki dengan fitur-fitur berikut:

- **Structured JSON Logging**: Semua log disimpan dalam format JSON untuk parsing yang mudah
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Request ID Tracking**: Setiap request memiliki ID unik untuk tracing
- **Enhanced Error Logging**: Error logging dengan context yang lengkap
- **File and Console Output**: Log ditulis ke file dan console secara bersamaan

## Log Levels

### DEBUG
- Database operations
- Request completion details
- Development debugging information

### INFO
- Incoming requests
- Successful operations
- General application flow

### WARN
- Validation errors
- 4xx HTTP errors
- Expected errors that don't break functionality

### ERROR
- 5xx HTTP errors
- Unexpected errors
- System failures

### FATAL
- Critical errors that cause application shutdown

## Configuration

Set log level menggunakan environment variable:

```bash
LOG_LEVEL=DEBUG  # Shows all logs
LOG_LEVEL=INFO   # Shows INFO, WARN, ERROR, FATAL
LOG_LEVEL=WARN   # Shows WARN, ERROR, FATAL
LOG_LEVEL=ERROR  # Shows ERROR, FATAL only
LOG_LEVEL=FATAL  # Shows FATAL only
```

## Usage Examples

### Basic Logging

```go
import "api.mijkomp.com/helpers/logger"

// Simple messages
logger.LogInfo("User logged in successfully")
logger.LogWarning("Invalid input detected")
logger.LogError("Database connection failed")
```

### Logging with Data

```go
// With additional context data
logger.LogInfoWithData("User created", map[string]interface{}{
    "user_id": 123,
    "email": "user@example.com",
    "role": "admin",
})

logger.LogErrorWithData("Payment failed", map[string]interface{}{
    "order_id": "ORD-123",
    "amount": 100.50,
    "error_code": "INSUFFICIENT_FUNDS",
})
```

### Database Operations

```go
logger.LogDBOperation("SELECT", "SELECT * FROM users WHERE id = ?", userID)
```

## Log Format

Semua log disimpan dalam format JSON:

```json
{
  "timestamp": "2024-01-15 10:30:45.123",
  "level": "INFO",
  "message": "User logged in successfully",
  "file": "auth_controller.go",
  "line": 45,
  "function": "Login",
  "data": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": 123,
    "ip": "192.168.1.1"
  }
}
```

## Request Tracking

Setiap HTTP request memiliki unique request ID yang dapat digunakan untuk tracing:

- Request ID ditambahkan ke response header `X-Request-ID`
- Request ID disertakan dalam semua log terkait request tersebut
- Memudahkan debugging dan monitoring

## Log Files

- **Location**: `logs/app.log`
- **Rotation**: Manual (belum ada auto-rotation)
- **Format**: JSON per line
- **Output**: File dan console bersamaan

## Best Practices

1. **Gunakan level yang tepat**:
   - DEBUG untuk development debugging
   - INFO untuk flow normal aplikasi
   - WARN untuk kondisi yang perlu perhatian
   - ERROR untuk error yang perlu investigasi

2. **Sertakan context data**:
   ```go
   // Good
   logger.LogInfoWithData("Order processed", map[string]interface{}{
       "order_id": orderID,
       "user_id": userID,
       "amount": amount,
   })
   
   // Avoid
   logger.LogInfo("Order processed")
   ```

3. **Jangan log sensitive data**:
   - Password
   - Credit card numbers
   - Personal identification numbers

4. **Gunakan structured data**:
   - Lebih baik menggunakan map[string]interface{} daripada string formatting
   - Memudahkan parsing dan analysis

## Monitoring dan Analysis

Log JSON dapat diintegrasikan dengan tools seperti:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Fluentd
- Grafana Loki
- CloudWatch Logs

## Troubleshooting

### Log tidak muncul
1. Periksa LOG_LEVEL environment variable
2. Pastikan direktori `logs/` dapat ditulis
3. Periksa permission file log

### Performance impact
- JSON marshaling memiliki overhead kecil
- Untuk production, pertimbangkan menggunakan level WARN atau ERROR
- Monitor disk space untuk log files