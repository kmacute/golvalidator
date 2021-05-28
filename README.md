# golvalidator
A golang laravel inspired validation


### Installation

Install the package using
```go
$ go get github.com/kmacute/golvalidator
```

### Usage

To use the package import it in your `*.go` code
```go
import "github.com/kmacute/golvalidator"
```

### Example
***Validate `struct` only at the moment***

```
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmacute/golvalidator"
)

func main() {
	app := fiber.New()

	app.Get("/", Test)

	app.Listen(":3002")
}

type UserType struct {
	FirstName  string `json:"first_name" validate:"required_with:last_name"`
	MiddleName string `json:"middle_name" validate:"required|string|min:3|same:last_name"`
	LastName   string `json:"last_name" validate:"required|string|min:3"`
}

func Test(c *fiber.Ctx) error {
	user := UserType{
		FirstName:  "",
		LastName:   "e1",
		MiddleName: "3",
	}

	errors := golvalidator.ValidateStructs(user)
	if len(errors) > 0 {
		return c.JSON(fiber.Map{
			"errors": errors,
		})
	}

	return c.SendString("No Errors")
}

```

### Errors
```
{
    "errors": {
        "first_name": [
            "The first name field is required when last name is present."
        ],
        "last_name": [
            "The last name must only contain letters.",
            "The last name must be at least 3 characters."
        ],
        "middle_name": [
            "The middle name must only contain letters.",
            "The middle name must be at least 3 characters.",
            "The middle name and last name must match."
        ]
    }
}
```

### Available Validation
```
alpha
string
numeric
alpha_num
alpha_space
alpha_dash
date
email
same
digits
digits_between

min,max,between:
-string length
-numeric

lt,gt,lte,gte:
-string length
-numeric

nullable
required
required_if
required_with
ip
ipv4
ipv6
url
credit_card
```

### Pending
```
Unique -> (Gorm Dependency)
File Size
File Extension
```
