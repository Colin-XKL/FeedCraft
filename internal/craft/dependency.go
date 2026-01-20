package craft

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"fmt"
	"strings"
)

type DependencyNode struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"` // recipe, flow, atom, built-in, missing, cycle
	Exists   bool              `json:"exists"`
	Children []*DependencyNode `json:"children,omitempty"`
	Details  string            `json:"details,omitempty"`
	Key      string            `json:"key"` // Unique key for tree
}

func AnalyzeDependencies() ([]*DependencyNode, error) {
	db := util.GetDatabase()

	// 1. Load Data
	recipes, err := dao.ListCustomRecipeV2(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load recipes: %w", err)
	}

	flows, err := dao.GetAllCraftFlows(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load flows: %w", err)
	}
	flowMap := make(map[string]dao.CraftFlow)
	for _, f := range flows {
		flowMap[f.Name] = f
	}

	atoms, err := dao.GetAllCraftAtoms(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load atoms: %w", err)
	}
	atomMap := make(map[string]dao.CraftAtom)
	for _, a := range atoms {
		atomMap[a.Name] = a
	}

	sysTemplates := GetSysCraftTemplateDict()

	// 2. Build Tree
	var roots []*DependencyNode

	for _, r := range recipes {
		root := &DependencyNode{
			Name:    r.ID,
			Type:    "recipe",
			Exists:  true,
			Key:     "recipe-" + r.ID,
			Details: r.Description,
		}

		stack := map[string]bool{} // checking cycles in names

		// r.Craft can be "flow1, flow2"
		parts := strings.Split(r.Craft, ",")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			child := analyzeNode(part, stack, flowMap, atomMap, sysTemplates, fmt.Sprintf("%s-%d", root.Key, i))
			root.Children = append(root.Children, child)
		}

		roots = append(roots, root)
	}

	return roots, nil
}

func analyzeNode(name string, stack map[string]bool, flowMap map[string]dao.CraftFlow, atomMap map[string]dao.CraftAtom, sysTemplates map[string]CraftTemplate, parentKey string) *DependencyNode {
	node := &DependencyNode{
		Name: name,
		Key:  parentKey,
	}

	// Cycle detection
	if stack[name] {
		node.Type = "cycle"
		node.Exists = false
		node.Details = "Cycle detected"
		return node
	}

	newStack := make(map[string]bool)
	for k, v := range stack {
		newStack[k] = v
	}
	newStack[name] = true

	// Check Built-in
	if _, ok := sysTemplates[name]; ok {
		node.Type = "built-in"
		node.Exists = true
		return node
	}

	// Check Atom
	if atom, ok := atomMap[name]; ok {
		node.Type = "atom"
		node.Exists = true
		// Check if atom references a valid template
		if _, ok := sysTemplates[atom.TemplateName]; !ok {
			node.Details = "Template missing: " + atom.TemplateName
		}
		return node
	}

	// Check Flow
	if flow, ok := flowMap[name]; ok {
		node.Type = "flow"
		node.Exists = true

		for i, item := range flow.CraftFlowConfig {
			child := analyzeNode(item.CraftName, newStack, flowMap, atomMap, sysTemplates, fmt.Sprintf("%s-%d", parentKey, i))
			node.Children = append(node.Children, child)
		}
		return node
	}

	// Missing
	node.Type = "missing"
	node.Exists = false
	return node
}
