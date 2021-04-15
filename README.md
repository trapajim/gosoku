# GOsoku

Goal of this repository is to have an easy way to generate and administrate JSON API's. 

## Commands

*new*
setup a new project 
```bash
    # example
    gosoku new
```

*generate scaffold (shortcut g s)* which generates your domain, json http handler, usecases and repository
```bash
    # example
    gosoku g s User Name:string Age:int
    # default the form builder tags are set automatically to input fields
    # if you would like to change that you can append a supported form type after the datatype
    gosoku g s User Name:string:textarea Age:int
```

the above command will generate following structure
```bash
app
    domain
        user.go #interface and model definitions
    user
        delivery #handles http requests and route group is defined here
            json.go
        repository #database transactions
            repository.go
        usecase # business logic here
            usecase.go 
system
    router
        user.go
```

*clean*
remove a generated type

```bash
    # example
    gosoku clean user
    #or 
    gosoku c user 
    # this will remove domain/user.go, router/user.go and the user directory
```

*sentry*
add sentry and the sentry echo middleware

see https://docs.sentry.io/platforms/go/guides/echo/ for more information 

```bash
    gosoku sentry
    #or 
    gosoku s
```

## Formbuilder

Forms can be automatically created based on struct tags.

**example**

```go
type User struct {
	Name      string    `json:"name" form:"input,type=text,name=Name"`
	Biography string    `json:"biography" form:"textarea,name=Biography"`
}

// this will generate a html form with an input text field and a textarea
formbuilder.Build(domain.User{})

// the generated form will have the values set from fields.
formbuilder.Build(domain.User{Name:"Gosoku", Biography: "Test"})
```


## todo 

* [ ] sql migrations
* [ ] admin dashboard
* [ ] support for multiple datasources
