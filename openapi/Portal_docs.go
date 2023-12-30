// Package openapi Code generated by swaggo/swag. DO NOT EDIT
package openapi

import "github.com/swaggo/swag"

const docTemplatePortal = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://kanthorlabs.com/terms/",
        "contact": {
            "name": "Kanthor Support",
            "url": "http://kanthorlabs.com/support",
            "email": "support@kanthorlabs.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/account": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "account"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/AccountGetRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "account"
                ],
                "parameters": [
                    {
                        "description": "setup options",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/AccountSetupReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/AccountSetupRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/application/{app_id}/message": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "application"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "limit returning records",
                        "name": "_limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 1669914060000,
                        "description": "starting time to scan in milliseconds",
                        "name": "_start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 1985533260000,
                        "description": "ending time to scan in milliseconds",
                        "name": "_end",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ApplicationListMessageRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/application/{app_id}/message/{msg_id}": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "application"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ApplicationGetMessageRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/credentials": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "credentials"
                ],
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "list by ids",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search keyword",
                        "name": "_q",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "limit returning records",
                        "name": "_limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "requesting page",
                        "name": "_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsListRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "credentials"
                ],
                "parameters": [
                    {
                        "description": "credentials properties",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsCreateReq"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/credentials/{wsc_id}": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "credentials"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "credentials id",
                        "name": "wsc_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsGetRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "credentials"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "credentials id",
                        "name": "wsc_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "credentials properties",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsUpdateRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/credentials/{wsc_id}/expiration": {
            "put": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "credentials"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "credentials id",
                        "name": "wsc_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "credentials properties",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsExpireReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCredentialsExpireRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/endpoint/{ep_id}/message": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "endpoint"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "limit returning records",
                        "name": "_limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 1669914060000,
                        "description": "starting time to scan in milliseconds",
                        "name": "_start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 1985533260000,
                        "description": "ending time to scan in milliseconds",
                        "name": "_end",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/EndpointListMessageRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/endpoint/{ep_id}/message/{msg_id}": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    },
                    {
                        "WorkspaceId": []
                    }
                ],
                "tags": [
                    "endpoint"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/EndpointGetMessageRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/workspace": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceGetRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "parameters": [
                    {
                        "description": "credentials properties",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceCreateReq"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        },
        "/workspace/{ws_id}": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "workspace id",
                        "name": "ws_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceGetRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "workspace id",
                        "name": "ws_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "credentials properties",
                        "name": "props",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/WorkspaceUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WorkspaceUpdateRes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/gateway.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "AccountGetRes": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/authenticator.Account"
                },
                "workspaces": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Workspace"
                    }
                }
            }
        },
        "AccountSetupReq": {
            "type": "object",
            "properties": {
                "workspace_name": {
                    "type": "string"
                }
            }
        },
        "AccountSetupRes": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/authenticator.Account"
                },
                "workspace": {
                    "$ref": "#/definitions/Workspace"
                }
            }
        },
        "ApplicationGetMessageRes": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "body": {
                    "type": "string"
                },
                "headers": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "ApplicationListMessageRes": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Message"
                    }
                }
            }
        },
        "EndpointGetMessageRes": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "body": {
                    "type": "string"
                },
                "headers": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "string"
                },
                "request_count": {
                    "type": "integer"
                },
                "request_latest_ts": {
                    "type": "integer"
                },
                "response_count": {
                    "type": "integer"
                },
                "response_latest_ts": {
                    "type": "integer"
                },
                "success_id": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "EndpointListMessageRes": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/EndpointMessage"
                    }
                }
            }
        },
        "EndpointMessage": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "body": {
                    "type": "string"
                },
                "headers": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "string"
                },
                "request_count": {
                    "type": "integer"
                },
                "request_latest_ts": {
                    "type": "integer"
                },
                "response_count": {
                    "type": "integer"
                },
                "response_latest_ts": {
                    "type": "integer"
                },
                "success_id": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "Message": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "body": {
                    "type": "string"
                },
                "headers": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "Workspace": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "tier": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "WorkspaceCreateReq": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "tier": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "WorkspaceCredentials": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "expired_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                },
                "ws_id": {
                    "type": "string"
                }
            }
        },
        "WorkspaceCredentialsCreateReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "WorkspaceCredentialsExpireReq": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer",
                    "default": 1800000
                }
            }
        },
        "WorkspaceCredentialsExpireRes": {
            "type": "object",
            "properties": {
                "expired_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "WorkspaceCredentialsGetRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "expired_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                },
                "ws_id": {
                    "type": "string"
                }
            }
        },
        "WorkspaceCredentialsListRes": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/WorkspaceCredentials"
                    }
                }
            }
        },
        "WorkspaceCredentialsUpdateReq": {
            "type": "object",
            "properties": {
                "expired_at": {
                    "type": "integer",
                    "default": 1893456000000
                },
                "name": {
                    "type": "string",
                    "default": "swagger demo update"
                }
            }
        },
        "WorkspaceCredentialsUpdateRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "expired_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                },
                "ws_id": {
                    "type": "string"
                }
            }
        },
        "WorkspaceGetRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/authorizator.Permission"
                    }
                },
                "tier": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "WorkspaceUpdateReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "another name"
                }
            }
        },
        "WorkspaceUpdateRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "tier": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "authenticator.Account": {
            "type": "object",
            "properties": {
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "sub": {
                    "type": "string"
                }
            }
        },
        "authorizator.Permission": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "object": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "gateway.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Authorization": {
            "description": "[Bearer JWT_TOKEN] or [Basic base64(key:secret)]",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "WorkspaceId": {
            "description": "The selected workspace id you are working with",
            "type": "apiKey",
            "name": "x-authorization-workspace",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfoPortal holds exported Swagger Info so clients can modify it
var SwaggerInfoPortal = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Kanthor Portal API",
	Description:      "Portal API",
	InfoInstanceName: "Portal",
	SwaggerTemplate:  docTemplatePortal,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfoPortal.InstanceName(), SwaggerInfoPortal)
}
