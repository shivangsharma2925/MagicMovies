package utilities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
)

func IsValidImdbID(id string) bool {
	if len(id) < 3 {
		return false
	}
	return strings.HasPrefix(id, "tt")
}

func WithRetry(attempts int, ctx context.Context, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return err
}

func DoRequest(ctx context.Context, url string, target any, httpClient *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tmdb error: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func MapGenres(tmdbGenres []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}, genreMap map[string]models.Genre) []models.Genre {

	// var result []models.Genre
	result := make([]models.Genre, 0, len(tmdbGenres))

	for _, g := range tmdbGenres {
		key := strings.ToLower(g.Name)

		if genre, ok := genreMap[key]; ok {
			result = append(result, genre)
		}
	}

	return result
}

func GetRanking(rating float64) models.Ranking {

	var ranking models.Ranking

	// ratingValue, err := strconv.ParseFloat(rating, 64)
	// if err != nil {
	// 	return models.Ranking{
	// 		RankingName:  "Not_Ranked",
	// 		RankingValue: 999,
	// 	}
	// }

	if rating > 8 {
		ranking.RankingName = "Excellent"
		ranking.RankingValue = 1
	} else if rating > 7 {
		ranking.RankingName = "Good"
		ranking.RankingValue = 2
	} else if rating > 6 {
		ranking.RankingName = "Okay"
		ranking.RankingValue = 3
	} else if rating > 5 {
		ranking.RankingName = "Bad"
		ranking.RankingValue = 4
	} else {
		ranking.RankingName = "Terrible"
		ranking.RankingValue = 5
	}

	return ranking
}

func MapToSchema(data *models.MovieData) *models.Movie {

	return &models.Movie{
		ImdbID:      data.ImdbID,
		Title:       data.Title,
		PosterPath:  data.Poster,
		YouTubeID:   data.YoutubeID,
		Genre:       data.Genres,
		AdminReview: "", // default empty
		Ranking:     data.Ranking,
	}
}
