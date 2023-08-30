package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (ac *UserController) GetUsers(ctx *gin.Context) {
	var userResponse []models.User
	ac.DB.Find(&userResponse)

	// TODO: debug segments
	//for _, user := range userResponse {
	//	var records []models.UserSegment
	//	ac.DB.Where(models.UserSegment{UserUUID: user.UserUUID}).Find(&records)
	//	for _, record := range records {
	//		var foundSegment *models.Segment
	//		resFind := ac.DB.Where(models.Segment{SegmentID: record.SegmentID}).First(&foundSegment)
	//		if resFind.RowsAffected == 0 {
	//			// Skip, not found
	//			continue
	//		}
	//		//fmt.Println(record.SegmentID, foundSegment.SegmentSlug)
	//		userResponse.Segments = append(userResponse.Segments, foundSegment.SegmentSlug)
	//	}
	//}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"users": userResponse}})
}

func (ac *UserController) CreateUser(ctx *gin.Context) {
	var payload *models.UserPayload
	var newUser models.User
	var userResponse models.UserResponse
	var records []models.UserSegment

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	status := http.StatusCreated

	newUUID := uuid.New()
	result := ac.DB.Where(models.User{UserID: payload.UserID}).Attrs(models.User{UserUUID: newUUID}).FirstOrCreate(&newUser)
	if result.RowsAffected == 0 {
		// User exists
		status = http.StatusOK
		newUUID = newUser.UserUUID
	}
	//fmt.Printf("%+v\n", newUser)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Bad gateway"})
		return
	}

	userResponse.UserID = newUser.UserID

	// Add segments
	for _, segmentAdd := range payload.SegmentsAdd {
		// Here if no segment with this slug exists, the system skips it (it is invalid)
		var foundSegment *models.Segment
		var foundRecord *models.UserSegment
		// TODO: check if record exists
		if err := ac.DB.First(&foundSegment, models.Segment{SegmentSlug: segmentAdd}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Skip, not found
			continue
		}
		segmentUUID := foundSegment.SegmentID
		if err := ac.DB.First(&foundRecord, models.UserSegment{UserUUID: newUUID, SegmentID: segmentUUID}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			ac.DB.Create(models.UserSegment{
				UserUUID:  newUUID,
				SegmentID: segmentUUID,
			})
		}
	}

	// Delete segments
	for _, segmentDel := range payload.SegmentsDelete {
		// Here if no segment with this slug exists, the system skips it (it is invalid)
		var foundSegment *models.Segment
		if err := ac.DB.First(&foundSegment, models.Segment{SegmentSlug: segmentDel}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Skip, not found
			continue
		}
		segmentUUID := foundSegment.SegmentID
		ac.DB.Where(models.UserSegment{SegmentID: segmentUUID}).Delete(&models.UserSegment{})
	}

	ac.DB.Find(&records, models.UserSegment{UserUUID: newUUID})
	fmt.Println(len(records))

	segments := make([]string, 0)
	for _, record := range records {
		var foundSegment *models.Segment
		if err := ac.DB.First(&foundSegment, models.Segment{SegmentID: record.SegmentID}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Skip, not found
			continue
		}
		segments = append(segments, foundSegment.SegmentSlug)
	}
	userResponse.Segments = segments

	ctx.JSON(status, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *UserController) GetUserByID(ctx *gin.Context) {
	var payload *models.UserPayload
	var foundUser models.User
	var userResponse models.UserResponse
	var records []models.UserSegment

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//fmt.Printf("%+v\n", payload)

	if err := ac.DB.First(&foundUser, models.User{UserID: payload.UserID}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// User does not exist
		fmt.Println("USER DOES NOT EXIST")
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User with this ID not found"})
		return
	}

	//if result.Error != nil {
	//	ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Bad gateway"})
	//	return
	//}

	userResponse.UserID = foundUser.UserID

	resRecords := ac.DB.Where(models.UserSegment{UserUUID: foundUser.UserUUID})
	resRecords.Find(&records)

	// TODO: rows affected

	if resRecords.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
		return
	}
	fmt.Println(len(records))

	// MAYBE IT IS OK??? PROBLEM WITH DB...

	for _, record := range records {
		var foundSegment models.Segment
		if err := ac.DB.First(&foundSegment, models.Segment{SegmentID: record.SegmentID}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Skip, not found
			continue
		}
		userResponse.Segments = append(userResponse.Segments, foundSegment.SegmentSlug)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *UserController) GetRecords(ctx *gin.Context) {
	var segmentResponse []models.UserSegment
	ac.DB.Find(&segmentResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": segmentResponse}})
}
