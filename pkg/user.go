package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) UserAdd(name, password string) (*http.Response, error) {
	if name == "" || password == "" {
		return nil, errors.New("empty name or password")
	}

	user := api.CreateUserByName{Name: *api.NewNullableString(&name), Password: *api.NewNullableString(&password)}

	_, resp, err := c.client.UserAPI.CreateUserByName(c.ctx).CreateUserByName(user).Execute()
	return resp, err
}

func (c *Controller) UserDel(userID string) (*http.Response, error) {
	if userID == "" {
		return nil, errors.New("empty id")
	}

	return c.client.UserAPI.DeleteUser(c.ctx, userID).Execute()
}

func (c *Controller) UserList() (*http.Response, error) {
	users, resp, err := c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return resp, err
	}

	fmt.Printf("ID\t\t\t\t\tAdmin\tHidden\tDisabled\tLast Active\tName\n")
	fmt.Println("-----------------------------------|----------|-------|-------|---------------------|------")
	for _, user := range users {
		fmt.Printf("%s\t%t\t%t\t%t\t%s\t%s\n",
			user.GetId(),
			user.Policy.Get().GetIsAdministrator(),
			user.Policy.Get().GetIsHidden(),
			user.Policy.Get().GetIsDisabled(),
			user.GetLastActivityDate().Local().Format(time.DateTime),
			//time.Since(user.GetLastActivityDate()).Round(1*time.Minute),
			user.GetName(),
		)
	}

	return resp, err
}
