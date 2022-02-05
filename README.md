## True Tickets  Metrics System

This project aims to create a metrics system to store and retrieve information related to metrics.

### Dependencies

To fully use the project's capabilities, the following dependencies are needed:

- Golang 1.17+
- Swagger UI: https://github.com/wordnik/swagger-ui.git (you should clone this)

### Usage

#### Test

```
$ make test
```

#### Get deps

```
$ make deps
```

#### Build

```
$ make build
```

#### Run without swagger

```
$ make run-without-swagger
```

#### Run with swagger

For this case, you can get only the openapi specification, or you can make use of swagger-ui capabilities. Using swagger-ui capabilities require the following:

- Clone https://github.com/wordnik/swagger-ui.git
- The project contains a "dist" folder. Its contents has all the Swagger UI files you need.
- The index.html has an url set to http://petstore.swagger.wordnik.com/api/api-docs. You need to change that to match your WebService JSON endpoint e.g. http://localhost:8080/apidocs.json
- Create an environment variable setting swagger-ui `dist` folder. Eg: `export SWAGGER_UI=/Users/viniciusnordi/Desktop/Jobs/swagger-ui/dist`


```
$ make run 
```

- Accessing only specification: http://localhost:8080/apidocs.json
- Accessing swagger-ui: http://localhost:8080/apidocs/