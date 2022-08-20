// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "JeongMinSik",
            "url": "https://github.com/JeongMinSik/go-leaderboard",
            "email": "jms6025a@naver.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "테스트용",
                "tags": [
                    "test"
                ],
                "responses": {
                    "200": {
                        "description": "Hello go-leaderboard",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/teapot": {
            "get": {
                "description": "테스트용",
                "tags": [
                    "test"
                ],
                "responses": {
                    "418": {
                        "description": "I'm a teapot",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "name으로 User의 score와 rank를 얻습니다.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Show a user info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/leaderboard.UserRank"
                        }
                    },
                    "400": {
                        "description": "name query param 확인 필요",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            },
            "post": {
                "description": "신규 user를 추가합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Add a user",
                "parameters": [
                    {
                        "description": "New User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/leaderboard.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/leaderboard.UserRank"
                        }
                    },
                    "400": {
                        "description": "request body 확인 필요",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            },
            "delete": {
                "description": "기존 user를 삭제합니다.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.deleteData"
                        }
                    },
                    "400": {
                        "description": "name 확인 필요",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            },
            "patch": {
                "description": "기존 user를 수정합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update a user",
                "parameters": [
                    {
                        "description": "Updated User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/leaderboard.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/leaderboard.UserRank"
                        }
                    },
                    "400": {
                        "description": "request body 확인 필요",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            }
        },
        "/users/count": {
            "get": {
                "description": "전체 유저 수",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user count",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.userCountData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            }
        },
        "/users/{start}/to/{stop}": {
            "get": {
                "description": "user list 를 받아옵니다.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "start index",
                        "name": "start",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "stop index",
                        "name": "stop",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/leaderboard.UserRank"
                            }
                        }
                    },
                    "400": {
                        "description": "param 확인 필요",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    },
                    "500": {
                        "description": "서버에러",
                        "schema": {
                            "$ref": "#/definitions/handler.messageData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.deleteData": {
            "type": "object",
            "properties": {
                "is_deleted": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.messageData": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.userCountData": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                }
            }
        },
        "leaderboard.User": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                }
            }
        },
        "leaderboard.UserRank": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "rank": {
                    "type": "integer"
                },
                "score": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:6025",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Leaderboard API",
	Description:      "go언어로 만든 리더보드 api 토이프로젝트입니다.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
