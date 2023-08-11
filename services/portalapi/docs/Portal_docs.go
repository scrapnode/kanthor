// Code generated by swaggo/swag. DO NOT EDIT.

package docs

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
                            "$ref": "#/definitions/portalapi.WorkspaceGetRes"
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
                        "minLength": 29,
                        "type": "string",
                        "description": "current query cursor",
                        "name": "_cursor",
                        "in": "query"
                    },
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
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "only return records with selected ids",
                        "name": "_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsListRes"
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
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsCreateRes"
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
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsGetRes"
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
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsUpdateRes"
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
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsExpireReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/portalapi.WorkspaceCredentialsExpireRes"
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
        "entities.WorkspaceCredentials": {
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
                "workspace_id": {
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
        },
        "portalapi.WorkspaceCredentialsCreateReq": {
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
        "portalapi.WorkspaceCredentialsCreateRes": {
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
        "portalapi.WorkspaceCredentialsExpireReq": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                }
            }
        },
        "portalapi.WorkspaceCredentialsExpireRes": {
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
        "portalapi.WorkspaceCredentialsGetRes": {
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
                "workspace_id": {
                    "type": "string"
                }
            }
        },
        "portalapi.WorkspaceCredentialsListRes": {
            "type": "object",
            "properties": {
                "cursor": {
                    "type": "string"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.WorkspaceCredentials"
                    }
                }
            }
        },
        "portalapi.WorkspaceCredentialsUpdateReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "portalapi.WorkspaceCredentialsUpdateRes": {
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
                "workspace_id": {
                    "type": "string"
                }
            }
        },
        "portalapi.WorkspaceGetRes": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "tier_name": {
                    "type": "string"
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
