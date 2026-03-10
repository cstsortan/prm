package service

import (
	"fmt"
	"path/filepath"

	"github.com/cstsortan/prm/internal/model"
)

// TreeNode represents an entity and its children in a hierarchy.
type TreeNode struct {
	Entity   *model.Entity
	Children []*TreeNode
}

// Tree builds the hierarchy tree for a given epic, or all epics if ref is empty.
func (svc *Service) Tree(ref string) ([]*TreeNode, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}

	if ref != "" {
		result, err := svc.Store.Resolve(idx, ref)
		if err != nil {
			return nil, err
		}
		node, err := svc.buildTreeNode(result.Dir, result.Entity)
		if err != nil {
			return nil, err
		}
		return []*TreeNode{node}, nil
	}

	// All epics
	epicDir := filepath.Join(svc.Store.Root(), "epics")
	slugs, err := svc.Store.ListDirs(epicDir)
	if err != nil {
		return nil, fmt.Errorf("listing epics: %w", err)
	}

	var nodes []*TreeNode
	for _, slug := range slugs {
		dir := filepath.Join(epicDir, slug)
		entity, err := svc.Store.ReadEntity(dir)
		if err != nil {
			continue
		}
		node, err := svc.buildTreeNode(dir, entity)
		if err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (svc *Service) buildTreeNode(dir string, entity *model.Entity) (*TreeNode, error) {
	node := &TreeNode{Entity: entity}

	childDirName := entity.ChildDir()
	if childDirName == "" {
		return node, nil
	}

	childrenDir := filepath.Join(dir, childDirName)
	slugs, err := svc.Store.ListDirs(childrenDir)
	if err != nil {
		return node, nil
	}

	for _, slug := range slugs {
		childDir := filepath.Join(childrenDir, slug)
		childEntity, err := svc.Store.ReadEntity(childDir)
		if err != nil {
			continue
		}
		childNode, err := svc.buildTreeNode(childDir, childEntity)
		if err != nil {
			continue
		}
		node.Children = append(node.Children, childNode)
	}

	return node, nil
}
