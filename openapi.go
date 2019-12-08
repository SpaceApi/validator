package main

var openapi = `{
  "openapi": "3.0.2",
  "info": {
    "description": "This is the SpaceApi Validator api",
    "version": "0.0.1",
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
        "tags": [
          "v1"
        ],
        "summary": "validate an input against the SpaceApi schema",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "data": {
                    "type": "object"
                  }
                },
                "required": [
                  "data"
                ]
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
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "url": {
                    "type": "string"
                  }
                },
                "required": [
                  "url"
                ]
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
          "content": {
            "application/json": {
              "schema": {
                "type": "object"
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
      }
    }
  }
}`
