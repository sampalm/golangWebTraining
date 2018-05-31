package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("cellar", func() {
	Title("The virtual wine cellar")
	Description("A simple goa service")
	Scheme("http")
	Host("localhost:8080")
})

var Bottle = MediaType("application/vnd.bottle+json", func() {
	Description("A bottle of distilled")
	Attributes(func() {
		Attribute("name", String, "Bottle name")
		Attribute("brand", String, "Bottle brand")
		Required("name", "brand")
	})
	View("default", func() {
		Attribute("name")
		Attribute("brand")
	})
})

var _ = Resource("bottle", func() {
	BasePath("/bottles")
	DefaultMedia(Bottle)

	Action("show", func() {
		Description("Get bottle by id")
		Routing(GET("/:bottleID"))
		Params(func() {
			Param("bottleID", Integer, "Bottle ID")
		})
		Response(OK)
		Response(NotFound)
	})

	Action("list", func() {
		Description("Get all bottles")
		Routing(GET("/"))
		Response(OK, CollectionOf(Bottle))
	})

	Action("create", func() {
		Routing(POST("/"))
		Payload(Bottle)
		Response(Created)
	})

	Action("delete", func() {
		Routing(DELETE("/:bottleID"))
		Params(func() {
			Param("bottleID", Integer, "Bottle ID")
		})

		Response(OK)
		Response(NotFound)
	})
})
