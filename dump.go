package main

import (
	"log"

	"gorm.io/gorm"
)

func Dump() {
	var holes []Hole
	var floors Floors
	result := DB.
		Where("hidden = false").
		FindInBatches(&holes, 1000, func(tx *gorm.DB, batch int) error {
			if len(holes) == 0 {
				return nil
			}

			err := tx.
				Table("floor").
				Select("id", "content", "updated_at").
				Where("hole_id in (?) and deleted = 0 and ((is_actual_sensitive IS NOT NULL AND is_actual_sensitive = false) OR (is_actual_sensitive IS NULL AND is_sensitive = false))", holeIDs).
				Scan(&floors).Error
			if err != nil {
				return err
			}
			if len(floors) == 0 {
				return nil
			}

			err = BulkInsert(floors)
			if err != nil {
				return err
			}

			log.Printf("insert holes [%d, %d]\n", holes[0].ID, holes[len(holes)-1])
			return nil
		})

	if result.Error != nil {
		log.Fatalf("dump err: %s", result.Error)
	}
}
