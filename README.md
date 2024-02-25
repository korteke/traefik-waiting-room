# Traefik Waiting room Middleware Plugin
This Traefik plugin tries to mimic Cloudlares implementation [Cloudflare Waiting Room](https://www.cloudflare.com/application-services/products/waiting-room/).

The idea behind it is that each request gets a "sticky" cookie from the Traefik load balancer. Once the cookie is received, the value of the cookie is stored in a cache for xx minutes. Once the cache is full, subsequent requests are redirected to a waiting page that updates every xx seconds, checking if there is capasity again.

Currently still very incomplete and under development. I will first try to see if it is even possible / feasible to implement it.

### Installation
WIP

### Usage
WIP

### Performance tests
WIP

### Docs
