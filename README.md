# General

This is a server rendered application that provides url short features similar with bit.ly and other 
internet services, with the following differences:

- Can be installed on one's server and administrate the shortened urls the way the owner wants
- Any operation requires valid authentication
- All the endpoints require authentication, except /s url that is used solely to retrieve the real url from the shortened one, and does so with a Location header
