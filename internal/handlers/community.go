package handlers

import (
	"diplomIshi/internal/models"
	"diplomIshi/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommunityHandler struct {
	repo *repository.UserRepository
}

func NewCommunityHandler(repo *repository.UserRepository) *CommunityHandler {
	return &CommunityHandler{repo: repo}
}

func (h *CommunityHandler) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.UserID = userID

	if err := h.repo.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *CommunityHandler) GetPosts(c *gin.Context) {
	posts, err := h.repo.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Optionally enrich with comment count or user details
	type PostResponse struct {
		models.Post
		CommentCount int `json:"comment_count"`
	}
	var response []PostResponse
	for _, post := range posts {
		comments, _ := h.repo.GetComments(post.ID)
		response = append(response, PostResponse{Post: post, CommentCount: len(comments)})
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommunityHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID := c.Param("post_id")
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pID, _ := strconv.Atoi(postID)
	comment.UserID = userID
	comment.PostID = uint(pID)

	_, err := h.repo.GetPostByID(comment.PostID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := h.repo.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommunityHandler) GetComments(c *gin.Context) {
	postID := c.Param("post_id")
	pID, _ := strconv.Atoi(postID)
	comments, err := h.repo.GetComments(uint(pID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}
