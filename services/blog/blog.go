package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type BlogResponse struct {
	BlogID    string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	Category  string    `json:"category,omitempty"`
	Author    string    `json:"author"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateBlog(req models.CreateBlogRequest, db *gorm.DB, userId string) (BlogResponse, error) {
	var user models.User
	blog := models.Blog{
		ID:       utility.GenerateUUID(),
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userId,
		Category: req.Category,
		Image:    req.ImageURL,
	}

	err := blog.Create(db)

	if err != nil {
		return BlogResponse{}, err
	}

	user, err = user.GetUserByID(db, userId)

	if err != nil {
		return BlogResponse{}, err
	}

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageURL:  blog.Image,
		Category:  blog.Category,
		Author:    user.Name,
		AuthorID:  user.ID,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
}

func DeleteBlog(blogId string, userId string, db *gorm.DB) error {
	blog, err := CheckBlogExists(blogId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog not found")
		}
		return err
	}

	if blog.AuthorID != userId {
		return errors.New("user not authorised to delete blog")
	}

	return blog.Delete(db)
}

func GetBlogs(db *gorm.DB, c *gin.Context) ([]BlogResponse, postgresql.PaginationResponse, error) {
	var (
		blog models.Blog
		user models.User
	)
	blogs, paginationResponse, err := blog.GetAllBlogs(db, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	var responses []BlogResponse

	for _, blog := range blogs {
		userId := blog.AuthorID
		user, _ = user.GetUserByID(db, userId)
		response := BlogResponse{
			BlogID:    blog.ID,
			Title:     blog.Title,
			Content:   blog.Content,
			ImageURL:  blog.Image,
			Category:  blog.Category,
			Author:    user.Name,
			AuthorID:  user.ID,
			CreatedAt: blog.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, paginationResponse, nil
}

func GetBlogById(blogId string, db *gorm.DB) (BlogResponse, error) {
	var user models.User
	blog, err := CheckBlogExists(blogId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BlogResponse{}, errors.New("blog not found")
		}
		return BlogResponse{}, err
	}

	userId := blog.AuthorID
	user, _ = user.GetUserByID(db, userId)

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageURL:  blog.Image,
		Category:  blog.Category,
		Author:    user.Name,
		AuthorID:  user.ID,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
}

func CheckBlogExists(blogId string, db *gorm.DB) (models.Blog, error) {
	var blog models.Blog

	blog, err := blog.GetBlogById(db, blogId)
	if err != nil {
		return blog, err
	}

	return blog, nil
}
