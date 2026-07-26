package parser

import (
	"fmt"
	"strings"
)

// PolicyInfo хранит полную информацию о политике
type PolicyInfo struct {
	Service    string
	Method     string
	Path       string
	Public     bool
	Permission string
	JWT        bool
}

// collectPolicies рекурсивно собирает все политики из дерева
func collectPolicies(node *Node, serviceName string, basePath string) []PolicyInfo {
	var policies []PolicyInfo
	// Формируем текущий путь
	currentPath := basePath + "/" + node.Path
	switch basePath {
	case "":
		currentPath = node.Path
	case "/":
		currentPath = "/" + node.Path
	}

	// Собираем политики для всех методов в узле
	for method, policy := range node.Methods {
		policies = append(policies, PolicyInfo{
			Service:    serviceName,
			Method:     method,
			Path:       currentPath,
			Public:     policy.Public,
			Permission: policy.Permission,
			JWT:        policy.JWT,
		})
	}

	// Рекурсивно обходим дочерние узлы
	for _, child := range node.Children {
		childPolicies := collectPolicies(child, serviceName, currentPath)
		policies = append(policies, childPolicies...)
	}

	return policies
}

// PrintPoliciesByService выводит политики сгруппированные по сервисам
func (s *ServiceSpecs) PrintPoliciesByService() {
	fmt.Println("\n=== POLICIES BY SERVICE ===")
	for serviceName, serviceInfo := range s.Services {
		fmt.Printf("\n📦 Service: %s\n", serviceName)
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("  %-8s %-40s %-10s %-15s %-5s\n", "METHOD", "PATH", "PUBLIC", "PERMISSION", "JWT")
		fmt.Println("  " + strings.Repeat("-", 76))

		policies := collectPolicies(serviceInfo.Tree.Root, serviceName, "")
		for _, p := range policies {
			publicStr := "🔒 private"
			if p.Public {
				publicStr = "🌍 public"
			}
			fmt.Printf("  %-8s %-40s %-15s %-15s %-5t\n",
				p.Method,
				p.Path,
				publicStr,
				p.Permission,
				p.JWT,
			)
		}
		fmt.Printf("  Total: %d endpoints\n", len(policies))
	}
}
