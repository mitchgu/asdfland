# asdf.land
This is the go and redis server backend for asdf.land, a url shortener.

I'm doing this because I want a self-hostable mashup of bit.ly and shoutkey with more link manipulation options. I also want to learn more go and redis. 

I like the idea of go and redis for this because:

* A single executable would be really easy to deploy anywhere
* The datastore doesn't need to be super relational. A link shortener is basically just a key value store
* Go and redis should be pretty blazing fast


## Types of links
* random, not expiring
* shoutkey style, expiring
* fully custom, with prefixes for each user