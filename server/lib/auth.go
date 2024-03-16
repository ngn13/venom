package lib

import "server/global"

func AddAgentCookie() string{
  cookie := MakeRandom(42)
  global.AgentCookies = append(global.AgentCookies, cookie)
  return cookie 
}

func CheckAgentCookie(cookie string) bool {
  for _, c := range global.AgentCookies {
    if c == cookie {
      return true
    }
  }
  return false
}

func RemoveAgentCookie(cookie string) {
  indx := -1
  for i, c := range global.AgentCookies {
    if c == c{
      indx = i
      break
    }    
  }

  if indx == -1 {
    return 
  }

  global.AgentCookies = append(
    global.AgentCookies[:indx], global.AgentCookies[indx+1:]...)
}

func AddUserCookie() string{
  cookie := MakeRandom(32)
  global.UserCookies = append(global.UserCookies, cookie)
  return cookie 
}

func CheckUserCookie(cookie string) bool {
  for _, c := range global.UserCookies {
    if c == cookie {
      return true
    }
  }
  return false
}

func RemoveUserCookie(cookie string) {
  indx := -1
  for i, c := range global.UserCookies {
    if c == c{
      indx = i
      break
    }    
  }

  if indx == -1 {
    return 
  }

  global.UserCookies = append(
    global.UserCookies[:indx], global.UserCookies[indx+1:]...)
}
