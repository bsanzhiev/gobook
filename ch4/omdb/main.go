package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Movie struct {
	Poster string
}

func main() {
	movieName := flag.String("movie", "", "Name of movie")
	flag.Parse()

	if *movieName == "" {
		fmt.Println("Please provide a movie name")
		os.Exit(1)
	}

	// Make request to OMDB API
	res, err := http.Get("https://wwww.omdbapi.com/?t=" + url.QueryEscape(*movieName) + "&apikey=8acc819d")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// Decode JSON response
	var movie Movie
	if err := json.NewDecoder(res.Body).Decode(&movie); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Download poster
	res, err = http.Get(movie.Poster)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// create the file
	out, err := os.Create("poster.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
