package op

import (
	"github.com/bestruirui/bestsub/internal/database/interfaces"
)

var repo interfaces.Repository

func SetRepo(repository interfaces.Repository) {
	repo = repository
}
func Close() error {
	return repo.Close()
}
