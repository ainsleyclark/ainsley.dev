package analytics

import (
	"context"

	"github.com/ainsleydev/website/api/_pkg/environment"

	"github.com/ainsleyclark/errors"
	"github.com/ainsleydev/website/api/_pkg/logger"
	adapter "github.com/axiomhq/axiom-go/adapters/logrus"
	"github.com/axiomhq/axiom-go/axiom"
	"github.com/sirupsen/logrus"
)

// InitAxiom initialises the Axiom client and returns a function to close it.
// Returns errors.INTERNAL the client failed to initialise.
func InitAxiom(config *environment.Config) error {
	const op = "Analytics.InitAxiom"

	if !config.IsProduction() {
		return nil
	}

	hook, err := adapter.New()
	if err != nil {
		return errors.NewInternal(err, "Could not initialise Axiom SDK", op)
	}

	err = axiom.ValidateCredentials(context.Background())
	if err != nil {
		return errors.NewInternal(err, "Validation of Axiom credentials failed", op)
	}

	logrus.RegisterExitHandler(hook.Close)
	logger.DefaultLogger.AddHook(hook)

	return nil
}
