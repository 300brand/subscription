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
               / Cookies \
               \  Good?  /
                 \     /
                   \ /
                    |
+------------------------+
| Load config for mysite |
+------------------------+

+-----------------------------------------+
| GET/POST login credentials to real site |
+-----------------------------------------+

+----------------+
| Ensure success |

