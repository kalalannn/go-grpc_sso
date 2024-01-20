# SSO microservice using gRPC and JWT.

### Overview

This microservice provides centralized authentication and authorization functionality,
using gRPC for communication and JWT (JSON Web Tokens) for secure transmission of user information.
It manages the lifecycle of JWT access and refresh tokens, including generation and validation.

The service exposes the following endpoints:

*   **`Register`**: Registers a new user.
*   **`Auth`**: Authenticates a user.
*   **`ValidateAccessToken`**: Validates an access token.
*   **`ValidateRefreshToken`**: Validates a refresh token.
*   **`RefreshTokens`**: Refreshes access and refresh tokens.

## Technologies

*   **Golang:** The primary programming language.
*   **gRPC:** A framework for building APIs.
*   **JWT:** The standard for creating tokens.
*   **SQLite:** Database for users management.
*   **Protocol Buffers:** For defining the gRPC API schema.

## gRPC API

The API is defined in the [`sso.proto`](https://github.com/kalalannn/go-grpc_sso/blob/master/proto/sso.proto) file.

### Service: `SSOService`

The service defines the following methods:

*   **`rpc Register (RegisterRequest) returns (RegisterResponse)`**
    *   Registers a new user.
    *   **Request**:
        *   `email` (string): User email.
        *   `username` (string): Username.
        *   `password` (string): User password.
    *   **Response**:
        *   `user_id` (string): Unique user identifier.

*   **`rpc Auth (AuthRequest) returns (AuthResponse)`**
    *   Authenticates a user.
    *   **Request**:
        *   `oneof auth_field`:
            *   `username` (string): Username.
            *   `email` (string): Email.
        *   `password` (string): User password.
    *   **Response**:
        *   `access_token` (string): JWT access token.
        *   `refresh_token` (string): JWT refresh token.

*   **`rpc ValidateAccessToken (ValidateTokenRequest) returns (ValidateTokenResponse)`**
    *   Validates an access token.
    *   **Request**:
        *   `token` (string): JWT access token.
    *   **Response**:
        *   `valid` (bool): Token is valid.
        *   `token_type` (string): Token type ("access").
        *   `expires_at` (string): Expiration date.
        *    `user_id` (string): User identifier.
        *   `username` (string): Username.

*   **`rpc ValidateRefreshToken (ValidateTokenRequest) returns (ValidateTokenResponse)`**
    *   Validates a refresh token.
    *   **Request**:
        *   `token` (string): JWT refresh token.
    *   **Response**:
       *   `valid` (bool): Token is valid.
       *   `token_type` (string): Token type ("refresh").
       *   `expires_at` (string): Expiration date.
       *    `user_id` (string): User identifier.
       *   `username` (string): Username.


*   **`rpc RefreshTokens (RefreshTokensRequest) returns (RefreshTokensResponse)`**
    *   Refreshes access and refresh tokens.
    *   **Request**:
        *   `refresh_token` (string): JWT refresh token.
    *   **Response**:
        *   `access_token` (string): New JWT access token.
        *   `refresh_token` (string): New JWT refresh token.

## Quick Start

(Instructions for setup and running will be added later)

<!-- ## Next Steps -->
<!-- *   Add setup and running instructions.
*   Expand functionality.
*   Include usage examples. -->