package controllers

import "github.com/lucaronca/wasa-homework/service/api/routes"

// Router defines the required methods for retrieving api routes
type Controller interface {
	Routes() routes.Routes
}
