package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"net/http"
	"strings"
)

type SegmentController struct {
	DB *gorm.DB
}

func NewSegmentController(DB *gorm.DB) SegmentController {
	return SegmentController{DB}
}

func (ac *SegmentController) CreateSegment(ctx *gin.Context) {
	var payload *models.SegmentPayload
	var foundSegment models.Segment

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.AutoAdd && (payload.UserPercentage <= 0 || payload.UserPercentage > 100) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Automatic addition impossible, invalid percentage"})
		return
	}

	newSegment := models.Segment{
		//now := time.Now()
		//SegmentID:   payload.SegmentID,
		SegmentSlug: payload.SegmentSlug,
	}
	// TODO: change to `var newSegment *models.Segment`

	newID := uuid.New()
	result := ac.DB.Where(models.Segment{SegmentSlug: payload.SegmentSlug}).Attrs(models.Segment{SegmentID: newID}).FirstOrCreate(&foundSegment)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Segment with that slug already exists"})
		return
	}
	newSegment.SegmentID = newID

	// TODO: maybe remove the following? It's not specified in swagger
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Segment with that ID already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Bad gateway"})
		return
	}

	if payload.AutoAdd {
		// Add users to this segment
		var allUsers []models.User
		ac.DB.Find(&allUsers)
		shuffledUsers := make([]uuid.UUID, len(allUsers))
		perm := rand.Perm(len(allUsers))
		for i, v := range perm {
			shuffledUsers[v] = allUsers[i].UserUUID
		}

		lastIdx := payload.UserPercentage * float64(len(allUsers)) / 100
		shuffledUsers = shuffledUsers[:int(math.Round(lastIdx))]

		for _, user := range shuffledUsers {
			ac.DB.Create(models.UserSegment{
				UserUUID:  user,
				SegmentID: foundSegment.SegmentID,
			})
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"segment": newSegment}})
}

func (ac *SegmentController) GetSegments(ctx *gin.Context) {
	var segmentResponse []models.Segment
	ac.DB.Find(&segmentResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"segments": segmentResponse}})
}

func (ac *SegmentController) DeleteSegment(ctx *gin.Context) {
	var payload *models.Segment

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	result := ac.DB.Where(models.Segment{SegmentSlug: payload.SegmentSlug})
	result.First(&payload)

	if result.Delete(&models.Segment{}).RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No segment with this slug"})
		return
	}

	ac.DB.Where(models.UserSegment{SegmentID: payload.SegmentID}).Delete(&models.UserSegment{})

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"slug": payload.SegmentSlug}})
}
