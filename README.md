# Cleopatchra

Cleopatchra downloads and stores GitHub state in a Postgres database, powers a Golang server that renders a structure with JSON API relationships, and re-displays the state quickly and easily.

## Running the Project

The project is written in Golang, JavaScript, and Ruby. To start, you'll need to download the repository into your `GOPATH`. If you don't know what that is, you can get started with Go [here](https://golang.org/doc/install).

Golang
- `cd` into the root of the directory
- Install dependent packages with `glide install`. If you don't have glide installed, you can install it [here](https://glide.sh/)
- Compile the program by running `go build`
- Run the server with `./cleopatchra`
- You'll need to setup a database and seed it with a small Ruby app before the app will fully function (see below)

Ruby
- `cd` into the `seed_db` directory from root
- Install [PostgreSQL](https://www.postgresql.org/docs/10/static/tutorial-install.html)
- Install Ruby version 2.4.1 with a tool such as [RVM](https://rvm.io/)
- Install dependencies with `bundle install`
- Seed the database for a given repository with `ruby run.rb seed --organization 'foo_org' --repo 'bar_repo'`

JavaScript
- The front end uses the Ember.js framework. You will need a current version (engine 8+) of [Node](https://nodejs.org/) and [yarn](https://yarnpkg.com/) installed
- Install the Ember command-line tools with `yarn global add ember-cli`
- `cd` into the `front_end` directory from root
- Install our JavaScript dependencies with `yarn`
- Serve the app with `ember serve`