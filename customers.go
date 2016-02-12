package main

import (
	"./desk"
	"./zendesk"
	"fmt"
	"strconv"
	"strings"
)

func (migrator *Migrator) checkCustomer(c *desk.Case) error {
	if strings.TrimSpace(c.Customer.Email()) == "" {
		return fmt.Errorf("Can't migrate customer. Blank Email. Aborting.")
	}
	return nil
}

func (migrator *Migrator) migrateCustomer(deskCustomer *desk.Customer) (*zendesk.Customer, error) {
	phone := ""
	if len(deskCustomer.Phones) > 0 {
		phone = deskCustomer.Phones[0].Value
	}

	name := deskCustomer.FullName()
	if len(strings.TrimSpace(deskCustomer.FullName())) == 0 {
		name = deskCustomer.Email()
	}

	return migrator.getOrCreateCustomer(&zendesk.Customer{
		ExternalId: strconv.FormatInt(deskCustomer.Id, 10),
		Name:       name,
		CreatedAt:  deskCustomer.CreatedAt,
		UpdatedAt:  deskCustomer.UpdatedAt,
		Email:      deskCustomer.Email(),
		Phone:      phone,
		Role:       "end-user",
		Verified:   true,
	})
}

func (migrator *Migrator) getOrCreateCustomer(customer *zendesk.Customer) (*zendesk.Customer, error) {

	cus, err := migrator.zenDeskClient.GetCustomerByEmail(customer.Email)
	if err != nil {
		if err == zendesk.ErrMissingEmail {
			newCustomer, err := migrator.zenDeskClient.CreateCustomer(customer)
			if err != nil {
				return nil, err
			}
			return newCustomer, nil
		}
		return nil, err
	}
	return cus, nil
}
