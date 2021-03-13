# Fold

Fold is a platform and a set of tools for developing and deploying backends. 
The goal is to make running a set of microservices on the cloud as easy as
running them locally.

**WARNING: This is seriously alpha...**

## Overview

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

