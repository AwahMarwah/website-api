package router

import "github.com/gin-contrib/cors"

var corsConfig = cors.Config{
	//AllowAllOrigins:  true,
	//AllowHeaders:     []string{"Authorization", "Content-Type"},
	//AllowMethods:     []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
	//AllowCredentials: true,

	AllowAllOrigins: true,
	AllowHeaders: []string{
		"Authorization",
		"Content-Type",
		"Content-Length",
		"Accept",
		"Accept-Encoding",
		"X-CSRF-Token",
		"X-Requested-With",
		"Origin",
		"User-Agent",
		"Referer",
		"Cache-Control",
	},
	AllowMethods: []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
		"OPTIONS", // ✅ WAJIB untuk preflight request
		"HEAD",
	},
	AllowCredentials: true,  // ✅ Tambahkan ini untuk Bearer token
	MaxAge:           86400, // Cache preflight response 24 jam
}
