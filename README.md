# Cleopatchra &middot; [![Build Status](https://travis-ci.org/jfo84/cleopatchra.svg?branch=master)](https://travis-ci.org/jfo84/cleopatchra) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Cleopatchra downloads and stores GitHub state in a PostgreSQL database and powers a Golang server that renders JSON:API-spec data.

## Running the Project

The project is written in Golang and Ruby. To start, you'll need to download the repository into your `GOPATH`. If you don't know what that is, you can get started with Go [here](https://golang.org/doc/install).

Golang

- `cd` into the root of the directory
- Install dependent packages with `glide install`. If you don't have glide installed, you can install it [here](https://glide.sh/)
- Compile the program by running `go build`
- Run the server with `./cleopatchra`
- You'll need to setup a database and seed it with a small Ruby app before the app will fully function (see below)

Ruby

- `cd` into the `seed-db` directory from root
- Install [PostgreSQL](https://www.postgresql.org/docs/10/static/tutorial-install.html)
- Set environment variables for `DEFAULT_POSTGRES_USER` and `DEFAULT_POSTGRES_PASSWORD` to configure PostgreSQL access
- Install Ruby version 2.4.1 with a tool such as [RVM](https://rvm.io/)
- Install dependencies with `bundle install`
- Seed the database for a given repository with `ruby run.rb seed --organization 'foo_org' --repo 'bar_repo'`
