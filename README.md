# ![RealWorld Example App](logo.png)

> ### [Gorilla Mux] codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged backend service built with **[Gorilla Mux]** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **[Gorilla Mux]** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.

# How it works

> I'am using Gorilla Mux for http router, and sqlx for relational database interaction.

## Prerequisite
- Git [Download](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- Go [Download](https://go.dev/dl/)
- Docker [Download](https://docs.docker.com/get-docker/)

# Getting started

- Clone this repo.
  ```bash
  git clone https://github.com/ashalfarhan/realworld.git
  ```
- Install dependencies
  ```bash
  go mod tidy
  ```
- Start PostgreSQL and Adminer
  ```bash
  docker-compose up -d
  ```
- Start the server
  ```bash
  go run .
  ```
  Or using hot reload like [Air](https://github.com/cosmtrek/air)


## Testing
- Unit Testing
  ```bash
  make test
  ```
- E2E Testing with Conduit Spec
  ```bash
  make test-spec
  ```

## Todo
- [ ] Article Service Test
- [ ] List articles by favorited user
- [ ] List articles by author
- [ ] Containerize with Docker
- [ ] Caching with Redis
