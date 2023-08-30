package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
)

type SegmentRouteController struct {
	segmentController controllers.SegmentController
}

func NewSegmentRouteController(segmentController controllers.SegmentController) SegmentRouteController {
	return SegmentRouteController{segmentController}
}

func (rc *SegmentRouteController) SegmentRoute(rg *gin.RouterGroup) {
	router := rg.Group("segment")

	router.POST("/create", rc.segmentController.CreateSegment)
	router.GET("/check", rc.segmentController.GetSegments)
	router.DELETE("/delete", rc.segmentController.DeleteSegment)
}
