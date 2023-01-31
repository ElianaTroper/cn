# Server

This repo holds the server source code, and also acts as the default save location for a running server instance.

## Running a "server"

### Prerequisites
- you have `docker-compose` installed and enabled
- The following ports are not in use by another program:
	- 4001 (+ opened to the web using the same number)
	- 4002 (+ opened to the web using the same number)
	- 5001 (NOT opened to the web)
- You do not have docker containers named:
	- `ipfs`
	- `ipfs-js`
- The binary of the main program is in-place in the repo as cloned

(If I develop this program further, I will remove these restrictions. In fact, some are already adjustable by modifying 1-2 `config`/`env` files, but I am not going to provide support should you make modifications at the moment ☺)

### Spin up a server
- from this directory, `cd ipfs && docker-compose up && cd ..`
- from this directory, `./cn ipfs init`
- from this directory, `./cn start`
	- Congrats! You are now tracking and providing the root of the system.
- to see more functionality, run `./cn --help` and explore subcommands
- you can stop a running node by returning to this directory and entering `./cn stop`

#### Running apps
To add support to your node for an app, run `./cn app enable APPNAME`. Built in apps are benchmark, post, and library.

##### post

Runs a basic message board which displays messages sent by users.

NEXT: Add more details

##### library

Runs a mirror of project gutenberg.

NEXT: Add more details, link project gutenberg

##### benchmark

Allows your node to serve benchmark tests to clients. **WARNING:** This mode is currently unauthenticated, does not have any resource limitations, and relies on trusting clients to not exhaust your machines resources as a result.

NEXT: Add more details
