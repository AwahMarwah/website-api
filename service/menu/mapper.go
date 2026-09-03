package menu

import menuModel "website-api/model/menu"

func toResponse(m menuModel.Menu) menuModel.MenuResponse {
	return menuModel.MenuResponse{
		ID:          m.ID,
		ParentID:    m.ParentID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Icon:        m.Icon,
		Path:        m.Path,
		SortOrder:   m.SortOrder,
		IsActive:    m.IsActive,
	}
}

func buildTree(menus []menuModel.Menu) []menuModel.MenuTree {
	nodeMap := make(map[string]*menuModel.MenuTree)
	for _, m := range menus {
		if !m.IsActive {
			continue
		}
		m := m
		node := &menuModel.MenuTree{
			ID:          m.ID,
			ParentID:    m.ParentID,
			Name:        m.Name,
			DisplayName: m.DisplayName,
			Icon:        m.Icon,
			Path:        m.Path,
			SortOrder:   m.SortOrder,
			IsActive:    m.IsActive,
			Children:    []menuModel.MenuTree{},
		}
		nodeMap[m.ID] = node
	}

	var roots []menuModel.MenuTree
	for _, m := range menus {
		if !m.IsActive {
			continue
		}
		node := nodeMap[m.ID]
		if node == nil {
			continue
		}
		if m.ParentID == nil || *m.ParentID == "" {
			roots = append(roots, *node)
		} else if parent, ok := nodeMap[*m.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		} else {
			roots = append(roots, *node)
		}
	}
	return roots
}
