# Go web dev project
my little go project to learn the fundamentals of web development

## Network Setup
For hosting both mariadb and my go app are being hosted on a raspberry pi.
The database is closed from the wider web and can only be accessed from localhost,
while the API has port 8080 port forwarded through the router and is exposed to traffic.
This Go API then serves HTML to the user via HTTP

You should be able to access this yourself by visiting http://81.100.84.76:8080!
