#dns

## Domain Name Format

Contrary to the way we type a domain name into our browser, from left to right, like `www.example.com` , the names are interpreted by DNS the other way around, from right to left

### Root
The first component of a domain name, starting from right to left, is the root represented by a **dot**

This dot is a very special dot and not to be confused with any other dots.

This dot is in fact so dominant that nothing else can take it's place. i.e. the right most component of a domain name can only be the root represented by this dot and nothing else

Root can be thought of as a component in the domain name, just so we don't confuse the special dot with other dots being present in remaining part of the name

### TLD - Top level domain
e.g. : `com` is just one of a plethora of TLDs that we could come across
other examples could be `net`, `org`, `io`

### SLD - Second Level Domain
e.g. `example`

### Zone Apex or Naked Domain
the SLD along with TLD together

### Third-Level Domain
usually something like `www`
it's more of a naming convention rather than a necessity

### Labels:
Components of a domain name are known as labels, with a root component being said to have a label of null

Each label is a subdomain of it's parent domain
e.g.:  
www is a subdomain of example.com
while, example is a subdomain of the com top-level domain

### FQDN and PQDN
The entire domain name, stretching all the way from root, down to the third level domain, is known as the **FULLY QUALIFIED DOMAIN NAME**, abbreviated as FQDN, which is the absolute reference to a domain name i.e. `www.example.com.`

FQDN is not the same thing as a URL
as a URL contains the domain name of a site, as well as other information, including the transfer protocol and the path.

Example of a URL: `https://www.example.com`
Example of FQDN: `www.example.com`

A portion of a domain name, such as the host portion or the `www.example.com` domain is known as a partially qualified domain name, abbreviated as **PQDN**

A partially qualified domain name starts with a hostname, but it may not reach up to the root

### Dot representing root is different than other dots
The root dot is a special dot, while the other dots serve as delimiters separating the labels of the domain

