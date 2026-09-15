package flutter

import "fmt"

func RunCommands(session *Session, selected []CommandID) error {
	for _, id := range selected {
		cmd, ok := CommandByID(id)
		if !ok {
			return fmt.Errorf("unknown command: %s", id)
		}

		fmt.Printf("\nRunning: %s\n", cmd.Label)

		switch id {
		case CmdUpdatePackageName:
			if err := runUpdatePackageName(session); err != nil {
				return err
			}
		case CmdGenerateUploadKeystore:
			if err := runGenerateUploadKeystore(session); err != nil {
				return err
			}
		case CmdUpdateSigningConfig:
			if err := runUpdateSigningConfig(session); err != nil {
				return err
			}
		case CmdBuildRelease:
			if err := runBuildRelease(session); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported command: %s", id)
		}
	}

	return nil
}
