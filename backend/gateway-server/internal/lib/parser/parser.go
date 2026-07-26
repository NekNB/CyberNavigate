package parser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

// Policy представляет политику доступа для метода
type Policy struct {
	Public     bool
	Permission string
	JWT        bool
}

// Node представляет узел в дереве путей
type Node struct {
	Path      string
	Methods   map[string]*Policy // HTTP метод -> политика
	Children  map[string]*Node   // дочерние узлы
	IsParam   bool               // является ли параметром пути (например {id})
	ParamName string             // имя параметра, если IsParam = true
}

// PathTree представляет дерево всех путей
type PathTree struct {
	Root *Node
	mu   sync.RWMutex
}

// NewPathTree создает новое дерево путей
func NewPathTree() *PathTree {
	return &PathTree{
		Root: &Node{
			Path:     "/",
			Methods:  make(map[string]*Policy),
			Children: make(map[string]*Node),
		},
	}
}

// ParseSpec парсит OpenAPI спецификацию и строит дерево путей
func ParseSpec(spec *openapi3.T) (*PathTree, error) {
	tree := NewPathTree()

	for path, pathItem := range spec.Paths.Map() {
		if err := tree.addPath(path, pathItem); err != nil {
			return nil, fmt.Errorf("failed to add path %s: %w", path, err)
		}
	}

	return tree, nil
}

// addPath добавляет путь и его методы в дерево
func (t *PathTree) addPath(path string, pathItem *openapi3.PathItem) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Нормализуем путь
	normalizedPath := normalizePath(path)
	segments := splitPath(normalizedPath)

	currentNode := t.Root

	// Проходим по сегментам пути
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		isParam := strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
		paramName := ""
		if isParam {
			paramName = segment[1 : len(segment)-1]
		}

		// Ищем существующий дочерний узел
		child, exists := currentNode.Children[segment]
		if !exists {
			child = &Node{
				Path:      segment,
				Methods:   make(map[string]*Policy),
				Children:  make(map[string]*Node),
				IsParam:   isParam,
				ParamName: paramName,
			}
			currentNode.Children[segment] = child
		}

		currentNode = child
	}

	// Добавляем политики для каждого HTTP метода
	methods := map[string]*openapi3.Operation{
		"GET":     pathItem.Get,
		"POST":    pathItem.Post,
		"PUT":     pathItem.Put,
		"DELETE":  pathItem.Delete,
		"PATCH":   pathItem.Patch,
		"HEAD":    pathItem.Head,
		"OPTIONS": pathItem.Options,
	}

	for method, operation := range methods {
		if operation == nil {
			continue
		}

		policy, err := extractPolicy(operation)
		if err != nil {
			return fmt.Errorf("failed to extract policy for %s %s: %w", method, path, err)
		}

		currentNode.Methods[method] = policy
	}

	return nil
}

// extractPolicy извлекает политику из x-auth расширения
func extractPolicy(operation *openapi3.Operation) (*Policy, error) {
	if operation.Extensions == nil {
		return nil, fmt.Errorf("no extensions found")
	}

	xAuthRaw, exists := operation.Extensions["x-auth"]
	if !exists {
		// Если x-auth не указан, создаем политику по умолчанию
		return &Policy{
			Public:     false,
			Permission: "user",
		}, nil
	}

	// Преобразуем в map[string]interface{}
	xAuthMap, ok := xAuthRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid x-auth format")
	}

	policy := &Policy{
		Public:     false,
		Permission: "user",
		JWT:        true,
	}

	// Извлекаем public
	if publicRaw, exists := xAuthMap["public"]; exists {
		if public, ok := publicRaw.(bool); ok {
			policy.Public = public
		}
	}

	// Извлекаем permission
	if permRaw, exists := xAuthMap["permission"]; exists {
		if perm, ok := permRaw.(string); ok {
			policy.Permission = perm
		}
	}

	// Извлекаем jwt
	if jwtRaw, exists := xAuthMap["jwt"]; exists {
		if jwt, ok := jwtRaw.(bool); ok {
			policy.JWT = jwt
		}
	}
	return policy, nil
}

// PolicyFinder обходчик для поиска политики по методу и пути
type PolicyFinder struct {
	tree *PathTree
}

// NewPolicyFinder создает новый обходчик
func NewPolicyFinder(tree *PathTree) *PolicyFinder {
	return &PolicyFinder{tree: tree}
}

// FindPolicy находит политику для заданного метода и пути
func (f *PolicyFinder) FindPolicy(method, path string) (*Policy, error) {
	f.tree.mu.RLock()
	defer f.tree.mu.RUnlock()

	normalizedPath := normalizePath(path)
	segments := splitPath(normalizedPath)

	policy, err := f.findPolicyRecursive(f.tree.Root, segments, method, 0)
	if err != nil {
		return nil, err
	}

	if policy == nil {
		return nil, fmt.Errorf("no policy found for %s %s", method, path)
	}

	return policy, nil
}

// findPolicyRecursive рекурсивно ищет политику в дереве
func (f *PolicyFinder) findPolicyRecursive(node *Node, segments []string, method string, depth int) (*Policy, error) {
	// Достигли конечного узла
	if depth == len(segments) {
		if policy, exists := node.Methods[method]; exists {
			return policy, nil
		}
		return nil, nil
	}

	currentSegment := segments[depth]

	// Пробуем точное совпадение
	if child, exists := node.Children[currentSegment]; exists {
		policy, err := f.findPolicyRecursive(child, segments, method, depth+1)
		if err != nil {
			return nil, err
		}
		if policy != nil {
			return policy, nil
		}
	}

	// Пробуем параметризованные узлы
	for _, child := range node.Children {
		if child.IsParam {
			policy, err := f.findPolicyRecursive(child, segments, method, depth+1)
			if err != nil {
				return nil, err
			}
			if policy != nil {
				return policy, nil
			}
		}
	}

	return nil, nil
}

// Вспомогательные функции
func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimSuffix(path, "/")
}

func splitPath(path string) []string {
	segments := strings.Split(path, "/")
	// Удаляем пустые сегменты в начале и конце
	var result []string
	for _, segment := range segments {
		if segment != "" {
			result = append(result, segment)
		}
	}
	return result
}

// GetTreeInfo возвращает информацию о дереве для отладки
func (t *PathTree) GetTreeInfo() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var builder strings.Builder
	t.buildTreeInfo(t.Root, "", &builder)
	return builder.String()
}

func (t *PathTree) buildTreeInfo(node *Node, prefix string, builder *strings.Builder) {
	if node.Path != "" {
		fmt.Fprintf(builder, "%s%s", prefix, node.Path)
		if node.IsParam {
			builder.WriteString(" (param)")
		}
		builder.WriteString("\n")

		for method, policy := range node.Methods {
			fmt.Fprintf(builder, "%s  %s: public=%v, permission=%s\n",
				prefix, method, policy.Public, policy.Permission)
		}
	}

	for _, child := range node.Children {
		t.buildTreeInfo(child, prefix+"  ", builder)
	}
}
