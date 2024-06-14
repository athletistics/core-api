// Brief   : Global config constants
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package constants

const (
	LocalsUserId         = "userId"
	LocalsUser           = "user"    // deprecated, use GetContext().User instead
	LocalsHasVerified2FA = "is2faOk" // deprecated, use GetContext().Is2FAOk() instead
	LocalsHas2FA         = "has2fa"  // deprecated, use GetContext().Has2FA() instead
	LocalsContext        = "context"
	KeySecondsAgo        = "secondsAgo"
	DB_PORT_DEFAULT      = 5432
	UserSuffixLength     = 4
	TrainingEventTypeId  = 4
)

// defaults
const (
	DefaultPageSize = 25
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)
