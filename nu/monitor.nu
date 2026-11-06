
def main [] {
  show-status 
  loop {
    let key = (input listen --types [key])
    if ($key.code == 'r') and ($key.modifiers == []) {
      show-status
    } else if ($key.code == 'x') and ($key.modifiers == []) or ($key.code == 'c' and $key.modifiers == ["keymodifiers(control)"]) {
      break
    } 
  }
}

def show-status [] {
  clear
  docker stack ps -f desired-state=running -f desired-state=ready cyber-navigate
  print ""
  print "R — обновить"
  print "X — выйти"
}