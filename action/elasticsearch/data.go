package elasticsearch

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

type Artist struct {
	ID    int64  `csv:"id" json:"id"`
	Name  string `csv:"name" json:"name"`
	Genre string `csv:"genre" json:"genre"`
	Plays int64  `csv:"plays" json:"plays"`
}

type Movie struct {
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	TitleNgran         string `json:"title_ngram"`
	TitleShingle       string `json:"title_shingle"`
	Theme              string `json:"theme"`
	Director           string `json:"director"`
	DirectorNgran      string `json:"director_ngram"`
	DirectorShingle    string `json:"director_shingle"`
	TitleDirector      string `json:"title_director"`
	TitleDirectorNgran string `json:"title_director_ngram"`
	Year               string `json:"year"`
	Views              int    `json:"views"`
}

func newMovieElastic(m *MovieCSV) *Movie {
	return &Movie{
		ID:                 m.ID,
		Title:              m.Title,
		TitleNgran:         m.Title,
		TitleShingle:       m.Title,
		Theme:              m.Theme,
		Director:           m.Director,
		DirectorNgran:      m.Director,
		DirectorShingle:    m.Director,
		TitleDirector:      m.Title + " " + m.Director,
		TitleDirectorNgran: m.Title + " " + m.Director,
		Year:               m.Year,
		Views:              rand.Intn(10000), // Fake views number
	}
}

type MovieCSV struct {
	ID       int    `csv:"codgio"`
	Title    string `csv:"titulo"`
	Resume   string `csv:"sinopse"`
	Director string `csv:"diretor"`
	Year     string `csv:"ano"`
	Theme    string `csv:"tema"`
}

func InsertFakeData(ctx context.Context) error {
	moviesFile, err := os.OpenFile("data/movies-pt.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer moviesFile.Close()

	var movies []*MovieCSV

	if err := gocsv.UnmarshalFile(moviesFile, &movies); err != nil {
		return err
	}

	var artists []*Artist

	artistFile, err := os.OpenFile("data/artists-pt.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	if err := gocsv.UnmarshalFile(artistFile, &artists); err != nil {
		return err
	}

	fmt.Println(artists[0])

	client, err := GetConnection()
	if err != nil {
		return err
	}

	fmt.Printf("Inserting %d movies for tests...\n", len(movies))
	fmt.Printf("Inserting %d artists for tests...\n", len(artists))

	bulk := client.Bulk()
	for index, movie := range movies {
		entry := elastic.NewBulkIndexRequest().
			Index(viper.GetString("indexName")).
			Type("movies").
			Id(strconv.Itoa(index)).
			Doc(newMovieElastic(movie))
		bulk = bulk.Add(entry)
	}

	for index, artist := range artists {
		entry := elastic.NewBulkIndexRequest().
			Index(viper.GetString("indexName")).
			Type("artists").
			Id(strconv.Itoa(index)).
			Doc(artist)
		bulk = bulk.Add(entry)
	}
	fmt.Printf("Bulk actions = %d\n", bulk.NumberOfActions())
	_, err = bulk.Do(ctx)
	return err
}
