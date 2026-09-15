package flutter

type CommandID string

const (
	CmdUpdatePackageName      CommandID = "update-package-name"
	CmdGenerateUploadKeystore CommandID = "generate-upload-keystore"
	CmdUpdateSigningConfig    CommandID = "update-signing-config"
	CmdBuildRelease           CommandID = "build-release"
	CmdUpdateIcons            CommandID = "update-icons"
)

type PromptField struct {
	Label    string
	Secret   bool
	Default  string
	Required bool
}

type Command struct {
	ID     CommandID
	Label  string
	Flag   string
	Fields []PromptField
}

var AllCommands = []Command{
	{
		ID:    CmdUpdatePackageName,
		Label: "Update Gradle package name",
		Flag:  "update-package-name",
		Fields: []PromptField{
			{Label: "Old package name: ", Required: true},
			{Label: "New package name: ", Required: true},
		},
	},
	{
		ID:    CmdGenerateUploadKeystore,
		Label: "Generate upload keystore",
		Flag:  "generate-upload-keystore",
		Fields: []PromptField{
			{Label: "Keystore name: ", Default: "upload-keystore"},
			{Label: "Store password: ", Secret: true, Required: true},
			{Label: "Key password (blank = store password): ", Secret: true},
			{Label: "Key alias: ", Default: "upload"},
			{Label: "Certificate common name: ", Default: "Android Upload"},
		},
	},
	{
		ID:    CmdUpdateSigningConfig,
		Label: "Update signing config",
		Flag:  "update-signing-config",
	},
	{
		ID:    CmdBuildRelease,
		Label: "Build release APKs and AAB",
		Flag:  "build-release",
	},
	{
		ID:    CmdUpdateIcons,
		Label: "Update Android, iOS, and web icons",
		Flag:  "update-icons",
		Fields: []PromptField{
			{Label: "IconKitchen output path: ", Required: true},
		},
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
