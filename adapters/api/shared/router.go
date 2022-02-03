package shared

import "github.com/emicklei/go-restful/v3"

func DefDefaultResponse(base *restful.RouteBuilder) *restful.RouteBuilder {
	base.Returns(400, "Bad Request", &ErrorResponse{})
	base.Returns(403, "Forbidden", &ErrorResponse{})
	base.Returns(404, "Not Found", &ErrorResponse{})
	base.Returns(500, "Internal Server Error", &ErrorResponse{})
	return base
}
