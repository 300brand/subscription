Subscription Proxy Service
==========================

Stands up a webserver to listen for incoming requests on a subdomain of itself
and then connect to the configured website, login if necessary, then perform
the HTTP request as a logged in user.

How it Works
------------

         +--------------------------------+
         | mysite.subscription.net/secret |
         +--------------------------------+
                         |
                         |
                         v
                        / \
                      /     \
                    / Cookies \_____________.
                    \  Good?  / Yep         |
                      \     /               |
                        \ /                 |
                    Nope |                  |
                         |                  |
                         v                  |
             +------------------------+     |
             | Load config for mysite |     |
             +------------------------+     |
                         |                  |
                         v                  |
            +--------------------------+    |
            | GET/POST login           |    |
            | credentials to real site |    |
            +--------------------------+    |
                         |                  |
                         v                  |
                +----------------+          |
                | Update cookies |          |
                +----------------+          |
                         |                  |
                         v                  |
               +------------------+         |
               | Request /secret  |<--------+
               +------------------+
                         |
                         v
          +----------------------------+
          | Rewrite any absolute URLs  |
          | to mysite with URLs to     |
          | mysite.subscription.net    |
          +----------------------------+
                         |
                         v
          +----------------------------+
          | Transfer headers           |
          |                            |
          | io.Copy(resp.Wr, req.Body) |
          +----------------------------+
