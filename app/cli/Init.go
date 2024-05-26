package cli

import (
	"fmt"
	"gis_map_info/app/helper"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

const (
	FLAG_USER_CREATE = "user-create"
	FLAG_USER        = "user"
)

func InitCli(
	DB *gorm.DB,
) bool {
	var flag string
	app := &cli.App{
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        FLAG_USER_CREATE,
				Value:       "",
				Usage:       "create a your own user",
				Destination: &flag,
			},
			&cli.StringFlag{
				Name:        FLAG_USER,
				Value:       "",
				Usage:       "manage your own user",
				Destination: &flag,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.FlagNames() == nil {
				flag = "bypass"
				return nil
			}

			flagUsed := cCtx.FlagNames()[0]
			switch flagUsed {
			case FLAG_USER_CREATE:
				action := ""
				if cCtx.NArg() > 0 {
					action = cCtx.Args().Get(0)
				}
				userService := service.UserServiceConstruct(DB)
				password := helper.RandStringBytes(10)
				userData, err := userService.Add(service.AddPayload_UserService{
					Name:     flag,
					Username: helper.RandStringBytes(10),
					Password: &password,
					Email:    fmt.Sprint(helper.RandStringBytes(10), "@mail.com"),
					Status:   userService.GetStatus().ACTIVE,
				})
				if err != nil {
					log.Fatal("The personal access not found :(")
					return nil
				}
				displayUserGenerate(action, model.UserView{
					Uuid:     userData.Uuid,
					Name:     userData.Name,
					Username: userData.Username,
					Password: &password,
					Email:    userData.Email,
				})
			}
			fmt.Println("flag-name", cCtx.FlagNames())
			flag = "block"
			os.Exit(0)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return flag == "bypass"
}

func displayUserGenerate(action string, gg model.UserView) {
	fmt.Println("----------------------------------------------")
	if action == FLAG_USER_CREATE {
		fmt.Println("User Created:")
	}
	fmt.Println("----------------------------------------------")
	fmt.Println("Name :: ", gg.Name)
	fmt.Println("UUID :: ", gg.Uuid)
	fmt.Println("Username :: ", gg.Username)
	fmt.Println("Password  :: ", *gg.Password)
	fmt.Println("----------------------------------------------")
}
