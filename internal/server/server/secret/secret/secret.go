//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package secret

import (
	"context"
	"net/http"

	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefix = "api:handler:secret"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*models.Secret, *errors.Err) {

	sm := service.NewSecretModel(ctx, envs.Get().GetStorage())
	sct, err := sm.Get(namespace, name)

	if err != nil {
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("secret").InternalServerError(err)
	}

	if sct == nil {
		err := errors.New("namespace not found")
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("secret").NotFound()
	}

	return sct, nil
}

func Apply(ctx context.Context, ns *models.Namespace, mf *request.SecretManifest) (*models.Secret, *errors.Err) {

	if mf.Meta.Name == nil {
		return nil, errors.New("secret").BadParameter("meta.name")
	}

	sct, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
	if err != nil {
		if err.Code != http.StatusText(http.StatusNotFound) {
			return nil, errors.New("secret").InternalServerError()
		}
	}

	if sct == nil {
		return Create(ctx, ns, mf)
	}

	return Update(ctx, ns, sct, mf)
}

func Create(ctx context.Context, ns *models.Namespace, mf *request.SecretManifest) (*models.Secret, *errors.Err) {

	sm := service.NewSecretModel(ctx, envs.Get().GetStorage())
	if mf.Meta.Name != nil {

		sc, err := sm.Get(ns.Meta.Name, *mf.Meta.Name)
		if err != nil {
			log.Errorf("%s:create:> get secret by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("secret").InternalServerError()
		}

		if sc != nil {
			log.Warnf("%s:create:> secret name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("secret").NotUnique("name")
		}
	}

	sct := new(models.Secret)
	sct.Meta.SetDefault()
	sct.Meta.SelfLink = *models.NewSecretSelfLink(ns.Meta.Name, *mf.Meta.Name)
	sct.Meta.Namespace = ns.Meta.Name

	mf.SetSecretMeta(sct)
	mf.SetSecretSpec(sct)

	if _, err := sm.Create(ns, sct); err != nil {
		log.Errorf("%s:create:> create secret err: %s", logPrefix, ns.Meta.Name, err.Error())
		return nil, errors.New("secret").InternalServerError()
	}

	return sct, nil
}

func Update(ctx context.Context, ns *models.Namespace, sct *models.Secret, mf *request.SecretManifest) (*models.Secret, *errors.Err) {

	sm := service.NewSecretModel(ctx, envs.Get().GetStorage())

	mf.SetSecretMeta(sct)
	mf.SetSecretSpec(sct)

	if _, err := sm.Update(sct); err != nil {
		log.Errorf("%s:update:> update secret err: %s", logPrefix, err.Error())
		return nil, errors.New("secret").InternalServerError()
	}

	return sct, nil
}
