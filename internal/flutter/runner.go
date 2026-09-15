package flutter

import (
	"fmt"

	"martini/internal/flutter/icons"
	"martini/internal/flutter/keystore"
	"martini/internal/flutter/pkgrename"
	"martini/internal/flutter/release"
	"martini/internal/flutter/signing"
)

func RunCommands(session *Session, selected []CommandID) error {
	for _, id := range selected {
		cmd, ok := CommandByID(id)
		if !ok {
			return fmt.Errorf("unknown command: %s", id)
		}

		fmt.Printf("\nRunning: %s\n", cmd.Label)

		var err error
		switch id {
		case CmdUpdatePackageName:
			err = pkgrename.Run(session.ProjectPath, session.AnswersFor(CmdUpdatePackageName))
		case CmdGenerateUploadKeystore:
			err = keystore.Run(session.ProjectPath, session.AnswersFor(CmdGenerateUploadKeystore))
		case CmdUpdateSigningConfig:
			err = signing.Run(session.ProjectPath)
		case CmdBuildRelease:
			err = release.Run(session.ProjectPath)
		case CmdUpdateIcons:
			err = icons.Run(session.ProjectPath, session.AnswersFor(CmdUpdateIcons))
		default:
			return fmt.Errorf("unsupported command: %s", id)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
