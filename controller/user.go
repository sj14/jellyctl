package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) UserAdd(name, password string) error {
	if name == "" || password == "" {
		return errors.New("empty name or password")
	}

	user := api.CreateUserByName{Name: *api.NewNullableString(&name), Password: *api.NewNullableString(&password)}

	_, _, err := c.client.UserAPI.CreateUserByName(c.ctx).CreateUserByName(user).Execute()
	return err
}

func (c *Controller) UserDel(userID string) error {
	if userID == "" {
		return errors.New("empty id")
	}

	_, err := c.client.UserAPI.DeleteUser(c.ctx, userID).Execute()
	return err
}

func (c *Controller) UserList() error {
	users, _, err := c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("ID                                Admin   Disabled   Last Active       Name\n")
	fmt.Printf("---------------------------------|-------|------|---------------------|------\n")
	for _, user := range users {
		fmt.Printf("%s  %t\t  %t\t  %s  %s\n",
			user.GetId(),
			user.Policy.Get().GetIsAdministrator(),
			user.Policy.Get().GetIsDisabled(),
			user.GetLastActivityDate().Local().Format(time.DateTime),
			user.GetName(),
		)
	}

	return nil
}

func (c *Controller) userUpdatePolicy(userID string, policy api.UserPolicy) error {
	if userID == "" {
		return errors.New("empty id")
	}

	_, err := c.client.UserAPI.
		UpdateUserPolicy(c.ctx, userID).
		UserPolicy(policy).
		Execute()
	return err
}

func (c *Controller) UserPolicy(userID string) (*api.UserDtoPolicy, error) {
	user, _, err := c.client.UserAPI.GetUserById(c.ctx, userID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %w", err)
	}

	p := user.Policy.Get()

	setZeroValue(p)

	return p, nil
}

func (c *Controller) UserEnable(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsDisabled = pointer(false)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}

func (c *Controller) UserDisable(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsDisabled = pointer(true)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}

func (c *Controller) UserSetAdmin(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsAdministrator = pointer(true)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}

func (c *Controller) UserUnsetAdmin(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsAdministrator = pointer(false)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}

func (c *Controller) UserSetHidden(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsHidden = pointer(true)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}

func (c *Controller) UserUnsetHidden(userID string) error {
	policy, err := c.UserPolicy(userID)
	if err != nil {
		return err
	}

	policy.IsHidden = pointer(false)
	return c.userUpdatePolicy(userID, api.UserPolicy(*policy))
}
