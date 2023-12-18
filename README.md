# DBO-MANAGEMENT-APP

Brief description of your project.

## Prerequisites

- Docker: [Install Docker](https://docs.docker.com/get-docker/)

## Getting Started

Follow these instructions to get your project up and running.

### Clone the Repository

    git clone https://github.com/dakochan666/dbo-management-app.git
    cd dbo-management-app

### Build Docker Image

    docker build -t dbo-management-app .

### Run the Docker Container

    docker run -p 8080:8080 dbo-management-app

The application should now be accessible at http://localhost:8080.

## How to run without docker

Go Program can be operated with two ways

1. Doing Compilation (*compile*) at first, then execution file the product of its compilation.

    ```bash
    $ go mod tidy
    $ go build main.go
    $ ./main
    ```

2. Without compilation

    ```bash
    $ go mod tidy
    $ go run main.go
    ```

## ERD

```dbml
// Use DBML to define your database structure

Table products {
  id integer [primary key]
  name varchar
  stock integer
  description varchar
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}

Table users {
  id integer [primary key]
  name varchar
  email varchar
  password varchar
  role varchar
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}

Table orders {
  id integer [primary key]
  user_id integer
  product_id integer
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}

Ref: products.id < orders.product_id

Ref: users.id < orders.user_id
```

![DB Diagram](/src/erd.png)

## Documentation

https://documenter.getpostman.com/view/13235416/2s9Ykn8MkY

## Admin Account
    
    email = admin@mail.com
    password = admin123
    
