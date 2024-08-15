json = require "json"

-- Define another function to print a greeting
function serverless(name)

    input = json.decode(name)

    return json.encode({result = "Hello, " .. input['url'] .. "!" })
end



