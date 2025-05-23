package commands

import (
	_ "context"
	//"github.com/guodongq/quickstart/internal/role/ports"
	//"github.com/guodongq/quickstart/internal/role/service"
	_ "net/http"

	"github.com/guodongq/quickstart/pkg/provider/app"
	_ "github.com/guodongq/quickstart/pkg/provider/chi"
	"github.com/guodongq/quickstart/pkg/provider/logger/logrus"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
	"github.com/guodongq/quickstart/pkg/provider/probes"
	"github.com/guodongq/quickstart/pkg/stack"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "role-quickstart",
		Short:             "A role quickstart example",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, _ []string) {
			st := stack.New()
			defer st.MustClose()

			logProvider := logrus.New()
			st.MustInit(logProvider)

			appProvider := app.New()
			st.MustInit(appProvider)

			probesProvider := probes.New(appProvider)
			st.MustInit(probesProvider)

			mongoProvider := mongodb.New(appProvider, probesProvider)
			st.MustInit(mongoProvider)

			mongoRepoProvider := mongodb.NewMongoRepository(mongoProvider)
			st.MustInit(mongoRepoProvider)

			//application, cleanup := service.NewApplication(context.Background(), mongoRepoProvider)
			//defer cleanup()
			//
			//chiProvider := genericchi.New(
			//	genericchi.WithChiOptionsHandlerFromMux(func(router chi.Router) http.Handler {
			//		return genapis.HandlerFromMux(ports.NewHttpServer(application), router)
			//	}),
			//)
			//st.MustInit(chiProvider)

			st.MustRun()
		},
	}
	return command
}
