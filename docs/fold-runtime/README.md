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
