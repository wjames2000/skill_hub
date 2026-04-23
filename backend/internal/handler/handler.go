package handler

import "github.com/gin-gonic/gin"

func RegisterSkillRoutes(rg *gin.RouterGroup) {
	rg.GET("/skills", ListSkills)
	rg.GET("/skills/:id", GetSkill)
}

func RegisterSearchRoutes(rg *gin.RouterGroup) {
	rg.POST("/skills/search", SearchSkills)
}

func RegisterRouterRoutes(rg *gin.RouterGroup) {
	rg.POST("/router/match", MatchRouter)
	rg.POST("/router/execute", ExecuteRouter)
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", Login)
	rg.POST("/auth/register", Register)
	rg.GET("/auth/github", GitHubOAuth)
}

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/profile", GetProfile)
	rg.POST("/user/favorites", AddFavorite)
	rg.DELETE("/user/favorites/:id", RemoveFavorite)
	rg.POST("/user/reviews", AddReview)
}

func RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/sync/status", SyncStatus)
	rg.POST("/sync/trigger", TriggerSync)
	rg.GET("/stats", AdminStats)
}

func ListSkills(c *gin.Context)   { c.JSON(200, gin.H{"message": "ListSkills"}) }
func GetSkill(c *gin.Context)     { c.JSON(200, gin.H{"message": "GetSkill"}) }
func SearchSkills(c *gin.Context) { c.JSON(200, gin.H{"message": "SearchSkills"}) }

func MatchRouter(c *gin.Context)   { c.JSON(200, gin.H{"message": "MatchRouter"}) }
func ExecuteRouter(c *gin.Context) { c.JSON(200, gin.H{"message": "ExecuteRouter"}) }

func Login(c *gin.Context)      { c.JSON(200, gin.H{"message": "Login"}) }
func Register(c *gin.Context)   { c.JSON(200, gin.H{"message": "Register"}) }
func GitHubOAuth(c *gin.Context) { c.JSON(200, gin.H{"message": "GitHubOAuth"}) }

func GetProfile(c *gin.Context)      { c.JSON(200, gin.H{"message": "GetProfile"}) }
func AddFavorite(c *gin.Context)      { c.JSON(200, gin.H{"message": "AddFavorite"}) }
func RemoveFavorite(c *gin.Context)   { c.JSON(200, gin.H{"message": "RemoveFavorite"}) }
func AddReview(c *gin.Context)        { c.JSON(200, gin.H{"message": "AddReview"}) }

func SyncStatus(c *gin.Context)   { c.JSON(200, gin.H{"message": "SyncStatus"}) }
func TriggerSync(c *gin.Context)  { c.JSON(200, gin.H{"message": "TriggerSync"}) }
func AdminStats(c *gin.Context)   { c.JSON(200, gin.H{"message": "AdminStats"}) }
