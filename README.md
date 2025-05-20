# Roster Management System

## How to Run

### 1. Set Up PostgreSQL
- Ensure a PostgreSQL instance is running
- Alternatively, use `docker-compose` in the `/deployments` directory to run both PostgreSQL and service-api

### 2. Configure Environment
- Adjust `configs/local.env` to match your environment settings

### 3. Sync Dependencies
- Run `go mod tidy` to update and clean up dependencies
- Run `go mod vendor` to store dependencies locally

### 4. Run Migrations
- Execute `go run migrations/migrate.go up` to apply database migrations

### 5. Start the Service API
- Run `go run cmd/service-api/main.go` to start the service

### 6. Start the Frontend
- See `/cmd/web-app/README.md` for frontend setup instructions
- Note: The frontend service is not included in docker-compose due to CORS issues during development
- While admin login is available, the MVP for this role is not fully implemented yet. You may use `postman` for admin operations instead

## Admin Operations

This section describes operations exclusive to administrators. All requests require the `admin_token` environment variable, which is pre-configured in the Postman collection. Sample requests are available in the Postman collection.

### 1. Authentication
- **Login**: `POST /api/auth/login`
  - Admin credentials are pre-configured via SQL migrations

### 2. Shift Management
- **Create Shift**: `POST /api/shift`
- **Update Shift**: `PUT /api/shift/{shiftID}`
- **Delete Shift**: `DELETE /api/shift/{shiftID}`
- **Get Assigned Shifts**: `GET /api/shift/assigned`
  - Use this endpoint to view shifts for reassignment purposes

### 3. Shift Request Management
- **Update Shift Request**: `PUT /api/shift-request/{requestID}`
  - When a request is approved, all other requests for the same shift are automatically set to `not_selected`
  - This endpoint can also be used for shift reassignment
- **Get Pending Requests**: `GET /api/shift-request/pending`

#### For a detailed end-to-end workflow, please refer to the sequence diagram in `docs/worker-shift-journey.puml`.