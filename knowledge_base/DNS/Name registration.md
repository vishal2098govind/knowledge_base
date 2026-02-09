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

## Choosing TLD
We have a lot of options to choose from, among
- generic TLDs
- country-code TLDs
- infrastructure TLDs
- sponsored TLDs
- IDN TLDs
- geographic TLDs

There are more than 1500+ TLDs to choose from, as of June 2020s

Not all TLDs are suited for all types of scenarios
So there is definitely not a one size fits all

Choice of TLD is dictated by the specific use-case

Considerations when choosing a TLD
- Is there support for DNSSEC?
	- DNSSEC is a security mechanism that provide protection against certain attacks like DNS cache poisoning
	- TLDs such as `.ero`, `.pro` and `.travel` are only a few of the TLDs that do not support the security feature at this point in time
- IDNs support?
	- Is there support for internationalized domain names?
	- Domain names that include non-ascii characters like Arabic and Korean languages do
- Privacy protection?
	- To protect us from people finding our personal information in the registration records
	- Many country code TLDs such as `.cn`, `.us` and `.eu` do not have privacy
- Target audience?
	- If our organization is region specific, such as healthcare provider offering services only in Poland, using Poland's `.po` country-code TLD would make more sense than a generic TLD
- Relevant field?
	- If our website is going to advertise our cafe business, we might want to use an industry specific TLD such as `.cafe` or `.coffee`
- Local presence requirements

## Choosing a second-level domain
Consider `.com` being chosen as the TLD

Second level TLD can have any form we like, as long as the name we pick follows the rules of the domain name format

Considerations:
- Use keywords that reflect our industry
- localized keywords
- keep short within 10 characters
- easy to spell, pronounce and remember
- avoid hyphens, numbers, acronyms

Once we register say for instance `dnstips.com` as a domain name, we now own the entire `dnstips.com` namespace, which means that we would not need to register any subdomains that we might want to create later on.

We can also buy the same domain with multiple TLDs, such as 
- `dnstips.info`
- `dnstips.io`
- `dnstips.co.uk`
and point them all to `dnstips.com`

So if a user by mistake types another TLD other than `.com`, they will still be directed to our website

We might also want to by misspellings of the domain we originally purchase, in order to prevent others from registering them

### Typosquatting or URL hijacking
A common tactic among malicious actors is to create a fake version of our website with a deliberately misspelled domain name, that looks very similar to ours, so that if an end user makes a typo while typing our domain name, such as `dnsstrips.com`, which if matches the misspelled domain name the malicious actor has set up, the end user will then be directed to their website instead of yours.
https://dnstwist.it helps to find any phishing domain scanners

## Choosing a Registrar
Considerations:
- Pricing
	- Registration fee
	- renewal fee
	- other charges like domain-transfers
	- cost benefits like
		- bulk pricing options
		- promotion deals
- add on services
	- web hosting
	- website builders
	- email hosting
	- domain privacy
	- SSL protection
- Supported TLDs
	- not all registrars have license to sell all top level domains
- Policies
	- domain transfer policy
	- domain expiration policy
		- grace period allowing renew of expired domain name
- Recommended to register our domain name on a registrar, but host on another provider, so that it is easier to switch hosting companies if required later on, provided that domain name is hosted on a platform different to the one it was registered on

## EPP Status Code
EPP - Extensible Provisioning Protocol domain status codes
Also called domain name status codes, indicate the status of domain name registration

Every domain has at least one status code, if not more

Each EPP code provides useful information about a domain that comes in handy for operations such as
- troubleshooting domain-related issues
- domain renewals
- domain transfers between registrars

There are two different types of EPP status codes:
- client codes
- server codes

Client Status Codes:
- Set by **registrars**, either while requesting a domain or after registering domain
- Examples:
	- clientHold
		- enacted during legal disputes
		- non-payment
		- when domain is subject to deletion
	- clientTransferProhibited
		- tells registry to reject requests to transfer domain from current registrar to another
	- clientUpdateProhibited
		- tells registry to reject requests to update domain

Server Status Codes:
- Server status codes, are set by **registries** and they take precedence over client code
- Examples:
	- ok
		- standard status for a domain, i.e. it has no pending operations or prohibition
	- autoRenewPeriod
		- grace period provided after domain name registration period expires, extended or renewed automatically by the registry
		- if the registrar deletes the domain name during this period, the registry provides a credit to the registrar for cost of the renewal
	- serverTransferProhibited
		- prevents our domain from being transferred from our current registrar
		- It is an uncommon status that is usually enacted during legal or other disputes at our request

We can inspect EPP codes of any domain we are interested in, simply by performing a WHOIS lookup against it

