#dns

The context in which the domain names exist

## Root Servers
The first recipient of a DNS query request is one of the root servers

Root servers sit at the top of the hierarchical DNS tree

Upon receiving a query to resolve a specific name, the root service job is to pass back a referral for a top level domain name server based on a TLD extension of the domain name request

For example, if we browse to `www.yahoo.com`, DNS request would be received by a root server which will pass back a referral with a `com` name server, so that the query is sent there for further processing

Thanks to the information received by the root server, we can now proceed with the resolution process.

Root servers are the first component of the DNS tree that receive a DNS query

They contribute to the resolution process by providing a referral to a top level domain name server that corresponds to the top level domain extension of the domain name being requested

There are 13 root name servers, each of which is named with letters A to M. More information about the root name servers can be found at https://root-servers.org

Because of the important role they play in DNS resolution process, root servers are one of the most critical infrastructures in the internet

Although they number 13, there are multiple copies of each one at over 130 locations all over the world, hosted in multiple secure sites with high bandwidth access to handle traffic load, with 938 root servers being in operation world wide as of Feb 28th 2019.

## TLD Servers
Thanks to the root server's referral, it is now up to the referred TLD name server to take on the mantle of name resolution and take us one step closer to having a DNS query resolved

Responsibility of the TLD name server is simply to provide a referral of it's own that will contain a server authoritative for the domain we want to browse to.

Root name servers don't know about server hosting `www.yahoo.com`. It only knows about the TLD server location

The referred TLD name server doesn't know about the host or system hosting `www.yahoo.com` either. It does know about the name server responsible for `yahoo` part of the domain

Thus, level of hierarchy or DNS tree provides a piece of the puzzle needed for the requested domain name to be resolved

Most operational aspects of DNS are hierarchical. From the way domain names are formatted or structured in DNS tree, down the way they are translated when browsing.

The main job of the TLD name server is to pass back a referral with a name server that has authority over the requested domain.

There are several categories of TLDs, of which the most important ones to know are 
- generic TLDs (gTLDs)
	- they are supposed to reflect the type of industry or space that organization in ownership of that domain, falls into
	- for example, it is standard practice for governmental agencies to use the `.gov` gTLD like www.usa.gov
	- For the universities to use the .edu gTLD as
	- The origin of gTLDs were 
		- `arpa`
		- `com`
		- `edu`
		- `gov`
		- `mil`
		- `net`
		- `org`
	- The list has grown bigger with additional gTLDs, such as
		- `aero`
		- `int`
		- `pro`
		- `info`
- country-code TLDs (ccTLDs)
	- have been created to allow countries to manage their own namespace
	- and make it easier for their citizens to find the resources they need in the languages they understand
	- example:
		- `gr` - Greece
		- `jp` - Japan
		- `pl` - Poland
	- Countries can also use organizational sub-domains within their country level top-level domains
		- example: \*.com.au used by Australia
- sTLDs
- IDNTLDs
- geoTLDs

Full list of TLDs can be found from https://iana.org/domains/root/db


## Authoritative Name servers
They are the most abstract component of the DNS tree
Hardest one to define

### Concept of Authority in DNS protocol
The DNS namespace is based on a hierarchy of domains
Linguistically, the term domain refers to a sphere of influence governed by an entity that exercises authority over that domain and is tasked with certain responsibilities
Nowadays, we often hear people of various professions saying that this is outside my domain, when asked questions not related to their trade, indicating that the requests made to them are outside of their scope of knowledge

Conversely, a well-known expert who is very good in their industry is referred to as the leading authority in that field, and so that person is responsible for educating rest of us on their specific subject matter

Along the same lines, the domain within the context of DNS is a particular slice of namespace where there is somebody with the authority to manage that segment and with the responsibility to provide, and with the responsibility to provide answers to DNS requests made for that very segment
![[Pasted image 20260209011837.png]]
DNS uses a globally distributed system of databases, and that system works in conjunction with an equally distributed system of authorities

This hierarchical authority structure complements the hierarchical name structure in DNS

It is not necessary for a different authority to exist at every level of hierarchy, as in many cases, a single authority may manage a section of namespace that spans more than one level of the structure