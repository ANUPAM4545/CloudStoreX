# CloudStoreX Storage REST API Control Plane

CloudStoreX exposes a unified REST API over `/api/v1/storage` for managing buckets and objects across any underlying cloud provider.

---

## 1. Authentication & Correlation

Every request must include:
- `Authorization: Bearer <jwt_token>` header.
- `X-Request-ID`: Automatically injected into every response header by middleware for distributed tracing.

---

## 2. Bucket Endpoints

| Method | Endpoint | Description | Response |
|---|---|---|---|
| `GET` | `/api/v1/storage/buckets` | List all buckets | `200 OK` (JSON array of `BucketDTO`) |
| `POST` | `/api/v1/storage/buckets` | Create a bucket | `201 Created` (`BucketDTO`) |
| `DELETE` | `/api/v1/storage/buckets/:bucket` | Delete a bucket | `204 No Content` |

---

## 3. Object Endpoints

| Method | Endpoint | Description | Response |
|---|---|---|---|
| `GET` | `/api/v1/storage/buckets/:bucket/objects` | List objects in bucket (optional `?prefix=`) | `200 OK` (JSON array of `ObjectDTO`) |
| `POST` | `/api/v1/storage/buckets/:bucket/objects` | Upload file (`multipart/form-data`) | `201 Created` (`ObjectMetadataDTO`) |
| `GET` | `/api/v1/storage/buckets/:bucket/objects/*key` | Download object file stream | `200 OK` (`application/octet-stream`) |
| `DELETE` | `/api/v1/storage/buckets/:bucket/objects/*key` | Delete object | `204 No Content` |

---

## 4. Standard Error Response Contract

All errors return structured JSON with domain error codes:
```json
{
  "error": {
    "code": "bucket_not_found",
    "message": "The specified bucket 'prod-media' does not exist."
  }
}
```
