# `Codex` Analytics on Current Transaction

You’re making the right refactor. Current code shows transaction and context ownership is mostly in repo impl, so app/usecase layer cannot define a real unit-of-work boundary yet.

## What I found in your code
* Repository interfaces are currently context-agnostic (e.g. ProductRepository, CartRepository, AccountRepository), so usecases cannot pass request context or tx scope.
* Most repo impl methods create their own timeout using context.WithTimeout(context.Background(), 5*time.Second), repeated across almost all persistence files.
* Transactions are started inside repo methods:
    * cart save flow in cart_repository_impl.go
    * payment account allocation in payment_account_repository_impl.go
    * product image create in product_image_repository_impl.go
* Multi-step business flows are in usecases (register_customer, verify_account, add_item, add_product_images) but they are not atomic across repositories.
* Handlers call usecases without r.Context(), so cancellation/deadline/auth-bound context cannot flow to persistence.

## Why this matters (current risks)
* Partial writes are possible in business flows (example: register_customer writes user/account/challenge separately, then mail send).
* You can’t enforce cross-repo atomicity from app layer.
* Timeout/cancellation policy is fragmented and duplicated in each repo.
* One critical bug surfaced: in product_image_repository_impl.go, method Create begins a tx but executes inserts via r.db.Exec(...) instead of tx.Exec(...); commit/rollback won’t protect those inserts.
## Recommended target architecture
* App/usecase owns:
    * incoming context.Context
    * transaction boundary (Begin/Commit/Rollback)
* Repository owns:
    * SQL mapping/query logic only
* Infra provides:
    * UnitOfWork / Transactor
    * query executor abstraction (pgx.Tx or pool-backed runner)

A practical shape:
* Usecase signature: Execute(ctx context.Context, input ...)
* Repository signature: Method(ctx context.Context, q Queryer, ...) or repository bound to tx-scoped executor.
* App layer does:
    1. tx := uow.Begin(ctx)
    2. call repos with tx
    3. commit/rollback in one place

## Next steps (safe migration sequence)
1. Introduce context first (no tx yet)
    * Update usecase Execute to accept ctx.
    * Pass r.Context() from handlers.
    * Update repository interfaces + impls to accept ctx.
    * Remove context.Background() from repo methods.
2. Add minimal transactor abstraction
    * In internal/infra/db, add Transactor (WithinTransaction(ctx, fn) or explicit begin/commit).
    * Keep old repo APIs temporarily to reduce blast radius.
3. Migrate one vertical slice as pilot
    * Best first candidate: authentication/register_customer or verify_account (clear multi-repo flow).
    * Move tx orchestration to usecase/app layer.
    * Repos stop opening transactions internally.
4. Refactor remaining tx-in-repo hotspots
    * cartRepo.Save
    * paymentAccountRepo.AcquireLeastLoaded
    * productImgRepo.Create (also fix tx.Exec issue)
5. Cleanup pass
    * Remove leftover repo-internal tx begin/commit.
    * Standardize timeout policy at app boundary (or middleware), not per query.
    * Add tests for rollback behavior in multi-repo usecases.

## Suggested first PR scope (small but high value)
* Propagate ctx from handler -> usecase -> repo for one module (authentication).
* Add transactor + migrate only RegisterCustomerUsecase transaction ownership.
* Include regression tests:
    * account insert fails => user/challenge not persisted
    * challenge insert fails => user/account not persisted

<br><br>



# `ChatGPT` Suggestion for Refactoring


## Step 1: Update usecase signature

Current:

```go
func (u *RegisterCustomerUsecase) Execute(
	params RegisterCustomerParams,
) (*uuid.UUID, error)
```

New:

```go
func (u *RegisterCustomerUsecase) Execute(
	ctx context.Context,
	params RegisterCustomerParams,
) (*uuid.UUID, error)
```

Handler:

```go
challengeID, err := uc.Execute(
	r.Context(),
	params,
)
```

Simple enough. No architectural drama yet.

## Step 2: Update repository contracts

Current:

```go
type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	CreateUser(props CreateUserProps) error
}
```

Becomes:

```go
type UserRepository interface {
	GetByUsername(
		ctx context.Context,
		username string,
	) (*domain.User, error)

	CreateUser(
		ctx context.Context,
		exec transaction.Executor,
		props CreateUserProps,
	) error
}
```

Likewise:

```go
type AccountRepository interface {
	GetByEmail(
		ctx context.Context,
		email string,
	) (*domain.Account, error)

	Create(
		ctx context.Context,
		exec transaction.Executor,
		account domain.Account,
	) error
}
```

and

```go
type VerificationChallengeRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		challenge domain.VerificationChallenge,
	) error
}
```

## Step 3: Repository implementation

Current:

```go
func (r *Repository) Create(
	account domain.Account,
) error {

	_, err := r.db.Exec(
		context.Background(),
		query,
		args...,
	)

	return err
}
```

Refactor:

```go
func (r *Repository) Create(
	ctx context.Context,
	exec transaction.Executor,
	account domain.Account,
) error {

	_, err := exec.Exec(
		ctx,
		query,
		args...,
	)

	return err
}
```

