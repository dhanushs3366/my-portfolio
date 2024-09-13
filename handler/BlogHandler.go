package handler

import (
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/db"
	"dhanushs3366/my-portfolio/utils"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func (h *Handler) createBlog(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	content := c.FormValue("content")
	if err != nil {
		return c.JSON(http.StatusNoContent, http.ErrNoCookie)
	}

	username, err := services.GetUsernameFromToken(cookie.Value)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	user, err := h.userStore.GetUser(username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	err = h.blogStore.CreateBlog(user, content)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog created")
}

func (h *Handler) editBlog(c echo.Context) error {
	blogID := c.QueryParam("ID")
	content := c.FormValue("content")

	err := h.blogStore.EditBlog(blogID, content)
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog edited")
}

func (h *Handler) deleteBlog(c echo.Context) error {
	blogID := c.QueryParam("ID")

	err := h.blogStore.DeleteBlog(blogID)

	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog Deleted")
}

func (h *Handler) getBlogs(c echo.Context) error {
	blogs, err := h.blogStore.GetBlogs()
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, blogs)
}

func (h *Handler) getBlog(c echo.Context) error {
	ID := c.Param("ID")
	blog, err := h.blogStore.GetBlogByID(ID)
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, blog)
}

func (h *Handler) debugS3Upload(c echo.Context) error {
	form, err := c.MultipartForm()

	if err != nil {
		return c.JSON(http.StatusNoContent, errors.New("no file upload"))
	}

	images := form.File["files"]
	// implement concurrent model mayhaps?
	var wg sync.WaitGroup
	wg.Add(len(images))
	errs := make(chan error, len(images))
	for _, image := range images {
		src, err := image.Open()

		if err != nil {
			log.Printf("Error %+v", err)
		}

		s3Key, err := utils.GenerateKeyForS3(image.Filename)

		if err != nil {

			log.Printf("Error %+v", err)
		}

		go func(file multipart.File, s3Key string, wg *sync.WaitGroup) {
			defer file.Close()
			defer wg.Done()
			err := h.aws.S3Upload(file, s3Key)
			if err != nil {
				log.Printf("Errors %+v", err)
				errs <- err
			} else {
				errs <- nil
			}

		}(src, s3Key, &wg)
	}

	wg.Wait()
	close(errs)

	var uploadErrors []error
	for err := range errs {
		if err != nil {
			uploadErrors = append(uploadErrors, err)
		}
	}
	if len(uploadErrors) > 0 {
		return c.JSON(http.StatusInternalServerError, uploadErrors)
	}
	return c.JSON(http.StatusOK, "All images uploaded successfully")
}
