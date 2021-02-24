# gin-swagger

Generates code from a Swagger spec for GIN Framework. Currently, only server code generation is supported.

_**Disclaimer:**_ This project uses code generation framework on https://github.com/go-swagger/go-swagger. Also, templates from https://github.com/mikkeloscar/gin-swagger have been used as a base for the templates in this project.

# Install

```sh
go get github.com/emreisikligil/gin-swagger
```

or install from the source

```sh
go install .
```

# GIN Server Generation

## CLI Usage

```sh
Usage:
  gin-swagger server [flags]

Flags:
  -p, --api-package string      api-package name (default "operations")
  -c, --client-package string   client-package name (default "client")
      --exclude-spec            Excludes spec if set
  -h, --help                    help for server
      --include-handler         Generates operation handlers if set (default true)
      --include-main            Generates main function if set
      --include-models          Generates models if set (default true)
      --include-support         Generates support docs if set (default true)
      --include-validation      Generates validators if set (default true)
  -m, --model-package string    model-package name (default "models")
  -n, --name string             name (default "name")
  -s, --server-package string   server-package name (default "restapi")
      --skip-operations         Skips operations if set
      --spec string             spec name (default "api/swagger.yml")
      --target string           target name (default "./")
  -t, --toggle                  Help message for toggle
```

A sample server is provided under [example](./example) folder. It has been generated using

```sh
cd example
gin-swagger server --name example --spec spec.yml 
```

## Authorization

API key authorization and oauth2 application flow are supported currently. Check [Swagger](https://swagger.io/docs/specification/2-0/authentication) for more details about the authentications methods.


Consider the following security definitions.

```yaml
securityDefinitions:
  oauth2:
    type: oauth2
    authorizationUrl: http://petstore.swagger.io/oauth/dialog
    flow: application
    scopes:
      read: read your pets
      write: modify pets in your account
  apiKey:
    type: apiKey
    name: X-API-Key
    in: header
    scopes:
      read: read your pets
      write: modify pets in your account
```

It adds the following 2 functions to the service interface

```go
APIKeyAuthenticate(*gin.Context, string, []string) (interface{}, error, int)
Oauth2Authenticate(*gin.Context, string, []string) (interface{}, error, int)
```

These functions are supposed to handle authorization. Both functions return the same triplet which is (principal, error, status_code). status_code is used as http response code when error is not nil. Otherwise, it is not used. principal is the entity who makes the request. If the authorization succeeds these functions should return an object identifying the caller. Object returned by one of these functions is added to the context with "Principal" key.

The first and the third arguments of this functions are the gin context and the required scopes. The second argument is the API key for API key authorization and Authorization header for oauth2 authorization.

Assume an operation is tagged with

```yaml
...
security:
  - apiKey: [read]
  - oauth2: [read]
...
```

Authenticator functions will be called in the given order. If any of the functions returns success the following authenticators are skipped and the request is passed to the handler. If none of the authenticators returns success status_code and error returned by the last authenticator will be returned as HTTP response.

```go
APIKeyAuthenticate(ctx, "<api_key>", []string{"read"})
Oauth2Authenticate(ctx, "<authorization_token>", []string{"read"})
```

## Service Interface

API operations are added to the service interface along with the authenticators. All you have to do is to implement functions in this interface.

```go
type ExampleService interface {
  // Authenticator Handlers
  APIKeyAuthenticate(*gin.Context, string, []string) (interface{}, error, int)
  Oauth2Authenticate(*gin.Context, string, []string) (interface{}, error, int)

  // API Operation Handlers
  AddPet(ctx *gin.Context, params *pet.AddPetParams) api.APIResponse
  GetPetByID(ctx *gin.Context, params *pet.GetPetByIDParams) api.APIResponse
```

## Cors

Generated code will have the following CORS configuration. Currently, it is not possible to modify this configuration through CLI. Templates may be edited in order to modify this behaviour.

```go
import "github.com/gin-contrib/cors"
...

config := cors.DefaultConfig()
config.AllowHeaders = append(config.AllowHeaders, "Authorization")
config.AllowAllOrigins = true
server.Router.Use(cors.New(config))
```

# Editing Templates

Templates can be edited using the following steps.

1. Edit templates
1. Use `go-bindata` to generate template binaries again.

   ```
   go get -u github.com/jteeuwen/go-bindata

   go-bindata -o cmd/bindata.go -pkg cmd templates

   go install .
   ```