Notice:

```go
r.db.Exec(...)
```

becomes

```go
exec.Exec(...)
```

That's the entire trick.

The repository no longer cares whether `exec` is:

```go
*pgxpool.Pool
```

or

```go
pgx.Tx
```

## Step 4: Inject Transactor

Usecase:

```go
type RegisterCustomerUsecase struct {
	transactor transaction.Transactor

	userRepo      UserRepository
	accountRepo   AccountRepository
	challengeRepo ChallengeRepository

	hasher Hasher
	otpGen OTPGenerator
	mailer Mailer
}
```

## Step 5: Wrap persistence in transaction

Everything before persistence remains unchanged.

You still do:

```go
user := ...
account := ...
challenge := ...
```

outside transaction.

Then:

```go
err = u.transactor.WithinTransaction(
	ctx,
	func(exec transaction.Executor) error {

		if err := u.userRepo.CreateUser(
			ctx,
			exec,
			user,
		); err != nil {
			return fmt.Errorf(
				"failed to create user: %w",
				err,
			)
		}

		if err := u.accountRepo.Create(
			ctx,
			exec,
			acc,
		); err != nil {
			return fmt.Errorf(
				"failed to create account: %w",
				err,
			)
		}

		if err := u.challengeRepo.Save(
			ctx,
			exec,
			challenge,
		); err != nil {
			return fmt.Errorf(
				"failed to create challenge: %w",
				err,
			)
		}

		return nil
	},
)

if err != nil {
	return nil, err
}
```

## Step 6: Email remains outside

After commit succeeds:

```go
mail := mailer.SendInput{
	To:      params.Email,
	Subject: "Verify your account",
	Text:    fmt.Sprintf(
		"Your OTP is %s",
		otp,
	),
}

if err := u.mailer.Send(mail); err != nil {
	return nil, fmt.Errorf(
		"failed to send otp: %w",
		err,
	)
}
```

<br><br>

# `ChatGPT` Suggestion for Which File to Refactor

Looking at the names alone, I'd split them into **3 buckets**.

## Bucket 1: Almost certainly need transaction ownership

These are the first ones I'd migrate.

### Authentication

```text
register_customer.go
verify_account.go
```

Why?

Register:

```text
Create User
Create Account
Create Challenge
```

Verify:

```text
Find Challenge
Activate Account
Delete Challenge / Mark Used
```

Multi-repository write flows.

### Product

```text
add_product_images.go
```

You already found a repository starting transactions internally.

This one is a prime candidate.

Likely:

```text
Create Image
Create Image
Create Image
```

should be atomic.

### Payment

```text
allocate_payment_method.go
select_payment_account.go
```

These names scream:

```text
Find Account
Update Account
Update Allocation
```

which usually means consistency matters.

I'd inspect these closely.

## Bucket 2: Maybe transaction

Need to read code.

### Cart

```text
add_item.go
remove_item.go
update_item.go
```

Depends heavily on implementation.

Example A:

```text
Insert CartItem
```

Single write.

No transaction needed.

Example B:

```text
Insert CartItem
Update Cart Total
Update Cart UpdatedAt
```

Multiple writes.

Transaction useful.

### Shop

```text
create_shop.go
```

If it only creates:

```text
shops
```

No.

If it creates:

```text
shop
address
merchant account
```

Yes.

### Inventory

```text
create_inventory.go
```

Need to see code.

Inventory systems often evolve into transaction-heavy code later.

### Courier

```text
configure_shop_courier.go
```

Could be:

```text
Delete existing config
Insert new config
```

which should be atomic.

## Bucket 3: Probably only context propagation

These likely don't need transaction management.

### Reads

```text
get_account.go
get_cart.go
get_product.go
get_shop.go
get_user.go

list_shop_addresses.go
list_user_addresses.go

list_payment_accounts.go
list_payment_methods.go

find_products.go
list_location.go

estimate_shipping_cost.go
get_shop_address.go
```

These should still receive:

```go
ctx context.Context
```

but almost certainly don't need:

```go
transaction.Executor
```

## My ranking for YOUR codebase

If I were doing this migration Friday night while questioning my life choices and PostgreSQL's locking semantics:

### Phase 1

```text
register_customer.go
verify_account.go
add_product_images.go
```

These are obvious wins.

### Phase 2

```text
add_item.go
update_item.go
remove_item.go

allocate_payment_method.go
select_payment_account.go

create_shop.go
create_inventory.go
```

Inspect and migrate if multiple writes exist.

### Phase 3

Everything else.

Mostly:

```go
Execute(ctx context.Context, ...)
```

and pass context through repositories.

A quick rule that works surprisingly well:

**Search for usecases that call 2+ repository write methods.**

Example:

```go
userRepo.Create(...)
accountRepo.Create(...)
challengeRepo.Save(...)
```

Transaction.

Example:

```go
productRepo.Create(...)
```

Probably not.

The moment you see:

```go
repoA.Create(...)
repoB.Update(...)
repoC.Delete(...)
```

your transaction senses should start tingling. That's usually where `Transactor` belongs.
