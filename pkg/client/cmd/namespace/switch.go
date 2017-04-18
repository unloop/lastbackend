//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package namespace

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func SwitchCmd(name string) {

	var (
		log = c.Get().GetLogger()
	)

	namespace, err := Switch(name)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("The namespace `%s` was selected as the current", namespace.Meta.Name)
}

func Switch(name string) (*types.Namespace, error) {

	var (
		er        = new(errors.Http)
		http      = c.Get().GetHttpClient()
		storage   = c.Get().GetStorage()
		namespace = new(types.Namespace)
	)

	_, _, err := http.
		GET("/namespace/"+name).
		AddHeader("Content-Type", "application/json").
		Request(&namespace, er)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	err = storage.Set("namespace", namespace)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return namespace, nil
}