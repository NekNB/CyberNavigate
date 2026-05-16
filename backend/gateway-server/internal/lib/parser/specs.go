package parser

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// SpecInfo хранит информацию о загруженной спецификации
type SpecInfo struct {
	Name string      // имя сервиса (например "article-service")
	Spec *openapi3.T // распарсенная спецификация
	Tree *PathTree   // дерево путей с политиками
}

// ServiceSpecs содержит все загруженные спецификации
type ServiceSpecs struct {
	Services map[string]*SpecInfo
}

// LoadSpecsFromFS загружает все спецификации из переданной embed.FS
// prefix - префикс для всех путей (например "/api/v1")
// rootDir - корневая директория в FS где лежат папки сервисов (например "docs")
func LoadSpecsFromFS(specsFS embed.FS, prefix string, rootDir string) (*ServiceSpecs, error) {
	serviceSpecs := &ServiceSpecs{
		Services: make(map[string]*SpecInfo),
	}

	err := fs.WalkDir(specsFS, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Обрабатываем только YAML/YML файлы
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		// Извлекаем имя сервиса из пути
		// rootDir/service-name/spec.yaml -> service-name
		serviceName := extractServiceNameFromPath(path, rootDir)

		// Читаем файл
		data, err := specsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Парсим спецификацию
		loader := openapi3.NewLoader()
		spec, err := loader.LoadFromData(data)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Строим дерево путей с префиксом
		tree := NewPathTree()
		for apiPath, pathItem := range spec.Paths.Map() {
			fullPath := prefix + apiPath
			if err := tree.addPath(fullPath, pathItem); err != nil {
				return fmt.Errorf("failed to add path %s: %w", fullPath, err)
			}
		}

		// Если сервис уже существует (несколько файлов спецификаций), объединяем
		if existing, ok := serviceSpecs.Services[serviceName]; ok {
			// Объединяем деревья
			if err := mergePathTrees(existing.Tree, tree); err != nil {
				return fmt.Errorf("failed to merge trees for %s: %w", serviceName, err)
			}
		} else {
			serviceSpecs.Services[serviceName] = &SpecInfo{
				Name: serviceName,
				Spec: spec,
				Tree: tree,
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk FS: %w", err)
	}

	return serviceSpecs, nil
}

// extractServiceNameFromPath извлекает имя сервиса из пути
// Пример: "docs/user-service/user.swagger.yaml" -> "user-service"
func extractServiceNameFromPath(path string, rootDir string) string {
	// Убираем rootDir из начала пути
	relativePath := strings.TrimPrefix(path, rootDir)
	relativePath = strings.TrimPrefix(relativePath, "/")

	// Разбиваем путь и берем первую директорию
	parts := strings.Split(filepath.ToSlash(relativePath), "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	// Если не получилось извлечь, используем имя файла без расширения
	baseName := filepath.Base(path)
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}

// BuildPathTree строит дерево путей из спецификации с префиксом
func BuildPathTree(spec *openapi3.T, prefix string) (*PathTree, error) {
	tree := NewPathTree()

	for path, pathItem := range spec.Paths.Map() {
		// Добавляем префикс ко всем путям
		fullPath := prefix + path

		if err := tree.addPath(fullPath, pathItem); err != nil {
			return nil, fmt.Errorf("failed to add path %s: %w", fullPath, err)
		}
	}

	return tree, nil
}

// FindPolicy ищет политику по всем сервисам
func (s *ServiceSpecs) FindPolicy(method, path string) (*Policy, string, error) {
	for serviceName, serviceInfo := range s.Services {
		finder := NewPolicyFinder(serviceInfo.Tree)
		policy, err := finder.FindPolicy(method, path)
		if err == nil && policy != nil {
			return policy, serviceName, nil
		}
	}

	return nil, "", fmt.Errorf("no policy found for %s %s", method, path)
}

// GetServiceByName возвращает информацию о сервисе по имени
func (s *ServiceSpecs) GetServiceByName(name string) (*SpecInfo, error) {
	service, exists := s.Services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

// GetAllServices возвращает список всех загруженных сервисов
func (s *ServiceSpecs) GetAllServices() []string {
	services := make([]string, 0, len(s.Services))
	for name := range s.Services {
		services = append(services, name)
	}
	return services
}

// mergePathTrees объединяет два дерева путей
func mergePathTrees(target, source *PathTree) error {
	for path, sourceNode := range source.Root.Children {
		if targetNode, exists := target.Root.Children[path]; exists {
			if err := mergeNodes(targetNode, sourceNode); err != nil {
				return err
			}
		} else {
			target.Root.Children[path] = sourceNode
		}
	}
	return nil
}

// mergeNodes рекурсивно объединяет узлы
func mergeNodes(target, source *Node) error {
	// Объединяем методы
	for method, policy := range source.Methods {
		target.Methods[method] = policy
	}

	// Объединяем дочерние узлы
	for path, sourceChild := range source.Children {
		if targetChild, exists := target.Children[path]; exists {
			if err := mergeNodes(targetChild, sourceChild); err != nil {
				return err
			}
		} else {
			target.Children[path] = sourceChild
		}
	}

	return nil
}
