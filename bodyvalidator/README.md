# Bodyvalidator

`fiber-modules/bodyvalidator` is a module for validating body json
on request.

## Usage

```go
import "github.com/gkampitakis/fiber-modules/bodyvalidator"

validator := bodyvalidator.New()

app.Post("/route",validator(bodyvalidator.Config{
	SchemaPath: "path/from/root/to/file.json",
},handler)
```

### Options

`fiber-modules/bodyvalidator` can be configured globally when calling `New()`. 
The two options you can pass are

- `bodyvalidator.ExposeErrors()` an option for exposing a description
on what validation errors occurred
- `bodyvalidator.RegisterKeywords(keywords)` with this options you can pass 
custom keywords to the underlying library used 
[jsonschema](https://github.com/qri-io/jsonschema). For more information about
[Custom Keywords](https://github.com/qri-io/jsonschema#custom-keywords)
- `bodyvalidator.SetResponse()` You can override the default bad request response.

  Example 

  ```go
  response := func(ke []jsonschema.KeyError) interface{} {
    return map[string]interface{}{
      "message": "oops you messed it up there",
    }
  }
  validator := bodyvalidator.New(bodyvalidator.SetResponse(response))
  ```

  > Please be mindful that if you override the default response the option
  `ExposeErrors` is not going to be in effect as the provided function is in 
  control of what's returned in the response.

The second way of configuring `fiber-modules/bodyvalidator` is per route level
by providing the `json schema` to validate the objects with. It can be either a
json file or a string literal representing a json.

- `SchemaPath` path to json file containing json schema
- `SchemaLiteral` string literal representing json schema

If both options are present `SchemaLiteral` takes precedence.


## Examples

### Simple example

```go
import "github.com/gkampitakis/fiber-modules/bodyvalidator"

validator := bodyvalidator.New(bodyvalidator.ExposeErrors())

app.Post("/route",validator(bodyvalidator.Config{
	SchemaPath: "path/from/root/to/file.json",
},handler)
```


### Custom Keyword

Registering a new keyword `email`.

```go
var emailRgx, _ = regexp.Compile("^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$")

type IsEmail bool

func newIsEmail() jsonschema.Keyword {
	return new(IsEmail)
}

func (e *IsEmail) Register(uri string, registry *jsonschema.SchemaRegistry) {}

func (e *IsEmail) Resolve(pointer jptr.Pointer, uri string) *jsonschema.Schema {
	return nil
}

// ValidateKeyword implements jsonschema.Keyword
func (f IsEmail) ValidateKeyword(ctx context.Context, currentState *jsonschema.ValidationState, data interface{}) {
	if str, ok := data.(string); ok {
		if !emailRgx.Match([]byte(str)) {
			currentState.AddError(data, fmt.Sprintf("should be comform to email. plz make '%s' == type email. plz", str))
		}
	}
}

validator := bodyvalidator.New(
	bodyvalidator.ExposeErrors(),
	bodyvalidator.RegisterKeywords(bodyvalidator.Keywords{"email": newIsEmail}),
)

...

app.Post("/route",validator(bodyvalidator.Config{
	SchemaLiteral: `{
  "type": "object",
  "properties": {
    "user": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "email": {
          "type": "string",
          "email": true
        }
      }
    }
  }
}`,
}),handler)
```

So in the above example if we do a post request 

```bash
curl --header 'Content-Type: application/json' --request POST http://localhost:8080/route --data '{"user":{"name":"george","email":1}}'
```

with `"email":1` we will get a response of `Bad request` and `status_code` 400

```json
{
    "description": [
        {
            "propertyPath": "/user/email",
            "invalidValue": "1",
            "message": "should be comform to email. plz make '1' == type email. plz"
        }
    ],
    "statusCode": 400,
    "message": "Bad request"
}
```

Also we can notice the description field that's enabled by calling `bodyvalidator.ExposeErrors()`.