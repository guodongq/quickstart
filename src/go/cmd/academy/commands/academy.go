package commands

import (
	generic_gin "github.com/gin-gonic/gin"
	application "github.com/guodongq/quickstart/internal/academy/app"
	"github.com/guodongq/quickstart/internal/academy/interfaces/http"
	openapi "github.com/guodongq/quickstart/pkg/api/codegen"
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

			appProvider := app.New(app.WithAppOptionsBasePath("/academy"))
			st.MustInit(appProvider)

			probesProvider := probes.New(appProvider)
			st.MustInit(probesProvider)

			mongoProvider := mongodb.New(appProvider, probesProvider)
			st.MustInit(mongoProvider)

			mongoRepository := mongodb.NewMongoRepository(mongoProvider, mongodb.WithDatabaseNameGetter(func() string { return "Zone" }))
			newApplication := application.NewApplication(mongoRepository)

			ginProvider := gin.New(gin.WithGinOptionsRegisterHandlers(func(engine *generic_gin.Engine) error {
				//swagger, err := openapi.GetSwagger()
				//if err != nil {
				//	return err
				//}
				//
				//options := ginmiddleware.Options{
				//	Options: openapi3filter.Options{
				//		ExcludeRequestBody:    false,
				//		ExcludeResponseBody:   false,
				//		IncludeResponseStatus: true,
				//		MultiError:            true,
				//	},
				//}
				//
				//engine.Use(ginmiddleware.OapiRequestValidatorWithOptions(swagger, &options))

				openapi.RegisterHandlersWithOptions(
					engine,
					openapi.NewStrictHandler(http.NewHttpServer(newApplication), nil),
					openapi.GinServerOptions{
						BaseURL: appProvider.ParsePath(),
					},
				)

				return nil
			}))
			st.MustInit(ginProvider)

			st.MustRun()
		},
	}
	return command
}
