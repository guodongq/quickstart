package commands

import (
	_ "context"
	"github.com/go-chi/chi/v5"
	//"github.com/guodongq/quickstart/internal/role/ports"
	//"github.com/guodongq/quickstart/internal/role/service"
	"github.com/guodongq/quickstart/pkg/util/provider/app"
	_ "github.com/guodongq/quickstart/pkg/util/provider/chi"
	"github.com/guodongq/quickstart/pkg/util/provider/logger/logrus"
	"github.com/guodongq/quickstart/pkg/util/provider/mongodb"
	"github.com/guodongq/quickstart/pkg/util/provider/probes"
	"github.com/guodongq/quickstart/pkg/util/stack"
	"github.com/spf13/cobra"
	_ "net/http"
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
