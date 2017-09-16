package cmd

import "fmt"

func banner() string {
	return fmt.Sprintf(`
   __                       ____                          __  __
  / /____  ______________ _/ __/___  _________ ___  _____/ /_/ /  v%s
 / __/ _ \/ ___/ ___/ __  / /_/ __ \/ ___/ __  __ \/ ___/ __/ /   Kris Nova ⚧
/ /_/  __/ /  / /  / /_/ / __/ /_/ / /  / / / / / / /__/ /_/ /    kris@nivenly.com
\__/\___/_/  /_/   \__,_/_/  \____/_/  /_/ /_/ /_/\___/\__/_/     %s
----------------------------------------------------------------------------------------------------------

Thank you too Hashicorp and the Terraform community for making this software possible! This would not be
possible to have build without their continued hard work, and support for the open source infrastructure
community.

`, Version, GitSha)
}
