# Migrate routing from `net/http` to `chi` in `internal/bootstrap/router.go`

This plan details the migration of the application's endpoints from standard `net/http` routing (and `apphttp.HandleMethods` multiplexer) to a more expressive, grouped routing using `go-chi/chi/v5`. 

This will also involve standardizing the endpoint structure into a more RESTful convention (using plural nouns for resources, avoiding trailing slashes, and using explicit path parameters like `/{id}`).

## User Review Required

> [!WARNING]
> **API Changes:** The proposed changes will update the endpoints to follow standard REST conventions (e.g. changing `/product` to `/products`, dropping trailing slashes for IDs). If you have frontend applications or mobile apps consuming these endpoints, they will need to be updated. If backwards compatibility is strictly required, I can keep the exact previous paths instead.

## Open Questions

> [!IMPORTANT]
> 1. Should I add a global prefix for all routes like `/api/v1`? 
> 2. Currently, handlers extract IDs manually using `strings.Split(r.URL.Path, "/")`. With `chi`, you can use `chi.URLParam(r, "id")` inside your handlers. This plan currently proposes just updating the `router.go` endpoints, meaning the previous handlers will still work without any change. Would you like me to also update the handlers to use `chi.URLParam(r, "id")`?

## Proposed Changes

### Router Migration (`router.go`)

We will utilize `chi.Router` grouping (via `r.Route`) to cleanly separate feature modules and make the routing tree readable. We will eliminate `apphttp.HandleMethods` entirely because `chi` natively supports method-based bindings (`r.Get`, `r.Post`, `r.Put`, `r.Delete`).

#### [MODIFY] router.go (file:///d:/__Projects/kage/chia.florist/service-core/internal/bootstrap/router.go)

The route definitions will be restructured like so:

```go
func NewRouter(c *Container) *chi.Mux {
	// ... (dependencies initialization remains the same)

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		
		// Products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", core(productHandler.FindProducts))
			r.Post("/", core(productHandler.CreateProduct))
			r.Get("/{id}", core(productHandler.GetProduct))
			r.Post("/images", core(productHandler.AddProductImages))
		})

		// Inventory
		r.Route("/inventories", func(r chi.Router) {
			r.Post("/", core(inventoryHandler.CreateInventory))
		})

		// Auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", core(authHandler.SignInEmail))
			r.Post("/signup", core(authHandler.SignUpAccount))
			r.Post("/verify", core(authHandler.VerifyAccount))
		})

		// Cart
		r.Route("/carts", func(r chi.Router) {
			r.Get("/", core(cartHandler.GetCart))
			r.Route("/items", func(r chi.Router) {
				r.Post("/", core(cartHandler.AddItem))
				r.Put("/", core(cartHandler.UpdateItem))
				r.Delete("/", core(cartHandler.RemoveItem))
			})
		})

		// Locations
		r.Route("/locations", func(r chi.Router) {
			r.Get("/provinces", core(locationHandler.Province))
			r.Get("/cities", core(locationHandler.City))
			r.Get("/districts", core(locationHandler.District))
			r.Get("/villages", core(locationHandler.Village))
		})

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", core(userHandler.GetUserByID))
			r.Post("/addresses", core(addressHandler.CreateUserAddress))
		})

		// Shops
		r.Route("/shops", func(r chi.Router) {
			r.Post("/", core(shopHandler.CreateShop))
			r.Get("/{id}", core(shopHandler.GetShopByID))
			
			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", core(addressHandler.ListShopAddresses))
				r.Post("/", core(addressHandler.CreateShopAddress))
				r.Get("/{id}", core(addressHandler.GetShopAddress))
			})

			r.Post("/couriers", core(courierHandler.ConfigureCourierShop))
		})

		// Payments
		r.Route("/payments", func(r chi.Router) {
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentAccount))
				r.Post("/", core(paymentHandler.CreatePaymentAccount))
			})
			r.Route("/methods", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentMethod))
				r.Post("/", core(paymentHandler.CreatePaymentMethod))
			})
		})

		// Shipping
		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", core(shipmentHandler.EstimateShippingOptions))
		})
	})

	return r
}
```

## Verification Plan

### Manual Verification
1. I will replace the contents of `router.go` with the proposed changes.
2. I will run a `go build` to make sure `router.go` correctly compiles.
3. If requested, I will also update the handler code to cleanly extract URL parameters using `chi.URLParam`.
