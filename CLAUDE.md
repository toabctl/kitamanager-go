# KitaManager Go - Development Guidelines

## API Handler Documentation

All API handlers MUST be documented using swaggo annotations. This enables automatic OpenAPI/Swagger specification generation.

### Required Annotations

Every handler function must include the following annotations:

```go
// HandlerName godoc
// @Summary Short description of the endpoint
// @Description Detailed description of what the endpoint does
// @Tags tag-name
// @Accept json
// @Produce json
// @Security BearerAuth  // For protected endpoints
// @Param paramName path/query/body type required "Description"
// @Success statusCode {object/array} ResponseType
// @Failure statusCode {object} ErrorResponse
// @Router /api/v1/path [method]
func (h *Handler) HandlerName(c *gin.Context) {
    // implementation
}
```

### Example

```go
// Create godoc
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserCreate true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
    // implementation
}
```

### Request/Response Types

All request and response structs should include `example` tags for better documentation:

```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required" example:"John Doe"`
    Email    string `json:"email" binding:"required,email" example:"john@example.com"`
    Password string `json:"password" binding:"required,min=6" example:"secret123"`
}
```

### Generating Documentation

Run the following command to generate/update the OpenAPI specification:

```bash
swag init -g cmd/api/main.go
```

This will create/update files in the `docs/` directory.
