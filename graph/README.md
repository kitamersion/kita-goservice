# GraphQL Local Development Guide

This project uses [gqlgen](https://github.com/99designs/gqlgen) for GraphQL support in Go.

## Workflow Overview

Working directory: `graph`

1. **Edit the GraphQL schema** (`schema.graphqls`)
2. **Generate Go code** from the schema using `make generate`
3. **Implement resolver logic** in the generated resolver files
4. **Run the server locally** and test your API

---

## 1. Editing the GraphQL Schema

- The schema is defined in `schema.graphqls`
- Update types, queries, mutations, and inputs as needed.
- Example:

  ```graphql
  type User {
    id: ID!
    name: String!
    email: String!
  }

  type Query {
    user(id: ID!): User
    users: [User!]!
  }

  type Mutation {
    createUser(name: String!, email: String!): User!
  }
  ```

## 2. Generating Go Code

- Run the following command to generate Go types and resolver stubs:
  ```sh
  make generate
  # or
  go run github.com/99designs/gqlgen generate
  ```
- This will update files in `graph/` and create or update resolver stubs in `graph/`.

## 3. Implementing Resolvers

- Open the resolver files in `graph/` (e.g., `schema.resolvers.go`).
- Implement the business logic for each resolver function.
- Example:
  ```go
  func (r *queryResolver) User(ctx context.Context, id string) (*generated.User, error) {
      // Fetch user from your service or database
  }
  ```

## 4. Running the Server Locally

- Start your API server:
  ```sh
  go run cmd/api/main.go
  # or use Docker Compose if configured
  make run
  ```
- Access the GraphQL Playground at [http://localhost:8080/playground](http://localhost:8080/playground) (or the port you configured).

## 5. Regenerating Code After Schema Changes

- Any time you update `schema.graphql`, rerun `make generate` to update Go types and resolvers.
- Implement any new or changed resolvers as needed.

## 6. Tips

- Use the generated types from the `generated` package in your resolver implementations.
- Keep your schema and Go code in sync by always running code generation after schema changes.
- For more info, see the [gqlgen docs](https://gqlgen.com/getting-started/).
