# Suggested Shared OTP Generator Structure

OTP generation is a shared security utility.

Recommended location:

```text
internal/shared/otp
```

Recommended structure:

```text
internal/shared/otp
├── contract.go
├── generator.go
└── numeric.go      (optional later)
```

---

# contract.go

```go
package otp

// Generator generates temporary authentication secrets.
//
// Examples:
// - numeric OTP
// - magic link token
// - passwordless login token
// - verification challenge token
//
// Generation only creates raw secret values.
// Hashing/storage are separate responsibilities.
type Generator interface {
	Generate() (string, error)
}
```

---

# generator.go

Example implementation for secure numeric OTP generation.

```go
package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// NumericGenerator generates secure numeric OTP values.
type NumericGenerator struct {
	length int
}

func NewNumericGenerator(length int) *NumericGenerator {
	return &NumericGenerator{
		length: length,
	}
}

func (g *NumericGenerator) Generate() (string, error) {
	var builder strings.Builder

	for range g.length {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate otp digit: %w", err)
		}

		builder.WriteString(n.String())
	}

	return builder.String(), nil
}
```

---

# Recommended Usage

```go
otpValue, err := u.otpGenerator.Generate()
if err != nil {
	return fmt.Errorf("failed to generate otp: %w", err)
}
```

Then:

```go
hash, err := u.hasher.Hash(otpValue)
```

Store:

```text
hash only
```

Send:

```text
raw otp value
```

---

# Responsibility Separation

| Component      | Responsibility      |
| -------------- | ------------------- |
| otp.Generator  | create secret       |
| hasher         | protect secret      |
| challenge repo | persist secret hash |
| mailer         | transport secret    |

This separation keeps auth lifecycle very clean and scalable.

---

# Why Shared Layer?

OTP generation is NOT authentication-domain specific.

Future systems may also use it for:

* MFA
* device verification
* transaction verification
* email verification
* password reset
* merchant onboarding verification

So:

```text
otp generation = shared security utility
```

NOT:

```text
auth-only logic
```
