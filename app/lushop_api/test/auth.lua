-- auth.lua
function request()
    -- 替换为有效的 JWT Token
    local token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiTmlja05hbWUiOiIxOTgyMTIxNjgwNiIsIkF1dGhvcml0eUlkIjoxLCJEZXZpY2VJRCI6IiIsImV4cCI6MTc1ODI2ODg3MCwiaXNzIjoibHVjaWVuIiwibmJmIjoxNzU1Njc2ODcwfQ.RgwaJda_KPq8jeVjowoEqEKkiRoYCXlg8Lpk3qDiWk8"
    wrk.headers["Authorization"] = "Bearer " .. token
    return wrk.format("GET", "/v2/goods/list/es?pmin=&pmax=&ih=0&in=0&it=0&c=&pn=1&pnum=10&q=&b=632")  -- 替换为实际接口路径和参数
end