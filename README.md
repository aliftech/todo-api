# Go TODO API

🚀 Exciting news, folks! Just dropped my latest portfolio update featuring a sleek, efficient todo REST-API built from scratch using the power of Go programming language, along with the dynamic duo of Gin and Gorm frameworks, topped off with CompileDaemon for seamless development. Dive into the world of RESTful APIs with me as I unravel the magic behind crafting this robust solution. Check out the link below for all the juicy details! 💻✨

## Instalation

### Pre-requisites

Before installing this project, you should know about the project requirements. Here are the requirements that you need to prepared:

1. Go Programming Laguage
2. Knowledge about Go.
3. Knowledge about REST-API

### Dependencies

This project is build using some third parties or dependencies. Here are the required dependencies:

1. Gin (https://gin-gonic.com/)
2. Gorm (https://gorm.io/)
3. Go dotenv (https://github.com/joho/godotenv)
4. CompileDaemon (https://github.com/githubnemo/CompileDaemon)
5. MySQL driver (https://gorm.io/docs/connecting_to_the_database.html)

### Database Migration

In this project, we only use a single migration, because the table is very simple - only consists of one database table named tasks. Before we step to the migration, you have to make a database called todo in your XAMPP.

After that, you can run the database migration by the following command:

```bash
go run migrations/migrate.task.go
```

### Running The Project

To run the project, you can use a usual go command:

```bash
go run main.go
```

But, in this project we are not going to use it. Instead of using `go run` command, we will using CompileDaemon to run our application. Here is the command to run our application using CompileDaemon:

```bash
CompileDaemon -command="./todo-api"
```

## DEPLOYMENT

```bash
gcloud auth configure-docker asia-southeast2-docker.pkg.dev
```

```bash
docker compose -f docker-compose.yml build
```

```bash
docker tag todoapp asia-southeast2-docker.pkg.dev/todo-429911/todo-api-todo
```

```bash
docker push asia-southeast2-docker.pkg.dev/todo-429911/todoapp/todo-api-todo
```
