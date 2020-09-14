# GOsoku

Goal of this repository is to have an easy way to generate and administrate JSON API's with rails like generators. The Generated code should be modulary structured so a monolithic application could be easily refactored into an microservice approach.


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
```