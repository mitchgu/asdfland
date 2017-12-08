# asdfland URL Shortener

## The only URL shortener you'll ever need.

* **3 link types** for 3 use cases
  * **Random** (like bit.ly)
    * Generates a **short, completely random** link (high entropy) (Ex: `asdf.land/wHRbxt`)
    * Good for **maximal shortness**, or making a long link that is astronomically **hard to guess** (effectively private)
    * Best for **sending electronically** in messages, email, documents to reduce clutter
  * **Readable** (like shoutkey.com)
    * Generates a **readable** link composed of common words (Ex: `asdf.land/ankle`)
    * Good for **sharing verbally** for an audience or nearby collaborators
  * **Custom** (like bit.ly, tinyurl.com)
    * choose a **custom** link (Ex: `asdf.land/XYZ.party.signup`)
    * Good for **memorable links** to put in publicity or to reinforce branding
    * Also good for creating a **personal library** of bookmarks that are named memorably (`asdf.land/my_resume`, `asdf.land/catgifs`)
* **Endlessly parametrizable** links. Change everything about your link to your liking: length, word dictionary, expiration, custom slug, password protection, analytics
* **Live preview** links before creating them
* **Manage links** in the dashboard (no account required)
  * **Edit** any existing link's slug or destination
  * **Associate** multiple shortlinks with one destination
  * **Annotate** links with descriptions
  * **Track** link visits over time
* **Open source** and **easily self-hostable** on a custom domain
* **Ludicrously fast**: powered by **Go** and **Redis**

## Quick Install
First, **make sure you either have Redis installed locally or running on a remote server**. The [Redis Quick Start Guide](https://redis.io/topics/quickstart) has great instructions for installing Redis.

1. Download and extract the [latest prebuilt release](https://github.com/mitchgu/asdfland/releases/latest) (`asdfland_<version>_<os>_<arch>.tar.gz`) (`.zip` for Windows)
3. Run `asdfland` (`asdfland.exe` for Windows). The default server address is `localhost:9090`.

## Configuration
By default, `asdfland` will look for Redis running locally on the default port (`localhost:6379`) and use database `7`. To configure these settings and many others, create a file called `config.yaml` in one of three places:

* The working directory where `asdfland` will be run
* `$HOME/.config/asdfland/config.yaml`
* `/etc/asdfland/config.yaml`

The `config_sample.yaml` file in the root of this repository has all available settings and is recommended to use as a template. 

## Building from source
This (obviously) requires Go to be installed. [Here](https://golang.org/doc/install) are instructions for that

1. `go get github.com/mitchgu/asdfland`
2. `cd $GOPATH/src/github.com/mitchgu/asdfland` (this is the repo root)
3.  From here you can explore the code or run the project several ways:
	* Run current source: `go run *.go`
	* Compile to a binary: `go build` (creates `asdfland` binary)
	* Install to `$GOBIN`: `go install`

## Editing the frontend

The go code in this repository runs server-side to expose an API to any compatible frontend. The default frontend built on Vue.js for asdfland lives in a separate repository: [mitchgu/asdfland\_vue\_client](https://github.com/mitchgu/asdfland_vue_client). To provide a default frontend in the go project, the Vue.js frontend is compiled into a few production-ready files, which are embedded into this repository in the `bindata.go` file using the [go-bindata](https://github.com/jteeuwen/go-bindata) tool.

Currently, the best way to customize the frontend is to clone both the backend (this repo) and the frontend repo, compile the frontend, and re-embed it into `bindata.go`. There are plans in the future to allow specification of an alternate filesystem directory for the go backend to serve the frontend then.

To link the frontend to the go repo for development purposes:

1. Clone the repo in another directory: `git clone https://github.com/mitchgu/asdfland_vue_client.git`
2. `npm install` to install dependencies (requires `node.js` and `npm`)
3. `npm run build` to compile the frontend to the `dist` folder
4. Symlink the frontend repo `dist` folder to a folder named `frontend` in the backend go repo (this repo)
		
		ln -s <frontend repo path>/dist <backend repo path>/frontend
4. Make sure `go-bindata` is installed: `go get -u github.com/jteeuwen/go-bindata/...`	
5. Run `go-bindata wordlists/ frontend/...` to embed the frontend and wordlists into the `bindata.go` file. During development, it's recommended to add the `-debug` flag so the go app will directly read from the `frontend` (symlinked `dist`) folder instead of actually copying the data into `bindata.go`. Before committing or releasing, rerun the command without the `-debug` flag to copy the data back in.

Now you are free to make edits to the frontend repo. Every time the frontend is recompiled to the `dist` folder with `npm build`, the go backend will reflect the new changes. 

The frontend repo also comes with a development server that is able to point to a running backend (local or remote) for live development. See that repo's README for instructions.