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
}

// GetAllPolicies возвращает список всех политик из всех сервисов
func (s *ServiceSpecs) GetAllPolicies() []PolicyInfo {
	var policies []PolicyInfo

	for serviceName, serviceInfo := range s.Services {
		policies = append(policies, collectPolicies(serviceInfo.Tree.Root, serviceName, "")...)
	}

	return policies
}

// collectPolicies рекурсивно собирает все политики из дерева
func collectPolicies(node *Node, serviceName string, basePath string) []PolicyInfo {
	var policies []PolicyInfo

	// Формируем текущий путь
	currentPath := basePath + "/" + node.Path
	if basePath == "/" {
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
		})
	}

	// Рекурсивно обходим дочерние узлы
	for _, child := range node.Children {
		childPolicies := collectPolicies(child, serviceName, currentPath)
		policies = append(policies, childPolicies...)
	}

	return policies
}

// PrintAllPolicies выводит все политики в читаемом формате
func (s *ServiceSpecs) PrintAllPolicies() {
	policies := s.GetAllPolicies()

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("%-20s %-8s %-40s %-10s %-15s\n", "SERVICE", "METHOD", "PATH", "PUBLIC", "PERMISSION")
	fmt.Println(strings.Repeat("-", 100))

	for _, p := range policies {
		publicStr := "private"
		if p.Public {
			publicStr = "public"
		}

		fmt.Printf("%-20s %-8s %-40s %-10s %-15s\n",
			p.Service,
			p.Method,
			p.Path,
			publicStr,
			p.Permission,
		)
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Total: %d policies\n", len(policies))
}

// PrintPoliciesByService выводит политики сгруппированные по сервисам
func (s *ServiceSpecs) PrintPoliciesByService() {
	for serviceName, serviceInfo := range s.Services {
		fmt.Printf("\n📦 Service: %s\n", serviceName)
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("  %-8s %-40s %-10s %-15s\n", "METHOD", "PATH", "PUBLIC", "PERMISSION")
		fmt.Println("  " + strings.Repeat("-", 76))

		policies := collectPolicies(serviceInfo.Tree.Root, serviceName, "")
		for _, p := range policies {
			publicStr := "🔒 private"
			if p.Public {
				publicStr = "🌍 public"
			}

			fmt.Printf("  %-8s %-40s %-15s %-15s\n",
				p.Method,
				p.Path,
				publicStr,
				p.Permission,
			)
		}
		fmt.Printf("  Total: %d endpoints\n", len(policies))
	}
}

// GetPoliciesByService возвращает политики сгруппированные по сервисам
func (s *ServiceSpecs) GetPoliciesByService() map[string][]PolicyInfo {
	result := make(map[string][]PolicyInfo)

	for serviceName, serviceInfo := range s.Services {
		result[serviceName] = collectPolicies(serviceInfo.Tree.Root, serviceName, "")
	}

	return result
}

// FilterPolicies позволяет фильтровать политики по различным критериям
type PolicyFilter struct {
	Service    string // если пусто - все сервисы
	Method     string // если пусто - все методы
	Public     *bool  // nil - все, true - только public, false - только private
	Permission string // если пусто - все права
}

// FindPoliciesByFilter возвращает политики соответствующие фильтру
func (s *ServiceSpecs) FindPoliciesByFilter(filter PolicyFilter) []PolicyInfo {
	allPolicies := s.GetAllPolicies()
	var filtered []PolicyInfo

	for _, p := range allPolicies {
		if filter.Service != "" && p.Service != filter.Service {
			continue
		}
		if filter.Method != "" && p.Method != filter.Method {
			continue
		}
		if filter.Public != nil && p.Public != *filter.Public {
			continue
		}
		if filter.Permission != "" && p.Permission != filter.Permission {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered
}

// PrintFilteredPolicies выводит отфильтрованные политики
func (s *ServiceSpecs) PrintFilteredPolicies(filter PolicyFilter) {
	policies := s.FindPoliciesByFilter(filter)

	if len(policies) == 0 {
		fmt.Println("No policies found matching the filter")
		return
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Filtered Policies:")
	if filter.Service != "" {
		fmt.Printf("Service: %s\n", filter.Service)
	}
	if filter.Method != "" {
		fmt.Printf("Method: %s\n", filter.Method)
	}
	if filter.Public != nil {
		fmt.Printf("Public: %v\n", *filter.Public)
	}
	if filter.Permission != "" {
		fmt.Printf("Permission: %s\n", filter.Permission)
	}
	fmt.Println(strings.Repeat("-", 100))

	fmt.Printf("%-20s %-8s %-40s %-10s %-15s\n", "SERVICE", "METHOD", "PATH", "PUBLIC", "PERMISSION")
	fmt.Println(strings.Repeat("-", 100))

	for _, p := range policies {
		publicStr := "private"
		if p.Public {
			publicStr = "public"
		}

		fmt.Printf("%-20s %-8s %-40s %-10s %-15s\n",
			p.Service,
			p.Method,
			p.Path,
			publicStr,
			p.Permission,
		)
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Found: %d policies\n", len(policies))
}

// ExportPoliciesToJSON экспортирует все политики в JSON формат (для API или отладки)
func (s *ServiceSpecs) ExportPoliciesToJSON() string {
	policies := s.GetAllPolicies()

	var json strings.Builder
	json.WriteString("[\n")

	for i, p := range policies {
		json.WriteString("  {\n")
		json.WriteString(fmt.Sprintf("    \"service\": \"%s\",\n", p.Service))
		json.WriteString(fmt.Sprintf("    \"method\": \"%s\",\n", p.Method))
		json.WriteString(fmt.Sprintf("    \"path\": \"%s\",\n", p.Path))
		json.WriteString(fmt.Sprintf("    \"public\": %v,\n", p.Public))
		json.WriteString(fmt.Sprintf("    \"permission\": \"%s\"\n", p.Permission))

		if i < len(policies)-1 {
			json.WriteString("  },\n")
		} else {
			json.WriteString("  }\n")
		}
	}

	json.WriteString("]")
	return json.String()
}
