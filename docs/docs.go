// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/liveness": {
            "get": {
                "description": "Returns application's status info if the service is alive",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "App Liveness",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/systemapp.Info"
                        }
                    }
                }
            }
        },
        "/readiness": {
            "get": {
                "description": "Check application's readiness",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "App Readiness",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/systemapp.Status"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errs.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errs.ErrCode": {
            "type": "object"
        },
        "errs.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/errs.ErrCode"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "systemapp.Info": {
            "type": "object",
            "properties": {
                "GOMAXPROCS": {
                    "type": "integer"
                },
                "build": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "namespace": {
                    "type": "string"
                },
                "node": {
                    "type": "string"
                },
                "podIP": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "systemapp.Status": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:4000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Backend System",
	Description:      "This is a backend system with various technologies.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
