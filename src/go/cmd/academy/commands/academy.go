package commands

import (
	generic_gin "github.com/gin-gonic/gin"
	ddd_app "github.com/guodongq/quickstart/internal/academy/app"
	ports "github.com/guodongq/quickstart/internal/academy/interfaces/http"
	genapi "github.com/guodongq/quickstart/pkg/api/genapi/server"
	"github.com/guodongq/quickstart/pkg/provider/app"
	"github.com/guodongq/quickstart/pkg/provider/gin"
	"github.com/guodongq/quickstart/pkg/provider/logger/zap"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
	"github.com/guodongq/quickstart/pkg/provider/probes"
	"github.com/guodongq/quickstart/pkg/stack"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "academy",
		Short:             "A academy example",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, _ []string) {
			st := stack.New()
			defer st.MustClose()

			loggerProvider := zap.New()
			st.MustInit(loggerProvider)

			appProvider := app.New()
			st.MustInit(appProvider)

			probesProvider := probes.New(appProvider)
			st.MustInit(probesProvider)

			mongoProvider := mongodb.New(appProvider, probesProvider)
			st.MustInit(mongoProvider)

			ginProvider := gin.New(
				gin.WithGinOptionsGinOptionFunc(func(ginEngine *generic_gin.Engine) {
					mongoRepository := mongodb.NewMongoRepository(mongoProvider)
					application := ddd_app.NewApplication(mongoRepository)
					genapi.RegisterHandlersWithOptions(ginEngine, ports.NewHttpServer(application), genapi.GinServerOptions{
						BaseURL: appProvider.ParseEndpoint(),
					})
				}),
			)
			st.MustInit(ginProvider)

			st.MustRun()
		},
	}
	return command
}
