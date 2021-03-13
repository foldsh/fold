# Installation

You just need a binary for the command line tool, `foldctl`. Head over to 
the [releases page](https://github.com/foldsh/fold/releases/) to grab the latest 
release for your system. Simply download and extract the binary, make it 
executable and place it somewhere on your `$PATH`.

The only dependency for local development is docker, which `foldctl` uses to
build containers and run them for you locally. You can get docker from [the
docker website](https://docs.docker.com/get-docker/).

## Linux

You will need to use the `foldctl-$VERSION-linux-amd64.tar.gz` release.

For example:

```
tar xzvf foldctl-$VERSION-darwin-amd64.tar.gz
chmod +x ./foldctl-$VERSION-linux-amd64
sudo mv ./foldctl-$VERSION-linux-amd64 /usr/local/bin/foldctl
```

## Mac

You will need to use the `foldctl-$VERSION-darwin-amd64.tar.gz` release.

For example:

```
tar xzvf foldctl-$VERSION-darwin-amd64.tar.gz
chmod +x ./foldctl-$VERSION-darwin-amd64
sudo mv ./foldctl-$VERSION-darwin-amd64 /usr/local/bin/foldctl
```

## Build from Source

Alternatively you can easily build the tool from source if you have the go
compiler available. Simply check out the source at the tag you want and run:

`go install ./cmd/foldctl/`

# Speed Run

To get started you will need to use the `new` command, which is used to create
new fold projects and services. They will both guide you through a few set up
options.

The `up` command will then start your new service. You will have to change the
path depending on the name you chose for your service in the setup. The below
example assumes you called your service `new-service`

```
foldctl new project
foldctl new service
foldctl up new-service/
```

Once the `up` command has run succesfully it will be available as an HTTP 
service on your machine. The command should give you some information about 
how to reach your new service and which routes have been registered:

```
Fold gateway is available at http://localhost:6123

    new-service is available at http://localhost:6123/service-name
    new-service routes:
        GET http://localhost:6123/new-service/hello/:name
```

Note the `new-service` in the URL. The gateway creates a path for every 
service, based on its name, and you must include that in the URL to contact 
the service you are interested in.

This is because fold supports running multiple services (as many as you'd like)
and you need to be able to identify the right one through the gateway. You can
try this out by creating a new service and starting it:

```
foldctl new service
foldctl up second-service/
```

The new service will be available on:

`localhost:6123/second-service`

When you're done run `foldct down` from the project root to bring down the
services.
