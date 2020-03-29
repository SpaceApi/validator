package main

var openapi = `{
  "openapi": "3.0.2",
  "info": {
    "description": "This is the SpaceApi Validator api",
    "version": "1.1.0",
    "title": "SpaceApi Validator"
  },
  "servers": [
    {
      "url": "https://validator.spaceapi.io",
      "description": "The SpaceApi Validator Service"
    }
  ],
  "paths": {
    "/v1": {
      "get": {
        "deprecated": true,
        "tags": [
          "v1"
        ],
        "responses": {
          "200": {
            "description": "get default information about the server",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ServerInformation"
                }
              }
            }
          }
        }
      }
    },
    "/v1/validate/": {
      "post": {
        "deprecated": true,
        "tags": [
          "v1"
        ],
        "summary": "validate an input against the SpaceApi schema",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ValidateV1"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ValidateV1Response"
                }
              }
            }
          },
          "400": {
            "description": "request body is malformed"
          },
          "500": {
            "description": "something went wrong"
          }
        }
      }
    },
    "/v2": {
      "get": {
        "tags": [
          "v2"
        ],
        "responses": {
          "200": {
            "description": "get default information about the server",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ServerInformation"
                }
              }
            }
          }
        }
      }
    },
    "/v2/validateURL": {
      "post": {
        "tags": [
          "v2"
        ],
        "summary": "validate an input against the SpaceApi schema",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ValidateUrlV2"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ValidateUrlV2Response"
                }
              }
            }
          },
          "400": {
            "description": "request body is malformed"
          },
          "500": {
            "description": "something went wrong"
          }
        }
      }
    },
    "/v2/validateJSON": {
      "post": {
        "tags": [
          "v2"
        ],
        "summary": "validate an input against the SpaceApi schema",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ValidateJsonV2"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ValidateJsonV2Response"
                }
              }
            }
          },
          "400": {
            "description": "request body is malformed"
          },
          "500": {
            "description": "something went wrong"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "ServerInformation": {
        "properties": {
          "description": {
            "type": "string"
          },
          "usage": {
            "type": "string"
          },
          "version": {
            "type": "string"
          }
        },
        "required": [
          "description",
          "usage",
          "version"
        ]
      },
      "ValidateV1": {
        "properties": {
          "data": {
            "type": "object"
          }
        },
        "required": [
          "data"
        ]
      },
      "ValidateUrlV2": {
        "properties": {
          "url": {
            "type": "string",
            "pattern": "uri"
          }
        },
        "required": [
          "url"
        ]
      },
      "ValidateJsonV2": {
        "type": "object"
      },
      "ValidateV1Response": {
        "properties": {
          "valid": {
            "type": "boolean"
          },
          "message": {
            "type": "string"
          }
        },
        "required": [
          "valid",
          "message"
        ]
      },
      "ValidateUrlV2Response": {
        "properties": {
          "valid": {
            "type": "boolean"
          },
          "message": {
            "type": "string"
          },
          "isHttps": {
            "type": "boolean"
          },
          "httpsForward": {
            "type": "boolean"
          },
          "reachable": {
            "type": "boolean"
          },
          "cors": {
            "type": "boolean"
          },
          "contentType": {
            "type": "boolean"
          },
          "certValid": {
            "type": "boolean"
          },
          "validatedJson": {
            "type": "object"
          },
          "schemaErrors": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/SchemaError"
            }
          }
        },
        "required": [
          "valid",
          "isHttps",
          "httpsForward",
          "reachable",
          "cors",
          "contentType",
          "certValid"
        ]
      },
      "ValidateJsonV2Response": {
        "properties": {
          "valid": {
            "type": "boolean"
          },
          "message": {
            "type": "string"
          },
          "validatedJson": {
            "type": "object"
          },
          "schemaErrors": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/SchemaError"
            }
          }
        },
        "required": [
          "valid",
          "message"
        ]
      },
      "SchemaError": {
        "properties": {
          "field": {
            "type": "string"
          },
          "message": {
            "type": "string"
          }
        },
        "required": [
          "field",
          "message"
        ]
      }
    }
  }
}`
