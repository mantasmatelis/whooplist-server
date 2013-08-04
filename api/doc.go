/*
Whooplist Api Server

Listens on /tmp/whooplist.socket for HTTP requests. (meant to be reverse-proxied by nginx or similar)

Integrates with:

Login APIs
 * Facebook
 * Google
  - API Access Token is passed from client to server
Note: with Login APIs PUT /users (i.e. CreateUser) is never used, an account is created on the fly on /login if a login API isused.

Place API
 * Google Places
  - A reference id is used to communicate place information between client and server. Read the Places API spec well, reference ids are not constant across the life of a location. For information consolidation server-side, the actual place id is used (thank Google for not making it something we can search on)

REST requests and responses follow a standard format.

Lack of existance is shown with a 404 response code.
A forbidden request (i.e. updating a user that is not your own) returns a 403.
A general bad request returns a 400.
Errors do not come with a response body unless otherwise stated.
(Server errors in the form of a 50x from nginx are also possible. Handle them.)

For "Show" type routes, the corresponding (array of) structure(s) in models.go is returned in JSON format. All parameters are in the URL to facilitate ease of caching.
For "Create" type routes, POST data for each field in the corresponding structure in models.go is required. The information that would be returned from the following corresponding "Show" request is returned.
For "Update" type routes, POST data for the updated fields in the corresponding structure in models.go is required. The information that would be returned from the following corresponding "Show" request is returned.
For "Delete" type routes, no POST data or response is returned.

Client/server authentication occurs in /users/login, where a user object as well as a key is returned. For all requests following login, the key should be provided in the request body (Yes it is possible to have a request body for GET).
*/
package main
