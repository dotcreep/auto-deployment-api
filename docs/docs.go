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
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/deploy/remove": {
            "delete": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Remove all data of user by username and domain used",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deploy"
                ],
                "summary": "Undeploy user data",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/deploy_api.RequestInputAddDomain"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/deploy/start": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Deployment to Cloudflare, Portainer, and Jenkins (Support rollback action if failed)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deploy"
                ],
                "summary": "Deploy All Third Party Environment",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/deploy_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessDeploy"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/add": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Add domain to cloudflare tunnel and dns record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "Add domain to cloudflare",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInputAddDomain"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/check": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Get domain is available or unavailable, only support for cloudflare providers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "Get domain is available",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/is-not-exists": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "True if domain is not exist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "True if domain is not exist",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "302": {
                        "description": "Found",
                        "schema": {
                            "$ref": "#/definitions/utils.FoundFail"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/nameserver": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Get cloudflare nameserver of base domain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "Get cloudflare nameserver",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/register-status": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Check status register domain with status ` + "`" + `pending` + "`" + ` and ` + "`" + `active` + "`" + `",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "Get status register domain",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/domain/status": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Are domain is accessible or still cannot access",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Domain"
                ],
                "summary": "Get status domain",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cloudflare_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/mobile/is-not-exists": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "True if item is not exist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Mobile"
                ],
                "summary": "True if item is not exist",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/jenkins_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/mobile/status": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Check status of mobile builder with return 'success', 'no build', 'failed', 'unknown'",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Mobile"
                ],
                "summary": "Check status of mobile builder",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/jenkins_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/system/is-not-exists": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "True if stack is not exist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "True if stack is not exist",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/portainer_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/system/stack": {
            "get": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Get all stack from portainer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Get all system stack",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/system/status": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "Get status of stack",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Get status of stack",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/portainer_api.RequestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        },
        "/api/v1/system/update": {
            "post": {
                "security": [
                    {
                        "X-Token": []
                    }
                ],
                "description": "You can update stack only using name of stack",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Update stack by name",
                "parameters": [
                    {
                        "description": "Body Input",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/portainer_api.RequestInputUpdateStack"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/utils.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.InternalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "cloudflare_api.RequestInput": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "sub.example.com"
                }
            }
        },
        "cloudflare_api.RequestInputAddDomain": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "example.com"
                },
                "username": {
                    "type": "string",
                    "example": "exampleusername"
                }
            }
        },
        "deploy_api.RequestInput": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "example.com"
                },
                "email": {
                    "type": "string",
                    "example": "sample@example.com"
                },
                "merchant_name": {
                    "type": "string",
                    "example": "Example Name"
                },
                "paket_merchant": {
                    "type": "string",
                    "example": "starter"
                },
                "username": {
                    "type": "string",
                    "example": "exampleusername"
                }
            }
        },
        "deploy_api.RequestInputAddDomain": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string",
                    "example": "example.com"
                },
                "username": {
                    "type": "string",
                    "example": "exampleusername"
                }
            }
        },
        "jenkins_api.RequestInput": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string",
                    "example": "exampleusername"
                }
            }
        },
        "portainer_api.RequestInput": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string",
                    "example": "exampleusername"
                }
            }
        },
        "portainer_api.RequestInputUpdateStack": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "mystack"
                }
            }
        },
        "utils.BadRequest": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "string"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "result": {
                    "type": "string",
                    "example": "null"
                },
                "status": {
                    "type": "integer",
                    "example": 400
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "utils.FoundFail": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "result": {
                    "type": "string",
                    "example": "null"
                },
                "status": {
                    "type": "integer",
                    "example": 302
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "utils.InternalServerError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "result": {
                    "type": "string",
                    "example": "null"
                },
                "status": {
                    "type": "integer",
                    "example": 500
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "utils.Success": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "null"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "result": {
                    "type": "string",
                    "example": "message"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "utils.SuccessDeploy": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "null"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "result": {
                    "type": "object",
                    "properties": {
                        "cloudflare": {
                            "type": "string",
                            "example": "success add domain sub.example.com"
                        },
                        "jenkins": {
                            "type": "string",
                            "example": "success deploy jenkins with status build in proccess"
                        },
                        "portainer": {
                            "type": "string",
                            "example": "success deploy portainer"
                        }
                    }
                },
                "status": {
                    "type": "integer",
                    "example": 200
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        }
    },
    "securityDefinitions": {
        "X-Token": {
            "description": "Input your token authorized",
            "type": "apiKey",
            "name": "X-Token",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Automate Deployment API",
	Description:      "Documentation for Automate Deployment Restful API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
