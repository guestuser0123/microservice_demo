request = function()
    voucherID = math.random(1, 10)
    itemID = math.random(1, 10)
    itemQty = math.random(1, 10)
    requestBody = "{\"item_id\":" .. itemID .. ",\"item_qty\":" .. itemQty .. ",\"voucher_id\":" .. voucherID .. "}"
    return wrk.format("POST", nil, nil, requestBody)
end
