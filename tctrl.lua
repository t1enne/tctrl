local json = require 'cjson'
local Date = require 'pl.Date'
local lapp = require 'pl.lapp'
local utils = require 'utils'

RT_BEARER_TOKEN = ''
RT_USER_ID = ''
API_BASE = ''

-- Function to get token
local function get_token()
  io.write 'Input username: '
  local username = io.read()
  io.write 'Input password: '
  local password = io.read()
  local body = { email = username .. '@raintonic.com', password = password }
  local data = utils.rt_req(API_BASE .. '/auth/login', body)
  return data
end

local function pick_tag()
  local tags = utils.rt_req(API_BASE .. '/hoursTags/fb', {})
  for i, c in ipairs(tags) do
    io.write(utils.escape(i .. ': ' .. c.name .. '\n', i % 2 == 0 and 'cyan' or 'green'))
  end
  io.write 'Pick a tag: '
  local picked = tags[io.read '*n']
  return picked
end

local function print_worked_hrs(input_date)
  local from_date = utils.date_to_string(utils.start_of_day(input_date))
  local to_date = utils.date_to_string(utils.end_of_day(input_date))
  local body = {
    relations = { 'release', 'release.project', 'release.project.customer', 'hoursTag' },
    where = {
      userId = RT_USER_ID,
      date = {
        _fn = 17,
        args = { from_date, to_date },
      },
    },
  }
  io.write 'Contacting server...\n'
  ---@type {hours: string, notes: string, release: table, hoursTag: table}[]
  local worked_data = utils.rt_req(API_BASE .. '/userHours/fb', body)
  local total = 0
  if #worked_data == 0 then
    io.write 'No hours worked\n'
  end
  for _, entry in ipairs(worked_data) do
    total = total + tonumber(entry.hours)
    print(utils.escape(entry.release.project.customer.name .. ': ' .. entry.release.name, 'yellow'))
    print(utils.escape(entry.hoursTag.name .. ' - ' .. entry.hours .. ' ore. Note: ' .. entry.notes, 'green'))
  end
  io.flush()
  io.write('Worked hours: ' .. total .. '\n')
end

---@param proj_id string
---@return { _id: string, name: string }
local function pick_release(proj_id)
  local body = { order = { name = 'ASC' }, where = { projectId = proj_id } }
  local releases = utils.rt_req(API_BASE .. '/releases/fb', body)
  for i, c in ipairs(releases) do
    io.write(utils.escape(i .. ': ' .. c.name .. '\n', i % 2 == 0 and 'cyan' or 'green'))
  end
  io.write 'Pick a release:\n'
  local picked = releases[io.read '*n']
  return picked
end

---@param client_id string
---@return { _id: string, name: string }
local function pick_project(client_id)
  local body = { order = { name = 'ASC' }, where = { customerId = client_id } }
  local projects = utils.rt_req(API_BASE .. '/projects/fb', body)
  for i, c in ipairs(projects) do
    io.write(utils.escape(i .. ': ' .. c.name .. '\n', i % 2 == 0 and 'cyan' or 'green'))
  end
  io.write 'Pick a project: '
  local picked = projects[io.read '*n']
  return picked
end

---@return { _id: string, name: string }
local function pick_client()
  local body = { order = { name = 'ASC' } }
  local clients = utils.rt_req(API_BASE .. '/customers/fb', body)
  for i, c in ipairs(clients) do
    io.write(utils.escape(i .. ': ' .. c.name .. '\n', i % 2 == 0 and 'cyan' or 'green'))
  end
  io.write 'Pick a client: '
  local picked = clients[io.read '*n']
  return picked
end

local function handle_auth()
  local CACHE = os.getenv 'HOME' .. '/.cache/tcontrol.json'
  local cache_file = io.open(CACHE, 'r')
  if cache_file then
    local contents = cache_file:read '*all'
    local cached_token = json.decode(contents)
    RT_USER_ID = cached_token.RT_USER_ID
    RT_BEARER_TOKEN = cached_token.RT_BEARER_TOKEN
    API_BASE = cached_token.API_BASE
    cache_file:close()
  else
    local response = get_token()
    cache_file = io.open(CACHE, 'w')
    if response == nil then
      error 'error getting token'
    end
    if cache_file == nil then
      error("wasn't able to open cache file" .. CACHE)
    end
    local cache_content = { RT_USER_ID = response.user._id, RT_BEARER_TOKEN = 'Bearer ' .. response.token }
    cache_file:write(json.encode(cache_content))
    cache_file:close()
  end
end

-- Function to main logic
local function main()
  ---@type { date: string }
  local args = lapp [[
	CLI interface for Raintonic's Tcontrol webapp.
		
		-d,--date (optional string)
	]]
  local df = Date.Format 'dd/mm/yy'
  local input_date = args.date and df:parse(args.date) or Date()
  handle_auth()
  print_worked_hrs(input_date)
  io.write '-------\n'
  local client = pick_client()
  local project = pick_project(client._id)
  local release = pick_release(project._id)
  io.flush()

  local notes = nil
  io.write 'Notes: '
  while notes == nil or notes == '' do
    notes = io.read()
  end

  local hours = nil
  io.write 'Hours: '
  while hours == nil or hours == '' do
    hours = io.read()
  end
  local tag = pick_tag()

  local isodate = utils.date_to_string(utils.start_of_day(input_date))
  local body = {
    notes = notes,
    hours = tonumber(hours),
    date = isodate,
    releaseId = release._id,
    hoursTagId = tag._id,
    userId = RT_USER_ID,
  }
  local _ = utils.rt_req(API_BASE .. '/userHours', body)
end

-- Run main function
main()
print 'All done!'
os.exit(0)
