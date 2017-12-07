# asdf.land URL Shortener - Go Backend
## Introduction

asdf.land is a URL shortener that brings many new innovations. Learn about them below, or try them out now on the main deployment: [https://asdf.land](https://asdf.land)

* **3 different link types** for 3 different use cases
	* **Random** (like bit.ly)
		* Generates a short, completely random link (high entropy) (Ex: `asdf.land/wHRbxt`)
		* Good for maximal shortness, or making a long link that is astronomically hard to guess (effectively private)
		* Best for sending electronically in messages, email, documents to reduce clutter
	* **Readable** (like shoutkey.com)
		* Generates a readable link composed of common words (Ex: `asdf.land/ankle`)
		* Great for sharing verbally for an audience or nearby collaborators
	* **Custom** (like bit.ly, tinyurl.com)
		* Completely custom shortlink (Ex: `asdf.land/XYZ.party.signup`)
		* Great for memorable links to put in publicity or to reinforce branding
		* Also good for creating a personal library of bookmarks that are named memorably (`asdf.land/my_resume`, `asdf.land/catgifs`)

* Many options for link **customization** (length, word dictionary, password)
* **Live preview** the link before creating it
* **Manage links** in a dashboard (no account required at first)
	* **Edit** the shortlink or destination
	* **Associate** multiple shortlinks with one destination
	* **Annotate** links with descriptions
	* **Track** link visits over time
* **Open source** and **easily self-hostable** on a custom domain
* **Ludicrously fast**; powered by **go** and **redis**

## Deploying
asdf.land comes in two parts: the backend (this repository) and the [main Vue.JS frontend](https://github.com/mitchgu/asdfland_vue_client). Other frontends that adhere to the same API can also be used.

### Build Environment
Building asdf.land requires **Go**, **Redis**, and **node.js/npm**. Each can be installed fairly easily by following their official instructions

### Deploy

1. Find a suitable place to clone the files (Ex: `/srv/asdf`)
1. Clone both repositories:

		git clone https://github.com/mitchgu/asdfland
		git clone https://github.com/mitchgu/asdfland_vue_client
		
1. Build the frontend:
	1. `cd asdfland_vue_client`
	2. `npm install`
	3. `npm run build`
4. Configure the backend:
	1. Copy `asdfland/config_sample.yaml` to `asdfland/config.yaml` and open it in an editor of your choice
	2. Configure the backend:
		1. `server_addr`: The port to serve asdf.land at (Ex: `localhost:9090`)
		2. `db_kind`: Only supports `redis` currently
		3. `redis_addr`: Address and port of redis (Ex: `localhost:6379`)
		4. `redis_dbnum`: Database number to use for redis (Ex: `7`)
		5. `redis_pass`: Redis password, empty if no password
		6. `frontend_dir`: The path of the frontend directory (Ex: `/srv/asdf/asdfand_vue_client/dist`)
	7. Run the backend: `./asdfland &`
	8. Recommended: Use Nginx or Apache virtual hosts to proxy traffic from the desired domain to the go server. 