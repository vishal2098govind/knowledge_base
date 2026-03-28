#dns

Will be inspecting closely the DNS messages
The core function of name resolution is carried out by two different messages:
- The query - sent by the client
- The response - sent by the server

Each exchange of these DNS messages is referred to as a transaction, and is assigned a transaction identifier Transaction-ID, which is useful in identifying which DNS response corresponds to which DNS query request.

Can be useful
- when trouble shooting DNS issues
- when conducting a forensic network analysis

DNS is a layer-7 application protocol
thus, DNS has to rely on a transport protocol at layer-4 to transfer DNS messages from client to server and vice versa

Most network protocols residing at the application layer usually rely on either TCP or UDP, 

but DNS is interestingly one of the few protocols that use both, depending on the operation taking place

For name resolution, where speed is of the essence, UDP is used as a transport protocol to carry both requests and responses

Since UDP doesn't offer a reliable method of delivering messages, it is up to the client to keep track of the requests sent so that if a response is not received at a specific time interval, the corresponding request can be retransmitted

In the interest of preventing excessive DNS traffic on the network, retransmissions are usually sent at an interval ranging from 2 to 5 seconds

TCP, on the other hand, is used when data has to be delivered reliably, as is the case in scenarios such as 
- zone transfers
- or when DNS response requires more than 512 bytes of space, since all UDP DNS messages are limited to a payload of 512 bytes, though if a response message is larger than that, the message is truncated, and a special bit in the header is set to indicate this outcome
	- In such an event, the client would then be expected to retry the query over TCP instead, which doesn't have the same size limitation as UDP


Generally, transactions taking place over TCP remain a very small fraction of overall DNS, regardless of which transport protocol is used, the standard port that the server listens on by default is 53, while the port used by the client will always be an ephemeral one


Other types of DNS messages outside of the area of name resolution are
- notify 
	- allows master servers to inform slave servers when the zone changed. 
	- carried over UDP
- update
	- add or delete resource records from a specified zone
	- carried over UDP or TCP, depending on size of the request



### Response Codes (RCODES)
- no error R code is the only response code that indicates success
- 0 : NoError
- 1 : FormErr
- 2 : ServFail
- 3 : NXDomain Non-Existent Domain
- 4 : NotImp Not Implemented 