package services

import (
	"fmt"
	"strings"

	"greenride/internal/models"
	"greenride/internal/protocol"
	"gorm.io/gorm"
)

// HardDeleteCleanupSummary stores row counts for each cleanup action.
type HardDeleteCleanupSummary struct {
	DryRun               bool  `json:"dry_run"`
	ZombieUsersFound      int64 `json:"zombie_users_found"`
	ZombieUsersDeleted    int64 `json:"zombie_users_deleted"`
	VehicleUnboundCount   int64 `json:"vehicle_unbound_count"`
	FCMTokensDeletedCount int64 `json:"fcm_tokens_deleted_count"`
	VehiclesDeletedCount  int64 `json:"vehicles_deleted_count"`
	PriceRulesDeleted     int64 `json:"price_rules_deleted"`
	PromotionsDeleted     int64 `json:"promotions_deleted"`
	RatingsDeleted        int64 `json:"ratings_deleted"`
}

// RunHardDeleteCleanup permanently removes legacy soft-deleted rows and zombie records.
func RunHardDeleteCleanup() (*HardDeleteCleanupSummary, protocol.ErrorCode) {
	return RunHardDeleteCleanupWithOptions(false)
}

// RunHardDeleteCleanupWithOptions supports dry-run preview mode.
func RunHardDeleteCleanupWithOptions(dryRun bool) (*HardDeleteCleanupSummary, protocol.ErrorCode) {
	db := models.GetDB()
	if db == nil {
		return nil, protocol.SystemError
	}

	summary := &HardDeleteCleanupSummary{DryRun: dryRun}
	err := db.Transaction(func(tx *gorm.DB) error {
		// NOTE: using raw SQL intentionally to keep cleanup deterministic even if model structs evolve.
		// 1) Find zombie users
		var userRows []struct {
			UserID string
		}
		rawUsers := tx.Raw(`
			SELECT user_id
			FROM t_users
			WHERE status = ?
			   OR (deleted_at IS NOT NULL AND deleted_at > 0)
			   OR phone LIKE 'deleted_%'
			   OR email LIKE 'deleted_%'
		`, protocol.StatusDeleted)
		if err := rawUsers.Scan(&userRows).Error; err != nil {
			return err
		}
		summary.ZombieUsersFound = int64(len(userRows))

		userIDs := make([]string, 0, len(userRows))
		for _, row := range userRows {
			if strings.TrimSpace(row.UserID) != "" {
				userIDs = append(userIDs, row.UserID)
			}
		}

		if len(userIDs) > 0 {
			type countRow struct {
				Count int64 `gorm:"column:cnt"`
			}
			var row countRow
			if err := tx.Raw(`SELECT COUNT(1) AS cnt FROM t_vehicles WHERE driver_id IN ?`, userIDs).Scan(&row).Error; err != nil {
				return err
			}
			summary.VehicleUnboundCount = row.Count

			if err := tx.Raw(`SELECT COUNT(1) AS cnt FROM t_fcm_tokens WHERE user_id IN ?`, userIDs).Scan(&row).Error; err != nil {
				return err
			}
			summary.FCMTokensDeletedCount = row.Count
			summary.ZombieUsersDeleted = int64(len(userIDs))

			if !dryRun {
				unbind := tx.Exec(`UPDATE t_vehicles SET driver_id = '' WHERE driver_id IN ?`, userIDs)
				if unbind.Error != nil {
					return unbind.Error
				}

				delFCM := tx.Exec(`DELETE FROM t_fcm_tokens WHERE user_id IN ?`, userIDs)
				if delFCM.Error != nil {
					return delFCM.Error
				}

				delUsers := tx.Exec(`DELETE FROM t_users WHERE user_id IN ?`, userIDs)
				if delUsers.Error != nil {
					return delUsers.Error
				}
			}
		}

		// 2) Purge vehicles that were marked deleted in legacy flows
		type countRow struct {
			Count int64 `gorm:"column:cnt"`
		}
		var row countRow
		if err := tx.Raw(`SELECT COUNT(1) AS cnt FROM t_vehicles WHERE status = ?`, protocol.StatusDeleted).Scan(&row).Error; err != nil {
			return err
		}
		summary.VehiclesDeletedCount = row.Count
		if !dryRun {
			delVehicles := tx.Exec(`DELETE FROM t_vehicles WHERE status = ?`, protocol.StatusDeleted)
			if delVehicles.Error != nil {
				return delVehicles.Error
			}
		}

		// 3) Purge legacy soft-deleted price rules
		if err := tx.Raw(`SELECT COUNT(1) AS cnt FROM t_price_rules WHERE status = ?`, protocol.StatusDeleted).Scan(&row).Error; err != nil {
			return err
		}
		summary.PriceRulesDeleted = row.Count
		if !dryRun {
			delPriceRules := tx.Exec(`DELETE FROM t_price_rules WHERE status = ?`, protocol.StatusDeleted)
			if delPriceRules.Error != nil {
				return delPriceRules.Error
			}
		}

		// 4) Purge legacy soft-deleted promotions
		if err := tx.Raw(`
			SELECT COUNT(1) AS cnt
			FROM t_promotions
			WHERE status = ?
			   OR (deleted_at IS NOT NULL AND deleted_at > 0)
		`, protocol.StatusDeleted).Scan(&row).Error; err != nil {
			return err
		}
		summary.PromotionsDeleted = row.Count
		if !dryRun {
			delPromos := tx.Exec(`
				DELETE FROM t_promotions
				WHERE status = ?
				   OR (deleted_at IS NOT NULL AND deleted_at > 0)
			`, protocol.StatusDeleted)
			if delPromos.Error != nil {
				return delPromos.Error
			}
		}

		// 5) Purge soft-deleted ratings (gorm.DeletedAt rows)
		if err := tx.Raw(`SELECT COUNT(1) AS cnt FROM t_order_ratings WHERE deleted_at IS NOT NULL`).Scan(&row).Error; err != nil {
			return err
		}
		summary.RatingsDeleted = row.Count
		if !dryRun {
			delRatings := tx.Exec(`DELETE FROM t_order_ratings WHERE deleted_at IS NOT NULL`)
			if delRatings.Error != nil {
				return delRatings.Error
			}
		}

		return nil
	})
	if err != nil {
		return nil, protocol.DatabaseError
	}

	return summary, protocol.Success
}

func (s *HardDeleteCleanupSummary) String() string {
	if s == nil {
		return "no cleanup summary"
	}
	return fmt.Sprintf(
		"dry_run=%t, zombie_users_found=%d, zombie_users_deleted=%d, vehicle_unbound=%d, fcm_tokens_deleted=%d, vehicles_deleted=%d, price_rules_deleted=%d, promotions_deleted=%d, ratings_deleted=%d",
		s.DryRun,
		s.ZombieUsersFound,
		s.ZombieUsersDeleted,
		s.VehicleUnboundCount,
		s.FCMTokensDeletedCount,
		s.VehiclesDeletedCount,
		s.PriceRulesDeleted,
		s.PromotionsDeleted,
		s.RatingsDeleted,
	)
}
