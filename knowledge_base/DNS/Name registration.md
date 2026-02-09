#dns

## Domain Name Registration Hierarchy
Name registration function of DNS involves the registration of domains

A domain represents a public identity on the internet, and is used to identify IP address of the computer system hosting the web content

To ensure that each domain is unique, name registration has to be processed within a globally distributed framework, called domain name registration hierarchy, designed to enforce a certain set of rules

### ICANN
At the top of the hierarchy, rules ICANN, a non-profit, internationally organized corporation that stands for **Internet-Corporation-for-Assigned-Names-And-Numbers**

It's main role is to oversee the huge complex interconnected network of unique identifiers that allow computers on the internet to find one another

ICANN is responsible for 
- managing generic TLDs and country-code TLDs
- managing how root name server systems function
- coordinating how IP addresses are supplied to avoid repetition or clashes
- maintaining a central repository of IP addresses

![[Pasted image 20260210003853.png]]
#### Registry
Below ICANN, are the 5 regional internet registries.
Each registry is responsible for obtaining IP ranges from ICANN to allocate them to internet service providers across a specific geographic region
#### Registrar
Subordinate registries are the registrars, which are ICANN accredited organizations responsible for processing registration of domain names.

Examples of registrars are 
- GoDaddy
- Namecheap
- Bluehost
#### Reseller
After registrars, come the resellers, which are third party companies that offer domain name registration services through registrars

Example:
- Route53 - The dedicated DNS and registration as-as-service provided by AWS
	- Route53 is the reseller of two registrars
		- Amazon Registrar for generic TLDs
		- Gandi for all other top level domains

#### Registrant
At the bottom of the hierarchy are the registrants, 
the people or organizations who register domains through a registrar or a reseller

## Domain Name Registration Process

- Registrant chooses a domain name and submits a request to register it with a **reseller** or ICANN accredited **registrar**
- Provided that the domain name is available, the registrar registers the name and then it creates a **WHOIS record**
- The WHOIS record contains 
	- registrant's name and contact information, 
	- registration date 
	- the name servers
	- the most recent update
	- expiration date
- WHOIS records may also provide administrative and technical contact information of the registrant
- The registrar will then send domain name request, along with contact and technical information of the domain name to the appropriate registry
- the registry will file all information provided and it will add domain zone file to master servers, which will tell other servers on the internet where our website is located