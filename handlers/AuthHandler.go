package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"
)

// var JWT_SECRET string

// Load environment variables
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("berhasil Lod ENV")
	// JWT_SECRET = os.Getenv("JWT_SECRET")

}
func GetUserDataAuth(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}
	id := userID.(uint)
	var data models.ParticipantRegistration
	var user models.DaftarUser
	// Fetch user by ID and include related ImageUser records
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data user Tidak Ditemukan"})

		return
	}
	users := gin.H{"nama_depan": user.Aset, "role": user.Role, "email": user.Email}
	c.JSON(http.StatusOK, gin.H{"user": users, "data-asset": data})
}
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var user models.DaftarUser
	if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau Password Salah!"})
		return
	}

	if !user.ComparePassword(loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau Password Salah!"})
		return
	}
	// Update LastLogin sebelum membuat sesi baru
	now := time.Now()
	config.DB.Model(&user).Update("last_login", now)

	// Check if the user has an existing session
	var existingSession models.SessionLogin
	if err := config.DB.Where("id = ?", user.ID).First(&existingSession).Error; err == nil {
		// Invalidate the existing session if found
		config.DB.Delete(&existingSession)
	}

	// Create a new session
	newSession := models.SessionLogin{
		UserID:    user.ID,
		SessionID: generateSessionID(), // Generate a unique session ID
	}
	config.DB.Create(&newSession)

	// Generate JWT
	token, err := generateJWT(user.ID, newSession.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Set Cookie JWT"})
		return
	}
	httpOnly, err := strconv.ParseBool(os.Getenv("COOKIE_HTTP_ONLY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Set Cookie JWT"})
		return
	}
	// Set token in secure cookie
	// c.SetCookie(os.Getenv("COOKIE_NAME"), token, 3600, os.Getenv("COOKIE_PATH"), os.Getenv("HOST"), secure, httpOnly)
	expirationTime := time.Now().Add(12 * time.Hour)
	// Set cookie dengan `SameSite=Lax`
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    token,
		Path:     os.Getenv("COOKIE_PATH"),
		Domain:   os.Getenv("HOST"),    // Sesuaikan dengan domain Anda
		Expires:  expirationTime,       // Berlaku 1 jam
		HttpOnly: httpOnly,             // Mencegah akses JavaScript (XSS Protection)
		Secure:   secure,               // Hanya dikirim melalui HTTPS
		SameSite: http.SameSiteLaxMode, // SameSite=Lax untuk mencegah CSRF
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func Logout(c *gin.Context) {
	// Ambil token JWT dari cookie
	tokenString, err := c.Cookie(os.Getenv("COOKIE_NAME"))
	if err != nil || tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Kembalikan kunci rahasia untuk validasi
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Ambil ID user dan session_id dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil || claims["session_id"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	userID := uint(claims["sub"].(float64))
	sessionID := claims["session_id"].(string)

	// Hapus session di database
	if err := config.DB.Where("user_id = ? AND session_id = ?", userID, sessionID).Delete(&models.SessionLogin{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		c.Abort()
		return
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Set Cookie JWT"})
		return
	}
	httpOnly, err := strconv.ParseBool(os.Getenv("COOKIE_HTTP_ONLY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Set Cookie JWT"})
		return
	}

	// c.SetCookie(os.Getenv("COOKIE_NAME"), "", -1, os.Getenv("COOKIE_PATH"), os.Getenv("HOST"), secure, httpOnly)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"), // Nama cookie yang dihapus
		Value:    "",                       // Set nilai kosong
		Path:     os.Getenv("COOKIE_PATH"), // Harus sama dengan path cookie saat login
		Domain:   os.Getenv("HOST"),        // Harus sama dengan domain saat login
		Expires:  time.Unix(0, 0),          // Set expired ke masa lalu
		MaxAge:   -1,                       // Hapus cookie segera
		HttpOnly: httpOnly,                 // Pastikan HttpOnly tetap true
		Secure:   secure,                   // Pastikan Secure tetap sesuai (HTTPS)
		SameSite: http.SameSiteLaxMode,     // Gunakan SameSite yang sama seperti saat login
	})
	// Berikan response sukses logout
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
func generateSessionID() string {
	// Implement your method to generate a unique session ID
	return fmt.Sprintf("%d-%s", time.Now().Unix(), randomString(10))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func generateJWT(userID uint, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID,
		"session_id": sessionID, // Include session_id in JWT claims
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateCSRFToken() string {
	// Generate a unique CSRF token (e.g., using UUID)
	return uuid.New().String() + randomString(500)
}
