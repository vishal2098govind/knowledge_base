#dns

DNS is globally distributed collection of databases.
This principle extends to DNS storage as well.
Since DNS data is stored in a database, that is known as a **zone**

## DNS Zone
### Forward Lookup Zone
- Used for resolving names to IP addresses
- Zone where most of the DNS data is stored
### Reverse Lookup Zone
- stores information used for reverse name resolution

Each of these zones is a collection of what we call **Resource Records**, commonly abbreviated as `RRs`

These `RRs` are sets of fields organized in rows.
There are many types of resource records in DNS, and each resource record type contains a specific set of data

### `A` record
contains a domain-name and it's associated IPv4 address

#### `AAAA` record
contains domain-name along with its corresponding IPv6 address

Despite their differences, all resource record types share a common resource record format

### Common Resource Record Format
- Name of record
- Type of the record
	- each resource record type has a different value
	- `A` record, for example, has a type value of `1`
	- while, `NS` record has a type value of `2`
- Class of the record
	- can be among `IN`/`CH`/`HS`
	- in 99% of cases, it should be `IN` standing for internet
- TTL
	- length of time the resource record remains cached for
- RDLength
	- Resource-Data length in bytes
- RData
	- actual data that the resource record stores


## SOA - Start of Authority Record
- Indicates the beginning of a zone
- It should be the first resource record to be specified in any zone file
- There can be only one start-of-authority record per zone
- Format:
```
<domain-name> <TTL> <class> SOA <m-name> <r-name> (
	<serial-number>
	<refresh-interval>
	<retry-interval>
	<expire-interval>
	<minimum>
)
```
- **m**-name: name of pri**m**ary authoritative name server for that zone
- **r**-name: email address of the administrator responsible for the zone
	- although the value of r-name represents an email address, there is no @ sign in it
	- i.e. `info.example.com` is same thing as `info@example.com`
- **\<serial-number>**
	- version number of zone
	- incremented by one every time a change is made into the zone
	- it should never be decreased
- **\<refresh-number>**, **\<retry-number>**, **\<expire-number>**
	- retry interval and expire interval have to do with primary and secondary name servers
- **\<minimum>**
	- historical purpose of this parameter was to hold the default TTL value for the records without explicit TTL value defined
	- nowadays, the minimum parameter represents the TTL value for negative caching

To query SOA record, we can use tools like `dig` or `nslookup`

```sh
$ nslookup -type=soa example.com
Server:		8.8.8.8
Address:	8.8.8.8#53

Non-authoritative answer:
example.com
	origin = elliott.ns.cloudflare.com
	mail addr = dns.cloudflare.com
	serial = 2395194487
	refresh = 10000
	retry = 2400
	expire = 604800
	minimum = 1800
```

## NS Record - Name Server
- NS : Name Server
- One of the most important resource records
- Point to authoritative name servers for a zone
- it is these name servers that hold the actual DNS information for a domain, so that, that domain can be accessible to the internet users
- If a domain such as example.com doesn't have any NS records configured in it's zone, there wouldn't be configured in it's zone, there wouldn't be any references to that domain's name servers, which means the `.com` TLD server contacted would not be able to return a list of name servers for `example.com`, which in turn means nobody will be able to browse to that domain
- It is the NS records that ensure the availability of the domain
- Because of that reason, every zone must have **at least** two NS records, of which each points to a different authoritative name server for redundancy
- This way, if one name server goes down or becomes unavailable, DNS queries can go to a different name server
- Name servers also typically reside in topologically separate networks for further resiliency
- while registering a domain, many DNS as a service providers by default provision a set of four name servers and they would be listed like so
	- ns1.example.com
	- ns2.example.com
	- ns3.example.com
	- ns4.example.com
- Format:
```
<domain-name> <TTL> <class> NS <nameserver's hostname>
```
example:
```sh
$ nslookup -type=ns knolia.ai
Server:		8.8.8.8
Address:	8.8.8.8#53

Non-authoritative answer:
knolia.ai	nameserver = ns63.domaincontrol.com.
knolia.ai	nameserver = ns64.domaincontrol.com.
```

## `A` record - Address Record
The resource record that stores the association between domain name and IPv4 address is the A-record or the address record

It is a primary record in DNS and is queried in forward lookup requests

Format:
```
<domain-name> <TTL> <class> A <IPv4 address>
```
### AAAA Record
The cousin of the A record is the AAAA record - aka Quad-A record, which stores a domain name, and IPv6 address

Reason for four As in AAAA-Record:
to signify that the value stored in it is four times as big as the one stored in A record
i.e. IPv4 has 32 bits, IPv6 has 128 bits

Format
```
<domain-name> <TTL> <class> AAAA <IPv6 address>
```

It is perfectly possible for an A record and a AAAA record to point to same domain name, in cases where dual stack is required

```sh
$ nslookup -type=A example.com
Server:		8.8.8.8
Address:	8.8.8.8#53

Non-authoritative answer:
Name:	example.com
Address: 104.18.27.120
Name:	example.com
Address: 104.18.26.120

$ nslookup -type=AAAA example.com
Server:		8.8.8.8
Address:	8.8.8.8#53

Non-authoritative answer:
example.com	has AAAA address 2606:4700:9a92:8c7c:33c5:0:ccc4:a209
```

## PTR Record - Pointer Record
- This resource record allows the reverse name resolution to happen
- Points an IP address to domain name

Format:
```
<reverse domain name> <class> PTR <domain name>
```

They are placed in reverse lookup zone

```sh
$ nslookup -type=ptr 52.214.181.141
Server:		8.8.8.8
Address:	8.8.8.8#53

Non-authoritative answer:
141.181.214.52.in-addr.arpa	name = ec2-52-214-181-141.eu-west-1.compute.amazonaws.com.
```

## Canonical Name Record CNAME
This appears very frequently in DNS configuration, troubleshooting and name resolution
CNAME = Canonical Name = real name of an object referenced by an alias
This record basically maps one domain to another

Format

```
<alias> <class> CNAME <TTL> <canonical name>
www.example.com CNAME example.com
```

### Use cases of CNAME records:
- Used to map subdomains to apex domains, such as www.example.com to example.com
- Redirect multiple TLDs to same second-level domain
- Example, if we have mydomain.com and mydomain.com.nz they get directed to mydomain.com
- Validate ownership or control of a domain

### Restrictions on CNAME
- Must always point to another domain name and never directly to an IP address
- Cannot point to an NS or MX record
- CNAME record cannot co-exist with another record for the same name
	- i.e. not possible to have both CNAME and TXT record for www.example.com
- CNAME can point to another CNAME, a mechanism that is known as CNAME chaining. Not recommended as it requires multiple DNS lookups before intended domain can be loaded, which slows down the name resolution process and in turn impacts user experience