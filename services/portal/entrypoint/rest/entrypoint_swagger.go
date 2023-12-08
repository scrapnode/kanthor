package rest

// @BasePath /api

// @title Kanthor Portal API
// @version 1.0
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description [Bearer <JWT token>] or [Basic base64(key:secret)]
// @securityDefinitions.apikey WsId
// @in header
// @name kanthor-ws-id
// @description The selected workspace id you are working with
