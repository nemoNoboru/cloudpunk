# Roadmap:
Adding and running functions dynamically is DONE AND WORKING.
Design Decisions:
all information should be in nats. HTML, Images, Code
Next:
- [] Implement a bridge nats-lua to be able to call remote functions
- [x] Implement a HTTP endpoint to call functions
- - [] adapt common HTTP params into a map to be sent to the lua function
- [] Research: Check how we could do templating here. Go templates vs client directly

node bootstraping:
Plan: we have three diferent type of data.
1. DB: mutable, runtime data that MUST be preserved
2. BLob: inmutable, on-disk data
3. Functions: runnable on-disk data.

the plan is to run cloudpunk from a folder that has three subfolders:
/functions - lua files containing functions. They will callable from nats or http://namespace.cloudpunk.org/api
/static - various files containing static data. They will be callable also from nats or http://namespace.cloudpunk.org/*

the idea is to have a "LOAD" message with certains metadata
- data ([]byte the actual data we want to load)
- handler (string, luavm, just a http response etc. )
- topic
- optional fields...

then, if the handler accepts the LOAD
it suscribes to the topic in order to RUN the HANDLER

handlers could be:
DB
LUA
WASM
FILE - for serving static files


