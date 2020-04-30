# Learning Go
Learning Golang - with Security in mind

Basically, I'm adding to this from the top down so you can consider the lowest file structure the most up to date/ 'secure' that I've got for now.

We're going to build a generic framework for logging in as a user to a public area (customised) and a series of private areas). The main functions will thus be:

- [x] Basic HTTP Server, with OWASP Zap! Passing baseline
- [x] X-Clacks-Overhead everywhere
- [x] Deliver static files such as CSS and JS
- [x] Deliver HTML based on programmed templates
- [x] Log in functionality
- [ ] Log out functionality
- [ ] RBAC: Public
- [ ] RBAC: Authenticated
- [ ] RBAC: Administrator
- [ ] Simple API
- [ ] Simple API - versioned
- [ ] Simple API - token authentication

After I know enough to do all of these things, I will probably abandon the project for something else I'm working on :)

## Beginning Server

1. [Basic Server with HTTP Tests and No Security](begserver/begserver.md)
2. [Basic Server with HTTP Tests and Some Security](begserver/hardened/begserverhardened.md)

## Adding static file delivery

1. [Basic Server - Static File Delivery](fileserver/fileserver.md)

## Porting to Gorilla

1. [Basic Server - Gorilla Mux](gorilla/gorilla.md)

## HTML Templates

1. [HTML Templates](templates/templates.md)

## Log in functionality

1. [Login Functionality - Not yet secure](login/login.md)
2. [Login Functionality - More secure](login/logincsrf.md)
