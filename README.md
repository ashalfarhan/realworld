# ![RealWorld Example App](logo.png)

> ### [Gorilla Mux] codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged backend service built with **[Gorilla Mux]** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **[Gorilla Mux]** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

> I'am using Gorilla Mux for http router, and sqlx for relational database interaction.

# Getting started

- Clone this repo.

# Features
- User
  - [x] Registration
  - [x] Login
  - [x] Get Current User
  - [x] Update User

- Profile
  - [x] Get Profile
  - [x] Follow
  - [x] Unfollow

- Article
  - [x] Create Article
  - [x] Delete Article
  - [x] Get Tags
  - [x] Get Article by Slug
  - [x] List Articles
    - [x] Paginated
    - [x] Filter by tag
    - [ ] Filter by favorited of a user
    - [ ] Filter by author
  - [x] Feed Articles
    - [x] Paginated
  - [x] Update Articles
  - [x] Favorite Articles
  - [x] Unfavorite Articles

- Comment
  - [ ] Add Comment
  - [ ] Delete Comment
  - [ ] Get Comment of an Article