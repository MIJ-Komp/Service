//go:build swag_admin
// +build swag_admin

package main

import (
	_ "api.mijkomp.com/docs/admin"
)

// @title       Admin API
// @version     1.0
// @description Swagger documentation for admin
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {}
