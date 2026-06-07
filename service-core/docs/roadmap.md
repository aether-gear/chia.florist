# Service Core Roadmap

This roadmap is adjusted to match the current backend shape.

- Most complete flow: `cart`
- Least complete flow: `checkout process`
- Easiest short-term unlock: `user address`
- Largest missing dependency across flows: `product -> shop availability`

## Current Lifecycle Status

### 1. Product Catalog

Goal:
Products with their available shops.

- [X] products
- [X] with inventory
- [ ] with links of main image
- ~~[ ] with shops~~

Impact:
Catalog is already good enough to list products, but it still cannot answer which shop can fulfill a product. That missing piece blocks meaningful shop selection later in cart and checkout.

### 2. Product Detail

Goal:
Product with available shops.
If logged in, include courier calculation from all available shops to the user's default address.

- [X] product
- [X] with inventory
- [ ] with links of bunch images
- [ ] with shops
- [ ] with user address
- [ ] with shop addresses
- [ ] with calculation

Impact:
This is a high-value screen because it can become the first place where shipping estimation is visible. It depends on both address readiness and shop availability readiness.

### 3. Cart

Goal:
Logged-in user sees cart products with selected shops and can switch to other available shops.

- [X] cart product (item)
- [X] with inventory
- [X] add, update, remove from cart
- [ ] with links of main image
- [ ] with shops

Impact:
This is the strongest existing flow. It is the best anchor for incremental delivery because the item lifecycle already works. Adding shop selection here will have immediate downstream value for checkout.

### 4. Checkout Process

Goal:
Logged-in user sees products with selected shop, selected courier, selected address, and can change all three.

- [X] products
- [X] with inventory
- [ ] with links of main image
- [ ] with couriers
- [ ] with shops
- [ ] with user address
- [ ] with shop addresses
- [ ] with calculation

Impact:
This is the weakest flow. It should not be treated as the first delivery target. It depends on work from product-shop mapping, addresses, and shipping calculation first.