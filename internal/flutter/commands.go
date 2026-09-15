package flutter

type CommandID string

const (
	CmdUpdatePackageName      CommandID = "update-package-name"
	CmdGenerateUploadKeystore CommandID = "generate-upload-keystore"
	CmdUpdateSigningConfig    CommandID = "update-signing-config"
)

type Command struct {
	ID      CommandID
	Label   string
	Flag    string
	Prompts []string
}

var AllCommands = []Command{
	{
		ID:    CmdUpdatePackageName,
		Label: "Update Gradle package name",
		Flag:  "update-package-name",
		Prompts: []string{
			"Old package name: ",
			"New package name: ",
		},
	},
	{
		ID:    CmdGenerateUploadKeystore,
		Label: "Generate upload keystore",
		Flag:  "generate-upload-keystore",
	},
	{
		ID:    CmdUpdateSigningConfig,
		Label: "Update signing config",
		Flag:  "update-signing-config",
	},
}

func CommandByID(id CommandID) (Command, bool) {
	for _, c := range AllCommands {
		if c.ID == id {
			return c, true
		}
	}
	return Command{}, false
}

func CommandsFromFlags(flags map[string]bool) []CommandID {
	var ids []CommandID
	for _, c := range AllCommands {
		if flags[c.Flag] {
			ids = append(ids, c.ID)
		}
	}
	return ids
}
