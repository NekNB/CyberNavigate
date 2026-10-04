const services = [
  "user-service",
  "gateway-server",
  "frontend",
  "nginx",
  "article-service",
  "simulator-service",
  "postgres",
  "mongo"
]
$services | enumerate | each {
  |elt| echo $elt.item
}
