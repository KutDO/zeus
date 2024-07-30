// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/downloads": {
            "post": {
                "description": "Create a new download task",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "downloads"
                ],
                "summary": "Create download task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "URL of the file to download",
                        "name": "url",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Optional filename for the downloaded file",
                        "name": "filename",
                        "in": "formData"
                    },
                    {
                        "type": "boolean",
                        "description": "Decompress the file if true",
                        "name": "decompress",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Optional user ID for the download",
                        "name": "user_id",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.DownloadTask"
                        }
                    }
                }
            }
        },
        "/progress/{id}": {
            "get": {
                "description": "Get the progress of a download task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "downloads"
                ],
                "summary": "Get download progress",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Download Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "integer"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.DownloadTask": {
            "type": "object",
            "properties": {
                "decompress": {
                    "type": "boolean"
                },
                "filename": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Download Service API",
	Description:      "This is a service for downloading files with progress tracking.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
