package command

import (
	"github.com/Fraegdegjevar/Gator/internal/config"
	"github.com/Fraegdegjevar/Gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
