json = require "json"

-- Define another function to print a greeting
function api(name)

    input = json.decode(name)

    local comments = {
        {
            author="john doe",
            comment="this site sucks"
        },
        {
            author="armitage",
            comment="ping me"
        },
        {
            author="Molly",
            comment="not too bad for a shit pet project"
        },
    }

    return json.encode(comments)
end

-- things that we need
-- template either or go/lua
-- a way to import and render fragments
-- a way to import files (html, template)
-- first thing: a two way from and to nats from lua
function render(data)
    local template = require("template")

    local t = cloudpunk.storage.get("index")

    local comments = json.decode(api(json.encode({})))

    print(comments)

    local env = {
        pairs  = pairs,
        ipairs = ipairs,
        type   = type,
        table  = table,
        string = string,
        date   = os.date,
        math   = math,
        comments=comments
    }

   return template.compile(t, env)
end
