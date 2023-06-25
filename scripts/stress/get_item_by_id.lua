request = function()
  wrk.method="GET"
  param_value = math.random(1,10)
  path = "/item/get?id=" .. param_value
  return wrk.format("GET", path)
end
