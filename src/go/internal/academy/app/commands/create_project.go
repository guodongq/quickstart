package commands

import (
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/pkg/decorator"
)

type CreateProjectHandler decorator.CommandHandler[dto.CreateProject]
