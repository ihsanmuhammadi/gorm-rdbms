package controller

import (
	"encoding/json"
	"fmt"
	"gorm-rdbms/database"
	"gorm-rdbms/models"
	"gorm-rdbms/request"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func CreateNews(c *fiber.Ctx) error  {
	// Start time
	startTime := time.Now()

	// Get the API Key from the environment
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "API key not set",
		})
	}

	// URL
	topic := "tesla"
	url := "https://newsapi.org/v2/top-headlines?" + "q=" + topic + "&apiKey=" + apiKey

	// Request duration
	startRequest := time.Now()

	// Get the data
	agent := fiber.Get(url)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	// Request duration
	requestDuration := time.Since(startRequest).Milliseconds()
	fmt.Printf("Total time taken for request to third party: %v ms\n", requestDuration)

	// DB duration
	startDB := time.Now()

	// Put data in the new struct
	var news request.News
	err := json.Unmarshal(body, &news)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	// Create with association
	newNews := models.News{
		Topic: topic,
		Status: news.Status,
		TotalResults: news.TotalResults,
		Articles: []models.Articles{},
	}
	for _, art := range news.Articles {
		newNews.Articles = append(newNews.Articles, models.Articles{
			Source:	models.Source{
				Name: art.Source.Name,
			},
			Author: art.Author,
			Title: art.Title,
			Description: art.Description,
			Url: art.Url,
			UrlToImage: art.UrlToImage,
			PublishedAt: art.PublishedAt,
			Content: art.Content,
		})
	}

	// Save to database concurrently
	go func(news models.News) {
		resultPost := database.DB.Create(&newNews)
		if resultPost.Error != nil {
			c.Status(500).JSON(fiber.Map{
				"message": "failed to make post!",
				"error": resultPost.Error.Error(),
			})
		}
	}(newNews)

	// Save to database
	// resultPost := database.DB.Create(&newNews)
	// if resultPost.Error != nil {
	// 	c.Status(500).JSON(fiber.Map{
	// 		"message": "failed to make post!",
	// 		"error": resultPost.Error.Error(),
	// 	})
	// }

	// Duration save to db
	dbDuration := time.Since(startDB).Milliseconds()
	fmt.Printf("Total time taken for database save: %v ms\n", dbDuration)

	// Total duration
	totalDuration := time.Since(startTime).Milliseconds()
	fmt.Printf("Total duration: %v ms\n", totalDuration)

	// Return the response with the created news
	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"data": newNews,
	})
}

// Make GET request from News API
func GetNewsApi(c *fiber.Ctx) (string, string, string, error) {
	// Get the API Key from the environment
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		return "", "", "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "API key not set",
		})
	}

	// URL
	url := "https://newsapi.org/v2/top-headlines?" + "q=tesla&apiKey=" + apiKey

	// Get the data
	agent := fiber.Get(url)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", "", "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	// Put data in the new struct
	var news request.News
	err := json.Unmarshal(body, &news)
	if err != nil {
		return "", "", "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	// Check if there is a data there
	if len(news.Articles) == 0 {
		return "", "", "", c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "no articles found",
		})
	}

	// Take the 'name' from Source data
	sourceData := news.Articles[0].Source.Name

	// Take the 'url' data
	urlData := news.Articles[0].Url

	// Take the 'author' data
	authorData := news.Articles[0].Author

	return sourceData, urlData, authorData, nil
}

func GetAllPosts(c *fiber.Ctx) error {
	// Retrieving all objects
	var posts []models.Post

	// Use preload to load the associatioted author, category & tags
	result := database.DB.Preload("Author").Preload("Category").Preload("Tags").Find(&posts)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "cannot retrieve posts",
			"error": result.Error.Error(),
		})
	}

	// Success retrieving data
	return c.Status(200).JSON(fiber.Map{
		"message": "posts retrieved successfully",
		"data": posts,
	})
}

func GetPostById(c *fiber.Ctx) error {
	// Take id from params
	postId := c.Params("id")

	// Retrieving object based on id
	var post models.Post
	result := database.DB.Preload("Author").Preload("Category").Preload("Tags").Find(&post, postId)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "cannot retrieve post",
			"error": result.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "post retrieved successfully",
		"data": post,
	})
}

func CreatePost(c *fiber.Ctx) error {
	// Take the data from the API
	sourceData, urlData, authorData, _ := GetNewsApi(c)

	// Make a new request
	p := new(request.CreatePost)

	// Parser request body to
	err := c.BodyParser(p);
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing",
			"error": err.Error(),
		})
	}

	// Validation
	validate := validator.New()
	errValidate := validate.Struct(p)
	if errValidate != nil {
		return errValidate
	}

	// Create with association
	newPost := models.Post{
		Title: p.Title,
		Content: p.Content,
		Comment: authorData,
		Author: models.Author{
			Name: p.Author.Name,
			Email: p.Author.Email,
		},
		Category: []models.Category{},
		Tags: []models.Tag{},
	}
	for _, cat := range p.Category {
		newPost.Category = append(newPost.Category, models.Category{
			Name:        cat.Name,
			Description: cat.Description,
		})
	}
	for _, tag := range p.Tags {
		newPost.Tags = append(newPost.Tags, models.Tag{
			Name: tag.Name,
		})
	}

	// Add a tag from request News API
	newPost.Tags = append(newPost.Tags, models.Tag{
		Name: sourceData,
	})
	newPost.Tags = append(newPost.Tags, models.Tag{
		Name: urlData,
	})

	// Save to database
	resultPost := database.DB.Create(&newPost)
	if resultPost.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to make post!",
			"error": resultPost.Error.Error(),
		})
	}

	// Return the response with the created post
	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"data": newPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	// Take id from params
	postId := c.Params("id")

	// Parser request body to a struct
	p := new(request.CreatePost)
	err := c.BodyParser(p);
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing",
			"error": err.Error(),
		})
	}

	// Validation
	validate := validator.New()
	errValidate := validate.Struct(p)
	if errValidate != nil {
		return errValidate
	}

	// Retrieve a single object with association
	var post models.Post
	result := database.DB.Preload("Author").Preload("Category").Preload("Tags").First(&post, postId)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "post not found",
			"error": result.Error.Error(),
		})
	}

	// Update data
	post.Title = p.Title
	post.Content = p.Content
	post.Comment = p.Comment
	post.Author = models.Author{
		Name: p.Author.Name,
		Email: p.Author.Email,
	}
	post.Category = []models.Category{}
	for _, cat := range p.Category {
		post.Category = append(post.Category, models.Category{
			Name: cat.Name,
			Description: cat.Description,
		})
	}
	post.Tags = []models.Tag{}
	for _, tag := range p.Tags {
		post.Tags = append(post.Tags, models.Tag{
			Name: tag.Name,
		})
	}

	//  Save updated post
	resultUpdate := database.DB.Save(&post)
	if resultUpdate.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update data",
			"error": resultUpdate.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "updated data success",
		"data": post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	// Take id from params
	postId := c.Params("id")

	// Retrieve a single object with association
	var post models.Post
	result := database.DB.Preload("Author").Preload("Category").Preload("Tags").First(&post, postId)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete data",
			"error": result.Error.Error(),
		})
	}

	// Delete data
	errDelete := database.DB.Select(clause.Associations).Delete(&post)
	if errDelete.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to delete data",
			"error": errDelete.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "success deleting data",
	})
}
