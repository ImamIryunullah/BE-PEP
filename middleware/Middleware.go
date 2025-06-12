package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/handlers"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Ambil token JWT dari cookie
		tokenString, err := c.Cookie(os.Getenv("COOKIE_NAME"))
		fmt.Println(tokenString)
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "StatusUnauthorized"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "StatusUnauthorized"})
			c.Abort()
			return
		}

		// Ambil ID user dan session_id dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == nil || claims["session_id"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "StatusUnauthorized"})
			c.Abort()
			return
		}

		userID := uint(claims["sub"].(float64))
		sessionID := claims["session_id"].(string)

		// Cari session berdasarkan userID dan sessionID
		var session models.SessionLogin
		if err := config.DB.Where("user_id = ? AND session_id = ?", userID, sessionID).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Cari data user dari userID
		var user models.DaftarUser
		var DataAnggota models.ParticipantRegistration
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if err := config.DB.Where("id = ?", userID).First(&DataAnggota).Error; err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			// c.Abort()
			// return
		}
		// Set user dalam konteks request
		fmt.Println(user.Email)
		// fmt.Println(user.)
		// fmt.Println(user.NamaBelakang)
		c.Set("userID", user.ID)
		c.Set("dataAnggotaID", DataAnggota.ID)
		c.Next()
	}
}

// CSRFTokenMiddleware mengecek apakah cookie "XSRF-TOKEN" ada.
// Jika tidak ada, middleware akan membuat token baru dan meng-set cookie tersebut.
func CSRFTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cek keberadaan cookie "XSRF-TOKEN"
		token, err := c.Cookie("CSRF-TOKEN")
		if err != nil || token == "" {
			// Jika cookie tidak ada, generate token baru
			token = handlers.GenerateCSRFToken() // Pastikan fungsi ini mengembalikan token (misalnya UUID)

			// Ambil pengaturan cookie dari environment
			secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
			if err != nil {
				secure = false
			}
			// httpOnly, err := strconv.ParseBool(os.Getenv("COOKIE_HTTP_ONLY"))
			// if err != nil {
			// 	httpOnly = false
			// }

			// Set cookie dengan token CSRF baru
			c.SetCookie("CSRF-TOKEN", token, 3600, os.Getenv("COOKIE_PATH"), os.Getenv("HOST"), secure, false)
			// Optionally, set token in response header juga
			c.Header("CSRF-TOKEN", token)
		}
		c.Next()
	}
}
func VerifyCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfToken := c.GetHeader("CSRF-TOKEN")
		fmt.Println(csrfToken)
		if csrfToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing"})
			c.Abort()
			return
		}
		// Ambil token dari cookie
		cookie, err := c.Cookie("CSRF-TOKEN")
		if err != nil || cookie != csrfToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Token valid, lanjutkan request
		c.Next()
	}
}

// bodyWriter digunakan untuk menangkap response body.
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Override method Write untuk menyimpan response body ke buffer.
func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type LogEntry struct {
	Time           string `json:"time"`
	RequestMethod  string `json:"request_method"`
	RequestURI     string `json:"request_uri"`
	RequestBody    string `json:"request_body"`
	ResponseStatus int    `json:"response_status"`
	Latency        string `json:"latency"`
	ResponseBody   string `json:"response_body"`
	Errors         string `json:"errors"`
}

// MonitoringMiddleware mencatat data request dan response untuk metode tertentu, dan menyimpannya ke file JSON.
func MonitoringMiddleware() gin.HandlerFunc {
	// Metode yang ingin dicatat
	allowedMethods := map[string]bool{
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"OPTIONS": true,
	}

	return func(c *gin.Context) {
		// Jika metode tidak termasuk, langsung lanjutkan handler berikutnya
		if !allowedMethods[c.Request.Method] {
			c.Next()
			return
		}

		// Tangkap request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Reset kembali request body agar bisa dibaca oleh handler selanjutnya
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		// Bungkus response writer untuk menangkap response body
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw

		// Catat waktu mulai
		startTime := time.Now()

		// Lanjutkan ke handler berikutnya
		c.Next()

		// Hitung latency dan ambil status code response
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Buat log entry
		entry := LogEntry{
			Time:           time.Now().Format(time.RFC3339),
			RequestMethod:  c.Request.Method,
			RequestURI:     c.Request.RequestURI,
			RequestBody:    string(requestBody),
			ResponseStatus: statusCode,
			Latency:        latency.String(),
			ResponseBody:   bw.body.String(),
			Errors:         c.Errors.String(),
		}

		// Nama file untuk menyimpan log
		fileName := "file.json"
		var entries []LogEntry

		// Baca file jika sudah ada isinya
		if fileData, err := ioutil.ReadFile(fileName); err == nil {
			// Jika file tidak kosong, unmarshal ke array
			if len(fileData) > 0 {
				if err := json.Unmarshal(fileData, &entries); err != nil {
					log.Printf("Error unmarshalling existing log file: %v", err)
				}
			}
		}

		// Tambahkan entry baru ke array
		entries = append(entries, entry)

		// Marshal array log dengan format yang rapi (indentasi)
		newData, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			log.Printf("Error marshalling log entries: %v", err)
			return
		}

		// Tulis ulang file dengan array log yang sudah terupdate
		if err := ioutil.WriteFile(fileName, newData, 0644); err != nil {
			log.Printf("Error writing log file: %v", err)
		}
	}
}

func StaticFileMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")
		userAgent := c.Request.Header.Get("User-Agent")

		// Pastikan referer atau origin valid
		if !isAllowed(origin) && !isAllowed(referer) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		// Cegah akses jika User-Agent kosong atau mencurigakan (opsional)
		if userAgent == "" || strings.Contains(strings.ToLower(userAgent), "crawler") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// **Tambahkan header untuk mencegah caching**
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
		c.Writer.Header().Set("ETag", "0")

		c.Next()
	}
}
func RequireAJAX() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedWith := c.GetHeader("X-Requested-With")

		if requestedWith != "XMLHttpRequest" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

var allowedOrigins = []string{
	"http://192.168.1.53:8080",
	"http://192.168.1.53:8080",
}

// Check if the request comes from an allowed origin
func isAllowed(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	for _, allowed := range allowedOrigins {
		if parsedURL.Host == strings.TrimPrefix(allowed, "http://") {
			return true
		}
	}
	return false
}
func CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya untuk file gambar
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	}
}
func ServeCompressedImage(c *gin.Context) {
	filename := c.Param("filename")
	path := filepath.Join("assets", filename)

	// Buka file
	file, err := os.Open(path)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	defer file.Close()

	// Decode gambar
	img, err := imaging.Decode(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal decode gambar")
		return
	}

	// Resize atau kompres (contoh: width max 800px)
	resized := imaging.Resize(img, 800, 0, imaging.Lanczos)

	// Kirim sebagai JPEG kualitas 80
	c.Header("Content-Type", "image/jpeg")
	c.Header("Cache-Control", "public, max-age=86400")
	_ = imaging.Encode(c.Writer, resized, imaging.JPEG, imaging.JPEGQuality(80))
}
