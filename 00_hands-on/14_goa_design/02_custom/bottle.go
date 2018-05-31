package main

import (
	"github.com/goadesign/goa"
	"github.com/sampalm/goa/cellar/app"
)

// BottleController implements the bottle resource.
type BottleController struct {
	*goa.Controller
}

var Store = []*app.Bottle{}

// NewBottleController creates a bottle controller.
func NewBottleController(service *goa.Service) *BottleController {
	return &BottleController{Controller: service.NewController("BottleController")}
}

// Create runs the create action.
func (c *BottleController) Create(ctx *app.CreateBottleContext) error {
	Store = append(Store, &app.Bottle{
		Name:  ctx.Payload.Name,
		Brand: ctx.Payload.Brand,
	})
	return ctx.Created()
}

// Delete runs the delete action.
func (c *BottleController) Delete(ctx *app.DeleteBottleContext) error {
	if ctx.BottleID >= len(Store) {
		return ctx.NotFound()
	}
	Store = append(Store[:ctx.BottleID], Store[ctx.BottleID+1:]...)
	res := &app.Bottle{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *BottleController) List(ctx *app.ListBottleContext) error {
	return ctx.OK(Store)
}

// Show runs the show action.
func (c *BottleController) Show(ctx *app.ShowBottleContext) error {
	return ctx.OK(Store[ctx.BottleID])
}
