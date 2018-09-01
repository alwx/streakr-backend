package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"streakr-backend/pkg/services"
	"streakr-backend/pkg/utils"
)

func UserRouter(data Data) {
	users := data.Router.Group("/users")
	{
		/**
		 * @api {post} /users/login Log in
		 * @apiGroup Users
		 * @apiName Login
		 * @apiVersion 0.1.0
		 * @apiDescription Logs the user in. Returns a token that should be used in `Authorization` header for
		 * all requests except `GET /config` and `POST /users/login`
		 *
		 * @apiExample {httpie} Example usage:
		 *     http -v --json POST <URL>/users/login username=admin@test.com password=testpass
		 *
		 * @apiSuccess {String} expire Expiration date and time (the example output: `2018-07-12T20:10:50Z`)
		 * @apiSuccess {String} token Newly created authorization token.
		 * @apiError 401 Some error with credentials (check `message` field)
		 */
		users.POST("/login", data.AuthMiddleware.LoginHandler)
		/**
		 * @api {post} /users/refresh-token Update authorization token
		 * @apiHeader {String} Authorization `Bearer <TOKEN>`
		 * @apiGroup Users
		 * @apiName RefreshToken
		 * @apiVersion 0.1.0
		 * @apiDescription Returns a fresh authorization token for already authorized user.
		 *
		 * @apiExample {httpie} Example usage:
		 *     http -v --json POST <URL>/users/refresh-token "Authorization:Bearer <TOKEN>"
		 *
		 * @apiSuccess {String} expire Expiration date and time (the example output: `2018-07-12T20:10:50Z`)
		 * @apiSuccess {String} token Newly created authorization token.
		 * @apiError 401 Some error (check `message` field)
		 */
		users.POST("/refresh-token", data.AuthMiddleware.RefreshHandler)

		/**
		 * @api {post} /users Create user
		 * @apiGroup Users
		 * @apiName Register
		 * @apiVersion 0.1.0
		 * @apiDescription Registers user and returns the id of it.
		 *
		 * @apiParam {User} user JSON object that contains an information about the user.
		 * @apiExample {httpie} Example usage:
		 *     http -v --json POST <URL>/users user:='{"username": "user", "email": "user@test.com", "password": "testpass"}'
		 *
		 * @apiSuccess {String} user_id User id
		 * @apiError 400 Bad Request.
		 * @apiError 403 Forbidden. You need to have an invite to register.
		 */
		users.POST("", func(c *gin.Context) {
			var registrationData services.RegistrationData
			if err := c.ShouldBindJSON(&registrationData); err == nil {
				privateKey, publicKey, err := utils.GenKeys()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				registrationData.User.PrivateKey = string(privateKey)
				registrationData.User.PublicKey = string(publicKey)

				token, err := services.BunqInstallation(registrationData.User)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				registrationData.User.Token = token

				_, err = services.BunqDeviceServer(registrationData.User)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				_, err = services.BunqSessionServer(registrationData.User)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				/*userId, err := registrationData.User.Insert(data.Database)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}*/

				c.JSON(http.StatusCreated, gin.H{"user_id": "123"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			}
		})

		secureArea := users.Group("")
		secureArea.Use(data.AuthMiddleware.MiddlewareFunc())
		{
			/**
			 * @api {get} /users Get all users
			 * @apiHeader {String} Authorization `Bearer <TOKEN>`
			 * @apiGroup Users
			 * @apiName GetUsers
			 * @apiVersion 0.1.0
			 * @apiDescription Returns a list of all users.
			 *
			 * @apiExample {httpie} Example usage:
			 *     http -v --json GET <URL>/users "Authorization:Bearer <TOKEN>"
			 *
			 * @apiSuccess {[]User} users List of users
			 */
			secureArea.GET("", func(c *gin.Context) {
				users, err := services.GetUsers(data.Database)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"users": users})
			})
		}
	}
}