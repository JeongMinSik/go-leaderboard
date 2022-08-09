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
