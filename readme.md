# beacon

Discovery service for designed for nsproxy but with global use.  This project is mainly taken from
[bitnuke](https://bitnuke.io), with code from
[here](https://git.unixvoid.com/mfaltys/bitnuke)  
Beacon allows for host discovery by uploading a hosts ip to the beacon server.
Once an ip is registered in beacon, clients are able to pull the ip globally.
The beacon api is available for public use globally at
https://beacon.unixvoid.com.  This service exposes an api where clients can
request a beacon id and their clients can fetch the data.  
A general usage flow is as follows:  
- A user requests a new id with `/provision`
- A beacon client updates the ip with `/update`
- Any client can fetch this ip with `/<beacon id>`  
The service can be run bare on the host or in a docker container
`unixvoid/beacon`

### api
- beacon exposes an api for provisioning id's, updating ip entries, viewing
  entries, and deleting id's.  The following  is the specification for endpoints
  and their protocols.
- `/<beacon id>` : `GET` : endpoint for getting a registered ip
  - example: `curl https://beacon.unixvoid.com/unixvoid`
- `/provision` : `POST` : endpoint for requesting a new beacon id
  - `id` : the intended ip
  - example: `curl -d id=unixvoid https://beacon.unixvoid.com/provision`
  - returns: `200` : client sec, a alphanumeric string for authorizing/removing
    entries
  - returns: `400` : the client id is already in use
- `/update` : `POST` : endpoint for updating client ip
  - `id` : registered beacon client id
  - `sec` : alphanumeric secret associated with registered beacon id
  - `address` : ip address to be updated to
  - example: `curl -d ip=unixvoid -d sec=yQHfXWrUMVDNaHoSkDhRhqG26 -d
    address=127.0.0.1 https://beacon.unixvoid.com/update`
  - returns: `200` : ip updated successfully
  - returns: `403` : client auth invalid
  - returns: `400` : client id does not exist
- `/remove` : `POST` : endpoint for remove a registered beacon id and
  associating metadata
  - `id` : registered beacon client id
  - `sec` : alphanumeric secret associated with registered beacon id
  - example: `curl -d ip=unixvoid -d sec=yQHfXWrUMVDNaHoSkDhRhqG26 https://beacon.unixvoid.com/remove`
  - returns: `200` : id and metadata removed successfully
  - returns: `403` : client auth invalid
  - returns: `400` : client id does not exist

### configuration
- beacon uses gcfg (INI-style config files for go structs).  The config uses some pretty sane defauls but the following fields are configurable:  
- `[beacon]`
  - `port:`  the port the API listens on.
  - `tokensize:` the length of authorization string returned from the client
  - `tokendictionary:` the allowed characters that make up the client
    authentication string
- `[redis]`
  - `host:`  this is the ip and port that the redis backend is running on  
  - `password:`  password to the redis database
