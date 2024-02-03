package rest

// @BasePath /api

// @title Kanthor Portal API
// @version 2022.1201.1701
// @description Portal API
// @termsOfService http://kanthorlabs.com/terms/

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @contact.name Kanthor Support
// @contact.url http://kanthorlabs.com/support
// @contact.email support@kanthorlabs.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @accept json
// @produce json

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
// @description [Bearer JWT_TOKEN] or [Basic base64(key:secret)]
// @securityDefinitions.apikey WorkspaceId
// @in header
// @name x-authorization-workspace
// @description The selected workspace id you are working with
