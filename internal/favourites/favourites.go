package favourites

// Add, remove, and check favourites for a playlist

func AddFavourite(favs map[string][]string, playlist, channel string) {
	for _, ch := range favs[playlist] {
		if ch == channel {
			return // already present
		}
	}
	favs[playlist] = append(favs[playlist], channel)
}

func RemoveFavourite(favs map[string][]string, playlist, channel string) {
	list := favs[playlist]
	for i, ch := range list {
		if ch == channel {
			favs[playlist] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func IsFavourite(favs map[string][]string, playlist, channel string) bool {
	for _, ch := range favs[playlist] {
		if ch == channel {
			return true
		}
	}
	return false
}
