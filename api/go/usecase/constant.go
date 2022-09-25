package usecase

import "time"

const (
	PERIOD    int = 7
	WEEK      int = 6
	COUNT_DAY int = 60 // 考慮日数
	FRAME_NUM int = 39
)

var (
	JST *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)
)
