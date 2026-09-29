package recommend

func DifficultyBand(rating int) (min, max int) {
	//TODO: ratingが低い時用に、AC済み問題からdifficultyを推測できるようにしたい
	return rating, rating + 200
}