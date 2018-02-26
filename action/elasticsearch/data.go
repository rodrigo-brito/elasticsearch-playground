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

type MovieElastic struct {
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	TitleNgran         string `json:"title_ngran"`
	TitleShingle       string `json:"title_shingle"`
	Theme              string `json:"theme"`
	Director           string `json:"director"`
	DirectorNgran      string `json:"director_ngran"`
	DirectorShingle    string `json:"director_shingle"`
	TitleDirector      string `json:"title_director"`
	TitleDirectorNgran string `json:"title_director_ngran"`
	Year               string `json:"year"`
	Views              int    `json:"views"`
}

func newMovieElastic(m *MovieCSV) *MovieElastic {
	return &MovieElastic{
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
	clientsFile, err := os.OpenFile("data/movies-pt.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer clientsFile.Close()

	var movies []*MovieCSV

	if err := gocsv.UnmarshalFile(clientsFile, &movies); err != nil {
		return err
	}

	client, err := GetConnection()
	if err != nil {
		return err
	}

	fmt.Printf("Inserting %d movies for tests...\n", len(movies))

	bulk := client.Bulk()
	for index, movie := range movies {
		entry := elastic.NewBulkIndexRequest().
			Index(viper.GetString("indexName")).
			Type("movies").
			Id(strconv.Itoa(index)).
			Doc(newMovieElastic(movie))
		bulk = bulk.Add(entry)
	}
	fmt.Printf("Bulk actions = %d\n", bulk.NumberOfActions())
	_, err = bulk.Do(ctx)
	return err
}
