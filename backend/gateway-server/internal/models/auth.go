package models

// RoutePolicy — результат парсинга OpenAPI
// Это "семантика endpoint'а"
type RoutePolicy struct {
	Method     string // GET, POST...
	Path       string // /users/{id}
	Public     bool   // можно без auth
	Permission string // users.list
}

// CompiledRoute — то, что попадёт в Trie
type CompiledRoute struct {
	Policy RoutePolicy
}
