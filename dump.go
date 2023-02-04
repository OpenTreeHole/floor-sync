package main

import (
	"gorm.io/gorm"
	"log"
)

func Dump() {
	var holes Holes
	var floors Floors
	result := DB.Model(&Hole{}).
		Select("id").
		Where("hidden = ?", false).
		FindInBatches(&holes, 1000, func(tx *gorm.DB, batch int) error {
			if len(holes) == 0 {
				return nil
			}

			var holeIDs = make([]int, 0, len(holes))
			for _, hole := range holes {
				holeIDs = append(holeIDs, hole.ID)
			}
			err := tx.
				Table("floor").
				Select("id", "content").
				Where("hole_id in (?) and deleted = 0", holeIDs).
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

			log.Printf("insert holes [%d, %d]\n", holeIDs[0], holeIDs[len(holeIDs)-1])
			return nil
		})

	if result.Error != nil {
		log.Fatalf("dump err: %s", result.Error)
	}
}
