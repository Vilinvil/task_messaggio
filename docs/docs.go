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
            "name": "Vladislav",
            "url": "https://t.me/Vilin0",
            "email": "ivn-15-07@mail.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/message": {
            "post": {
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "message"
                ],
                "summary": "добавить сообщение в систему",
                "parameters": [
                    {
                        "description": "текст сообщения",
                        "name": "orderChanges",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseSuccessful"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/myerrors.Error"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/message/statistic": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "message"
                ],
                "summary": "посмотреть статистику по сообщениям",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.MessageStatistic"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.MessageStatistic": {
            "type": "object",
            "properties": {
                "handled": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "myerrors.Error": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "responses.ResponseSuccessful": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8090",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger message API",
	Description:      "This api for message service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
