# Working with Fold

## Project Templates

You can find example projects over in the [templates repository](https://github.com/foldsh/templates). These templates are also what `foldctl` uses to create new services. The `foldctl new service` command will present you with options based on what is available there.

Feel free to submit your own via a PR! The intention over time is to create a repository of useful generic services which can be added to your project with a single command.

For example, there may be an `auth-service`, removing the need for you to build an integrate it yourself.

## Project Structure

A fold project has a very simple structure. It consists of a project root, which is identified by the `fold.yaml` file, and services.

Services are simply subdirectories within a fold project that container a valid `Dockerfile` and which are registered in your `fold.yaml` file. The only caveat is that the `Dockerfile` must build an image which implements the fold runtime interface.

## Project Config

A key aim of the fold platform is to eliminate all of the endless yaml that typically comes with looking after modern cloud based applications. That said, a small concession has been made in the form of the `fold.yaml` file which resides in your project root.

It looks something like this:

```text
email: test@test.com
maintainer: me
name: fold-project
repository: github.com/foo
services:
- name: user-service
  path: ./user-service
  mounts:
  - ./app
- name: score-service
  path: ./score-service
  mounts:
  - ./app
```

It is just used to configure a few bits of metadata and to register services. Additionally, there are a few options on the services that allow you to configure things for local development, for example which directories to mount on your containers so you can hot reload your changes.

## Hot Reloading

In the service config, there is a key called `mounts` which simply takes a list of paths to mount to your running development containers. The paths are relative to the service and will be mounted related to the `WORKDIR` in yoru container.

For example, if I have a service located at `./foo`, and mount `./bar/baz`, then the directory `./foo/bar/baz` \(relative to the project root\) will be mounted at `WORKDIR/bar/baz` in your development containers.

The fold runtime will watch for changes in all mounted directories during local development and reload the service if it detects any changes.

There are two key points to bear in mind if you are not familiar with how bind mounts work on a docker container.

1. If the directory does not exist on your host then the container will fail 

   to start. 

2. The directory on the host overwrites the directory on the container. This 

   means that if the the directory is empty on the host \(for example a build 

   output directory\) then it will also be empty in the container. Don't worry

   though, if you end up here the just build the service locally and the fold

   runtime will detect the changes and start up the service.

