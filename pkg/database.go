package pkg

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type Container struct {
	db *sql.DB
}

var DB *Container

func PrepareDB() {
	var dbUrl = "file:file.db"
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", dbUrl, err)
	}

	DB = &Container{db: db}
}

var base_table_sql string = `
    CREATE TABLE IF NOT EXISTS links (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    title TEXT NOT NULL,
	    description TEXT,
	    link TEXT NOT NULL,
	    thumbnail text
    );
`

func (c *Container) Initial() {
	_, err := c.db.Exec(base_table_sql)
	if err != nil {
		log.Fatalf("failed to create initial schema: %s", err)
	}
}

func (c *Container) Listen() {
	go func() {
		for true {
			time.Sleep(time.Second * 10)

			links, err := c.GetUnprocessedLinks()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, link := range links {
				context, cancel := chromedp.NewContext(context.Background())
				defer cancel()

				var filebyte []byte
				if err := chromedp.Run(context, chromedp.Tasks{
					chromedp.Navigate(link.Link),
					chromedp.Sleep(3 * time.Second),
					chromedp.CaptureScreenshot(&filebyte),
				}); err != nil {
					log.Println(err)
					continue
				}

				encodedImage := base64.StdEncoding.EncodeToString(filebyte)
				if _, err := c.db.Exec("UPDATE links SET thumbnail = ? WHERE id = ?", encodedImage, link.ID); err != nil {
					log.Println(err)
				}
			}

		}
	}()
}

func (c *Container) GetLinks() ([]Link, error) {
	rows, err := c.db.Query("SELECT * FROM links")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Description, &link.Link, &link.Thumbnail)
		if err != nil {
			log.Println(err)
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func (c *Container) GetUnprocessedLinks() ([]Link, error) {
	rows, err := c.db.Query("SELECT * FROM links WHERE thumbnail is null")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Description, &link.Link, &link.Thumbnail)
		if err != nil {
			log.Println(err)
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func (c *Container) CreateLink(body CreateLinkRequest) error {

	_, err := c.db.Exec("INSERT INTO links (title, description, link) VALUES (?, ?, ?)", body.Title, body.Description, body.Link)
	if err != nil {
		return err
	}

	return nil
}
