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
        "/account/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "account"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.AccountGetRes"
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
        "/account/setup": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
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
                            "$ref": "#/definitions/rest.AccountSetupReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.AccountSetupRes"
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
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceGetRes"
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
        "/workspace/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceGetRes"
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
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
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
                            "$ref": "#/definitions/rest.WorkspaceUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceUpdateRes"
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
        "/workspace/me/credentials": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "workspace"
                ],
                "parameters": [
                    {
                        "maxLength": 32,
                        "minLength": 2,
                        "type": "string",
                        "description": "search keyword",
                        "name": "_q",
                        "in": "query"
                    },
                    {
                        "maximum": 30,
                        "minimum": 5,
                        "type": "integer",
                        "description": "limit returning records",
                        "name": "_limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "requesting page",
                        "name": "_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceCredentialsListRes"
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
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
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
                            "$ref": "#/definitions/rest.WorkspaceCredentialsCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceCredentialsCreateRes"
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
        "/workspace/me/credentials/{wsc_id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "workspace"
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
                            "$ref": "#/definitions/rest.WorkspaceCredentialsGetRes"
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
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "workspace"
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
                            "$ref": "#/definitions/rest.WorkspaceCredentialsUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceCredentialsUpdateRes"
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
        "/workspace/me/credentials/{wsc_id}/expiration": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    },
                    {
                        "WsId": []
                    }
                ],
                "tags": [
                    "workspace"
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
                            "$ref": "#/definitions/rest.WorkspaceCredentialsExpireReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.WorkspaceCredentialsExpireRes"
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
        "entities.Workspace": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "I didn't find a way to disable automatic fields modify yet\nso, I use a tag to disable this feature here\nbut, we should keep our entities stateless if we can",
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
        },
        "rest.AccountGetRes": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/authenticator.Account"
                },
                "workspaces": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Workspace"
                    }
                }
            }
        },
        "rest.AccountSetupReq": {
            "type": "object",
            "properties": {
                "workspace_name": {
                    "type": "string"
                }
            }
        },
        "rest.AccountSetupRes": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/authenticator.Account"
                },
                "workspace": {
                    "$ref": "#/definitions/entities.Workspace"
                }
            }
        },
        "rest.WorkspaceCredentials": {
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
        "rest.WorkspaceCredentialsCreateReq": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "expired_at": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "rest.WorkspaceCredentialsCreateRes": {
            "type": "object",
            "properties": {
                "id": {
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
        "rest.WorkspaceCredentialsExpireReq": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                }
            }
        },
        "rest.WorkspaceCredentialsExpireRes": {
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
        "rest.WorkspaceCredentialsGetRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "I didn't find a way to disable automatic fields modify yet\nso, I use a tag to disable this feature here\nbut, we should keep our entities stateless if we can",
                    "type": "integer"
                },
                "expired_at": {
                    "type": "integer"
                },
                "hash": {
                    "type": "string"
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
        "rest.WorkspaceCredentialsListRes": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/rest.WorkspaceCredentials"
                    }
                }
            }
        },
        "rest.WorkspaceCredentialsUpdateReq": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "rest.WorkspaceCredentialsUpdateRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "I didn't find a way to disable automatic fields modify yet\nso, I use a tag to disable this feature here\nbut, we should keep our entities stateless if we can",
                    "type": "integer"
                },
                "expired_at": {
                    "type": "integer"
                },
                "hash": {
                    "type": "string"
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
        "rest.WorkspaceGetRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "I didn't find a way to disable automatic fields modify yet\nso, I use a tag to disable this feature here\nbut, we should keep our entities stateless if we can",
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
        "rest.WorkspaceUpdateReq": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "rest.WorkspaceUpdateRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "I didn't find a way to disable automatic fields modify yet\nso, I use a tag to disable this feature here\nbut, we should keep our entities stateless if we can",
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
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "[Bearer \u003cJWT token\u003e] or [Basic base64(key:secret)]",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "WsId": {
            "description": "The selected workspace id you are working with",
            "type": "apiKey",
            "name": "kanthor-ws-id",
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
