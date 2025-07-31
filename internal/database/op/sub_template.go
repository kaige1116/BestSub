package op

import "github.com/bestruirui/bestsub/internal/database/interfaces"

var subOutputTemplateRepo interfaces.SubTemplateRepository

func SubOutputTemplateRepo() interfaces.SubTemplateRepository {
	if subOutputTemplateRepo == nil {
		subOutputTemplateRepo = repo.SubTemplate()
	}
	return subOutputTemplateRepo
}
