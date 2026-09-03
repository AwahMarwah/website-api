package menu

import (
	"time"

	"github.com/google/uuid"
)

func (r *repo) AssignMenus(roleID string, menuIDs []string, now time.Time) error {
	if err := r.db.Exec("DELETE FROM role_menus WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}
	for _, menuID := range menuIDs {
		if err := r.db.Exec(
			"INSERT INTO role_menus (id, role_id, menu_id, created_at) VALUES (?, ?, ?, ?)",
			uuid.NewString(), roleID, menuID, now,
		).Error; err != nil {
			return err
		}
	}
	return nil
}
