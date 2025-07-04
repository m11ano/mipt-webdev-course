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
        "/orders": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Получить список заказов",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetOrdersOut"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Создать заказ",
                "parameters": [
                    {
                        "description": "JSON",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateOrderIn"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.CreateOrderOut"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/orders/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Получить заказ по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetOrderOut"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Обновить заказ",
                "parameters": [
                    {
                        "description": "JSON",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateOrderIn"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/orders/{id}/status": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Поменять статус заказу",
                "parameters": [
                    {
                        "description": "JSON",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.SetOrderStatusIn"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/orders/{id}/{secret_key}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Получить заказ по ID + secret_key",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetOrderOut"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorJSON"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.CreateOrderIn": {
            "type": "object",
            "required": [
                "details",
                "products"
            ],
            "properties": {
                "details": {
                    "$ref": "#/definitions/controller.CreateOrderInDetails"
                },
                "products": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/controller.CreateOrderInProduct"
                    }
                }
            }
        },
        "controller.CreateOrderInDetails": {
            "type": "object",
            "required": [
                "client_email",
                "client_name",
                "client_phone",
                "client_surname",
                "delivery_address"
            ],
            "properties": {
                "client_email": {
                    "type": "string"
                },
                "client_name": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "client_phone": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 1
                },
                "client_surname": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "delivery_address": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                }
            }
        },
        "controller.CreateOrderInProduct": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "minimum": 0
                },
                "quantity": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "controller.CreateOrderOut": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "secret_key": {
                    "type": "string"
                }
            }
        },
        "controller.GetOrderOut": {
            "type": "object",
            "properties": {
                "details": {
                    "$ref": "#/definitions/controller.GetOrderOutDetails"
                },
                "id": {
                    "type": "integer"
                },
                "order_sum": {
                    "type": "number"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controller.GetOrderOutProduct"
                    }
                },
                "secret_key": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "controller.GetOrderOutDetails": {
            "type": "object",
            "properties": {
                "client_email": {
                    "type": "string"
                },
                "client_name": {
                    "type": "string"
                },
                "client_phone": {
                    "type": "string"
                },
                "client_surname": {
                    "type": "string"
                },
                "delivery_address": {
                    "type": "string"
                }
            }
        },
        "controller.GetOrderOutProduct": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "controller.GetOrdersOut": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controller.GetOrdersOutItem"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "controller.GetOrdersOutItem": {
            "type": "object",
            "properties": {
                "details": {
                    "$ref": "#/definitions/controller.GetOrderOutDetails"
                },
                "id": {
                    "type": "integer"
                },
                "order_sum": {
                    "type": "number"
                },
                "secret_key": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "controller.SetOrderStatusIn": {
            "type": "object",
            "required": [
                "status"
            ],
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "controller.UpdateOrderIn": {
            "type": "object",
            "required": [
                "details",
                "products"
            ],
            "properties": {
                "details": {
                    "$ref": "#/definitions/controller.UpdateOrderInDetails"
                },
                "products": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/controller.UpdateOrderInProduct"
                    }
                }
            }
        },
        "controller.UpdateOrderInDetails": {
            "type": "object",
            "required": [
                "client_email",
                "client_name",
                "client_phone",
                "client_surname",
                "delivery_address"
            ],
            "properties": {
                "client_email": {
                    "type": "string"
                },
                "client_name": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "client_phone": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 1
                },
                "client_surname": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "delivery_address": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                }
            }
        },
        "controller.UpdateOrderInProduct": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "minimum": 0
                },
                "price": {
                    "type": "number",
                    "minimum": 0
                },
                "quantity": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "middleware.ErrorJSON": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "details": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "error": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Orders API",
	Description:      "API документация для orders",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
