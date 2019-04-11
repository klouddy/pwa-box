# PWA Box

**PWA Box** is trying to simplify getting started creating `PWA` or `SPA` applications.  

## Basic Goals

- Serve static html, javascript, css apps.
- Setup reverse proxy to external apis.
- Provide basic metrics out of box (pun intended, *always*) 
- Provide logging endpoint for use by PWA out of the box
- Provide some production ready items out of the box.  
Basically limit number of issues you might have because of any downstream services.  
- Designed with docker and kubernetes in mind.
- Config file should define as much as possible.
- Make it easy to reuse and light on resources so it can be used early and often.

Simple app to host a static js application i.e. Angular or React

## Usage

First crete a config file.  Example:

```json
{
  "reverseProxies": [
    {
      "route": "/api2",
      "remoteServer": "http://localhost:3000"
    }
  ],
  "staticApps": [
    {
      "route": "/app",
      "directory": "./public"
    }
  ],
  "port": 4200
}
```

The above when run with `pwa-box` will create a server at localhost:4200.  The route `/api2` will be a
reverse proxy to `http://lcoalhost:3000`.  There is also static apps configured for route `/app` and will serve
the files located at `./public`. 


The config file be default will be looked for relative to the `pwa-box` at `./config.json`.  
You can specify your own location by including the `-config` argument.
