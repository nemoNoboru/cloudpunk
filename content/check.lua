json = require "json"

-- Define a simple function that takes two numbers and returns their sum
function add(a, b)
    return a + b
end

-- Define another function to print a greeting
function serverless(name)
    return json.encode({result = "Hello, " .. name .. "!", sanity_check = add(2,2)})
end



