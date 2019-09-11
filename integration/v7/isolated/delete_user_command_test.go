package isolated

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("delete-user command", func() {
	Describe("help", func() {
		When("--help flag is set", func() {
			It("Displays command usage to output", func() {
				session := helpers.CF("delete-user", "--help")
				Eventually(session).Should(Say("NAME:"))
				Eventually(session).Should(Say("delete-user - Delete a user"))
				Eventually(session).Should(Say("USAGE:"))
				Eventually(session).Should(Say(`cf delete-user USERNAME \[-f\]`))
				Eventually(session).Should(Say(`cf delete-user USERNAME \[--origin ORIGIN\]`))
				Eventually(session).Should(Say("EXAMPLES:"))
				Eventually(session).Should(Say("   cf delete-user jsmith                   # internal user"))
				Eventually(session).Should(Say("   cf delete-user jsmith --origin ldap     # LDAP user"))
				Eventually(session).Should(Say("OPTIONS:"))
				Eventually(session).Should(Say(`-f\s+Prompt interactively for password`))
				Eventually(session).Should(Say(`--origin\s+Origin for mapping a user account to a user in an external identity provider`))
				Eventually(session).Should(Say("SEE ALSO:"))
				Eventually(session).Should(Say("org-users"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	When("the environment is not setup correctly", func() {
		It("fails with the appropriate errors", func() {
			helpers.CheckEnvironmentTargetedCorrectly(false, false, ReadOnlyOrg, "delete-user", "username")
		})
	})

	When("the environment is setup correctly", func() {
		var (
			someUser string
		)

		BeforeEach(func() {
			helpers.LoginCF()
			noobUser := helpers.NewUsername()
			noobPassword := helpers.NewPassword()
			session := helpers.CF("create-user", noobUser, noobPassword)
			Eventually(session).Should(Exit(0))

			helpers.LogoutCF()

			env := map[string]string{
				"CF_USERNAME": noobUser,
				"CF_PASSWORD": noobPassword,
			}
			session = helpers.CFWithEnv(env, "auth")
			Eventually(session).Should(Exit(0))
		})

		When("the logged in user is not authorized to delete new users", func() {
			It("fails with insufficient scope error", func() {
				session := helpers.CF("delete-user", someUser, "-f")
				Eventually(session).Should(Say(`Deleting user %s\.\.\.`, someUser)) // TODO: really? how about the create? using -f to get round this for now
				Eventually(session.Err).Should(Say(`You are not authorized to perform the requested action\.`))
				Eventually(session).Should(Say("FAILED"))
				Eventually(session).Should(Exit(1))
			})
		})

		When("the logged in user is authorized to delete new users", func() {
			BeforeEach(func() {
				helpers.LoginCF()
				someUser = helpers.NewUsername()
				somePassword := helpers.NewPassword()
				session := helpers.CF("create-user", someUser, somePassword)
				Eventually(session).Should(Exit(0))
			})

			When("the user to be deleted is found", func() {
				It("deletes the user", func() {
					session := helpers.CF("delete-user", someUser, "-f")
					Eventually(session).Should(Say("Deleting user %s...", someUser))
					Eventually(session).Should(Say("OK"))
					Eventually(session).Should(Exit(0))
				})
			})

			// 		When("passed invalid username", func() {
			// 			DescribeTable("when passed funkyUsername",
			// 				func(funkyUsername string) {
			// 					session := helpers.CF("delete-user", funkyUsername, helpers.NewPassword())
			// 					Eventually(session.Err).Should(Say("Username must match pattern: \\[\\\\p\\{L\\}\\+0\\-9\\+\\\\\\-_\\.@'!\\]\\+"))
			// 					Eventually(session).Should(Say("FAILED"))
			// 					Eventually(session).Should(Exit(1))
			// 				},

			// 				Entry("fails when passed an emoji", "😀"),
			// 				Entry("fails when passed a backtick", "`"),
			// 			)

			// 			When("the username is empty", func() {
			// 				It("fails with a username must be provided error", func() {
			// 					session := helpers.CF("delete-user", "", helpers.NewPassword())
			// 					Eventually(session.Err).Should(Say("A username must be provided."))
			// 					Eventually(session).Should(Say("FAILED"))
			// 					Eventually(session).Should(Exit(1))
			// 				})
			// 			})
			// 		})

			// 		When("the user passes in an origin flag", func() {
			// 			When("the origin is UAA", func() {
			// 				When("password is not present", func() {
			// 					It("errors and prints usage", func() {
			// 						newUser := helpers.NewUsername()
			// 						session := helpers.CF("delete-user", newUser, "--origin", "UAA")
			// 						Eventually(session.Err).Should(Say("Incorrect Usage: the required argument `PASSWORD` was not provided"))
			// 						Eventually(session).Should(Say("FAILED"))
			// 						Eventually(session).Should(Say("USAGE"))
			// 						Eventually(session).Should(Exit(1))
			// 					})
			// 				})
			// 			})
			// 			When("the origin is the empty string", func() {
			// 				When("password is not present", func() {
			// 					It("errors and prints usage", func() {
			// 						newUser := helpers.NewUsername()
			// 						session := helpers.CF("delete-user", newUser, "--origin", "")
			// 						Eventually(session.Err).Should(Say("Incorrect Usage: the required argument `PASSWORD` was not provided"))
			// 						Eventually(session).Should(Say("FAILED"))
			// 						Eventually(session).Should(Say("USAGE"))
			// 						Eventually(session).Should(Exit(1))
			// 					})
			// 				})
			// 			})

			// 			When("the origin is not UAA or empty", func() {
			// 				It("deletes the new user in the specified origin", func() {
			// 					newUser := helpers.NewUsername()
			// 					session := helpers.CF("delete-user", newUser, "--origin", "ldap")
			// 					Eventually(session).Should(Say("Creating user %s...", newUser))
			// 					Eventually(session).Should(Say("OK"))
			// 					Eventually(session).Should(Say("TIP: Assign roles with 'cf set-org-role' and 'cf set-space-role'"))
			// 					Eventually(session).Should(Exit(0))
			// 				})
			// 			})

			// 			When("argument for flag is not present", func() {
			// 				It("fails with incorrect usage error", func() {
			// 					session := helpers.CF("delete-user", helpers.NewUsername(), "--origin")
			// 					Eventually(session.Err).Should(Say("Incorrect Usage: expected argument for flag `--origin'"))
			// 					Eventually(session).Should(Exit(1))
			// 				})
			// 			})
			// 		})

			// 		When("the user passes in a password-prompt flag", func() {
			// 			It("prompts the user for their password", func() {
			// 				newUser := helpers.NewUsername()
			// 				buffer := NewBuffer()
			// 				_, err := buffer.Write([]byte(fmt.Sprintf("%s\n", "some-password")))
			// 				Expect(err).ToNot(HaveOccurred())
			// 				session := helpers.CFWithStdin(buffer, "delete-user", newUser, "--password-prompt")
			// 				Eventually(session).Should(Say("Password: "))
			// 				Eventually(session).Should(Say("Creating user %s...", newUser))
			// 				Eventually(session).Should(Say("OK"))
			// 				Eventually(session).Should(Say("TIP: Assign roles with 'cf set-org-role' and 'cf set-space-role'"))
			// 				Eventually(session).Should(Exit(0))
			// 			})
			// 		})

			// 		When("password is not present", func() {
			// 			It("fails with incorrect usage error", func() {
			// 				session := helpers.CF("delete-user", helpers.NewUsername())
			// 				Eventually(session.Err).Should(Say("Incorrect Usage: the required argument `PASSWORD` was not provided"))
			// 				Eventually(session).Should(Say("FAILED"))
			// 				Eventually(session).Should(Say("USAGE"))
			// 				Eventually(session).Should(Exit(1))
			// 			})
			// 		})

			// 		When("the user already exists", func() {
			// 			var (
			// 				newUser     string
			// 				newPassword string
			// 			)

			// 			BeforeEach(func() {
			// 				newUser = helpers.NewUsername()
			// 				newPassword = helpers.NewPassword()
			// 				session := helpers.CF("delete-user", newUser, newPassword)
			// 				Eventually(session).Should(Exit(0))
			// 			})

			// 			It("fails with the user already exists message", func() {
			// 				session := helpers.CF("delete-user", newUser, newPassword)
			// 				Eventually(session.Err).Should(Say("User '%s' already exists.", newUser))
			// 				Eventually(session).Should(Say("OK"))
			// 				Consistently(session).ShouldNot(Say("TIP: Assign roles with 'cf set-org-role' and 'cf set-space-role'"))
			// 				Eventually(session).Should(Exit(0))
			// 			})
			// 		})

			// 		When("the user does not already exist", func() {
			// 			It("deletes the new user", func() {
			// 				newUser := helpers.NewUsername()
			// 				newPassword := helpers.NewPassword()
			// 				session := helpers.CF("delete-user", newUser, newPassword)
			// 				Eventually(session).Should(Say("Creating user %s...", newUser))
			// 				Eventually(session).Should(Say("OK"))
			// 				Eventually(session).Should(Say("TIP: Assign roles with 'cf set-org-role' and 'cf set-space-role'"))
			// 				Eventually(session).Should(Exit(0))
			// 			})
			// 		})
		})
	})
})
