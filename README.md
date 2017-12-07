# asdf.land URL Shortener - Go Backend

## The only URL shortener you'll ever need. Here's why:

### **3 link types** for 3 use cases

#### **Random** (like bit.ly)
  * Generates a **short, completely random** link (high entropy) (Ex: `asdf.land/wHRbxt`)
  * Good for **maximal shortness**, or making a long link that is astronomically **hard to guess** (effectively private)
  * Best for **sending electronically** in messages, email, documents to reduce clutter
#### **Readable** (like shoutkey.com)
  * Generates a **readable** link composed of common words (Ex: `asdf.land/ankle`)
  * Great for **sharing verbally** for an audience or nearby collaborators
#### **Custom** (like bit.ly, tinyurl.com)
  * choose a **custom** link (Ex: `asdf.land/XYZ.party.signup`)
  * Great for **memorable links** to put in publicity or to reinforce branding
  * Also good for creating a **personal library** of bookmarks that are named memorably (`asdf.land/my_resume`, `asdf.land/catgifs`)

### Endlessly parametrizable links
Change everything about your link to your liking
* length
* word dictionary
* expiration
* custom slug
* password protection
* analytics

### **Live preview** links before creating them
The link you see is the link you'll get. Tweak it until you like how it looks.

### **Manage links** in the dashboard (no account required)
  * **Edit** any existing link's slug or destination
  * **Associate** multiple shortlinks with one destination
  * **Annotate** links with descriptions
  * **Track** link visits over time
### **Open source** and **easily self-hostable** on a custom domain
  * Go backend: [https://github.com/mitchgu/asdfland](https://github.com/mitchgu/asdfland)
  * Vue.js frontend: [https://github.com/mitchg/asdfland_vue_client](https://github.com/mitchgu/asdfland_vue_client)
### **Ludicrously fast**: powered by **Go** and **Redis**

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
	* `cd asdfland_vue_client`
	* `npm install`
	* `npm run build`
1. Configure the backend:
	* Copy `asdfland/config_sample.yaml` to `asdfland/config.yaml` and open it in an editor of your choice
	* Configure the backend:
		* `server_addr`: The port to serve asdf.land at (Ex: `localhost:9090`)
		*`db_kind`: Only supports `redis` currently
		* `redis_addr`: Address and port of redis (Ex: `localhost:6379`)
		* `redis_dbnum`: Database number to use for redis (Ex: `7`)
		* `redis_pass`: Redis password, empty if no password
		* `frontend_dir`: The path of the frontend directory (Ex: `/srv/asdf/asdfand_vue_client/dist`)
	* Run the backend: `./asdfland &`
	* Recommended: Use Nginx or Apache virtual hosts to proxy traffic from the desired domain to the go server. 