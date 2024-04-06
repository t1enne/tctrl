local http = require 'socket.http'
local json = require 'cjson'
local ltn12 = require 'ltn12'

---@alias colors
---| '"white"'
---| '"cyan"'
---| '"magenta"'
---| '"blue"'
---| '"yellow"'
---| '"green"'
---| '"red"'
---| '"black"'
---| '"reset"'

local colors = {
  reset = '\27[0m',
  black = '\27[30m',
  red = '\27[31m',
  green = '\27[32m',
  yellow = '\27[33m',
  blue = '\27[34m',
  magenta = '\27[35m',
  cyan = '\27[36m',
  white = '\27[37m',
}

---@class Utils
---@field start_of_day fun(date: `Date`): `Date`
---@field end_of_day fun(date: `Date`): `Date`
---@field date_to_string fun(date: `Date`): string
---@field escape fun(input: string, color: colors): string
---@field rt_req fun(url: string, data: table<string|any>): any
local M = {}
-- sets the time to 00:00 on the ISO string
M.start_of_day = function(date)
  date:hour(0)
  date:min(0)
  date:sec(0)
  return date
end
-- sets the time to 23:59 on the ISO string
M.end_of_day = function(date)
  date:hour(23)
  date:min(59)
  date:sec(59)
  return date
end
---@param date { tab: {year: number, month: number, day: number, hour: number, min: number, sec: number, msec: number}}
M.date_to_string = function(date)
  local year = date.tab.year
  local day = date.tab.day
  local min = date.tab.min
  local hour = date.tab.hour
  local month = date.tab.month
  local sec = date.tab.sec
  local msec = sec == 59 and 999 or 000
  return string.format('%04d-%02d-%02dT%02d:%02d:%02d.%03dZ', year, month, day, hour, min, sec, msec)
end

M.rt_req = function(url, body)
  local reqbody = json.encode(body)
  local resbody = {}
  local _, c = http.request {
    url = url,
    method = 'POST',
    source = ltn12.source.string(reqbody),
    sink = ltn12.sink.table(resbody),
    headers = {
      ['Content-Type'] = 'application/json',
      ['Accept'] = 'application/json, text/plain, */*',
      ['content-length'] = string.len(reqbody),
      ['Authorization'] = RT_BEARER_TOKEN,
    },
  }
  local response = json.decode(table.concat(resbody))
  if c < 200 and 299 < c then
    error('request' .. url .. 'failed with: ' .. c .. '. message: ' .. response.message)
  end
  return response.data
end

M.escape = function(input, color)
  return colors[color] .. input .. colors.reset
end
return M
