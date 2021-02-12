# Fold

Fold is a platform and a set of tools for developing and deploying backends. 
The goal is to make running a set of microservices on the cloud as easy as
running them locally.

**WARNING: This is seriously alpha...**

# Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Speed Run](#speed-run)
- [Working with Fold](#working-with-fold)
  * [Project templates](#project-templates)
  * [Project Structure](#project-structure)
  * [Project Config](#project-config)
  * [Hot reloading](#hot-reloading)
  * [SDKs](#sdks)
- [Deployment](#deployment)
- [Fold Runtime](#fold-runtime)
  * [Base Images](#base-images)
  * [Custom Images](#custom-images)
- [License](LICENSE)

# Overview

Fold aims to take cloud based infrastructure and serverless computing to its
ultimate conclusion. Infrastructure shouldn't just be easy, we shouldn't have to
think about it at all.

It has happened countless times already in our industry. Compilers removed the
need to think about machine code, garbage collection freed us from managing
memory, web frameworks liberated us from network programming and boilerplate,
and cloud computing unchained us from looking after servers in cupboards.

Fold aims to do the same for cloud based applications. No one wants to spend
hours fiddling with yaml and cursing configuration errors. No, we want to
develop applications simply and locally, and then magically have them deployed
around the globe at any scale we like.

Fold is trying to build that, here are the principles behind the experience 
we're trying to create:
- Local development is a first class citizen. If it works locally it should work
  remotely.
- You shouldn't have to think about where your application is running or how
  it will get there.
- You shouldn't have to think about creating and managing infrastructure
  such as databases or message brokers.
- You should be able to write the application in the same way whether you are
  serving 10 or 10,000,000 customers.
- You shouldn't have to spend any time iterating on complex configuration.

# Installation

You just need one binary to run the command line tool, `foldctl`. Head over to 
the releases page to grab the latest release for your system. Simply download 
the binary and place it somewhere on your `$PATH`.

The only dependency for local development is docker, which `foldctl` uses to
build containers and run them for you locally. You can get docker from [the
docker website](https://docs.docker.com/get-docker/).

# Speed Run

```
foldctl init  # Fill in some project details at the prompts
foldctl new basic js hello-service
foldctl up hello-service/
```

This will set up a fold project and create a new service from the basic
javascript template. Once the `up` command has run succesfully it will be
available as an HTTP service on your machine. By default it will be running on: 

`localhost:8080/hello-service`

Note the `hello-service` in the URL. The gateway creates a path for every 
service, based on its name, and you must include that in the URL to contact 
the service you are interested in.

This is because fold supports running multiple services (as many as you'd like)
and you need to be able to identify the right one through the gateway. You can
try this out by creating a new service and starting it:

```
foldctl new basic js goodbye-service
foldctl up goodbye-service/
```

The new service will be available on:

```
`localhost:8080/hello-service`
```

When you're done run `foldct down` from the project root to bring down the
services.

# Working with Fold

## Project templates

You can find example projects over in the [templates repository](https://github.com/foldsh/templates).
These templates are also what `foldctl` uses to create new services. If you want
to use the template located at github.com/foldsh/templates/basic/js, you can use
the command:

```
foldct new basic js <service-name>
```

Feel free to submit your own via a PR! The intention over time is to create a 
repository of useful generic services which can be added to your project with a
single command.

For example, there may be an `auth-service`, removing the need for you to build
an integrate it yourself.

## Project Structure

A fold projet has a very simple structure. It consists of a project root, which
is identified by the `fold.yaml` file, and services.

Services are simply subdirectories within a fold project that container a valid
`Dockerfile` and which are registered in your `fold.yaml` file. The only caveat 
is that the `Dockerfile` must build an image which implements the fold runtime 
interface.

## Project Config

A key aim of the fold platform is to eliminate all of the endless yaml that
typically comes with looking after modern cloud based applications. That said, 
a small concession has been made in the form of the `fold.yaml` file which 
resides in your project root. 

It looks something like this:

```
email: test@test.com
maintainer: me
name: fold-project
repository: github.com/foo
services:
- name: hello-service
  path: ./hello-service
  mounts:
  - ./app
- name: goodbye-service
  path: ./goodbye-service
  mounts:
  - ./app
```

It is just used to configure a few bits of metadata and to register services.
Additionally, there are a few options on the services that allow you to
configure things for local development, for example which directories to mount
on your containers so you can hot reload your changes.

## Hot reloading

In the service config, there is a key called `mounts` which simply takes a list
of paths to mount to your running development containers. The paths are relative
to the service and will be mounted related to the `WORKDIR` in yoru container.

For example, if I have a service located at `./foo`, and mount `./bar/baz`, then
the directory `./foo/bar/baz` (relative to the project root) will be mounted at
`WORKDIR/bar/baz` in your development containers.

The fold runtime will watch for changes in all mounted directories during local 
development and reload the service if it detects any changes.

There are two key points to bear in mind if you are not familiar with how bind
mounts work on a docker container.

1. If the directory does not exist on your host then the container will fail 
to start. 
2. The directory on the host overwrites the directory on the container. This 
means that if the the directory is empty on the host (for example a build 
output directory) then it will also be empty in the container. Don't worry
though, if you end up here the just build the service locally and the fold
runtime will detect the changes and start up the service.

## SDKs

Currently there are two fold SDKs, one for NodeJS, supporting both TypeScript
and JavaScript, and another for go.

There is no documentation for the SDKs yet but they are still very small and
simple to use. The examples from the templates show off more or less everything
they can do right now.

# Deployment

You can't yet, coming soon!

# Fold Runtime

The fold runtime acts as the interface between your services, implemented using
the fold sdk, and the infrastructure which the fold platform has provisioned for
you.

In cloud design pattern parlance it is akin to an ambassador, however it is
embedded directly in your container.

All of your services must implement the runtime interface. The easiest way is to
use a provided base image but you can create your own very easily too.

## Base Images

A valid fold service is simply an OCI image that implements the fold runtime 
interface. In order to make things easy there are some base images that you can
extend, check out the available tags on [docker hub](https://hub.docker.com/r/foldsh/foldrt).

In order to use one of these, you just have to extend the image and set the 
`CMD` in your `Dockerfile` to be the command for invoking your application. It's
important that you don't override the `ENTRYPOINT`.

## Custom Images

It's very easy to create your own image if you need a bit more control. You just
need to get a copy of the `foldrt` binary from the releases page. It is pretty
small (around 12MB) and has no dependencies (not even libc), so it is very easy
to work with.

All you need to do is set your container up to run the `foldrt` binary with your
application command and arguments passed as the arguments to `foldrt`. For
example, if you normally run your app with the command `node ./dist/index.js`, you
need to run `foldrt node ./dist/index.js`.
