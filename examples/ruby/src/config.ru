require 'json'

run do |env|
    [
        200,
        { "content-type" => "application/json" },
        [JSON.generate({status:"ok",app:"example-ruby-app"})]
    ]
end
