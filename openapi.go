package main

var openapi = `{
  "openapi": "3.0.2",
  "info": {
    "description": "This is the SpaceApi Validator api",
    "version": "1.3.0",
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
        "deprecated": true,
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
        "deprecated": true,
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
        "deprecated": true,
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
    },
    "/v3": {
      "get": {
        "tags": [
          "v3"
        ],
        "summary": "A short summary of the service and how to use it.",
        "responses": {
          "200": {
            "description": "A short summart of the service and how to use it.",
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
    "/v3/validateURL": {
      "post": {
        "tags": [
          "v3"
        ],
        "summary": "Validate a SpaceAPI endpoint URL against all SpaceApi schema versions it claims to support.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ValidateUrlV3"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful opperation. Whether the document passed validation can be taken from the 'valid' key in the response.",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ValidateUrlV3Response"
                }
              }
            }
          },
          "400": {
            "description": "The request body is malformed."
          },
          "429": {
            "description": "The request was rate limited."
          },
          "500": {
            "description": "Something went wrong."
          }
        }
      }
    },
    "/v3/validateJSON": {
      "post": {
        "tags": [
          "v3"
        ],
        "summary": "Validate a JSON document against all SpaceApi schema versions it claims to support.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ValidateJsonV3"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful opperation. Whether the document passed validation can be taken from the 'valid' key in the response.",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ValidateJsonV3Response"
                }
              }
            }
          },
          "400": {
            "description": "The request body is malformed."
          },
          "500": {
            "description": "Something went wrong."
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
            "type": "string",
            "description": "A short description of the service.",
            "example": "The SpaceAPI validator service."
          },
          "usage": {
            "type": "string",
            "description": "A short summary of how to use the service.",
            "example": "To validate a JSON document, send a POST request to /v3/validateJSON. To validate an endpoint URL, send a POST request to /v3/validateURL. See /swagger for an API documentation."
          },
          "version": {
            "type": "string",
            "description": "Server version information.",
            "example": "1.3.0"
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
          "checkedVersions": {
            "type": "array",
            "items": {
              "type": "string"
            }
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
          "checkedVersions": {
            "type": "array",
            "items": {
              "type": "string"
            }
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
      "ValidateJsonV3": {
        "type": "object",
         "description": "The JSON document to validate."
      },
      "ValidateJsonV3Response": {
        "properties": {
          "valid": {
            "type": "boolean",
            "description": "Whether the JSON document was valid according to all SpaceAPI schema versions it claims to support."
          },
          "message": {
            "type": "string"
          },
          "checkedVersions": {
            "type": "array",
            "description": "The list of SpaceAPI schema versions the document was validated against, both from the 'api' and 'api_compatibility' keys.",
            "items": {
              "type": "string"
            }
          },
          "validatedJson": {
            "type": "object",
            "description": "The JSON document that was validated."
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
          "checkedVersions"
        ]
      },
      "ValidateUrlV3": {
        "properties": {
          "url": {
            "type": "string",
            "description": "The URL of the endpoint to validate.",
            "pattern": "uri",
            "example": "http://example.org/spaceapi"
          },
          "extendedValidation": {
            "type": "boolean",
            "description": "Whether to perform some extended validation (e.g. IP dual-stack capability, User-Agent filtering). The exact set of validation parameters is implementation dependent, but is returned in the response. If false, the response is guaranteed to only contain a single result. If true, multiple results may be returned in the response.",
            "default": false
          },
          "withValidatedJson": {
            "type": "boolean",
            "description": "Whether the endpoint response should be included in the validator response.",
            "default": true
          }
        },
        "required": [
          "url"
        ]
      },
      "ValidateUrlV3Parameters": {
        "properties": {
          "addressFamily": {
            "type": "string",
            "description": "The IP address family restriction used to connect to the endpoint.",
            "enum": [
              "any",
              "ipv6",
              "ipv4"
            ]
          },
          "userAgent": {
            "type": "string",
            "description": "The User-Agent header used to connect to the endpoint.",
            "example": "Go-http-client/1.1"
          }
        },
        "required": [
          "addressFamily",
          "userAgent"
        ]
      },
      "ValidateUrlV3Result": {
        "properties": {
          "parameters": {
            "$ref": "#/components/schemas/ValidateUrlV3Parameters",
            "description": "The validation parameters (e.g. IP address family, HTTP User-Agent) this result applies to."
          },
          "originalUrl": {
            "type": "string",
            "description": "The URL as it was submitted to the validator, possibly with some normalizations (e.g. fragment removal) applied.",
            "pattern": "uri",
            "example": "http://example.org/spaceapi"
          },
          "redirectedUrl": {
            "type": "string",
            "description": "The final URL encountered during the request and subsequent redirects. If no redirect occurred, this is the same as 'originalUrl'.",
            "pattern": "uri",
            "example": "https://example.org/spaceapi"
          },
          "reachable": {
            "type": "boolean",
            "description": "True if the endpoint accepted the connection and responded to the HTTP request with a (not necessarily successful) HTTP response. (NOTE: This is different from v2, where 'reachable' also implied a HTTP 200 OK response.)",
            "example": true
          },
          "httpStatus": {
            "type": "integer",
            "description": "The HTTP status code to the final request. 0 if 'reachable' is false.",
            "minimum": 0,
            "maximum": 999,
            "example": 200
          },
          "isHttps": {
            "type": "boolean",
            "description": "Whether the final request was done using HTTPS.",
            "example": true
          },
          "httpsRedirect": {
            "type": "boolean",
            "description": "Whether the request was upgraded from HTTP to HTTPS.",
            "example": true
          },
          "permanentRedirect": {
            "type": "boolean",
            "description": "True if all HTTP redirects encountered were permanent (HTTP status code 301 or 308).  Absent if no redirect was encountered.",
            "example": true
          },
          "cors": {
            "type": "boolean",
            "description": "True if all HTTP responses (including redirects) set an 'Access-Control-Allow-Origin: *' header.",
            "example": true
          },
          "contentType": {
            "type": "boolean",
            "description": "Whether the final response serves the correct content type header: 'Content-Type: application/json'",
            "example": true
          },
          "invalidCert": {
            "type": "boolean",
            "description": "Whether there was a problem with HTTPS certificate validation. Always true for plain HTTP."
          },
          "validation": {
            "$ref": "#/components/schemas/ValidateJsonV3Response",
            "description": "The result of the SpaceAPI schema validation itself."
          }
        },
        "required": [
          "parameters",
          "originalUrl",
          "redirectedUrl",
          "reachable",
          "isHttps",
          "httpsRedirect",
          "cors",
          "contentType",
          "invalidCert",
          "validation"
        ]
      },
      "ValidateUrlV3Response": {
        "properties": {
          "results": {
            "type": "array",
            "description": "The list of paramter sets and their results.",
            "items": {
              "$ref": "#/components/schemas/ValidateUrlV3Result"
            }
          },
          "reachable": {
            "type": "boolean",
            "description": "True if the endpoint was reachable under at least one parameter set.",
            "example": true
          },
          "valid": {
            "type": "boolean",
            "description": "True if the endpoint was reachable under ALL parameter sets it was reachable under.",
            "example": true
          }
        },
        "required": [
          "results",
          "reachable",
          "valid",
          "inconsistencies"
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
